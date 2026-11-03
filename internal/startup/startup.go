package startup

import (
	"errors"
	"path/filepath"
	"strings"
)

var ErrForeignEntry = errors.New("Vorhandener Autostart-Eintrag gehört zu einer anderen Anwendung")

type Entry interface {
	Read() (string, error)
	Write(string) error
	Remove() error
}
type Manager struct {
	Entry      Entry
	Executable string
}

func (m Manager) command() (string, error) {
	if !filepath.IsAbs(m.Executable) || strings.ContainsAny(m.Executable, "\"\r\n") {
		return "", errors.New("Ungültiger Programmpfad")
	}
	return `"` + m.Executable + `" --autostart`, nil
}
func (m Manager) Enabled() (bool, error) {
	command, err := m.command()
	if err != nil {
		return false, err
	}
	old, err := m.Entry.Read()
	if err != nil {
		return false, err
	}
	if old != "" && old != command {
		return false, ErrForeignEntry
	}
	return old == command, nil
}
func (m Manager) Set(enabled bool) error {
	command, err := m.command()
	if err != nil {
		return err
	}
	old, err := m.Entry.Read()
	if err != nil {
		return err
	}
	if old != "" && old != command {
		return ErrForeignEntry
	}
	if enabled {
		return m.Entry.Write(command)
	}
	if old != "" {
		return m.Entry.Remove()
	}
	return nil
}
