package worker

import (
	"context"
	"errors"
	"path/filepath"
	"simklbridge/internal/config"
	"simklbridge/internal/supervisor"
	"testing"
)

func TestMissingConfigurationRequiresSignIn(t *testing.T) {
	err := Run(context.Background(), t.TempDir())
	if !errors.Is(err, supervisor.ErrAuthentication) {
		t.Fatal(err)
	}
}
func TestInvalidSettingsDoNotStartService(t *testing.T) {
	dir := t.TempDir()
	store := config.Store{Path: filepath.Join(dir, "settings.json")}
	c := config.Default()
	c.ClientID = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := store.Save(c); err != nil {
		t.Fatal(err)
	}
	if err := Run(context.Background(), dir); !errors.Is(err, supervisor.ErrConfiguration) {
		t.Fatal(err)
	}
}
