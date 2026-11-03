package desktop

import (
	"context"
	"simklbridge/internal/config"
	"strings"
	"testing"
	"time"
)

func TestUnavailableActionsGiveClearGuidance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	checked := false
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) {
		a.ShowSettings(config.Default(), SettingsActions{})
		a.Dispatch(func() {
			s := a.settings
			disabled := call(user32, "IsWindowEnabled", s.controls[settingBrowser]) == 0
			s.command(settingCopy)
			checked = disabled && strings.Contains(windowText(s.controls[settingStatus]), "auswählen")
			a.Close()
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checked {
		t.Fatal("Nicht verfügbare Aktion blieb ohne verständlichen Hinweis")
	}
}
