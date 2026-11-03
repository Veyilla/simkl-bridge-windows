package desktop

import (
	"context"
	"simklbridge/internal/config"
	"testing"
	"time"
)

func TestOpeningSettingsAfterLogKeepsWorkingSave(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	saved := false
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) {
		a.ShowLog(config.Default(), SettingsActions{ReadOnly: true, LogText: "Beispiel"})
		a.ShowSettings(config.Default(), SettingsActions{Save: func(config.Config) error { saved = true; return nil }})
		a.Dispatch(func() { a.settings.command(settingSave); a.Close() })
	})
	if err != nil {
		t.Fatal(err)
	}
	if !saved {
		t.Fatal("Wechsel vom Protokoll ließ Einstellungen funktionslos")
	}
}
