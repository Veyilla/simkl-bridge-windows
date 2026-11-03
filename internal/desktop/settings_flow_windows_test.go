package desktop

import (
	"context"
	"simklbridge/internal/config"
	"testing"
	"time"
)

func TestInlineListSelectionAndSaveReachBackend(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	discovered := make(chan struct{}, 1)
	inspected := make(chan string, 1)
	saved := make(chan config.Config, 1)
	c := config.Default()
	c.ClientID = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	actions := SettingsActions{Discover: func(context.Context, string) ([]ListChoice, error) {
		discovered <- struct{}{}
		return []ListChoice{{ID: 123, Name: "Beispiel", MediaType: "tv"}}, nil
	}, Inspect: func(_ context.Context, _ string, input, dest string) (config.Target, error) {
		inspected <- input + "/" + dest
		return config.Target{ID: 123, Name: "Beispiel", Destination: dest}, nil
	}, Save: func(c config.Config) error { saved <- c; return nil }}
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) {
		a.ShowSettings(c, actions)
		a.Dispatch(func() { a.settings.command(settingDiscover) })
		go func() {
			<-discovered
			deadline := time.Now().Add(2 * time.Second)
			for {
				ready := make(chan bool, 1)
				a.Dispatch(func() { ready <- a.settings != nil && !a.settings.busy && len(a.settings.choices) == 1 })
				if <-ready {
					break
				}
				if time.Now().After(deadline) {
					t.Error("Listen wurden nicht angezeigt")
					a.Close()
					return
				}
				time.Sleep(time.Millisecond)
			}
			a.Dispatch(func() {
				call(user32, "SendMessageW", a.settings.controls[settingAvailable], 0x186, 0, 0)
				a.settings.command(settingSelectSonarr)
			})
			if got := <-inspected; got != "123/sonarr" {
				t.Error("Falsche Zuordnung", got)
			}
			for {
				ready := make(chan bool, 1)
				a.Dispatch(func() { ready <- a.settings != nil && !a.settings.busy && len(a.settings.draft.Lists) == 1 })
				if <-ready {
					break
				}
				if time.Now().After(deadline) {
					t.Error("Zuordnung wurde nicht übernommen")
					a.Close()
					return
				}
				time.Sleep(time.Millisecond)
			}
			a.Dispatch(func() { a.settings.command(settingSave) })
			select {
			case got := <-saved:
				if len(got.Lists) != 1 || got.Lists[0].ID != 123 || got.Lists[0].Destination != "sonarr" {
					t.Error("Speichern veränderte Auswahl", got.Lists)
				}
			case <-time.After(time.Second):
				t.Error("Speichern ohne Funktion")
			}
			a.Close()
		}()
	})
	if err != nil {
		t.Fatal(err)
	}
}
func TestInlineProtocolIsSameWindow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	checked := false
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) {
		a.ShowLog(config.Default(), SettingsActions{ReadOnly: true, LogText: "Beispiel-Protokoll"})
		a.Dispatch(func() {
			checked = a.settings != nil && a.settings.hwnd == a.hwnd && windowText(a.settings.controls[207]) == "Beispiel-Protokoll"
			a.Close()
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checked {
		t.Fatal("Protokoll wurde nicht im selben Menü dargestellt")
	}
}
