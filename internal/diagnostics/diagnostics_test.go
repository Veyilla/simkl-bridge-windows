package diagnostics

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestStatusKeepsOtherListsFailureAndSanitizesUnknownErrors(t *testing.T) {
	dir := t.TempDir()
	var log bytes.Buffer
	s := New(dir, 1234, &log)
	if err := s.Record(10, errors.New("synthetic-secret")); err != nil {
		t.Fatal(err)
	}
	r, err := Load(dir, 1234)
	if err != nil || r.Code != "bridge_error" || strings.Contains(log.String(), "synthetic-secret") {
		t.Fatal(r, err, log.String())
	}
	s.Record(20, errors.New("TVDB-Zuordnung fehlt"))
	s.Record(10, nil)
	r, err = Load(dir, 1234)
	if err != nil || r.Code != "tvdb_missing" {
		t.Fatal("Fehler einer anderen Liste ging verloren", r, err)
	}
	s.Record(20, nil)
	r, err = Load(dir, 1234)
	if err != nil || r.Code != "" {
		t.Fatal("Behobener Fehler blieb bestehen", r, err)
	}
}
