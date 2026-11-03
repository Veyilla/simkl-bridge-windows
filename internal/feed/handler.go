package feed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"simklbridge/internal/config"
	"simklbridge/internal/simkl"
	"strconv"
	"strings"
	"time"
)

type ListService interface {
	Get(context.Context, int64, simkl.RefreshReason) (simkl.List, error)
	DetailIDs(context.Context, string, int64) (simkl.IDs, error)
}
type series struct {
	TVDBID int64  `json:"tvdbId"`
	Title  string `json:"title,omitempty"`
}
type movie struct {
	ID    int64  `json:"id"`
	Title string `json:"title,omitempty"`
}

func Format(ctx context.Context, s ListService, list simkl.List, destination string) (any, error) {
	if (destination == "sonarr" && list.MediaType == "movies") || (destination == "radarr" && list.MediaType == "tv") {
		return nil, errors.New("Medienart passt nicht zum Zielprogramm")
	}
	shows := []series{}
	movies := []movie{}
	seen := map[int64]bool{}
	for _, item := range list.Items {
		if err := simkl.ValidateItem(list.MediaType, item); err != nil {
			return nil, err
		}
		isMovie := list.MediaType == "movies"
		if list.MediaType == "anime" {
			if item.AnimeType == "" {
				return nil, errors.New("Anime-Art konnte nicht erkannt werden")
			}
			isMovie = item.AnimeType == "movie"
		}
		if (destination == "radarr") != isMovie {
			continue
		}
		ids := item.IDs
		if (isMovie && ids.TMDB <= 0) || (!isMovie && ids.TVDB <= 0) {
			detail, err := s.DetailIDs(ctx, list.MediaType, int64(ids.Simkl))
			if err != nil {
				return nil, err
			}
			if detail.TVDB > 0 {
				ids.TVDB = detail.TVDB
			}
			if detail.TMDB > 0 {
				ids.TMDB = detail.TMDB
			}
		}
		if isMovie {
			if ids.TMDB <= 0 {
				return nil, errors.New("TMDb-Zuordnung fehlt")
			}
			id := int64(ids.TMDB)
			if !seen[id] {
				movies = append(movies, movie{ID: id, Title: item.Title})
				seen[id] = true
			}
		} else {
			if ids.TVDB <= 0 {
				return nil, errors.New("TVDB-Zuordnung fehlt")
			}
			id := int64(ids.TVDB)
			if !seen[id] {
				shows = append(shows, series{TVDBID: id, Title: item.Title})
				seen[id] = true
			}
		}
	}
	if destination == "sonarr" {
		return shows, nil
	}
	return movies, nil
}
func Handler(service ListService, configured []config.Target) http.Handler {
	allowed := map[string]int64{}
	for _, target := range configured {
		allowed[fmt.Sprintf("/%s/%d", target.Destination, target.ID)] = target.ID
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		fail := func(code int, message string) {
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(map[string]string{"error": message})
		}
		if r.Method != "GET" {
			w.Header().Set("Allow", "GET")
			fail(405, "Nur lesender Zugriff erlaubt")
			return
		}
		host := r.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		if host != "127.0.0.1" && host != "localhost" {
			fail(403, "Nur lokaler Zugriff erlaubt")
			return
		}
		if r.URL.Path == "/healthz" {
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			return
		}
		id, ok := allowed[r.URL.Path]
		if !ok {
			fail(404, "Liste nicht eingerichtet")
			return
		}
		report := func(err error) {
			if observer, ok := service.(interface{ Report(int64, error) }); ok {
				observer.Report(id, err)
			}
		}
		list, err := service.Get(r.Context(), id, simkl.Requested)
		if err != nil {
			report(err)
			fail(503, "Simkl-Liste derzeit nicht abrufbar; Anmeldung und Protokoll prüfen")
			return
		}
		destination := strings.Split(r.URL.Path, "/")[1]
		out, err := Format(r.Context(), service, list, destination)
		if err != nil {
			report(err)
			fail(422, "Liste konnte nicht vollständig zugeordnet werden")
			return
		}
		report(nil)
		json.NewEncoder(w).Encode(out)
	})
}
func URL(port int, target config.Target) string {
	return fmt.Sprintf("http://127.0.0.1:%d/%s/%d", port, target.Destination, target.ID)
}
func Listen(port int) (net.Listener, error) {
	if port < 0 || port > 65535 {
		return nil, errors.New("Ungültiger Port")
	}
	l, err := net.Listen("tcp4", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
	if err != nil {
		return nil, errors.New("Bridge-Port ist belegt oder nicht verfügbar")
	}
	return l, nil
}
func Serve(ctx context.Context, l net.Listener, h http.Handler) error {
	server := &http.Server{Handler: h, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 2 * time.Minute, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 16 << 10}
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(shutdown)
		case <-done:
		}
	}()
	err := server.Serve(l)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
