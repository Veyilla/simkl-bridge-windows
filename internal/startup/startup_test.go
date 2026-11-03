package startup

import (
	"errors"
	"testing"
)

type memoryEntry struct{ value string }

func (e *memoryEntry) Read() (string, error) { return e.value, nil }
func (e *memoryEntry) Write(v string) error  { e.value = v; return nil }
func (e *memoryEntry) Remove() error         { e.value = ""; return nil }
func TestOwnStartupCanBeEnabledAndDisabled(t *testing.T) {
	s := &memoryEntry{}
	m := Manager{Entry: s, Executable: `C:\Program Files\Simkl Bridge\simkl-bridge.exe`}
	if err := m.Set(true); err != nil {
		t.Fatal(err)
	}
	if s.value != `"C:\Program Files\Simkl Bridge\simkl-bridge.exe" --autostart` {
		t.Fatal(s.value)
	}
	if enabled, err := m.Enabled(); err != nil || !enabled {
		t.Fatal(enabled, err)
	}
	if err := m.Set(false); err != nil || s.value != "" {
		t.Fatal(err)
	}
}
func TestForeignStartupEntryIsPreserved(t *testing.T) {
	s := &memoryEntry{value: `"C:\Other\program.exe"`}
	m := Manager{Entry: s, Executable: `C:\Simkl\simkl-bridge.exe`}
	for _, enabled := range []bool{true, false} {
		if err := m.Set(enabled); !errors.Is(err, ErrForeignEntry) {
			t.Fatal(err)
		}
		if s.value != `"C:\Other\program.exe"` {
			t.Fatal("Fremder Eintrag verändert")
		}
	}
}
