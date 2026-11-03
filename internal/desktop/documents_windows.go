package desktop

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

// OpenLocalDocument öffnet nur die mitgelieferten Lizenzunterlagen.
func OpenLocalDocument(name string) error {
	if name != "LICENSE" && name != "THIRD_PARTY_NOTICES.md" && name != "licenses" {
		return errors.New("Unbekannte Lizenzunterlage")
	}
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	path := filepath.Join(filepath.Dir(executable), name)
	if _, err = os.Stat(path); err != nil {
		return err
	}
	program := "notepad.exe"
	if name == "licenses" {
		program = "explorer.exe"
	}
	cmd := exec.Command(filepath.Join(os.Getenv("WINDIR"), "System32", program), path)
	if err = cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
