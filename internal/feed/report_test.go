package feed

import (
	"errors"
	"net/http/httptest"
	"simklbridge/internal/config"
	"simklbridge/internal/simkl"
	"testing"
)

type reportingSource struct {
	testSource
	reported error
	id       int64
}

func (s *reportingSource) Report(id int64, err error) { s.id = id; s.reported = err }
func TestMappingFailureReachesStatusReporter(t *testing.T) {
	s := &reportingSource{testSource: testSource{list: simkl.List{ListMeta: simkl.ListMeta{ID: 123, MediaType: "movies"}, Items: []simkl.Item{{Type: "tv", IDs: simkl.IDs{Simkl: 1, TMDB: 42}}}}}}
	r := httptest.NewRecorder()
	Handler(s, []config.Target{{ID: 123, Destination: "radarr"}}).ServeHTTP(r, httptest.NewRequest("GET", "http://127.0.0.1/radarr/123", nil))
	if r.Code != 422 || s.id != 123 || s.reported == nil {
		t.Fatalf("Zuordnungsfehler ohne Statusmeldung: %d %v", r.Code, s.reported)
	}
	s.err = errors.New("synthetic-network")
	r = httptest.NewRecorder()
	Handler(s, []config.Target{{ID: 123, Destination: "radarr"}}).ServeHTTP(r, httptest.NewRequest("GET", "http://127.0.0.1/radarr/123", nil))
	if r.Code != 503 || s.reported == nil {
		t.Fatal("Abruffehler ohne Statusmeldung")
	}
}
