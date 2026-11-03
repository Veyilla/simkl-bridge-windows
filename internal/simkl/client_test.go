package simkl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type testAuth struct {
	token     string
	refreshes int
}

func (a *testAuth) AccessToken(context.Context) (string, error) { return a.token, nil }
func (a *testAuth) Refresh(context.Context, string) (string, error) {
	a.refreshes++
	a.token = "refreshed"
	return a.token, nil
}
func testClient(t *testing.T, h http.HandlerFunc) (*Client, *httptest.Server, *testAuth) {
	t.Helper()
	s := httptest.NewServer(h)
	t.Cleanup(s.Close)
	a := &testAuth{token: "synthetic"}
	c := New("client-for-local-test", s.URL, s.Client(), a)
	c.wait = func(context.Context, time.Duration) error { return nil }
	return c, s, a
}
func TestReadsEveryPageAndExternalStringIDs(t *testing.T) {
	calls := 0
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Header.Get("Authorization") != "Bearer synthetic" || r.URL.Query().Get("client_id") == "" || r.URL.Query().Get("app-name") == "" || r.Header.Get("User-Agent") == "" {
			t.Error("Pflichtangaben fehlen")
		}
		page := r.URL.Query().Get("page")
		if page == "1" {
			w.Write([]byte(`{"id":123,"name":"Test","media_type":"tv","type":"regular","user":{"id":50},"updated_at":"2026-01-01T00:00:00Z","counts":{"items":2},"pagination":{"page":1,"limit":1,"total_items":2,"total_pages":2},"items":[{"title":"A","type":"tv","ids":{"simkl_id":1,"tvdb":"101"}}]}`))
		} else {
			w.Write([]byte(`{"id":123,"name":"Test","media_type":"tv","type":"regular","user":{"id":50},"updated_at":"2026-01-01T00:00:00Z","counts":{"items":2},"pagination":{"page":2,"limit":1,"total_items":2,"total_pages":2},"items":[{"title":"B","type":"tv","ids":{"simkl_id":2,"tvdb":"102"}}]}`))
		}
	})
	list, err := c.ListItems(context.Background(), 123)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 || len(list.Items) != 2 || list.Items[1].IDs.TVDB != 102 {
		t.Fatalf("Unvollständige Liste: %+v", list)
	}
}
func TestErrorsNeverBecomeEmptyLists(t *testing.T) {
	for _, body := range []string{`{"error":"premium_only","message":"secret-server-data"}`, `{}`, `{"id":123,"media_type":"tv","counts":{"items":2},"pagination":{"limit":500,"page":1,"total_items":2,"total_pages":1},"items":[]}`, `{"id":123,"media_type":"tv","counts":{"items":0},"pagination":{"limit":500,"page":1,"total_items":0,"total_pages":0}}`} {
		t.Run(body[:min(len(body), 30)], func(t *testing.T) {
			c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(body)) })
			_, err := c.ListItems(context.Background(), 123)
			if err == nil || strings.Contains(err.Error(), "secret-server-data") {
				t.Fatalf("Ungültige Antwort akzeptiert: %v", err)
			}
		})
	}
}
func TestRealEmptyListIsValid(t *testing.T) {
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":123,"media_type":"tv","counts":{"items":0},"pagination":{"limit":500,"page":1,"total_items":0,"total_pages":0},"items":[]}`))
	})
	list, err := c.ListItems(context.Background(), 123)
	if err != nil || len(list.Items) != 0 {
		t.Fatalf("Leere Liste: %+v %v", list, err)
	}
}
func TestUnauthorizedRefreshesOnlyOnce(t *testing.T) {
	calls := 0
	c, _, a := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(401)
		w.Write([]byte(`{"error":"user_token_required"}`))
	})
	_, err := c.ListItems(context.Background(), 123)
	if err == nil || calls != 2 || a.refreshes != 1 {
		t.Fatalf("Abrufe=%d Erneuerungen=%d Fehler=%v", calls, a.refreshes, err)
	}
}
func TestChangingPagesAreRejected(t *testing.T) {
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		fmt.Fprintf(w, `{"id":123,"media_type":"tv","updated_at":"version-%s","counts":{"items":2},"pagination":{"limit":500,"page":%s,"total_items":2,"total_pages":2},"items":[{"type":"tv","ids":{"simkl_id":%s}}]}`, page, page, page)
	})
	if _, err := c.ListItems(context.Background(), 123); err == nil {
		t.Fatal("Während des Abrufs geänderte Liste akzeptiert")
	}
}
func TestRateLimitHonorsRetryAfter(t *testing.T) {
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(429)
		w.Write([]byte(`{"error":"rate_limit_exceeded"}`))
	})
	c.ListItems(context.Background(), 123)
	var waited time.Duration
	c.wait = func(ctx context.Context, d time.Duration) error { waited = d; return context.Canceled }
	_, err := c.ListItems(context.Background(), 123)
	if !errors.Is(err, context.Canceled) || waited < 29*time.Second {
		t.Fatalf("Wartezeit %s, Fehler %v", waited, err)
	}
}
func TestDetailLookupUsesCorrectCatalog(t *testing.T) {
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/anime/77" || r.URL.Query().Get("extended") != "full" {
			t.Error("Falscher Detailabruf", r.URL)
		}
		json.NewEncoder(w).Encode(map[string]any{"ids": map[string]any{"simkl": 77, "tvdb": "999"}})
	})
	ids, err := c.DetailIDs(context.Background(), "anime", 77)
	if err != nil || ids.TVDB != 999 {
		t.Fatalf("Kennungen: %+v %v", ids, err)
	}
}
