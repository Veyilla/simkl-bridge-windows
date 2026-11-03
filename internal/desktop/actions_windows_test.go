package desktop

import (
	"context"
	"simklbridge/internal/config"
	"testing"
	"time"
	"unsafe"
)

func TestEveryHomeButtonReachesItsActionAtBothScales(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	events := make(chan string, 32)
	record := func(name string) func() { return func() { events <- name } }
	actions := Actions{Start: record("start"), Stop: record("stop"), Restart: record("restart"), OpenSettings: record("settings"), OpenLog: record("log"), OpenAbout: record("about"), Quit: record("quit"), SetStartup: func(bool) { events <- "startup" }, SetLanguage: func(s string) { events <- s }}
	err := RunView(ctx, View{Phase: "stopped", Language: "de"}, actions, func(a *App) {
		a.Dispatch(func() {
			for _, scale := range []float64{1, 2} {
				a.fonts(scale)
				call(user32, "SetWindowPos", a.hwnd, 0, 0, 0, uintptr(int(368*scale)), uintptr(int(552*scale)), 6)
				for i, want := range []string{"start", "restart", "startup", "settings", "log", "about", "quit"} {
					a.view.Phase = "stopped"
					r := HomeLayout().Rows[i]
					x, y := int(70*scale), int(float64(r.Y+20)*scale)
					call(user32, "SendMessageW", a.hwnd, 0x202, 0, uintptr(uint32(x)|uint32(y)<<16))
					select {
					case got := <-events:
						if got != want {
							t.Errorf("%s löst %s aus", want, got)
						}
					case <-time.After(time.Second):
						t.Errorf("Taste %s ohne Funktion", want)
					}
				}
				a.view.Phase = "running"
				r := HomeLayout().Rows[0]
				x, y := int(70*scale), int(float64(r.Y+20)*scale)
				call(user32, "SendMessageW", a.hwnd, 0x202, 0, uintptr(uint32(x)|uint32(y)<<16))
				select {
				case got := <-events:
					if got != "stop" {
						t.Error(got)
					}
				case <-time.After(time.Second):
					t.Error("Stoppen ohne Funktion")
				}
			}
			a.Close()
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSettingsBackReturnsToSameMenu(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	checked := false
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) {
		a.ShowSettings(config.Default(), SettingsActions{})
		a.Dispatch(func() {
			hwnd := a.hwnd
			a.settings.command(119)
			var r winRect
			call(user32, "GetClientRect", hwnd, uintptr(unsafe.Pointer(&r)))
			checked = a.settings == nil && a.hwnd == hwnd && r.Right > 0
			a.Close()
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checked {
		t.Fatal("Zurück öffnete nicht dasselbe Hauptmenü")
	}
}
