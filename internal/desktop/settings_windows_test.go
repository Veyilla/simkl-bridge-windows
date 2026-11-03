package desktop

import (
	"context"
	"simklbridge/internal/config"
	"testing"
	"time"
)

func TestSettingsWindowShowsSavedPortAndClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	checked := make(chan bool, 1)
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) {
		c := config.Default()
		c.Port = 9001
		c.ClientID = "synthetic-client"
		a.ShowSettings(c, SettingsActions{})
		a.Dispatch(func() {
			s := a.settings
			if s == nil {
				checked <- false
			} else {
				checked <- s.hwnd == a.hwnd && windowText(s.controls[settingClient]) == "synthetic-client" && windowText(s.controls[settingPort]) == "9001"
			}
			a.Close()
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !<-checked {
		t.Fatal("Gespeicherte Einstellungen fehlen in der Oberfläche")
	}
}
