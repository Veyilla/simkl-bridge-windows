package simkl

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestDeletedListCannotHideBehindAnotherListsActivity(t *testing.T) {
	deleted := false
	contentsA := 0
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/users/settings":
			fmt.Fprint(w, `{"account":{"id":50}}`)
		case "/sync/activities":
			fmt.Fprint(w, `{"custom_lists":{"lists":{"all":"new"}}}`)
		case "/lists/user/50":
			fmt.Fprint(w, `{"pagination":{"limit":500,"page":1,"total_items":1,"total_pages":1},"lists":[{"id":124,"user":{"id":50},"updated_at":"same"}]}`)
		case "/lists/123":
			contentsA++
			if deleted {
				w.WriteHeader(404)
				return
			}
			fmt.Fprint(w, `{"id":123,"media_type":"tv","type":"regular","user":{"id":50},"updated_at":"old","counts":{"items":0},"pagination":{"limit":500,"page":1,"total_items":0,"total_pages":0},"items":[]}`)
		case "/lists/124":
			fmt.Fprint(w, `{"id":124,"media_type":"tv","type":"regular","user":{"id":50},"updated_at":"same","counts":{"items":0},"pagination":{"limit":500,"page":1,"total_items":0,"total_pages":0},"items":[]}`)
		}
	})
	cache := NewCache(c)
	now := time.Unix(1000, 0)
	cache.now = func() time.Time { return now }
	cache.Get(context.Background(), 123, Requested)
	cache.Get(context.Background(), 124, Requested)
	deleted = true
	now = now.Add(2 * time.Minute)
	if _, err := cache.Get(context.Background(), 124, Requested); err != nil {
		t.Fatal(err)
	}
	if _, err := cache.Get(context.Background(), 123, Requested); err == nil {
		t.Fatal("Gelöschte Liste wurde erfolgreich aus Cache geliefert")
	}
	if contentsA != 2 {
		t.Fatal("Fehlende Liste wurde nicht nach ID überprüft")
	}
}
func TestMissingCountersAreNotExplicitZero(t *testing.T) {
	for _, body := range []string{`{"id":123,"media_type":"tv","counts":{},"pagination":{"limit":500,"page":1},"items":[]}`, `{"id":123,"media_type":"tv","counts":{"items":0},"pagination":{"limit":500,"page":1,"total_pages":0},"items":[]}`} {
		c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, body) })
		if _, err := c.ListItems(context.Background(), 123); err == nil {
			t.Fatal("Fehlende Zähler als erfolgreiche leere Liste akzeptiert")
		}
	}
}
func TestWrongItemTypeIsRejected(t *testing.T) {
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":123,"media_type":"movies","counts":{"items":1},"pagination":{"limit":500,"page":1,"total_items":1,"total_pages":1},"items":[{"type":"tv","ids":{"simkl_id":1,"tmdb":"77"}}]}`)
	})
	if _, err := c.ListItems(context.Background(), 123); err == nil {
		t.Fatal("TV-Eintrag als Film akzeptiert")
	}
}
