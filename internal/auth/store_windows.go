package auth

import (
	"encoding/json"
	"errors"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"simklbridge/internal/config"
	"unsafe"
)

type ProtectedStore struct{ Path, ClientID string }

func StorePath(dir, clientID string) string {
	return filepath.Join(dir, "authorization-"+clientID+".dat")
}

func blob(b []byte) *windows.DataBlob {
	v := &windows.DataBlob{Size: uint32(len(b))}
	if len(b) > 0 {
		v.Data = &b[0]
	}
	return v
}
func (s ProtectedStore) Save(t Tokens) error {
	plain, err := json.Marshal(t)
	if err != nil {
		return err
	}
	defer clear(plain)
	entropy := []byte("SimklBridge/" + s.ClientID)
	var out windows.DataBlob
	if err = windows.CryptProtectData(blob(plain), nil, blob(entropy), 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return errors.New("Windows konnte die Anmeldung nicht schützen")
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	return config.AtomicWrite(s.Path, unsafe.Slice(out.Data, out.Size))
}
func (s ProtectedStore) Load() (Tokens, error) {
	encrypted, err := os.ReadFile(s.Path)
	if err != nil {
		return Tokens{}, ErrAuthentication
	}
	if len(encrypted) == 0 || len(encrypted) > 1<<20 {
		return Tokens{}, ErrAuthentication
	}
	entropy := []byte("SimklBridge/" + s.ClientID)
	var out windows.DataBlob
	if err = windows.CryptUnprotectData(blob(encrypted), nil, blob(entropy), 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return Tokens{}, ErrAuthentication
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	plain := unsafe.Slice(out.Data, out.Size)
	defer clear(plain)
	var token Tokens
	if err = json.Unmarshal(plain, &token); err != nil {
		return Tokens{}, ErrAuthentication
	}
	return token, nil
}
