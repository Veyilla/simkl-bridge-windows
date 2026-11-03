package feed

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"simklbridge/internal/config"
	"simklbridge/internal/simkl"
	"strings"
	"testing"
)

type testSource struct {
	list    simkl.List
	err     error
	details int
}

func (s *testSource) Get(context.Context, int64, simkl.RefreshReason) (simkl.List, error) {
	return s.list, s.err
}
func (s *testSource) DetailIDs(context.Context, string, int64) (simkl.IDs, error) {
	s.details++
	return simkl.IDs{TVDB: 999, TMDB: 888}, nil
}
func TestFeedFormatsAndMapping(t *testing.T) {
	for _, v := range []struct {
		dest, media, want string
		ids               simkl.IDs
	}{{"sonarr", "tv", `[{"tvdbId":101,"title":"Test"}]`, simkl.IDs{Simkl: 1, TVDB: 101}}, {"radarr", "movies", `[{"id":202,"title":"Test"}]`, simkl.IDs{Simkl: 2, TMDB: 202}}, {"sonarr", "tv", `[{"tvdbId":999,"title":"Test"}]`, simkl.IDs{Simkl: 3}}} {
		s := &testSource{list: simkl.List{ListMeta: simkl.ListMeta{ID: 123, MediaType: v.media}, Items: []simkl.Item{{Title: "Test", Type: map[string]string{"tv": "tv", "movies": "movie"}[v.media], IDs: v.ids}}}}
		h := Handler(s, []config.Target{{ID: 123, Destination: v.dest}})
		r := httptest.NewRecorder()
		h.ServeHTTP(r, httptest.NewRequest("GET", "http://127.0.0.1/"+v.dest+"/123", nil))
		if r.Code != 200 || strings.TrimSpace(r.Body.String()) != v.want {
			t.Fatalf("%s: %d %s", v.dest, r.Code, r.Body.String())
		}
		if v.ids.Simkl == 3 && s.details != 1 {
			t.Fatal("Fehlende Kennung nicht nachgeladen")
		}
	}
}
func TestFeedErrorsAreNotEmptySuccess(t *testing.T) {
	for _, v := range []struct {
		path, media string
		err         error
		want        int
	}{{"/sonarr/999", "tv", nil, 404}, {"/radarr/123", "tv", nil, 404}, {"/sonarr/123", "movies", nil, 422}, {"/sonarr/123", "tv", errors.New("synthetic-secret"), 503}} {
		s := &testSource{list: simkl.List{ListMeta: simkl.ListMeta{ID: 123, MediaType: v.media}}, err: v.err}
		r := httptest.NewRecorder()
		Handler(s, []config.Target{{ID: 123, Destination: "sonarr"}}).ServeHTTP(r, httptest.NewRequest("GET", "http://127.0.0.1"+v.path, nil))
		if r.Code != v.want || strings.Contains(r.Body.String(), "synthetic-secret") {
			t.Fatalf("Antwort %d: %s", r.Code, r.Body.String())
		}
	}
}
func TestEmptyListAndMethods(t *testing.T) {
	s := &testSource{list: simkl.List{ListMeta: simkl.ListMeta{ID: 123, MediaType: "tv"}, Items: []simkl.Item{}}}
	h := Handler(s, []config.Target{{ID: 123, Destination: "sonarr"}})
	r := httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest("GET", "http://127.0.0.1/sonarr/123", nil))
	if r.Code != 200 || strings.TrimSpace(r.Body.String()) != "[]" {
		t.Fatalf("Leere Liste: %s", r.Body.String())
	}
	r = httptest.NewRecorder()
	h.ServeHTTP(r, httptest.NewRequest("POST", "http://127.0.0.1/sonarr/123", nil))
	if r.Code != http.StatusMethodNotAllowed {
		t.Fatal(r.Code)
	}
}
func TestLoopbackListenerAndBusyPort(t *testing.T) {
	l, err := Listen(0)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if !strings.HasPrefix(l.Addr().String(), "127.0.0.1:") {
		t.Fatal(l.Addr())
	}
	port := l.Addr().(*net.TCPAddr).Port
	if other, err := Listen(port); err == nil {
		other.Close()
		t.Fatal("Belegter Port akzeptiert")
	}
}
