package auth

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProtectedStoreRoundTripAndClientBinding(t *testing.T) {
	path := filepath.Join(t.TempDir(), "authorization.dat")
	store := ProtectedStore{Path: path, ClientID: testClient}
	want := Tokens{ClientID: testClient, AccessToken: "synthetic-access-secret", RefreshToken: "synthetic-refresh-secret", ExpiresAt: time.Unix(3000, 0)}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte(want.AccessToken)) || bytes.Contains(raw, []byte(want.RefreshToken)) {
		t.Fatal("Anmeldung im Klartext gespeichert")
	}
	got, err := store.Load()
	if err != nil || got.AccessToken != want.AccessToken || got.RefreshToken != want.RefreshToken {
		t.Fatalf("Anmeldung nicht lesbar: %v", err)
	}
	other := ProtectedStore{Path: path, ClientID: "other-client"}
	if _, err = other.Load(); err == nil {
		t.Fatal("Andere Registrierung konnte Anmeldung lesen")
	}
	raw[len(raw)/2] ^= 1
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = store.Load(); err == nil {
		t.Fatal("Beschädigte Anmeldung akzeptiert")
	}
}
