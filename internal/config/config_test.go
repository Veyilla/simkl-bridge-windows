package config

import (
	"path/filepath"
	"testing"
)

func TestSaveLoadPreservesSettings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	s := Store{Path: path}
	want := Default()
	want.Language = "en"
	want.Lists = []Target{{ID: 123, Destination: "sonarr"}}
	if err := s.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Language != "en" || got.Port != 8787 || len(got.Lists) != 1 || got.Lists[0].ID != 123 {
		t.Fatalf("Einstellungen verändert: %+v", got)
	}
	want.Port = 9000
	if err := s.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err = s.Load()
	if err != nil || got.Port != 9000 {
		t.Fatalf("Aktualisieren: %+v %v", got, err)
	}
}
func TestRejectInvalidSettings(t *testing.T) {
	for _, mutate := range []func(*Config){func(c *Config) { c.Port = 0 }, func(c *Config) { c.Language = "xx" }, func(c *Config) { c.ClientID = "ungueltig" }, func(c *Config) { c.Lists = []Target{{ID: -1, Destination: "radarr"}} }, func(c *Config) { c.Lists = []Target{{ID: 12, Destination: "other"}} }} {
		c := Default()
		mutate(&c)
		if c.Validate() == nil {
			t.Fatalf("Ungültige Einstellung akzeptiert: %+v", c)
		}
	}
}
