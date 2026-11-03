package diagnostics

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"simklbridge/internal/auth"
	"simklbridge/internal/config"
	"simklbridge/internal/simkl"
	"sort"
	"sync"
	"time"
)

type Report struct {
	PID     int       `json:"pid"`
	Code    string    `json:"code"`
	Updated time.Time `json:"updated"`
}
type Store struct {
	mu     sync.Mutex
	dir    string
	pid    int
	issues map[int64]string
	last   string
	logger *log.Logger
}

func New(dir string, pid int, writer io.Writer) *Store {
	if writer == nil {
		writer = io.Discard
	}
	return &Store{dir: dir, pid: pid, issues: map[int64]string{}, last: "initial", logger: log.New(writer, "", log.Ldate|log.Ltime)}
}
func file(dir string, pid int) string {
	return filepath.Join(dir, fmt.Sprintf("worker-status-%d.json", pid))
}
func Load(dir string, pid int) (Report, error) {
	b, err := os.ReadFile(file(dir, pid))
	if err != nil {
		return Report{}, err
	}
	var r Report
	if len(b) > 4096 || json.Unmarshal(b, &r) != nil || r.PID != pid {
		return Report{}, errors.New("Ungültige Statusmeldung")
	}
	return r, nil
}
func (s *Store) Record(id int64, err error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err == nil {
		delete(s.issues, id)
	} else {
		s.issues[id] = Code(err)
	}
	ids := []int64{}
	for id := range s.issues {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	code := ""
	if len(ids) > 0 {
		code = s.issues[ids[0]]
	}
	if code == s.last {
		return nil
	}
	r := Report{PID: s.pid, Code: code, Updated: time.Now().UTC()}
	b, _ := json.Marshal(r)
	if err := config.AtomicWrite(file(s.dir, s.pid), b); err != nil {
		return err
	}
	s.last = code
	if code != "" {
		s.logger.Print(Describe(code, "de"))
	}
	return nil
}
func Code(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, auth.ErrAuthentication) {
		return "authentication_required"
	}
	if errors.Is(err, auth.ErrInvalidClient) {
		return "configuration_required"
	}
	if errors.Is(err, auth.ErrScope) {
		return "scope_required"
	}
	var api *simkl.APIError
	if errors.As(err, &api) {
		if api.Status == 401 {
			return "authentication_required"
		}
		switch api.Code {
		case "premium_only":
			return "premium_required"
		case "oauth2_token_required":
			return "auth_v2_required"
		case "rate_limited":
			return "rate_limited"
		case "private_list":
			return "list_private"
		case "not_found":
			return "list_missing"
		default:
			return "simkl_error"
		}
	}
	switch err.Error() {
	case "Simkl-Anmeldung erforderlich":
		return "authentication_required"
	case "Bridge-Einstellungen prüfen":
		return "configuration_required"
	case "TVDB-Zuordnung fehlt":
		return "tvdb_missing"
	case "TMDb-Zuordnung fehlt":
		return "tmdb_missing"
	case "Simkl ist nicht erreichbar":
		return "simkl_unreachable"
	case "Medienart des Eintrags widerspricht der Liste", "Medienart passt nicht zum Zielprogramm", "Anime-Untertyp wird nicht unterstützt", "Anime-Art konnte nicht erkannt werden":
		return "media_invalid"
	case "Ungültige Simkl-Antwort", "Simkl-Liste unvollständig", "Simkl-Zähler fehlen oder sind ungültig", "Simkl-Seitendaten fehlen oder sind ungültig":
		return "response_incomplete"
	}
	return "bridge_error"
}
func Describe(code, lang string) string {
	de, en := "Bridge-Fehler; Protokoll öffnen", "Bridge error; open log"
	switch code {
	case "authentication_required":
		de, en = "Simkl-Anmeldung erforderlich", "Simkl sign-in required"
	case "configuration_required":
		de, en = "AUTH-V2-Client-ID prüfen", "Check AUTH V2 client ID"
	case "scope_required":
		de, en = "Lesende Simkl-Berechtigung erforderlich", "Read-only Simkl access required"
	case "premium_required":
		de, en = "Simkl PRO oder VIP erforderlich", "Simkl PRO or VIP required"
	case "auth_v2_required":
		de, en = "AUTH-V2-Anmeldung erforderlich", "AUTH V2 sign-in required"
	case "rate_limited":
		de, en = "Simkl-Abfragelimit erreicht", "Simkl request limit reached"
	case "list_private":
		de, en = "Kein Zugriff auf diese Liste", "No access to this list"
	case "list_missing":
		de, en = "Simkl-Liste nicht gefunden", "Simkl list not found"
	case "tvdb_missing":
		de, en = "TVDB-Zuordnung fehlt", "TVDB mapping missing"
	case "tmdb_missing":
		de, en = "TMDb-Zuordnung fehlt", "TMDb mapping missing"
	case "simkl_unreachable":
		de, en = "Simkl derzeit nicht erreichbar", "Simkl currently unreachable"
	case "media_invalid":
		de, en = "Medienart der Liste prüfen", "Check list media type"
	case "response_incomplete":
		de, en = "Simkl-Antwort unvollständig", "Incomplete Simkl response"
	case "simkl_error":
		de, en = "Simkl-Anfrage fehlgeschlagen", "Simkl request failed"
	}
	if lang == "en" {
		return en
	}
	return de
}
