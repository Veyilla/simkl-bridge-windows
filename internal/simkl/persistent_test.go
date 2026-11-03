package simkl

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"
)

func TestCacheSurvivesRestartWithoutContentDownload(t *testing.T) {
	calls := 0
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		fmt.Fprint(w, `{"id":123,"media_type":"tv","type":"regular","user":{"id":50},"counts":{"items":0},"pagination":{"limit":500,"page":1,"total_items":0,"total_pages":0},"items":[]}`)
	})
	path := filepath.Join(t.TempDir(), "lists.json")
	first, err := NewPersistentCache(c, path, 50)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = first.Get(context.Background(), 123, Requested); err != nil {
		t.Fatal(err)
	}
	second, err := NewPersistentCache(c, path, 50)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = second.Get(context.Background(), 123, Requested); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("Neustart verursachte %d Inhaltsabrufe", calls)
	}
	other, err := NewPersistentCache(c, path, 51)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = other.Get(context.Background(), 123, Requested); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatal("Cache eines anderen Kontos wurde übernommen")
	}
}
func TestDetailMappingsAreCached(t *testing.T) {
	calls := 0
	c, _, _ := testClient(t, func(w http.ResponseWriter, r *http.Request) { calls++; fmt.Fprint(w, `{"ids":{"tvdb":"101"}}`) })
	cache := NewCache(c)
	for i := 0; i < 2; i++ {
		ids, err := cache.DetailIDs(context.Background(), "tv", 1)
		if err != nil || ids.TVDB != 101 {
			t.Fatal(ids, err)
		}
	}
	if calls != 1 {
		t.Fatalf("%d identische Detailabrufe", calls)
	}
}
