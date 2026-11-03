package simkl

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheOnlyDownloadsChangedContents(t *testing.T) {
	var contents, indices, activities int
	version := "2026-01-01T00:00:00Z"
	fail := false
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/users/settings":
			fmt.Fprint(w, `{"account":{"id":50}}`)
		case "/sync/activities":
			activities++
			fmt.Fprintf(w, `{"custom_lists":{"lists":{"all":%q}}}`, version)
		case "/lists/user/50":
			indices++
			fmt.Fprintf(w, `{"pagination":{"limit":500,"page":1,"total_items":1,"total_pages":1},"lists":[{"id":123,"media_type":"tv","user":{"id":50},"updated_at":%q}]}`, version)
		case "/lists/123":
			contents++
			if fail {
				w.WriteHeader(500)
				return
			}
			fmt.Fprintf(w, `{"id":123,"media_type":"tv","type":"regular","user":{"id":50},"updated_at":%q,"counts":{"items":1},"pagination":{"limit":500,"page":1,"total_items":1,"total_pages":1},"items":[{"type":"tv","ids":{"simkl_id":1,"tvdb":"101"}}]}`, version)
		default:
			t.Error("Unerwarteter Abruf", r.URL.Path)
		}
	})
	cache := NewCache(c)
	now := time.Unix(1000, 0)
	cache.now = func() time.Time { return now }
	get := func() {
		t.Helper()
		if _, err := cache.Get(context.Background(), 123, Requested); err != nil {
			t.Fatal(err)
		}
	}
	get()
	now = now.Add(2 * time.Minute)
	get()
	now = now.Add(2 * time.Minute)
	get()
	if contents != 1 || activities != 2 || indices != 1 {
		t.Fatalf("Unverändert: Inhalte=%d Aktivitäten=%d Index=%d", contents, activities, indices)
	}
	version = "2026-01-02T00:00:00Z"
	now = now.Add(2 * time.Minute)
	fail = true
	if _, err := cache.Get(context.Background(), 123, Requested); err == nil {
		t.Fatal("Fehler als Liste ausgeliefert")
	}
	fail = false
	get()
	if contents != 3 {
		t.Fatalf("Fehlgeschlagenen Stand nicht erneut geladen: %d", contents)
	}
}
func TestConcurrentRequestsShareDownload(t *testing.T) {
	var calls atomic.Int32
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		fmt.Fprint(w, `{"id":123,"media_type":"tv","type":"regular","user":{"id":50},"counts":{"items":0},"pagination":{"limit":500,"page":1,"total_items":0,"total_pages":0},"items":[]}`)
	})
	cache := NewCache(c)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Go(func() {
			if _, err := cache.Get(context.Background(), 123, Requested); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if calls.Load() != 1 {
		t.Fatalf("%d Inhaltsabrufe", calls.Load())
	}
}
func TestAutoListWithoutTimestampDownloadsAtMostDaily(t *testing.T) {
	contents := 0
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/users/settings":
			fmt.Fprint(w, `{"account":{"id":50}}`)
		case "/sync/activities":
			fmt.Fprint(w, `{"custom_lists":{"lists":{"all":"same"}}}`)
		case "/lists/user/50":
			fmt.Fprint(w, `{"pagination":{"limit":500,"page":1,"total_items":1,"total_pages":1},"lists":[{"id":123,"type":"auto","media_type":"tv","user":{"id":50}}]}`)
		case "/lists/123":
			contents++
			fmt.Fprint(w, `{"id":123,"media_type":"tv","type":"auto","user":{"id":50},"counts":{"items":0},"pagination":{"limit":500,"page":1,"total_items":0,"total_pages":0},"items":[]}`)
		}
	})
	cache := NewCache(c)
	now := time.Unix(1000, 0)
	cache.now = func() time.Time { return now }
	for i := 0; i < 3; i++ {
		if _, err := cache.Get(context.Background(), 123, Requested); err != nil {
			t.Fatal(err)
		}
		now = now.Add(2 * time.Minute)
	}
	if contents != 1 {
		t.Fatalf("%d Abrufe innerhalb eines Tages", contents)
	}
	now = now.Add(24 * time.Hour)
	if _, err := cache.Get(context.Background(), 123, Requested); err != nil {
		t.Fatal(err)
	}
	if contents != 2 {
		t.Fatal("Tägliche Aktualisierung fehlt")
	}
}
