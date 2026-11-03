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
	cmd := exec.Command(documentProgram(name), path)
	if err = cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}

func documentProgram(name string) string {
	if name == "licenses" {
		return filepath.Join(os.Getenv("WINDIR"), "explorer.exe")
	}
	return filepath.Join(os.Getenv("WINDIR"), "System32", "notepad.exe")
}
