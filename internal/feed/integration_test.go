package feed

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"simklbridge/internal/config"
	"simklbridge/internal/simkl"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// Die API und beide Zielprogramme werden ausschließlich durch lokale Testserver ersetzt.
// Client, Seitenauswertung, Zwischenspeicher und Ausgabe sind die echten Komponenten.
type integrationTokens struct{}

func (integrationTokens) AccessToken(context.Context) (string, error) { return "test-access", nil }
func (integrationTokens) Refresh(context.Context, string) (string, error) {
	return "test-refreshed", nil
}

func integrationBridge(t *testing.T, api http.HandlerFunc, targets []config.Target) (*httptest.Server, *simkl.Cache) {
	t.Helper()
	upstream := httptest.NewServer(api)
	t.Cleanup(upstream.Close)
	cache := simkl.NewCache(simkl.New("test-client", upstream.URL, upstream.Client(), integrationTokens{}))
	bridge := httptest.NewServer(Handler(cache, targets))
	t.Cleanup(bridge.Close)
	return bridge, cache
}

func TestIntegrationPaginatedAnimeFeedsBothDestinations(t *testing.T) {
	var pages atomic.Int32
	bridge, _ := integrationBridge(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/lists/123" || r.Header.Get("Authorization") != "Bearer test-access" {
			t.Errorf("Unerwarteter API-Abruf: %s", r.URL.Path)
			http.Error(w, "unexpected request", 400)
			return
		}
		pages.Add(1)
		var items string
		switch r.URL.Query().Get("page") {
		case "1":
			items = `[{"title":"Serie A","type":"anime","anime_type":"tv","ids":{"simkl_id":1,"tvdb":"101"}},{"title":"Film A","type":"anime","anime_type":"movie","ids":{"simkl_id":2,"tmdb":"202"}}]`
		case "2":
			// Zwei Simkl-Einträge mit derselben Zielkennung dürfen keinen doppelten Import erzeugen.
			items = `[{"title":"Serie A erneut","type":"anime","anime_type":"ova","ids":{"simkl_id":3,"tvdb":"101"}},{"title":"Film B","type":"anime","anime_type":"movie","ids":{"simkl_id":4,"tmdb":"203"}}]`
		default:
			t.Errorf("Unerwartete Seite: %s", r.URL.RawQuery)
			http.Error(w, "unexpected page", 400)
			return
		}
		fmt.Fprintf(w, `{"id":123,"media_type":"anime","type":"regular","user":{"id":50},"counts":{"items":4},"pagination":{"limit":2,"page":%s,"total_items":4,"total_pages":2},"items":%s}`, r.URL.Query().Get("page"), items)
	}, []config.Target{{ID: 123, Destination: "sonarr"}, {ID: 123, Destination: "radarr"}})

	// Gleichzeitige Abrufe beider Ziele müssen denselben vollständigen Download teilen.
	var requests sync.WaitGroup
	for i := 0; i < 12; i++ {
		requests.Go(func() {
			for _, target := range []struct{ path, want string }{
				{"/sonarr/123", `[{"tvdbId":101,"title":"Serie A"}]`},
				{"/radarr/123", `[{"id":202,"title":"Film A"},{"id":203,"title":"Film B"}]`},
			} {
				res, err := bridge.Client().Get(bridge.URL + target.path)
				if err != nil {
					t.Error(err)
					return
				}
				var body strings.Builder
				_, err = io.Copy(&body, res.Body)
				res.Body.Close()
				if err != nil || res.StatusCode != 200 || strings.TrimSpace(body.String()) != target.want {
					t.Errorf("%s: Status=%d Antwort=%s Fehler=%v", target.path, res.StatusCode, body.String(), err)
				}
			}
		})
	}
	requests.Wait()
	if pages.Load() != 2 {
		t.Fatalf("Erwartet: zwei Seiten insgesamt; erhalten: %d", pages.Load())
	}
}

func TestIntegrationUpstreamFailuresCannotBecomeSuccessfulImports(t *testing.T) {
	for _, example := range []struct {
		name, body string
		status     int
	}{
		{"PRO fehlt", `{"error":"premium_only","message":"test-private-detail"}`, 503},
		{"Liste privat", `{"error":"private_list","message":"test-private-detail"}`, 503},
		{"Liste unvollständig", `{"id":123,"media_type":"tv","counts":{"items":2},"pagination":{"limit":500,"page":1,"total_items":2,"total_pages":1},"items":[{"type":"tv","ids":{"simkl_id":1,"tvdb":101}}]}`, 503},
		{"Zuordnung fehlt", `{"id":123,"media_type":"tv","counts":{"items":1},"pagination":{"limit":500,"page":1,"total_items":1,"total_pages":1},"items":[{"type":"tv","ids":{"simkl_id":1}}]}`, 422},
	} {
		t.Run(example.name, func(t *testing.T) {
			bridge, _ := integrationBridge(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/tv/1" {
					fmt.Fprint(w, `{"ids":{"simkl_id":1}}`)
					return
				}
				if r.URL.Path != "/lists/123" {
					t.Errorf("Unerwarteter API-Abruf: %s", r.URL.Path)
				}
				fmt.Fprint(w, example.body)
			}, []config.Target{{ID: 123, Destination: "sonarr"}})
			res, err := bridge.Client().Get(bridge.URL + "/sonarr/123")
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			var body strings.Builder
			_, err = io.Copy(&body, res.Body)
			if err != nil || res.StatusCode != example.status || !strings.Contains(body.String(), `"error"`) || strings.Contains(body.String(), "test-private-detail") {
				t.Fatalf("Status=%d Antwort=%s Fehler=%v", res.StatusCode, body.String(), err)
			}
		})
	}
}
