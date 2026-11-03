package config

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Target struct {
	ID          int64  `json:"id"`
	Name        string `json:"name,omitempty"`
	Destination string `json:"destination"`
}
type Config struct {
	Version  int      `json:"version"`
	ClientID string   `json:"client_id"`
	Port     int      `json:"port"`
	Language string   `json:"language"`
	Startup  bool     `json:"startup"`
	Lists    []Target `json:"lists"`
}

func Default() Config { return Config{Version: 1, Port: 8787, Language: "de", Lists: []Target{}} }
func ValidClientID(s string) bool {
	if len(s) != 64 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
func (c Config) Validate() error {
	if c.Version != 1 || c.Port < 1 || c.Port > 65535 || (c.Language != "de" && c.Language != "en") {
		return errors.New("Ungültige Bridge-Einstellungen")
	}
	if c.ClientID != "" && !ValidClientID(c.ClientID) {
		return errors.New("Client-ID prüfen")
	}
	seen := map[Target]bool{}
	for _, v := range c.Lists {
		if v.ID <= 0 || (v.Destination != "sonarr" && v.Destination != "radarr") {
			return errors.New("Listen-Zuordnung prüfen")
		}
		key := Target{ID: v.ID, Destination: v.Destination}
		if seen[key] {
			return errors.New("Doppelte Listen-Zuordnung")
		}
		seen[key] = true
	}
	return nil
}

type Store struct{ Path string }

func (s Store) Load() (Config, error) {
	b, err := os.ReadFile(s.Path)
	if errors.Is(err, os.ErrNotExist) {
		return Default(), nil
	}
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err = json.Unmarshal(b, &c); err != nil {
		return Config{}, errors.New("Einstellungen konnten nicht gelesen werden")
	}
	return c, c.Validate()
}
func (s Store) Save(c Config) error {
	if err := c.Validate(); err != nil {
		return err
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return AtomicWrite(s.Path, b)
}
func AtomicWrite(path string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".pending-*")
	if err != nil {
		return err
	}
	if err = f.Chmod(0600); err != nil {
		f.Close()
		return err
	}
	if _, err = f.Write(b); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
