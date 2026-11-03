package desktop

import (
	"context"
	"simklbridge/internal/config"
	"testing"
	"time"
)

func TestRepeatedOpenUsesCurrentSettingsCompletion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	closed := make(chan int, 2)
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) {
		a.ShowSettings(config.Default(), SettingsActions{Closed: func(bool) { closed <- 1 }})
		a.ShowSettings(config.Default(), SettingsActions{Closed: func(bool) { closed <- 2 }})
		a.Dispatch(func() {
			a.settings.hide()
			select {
			case current := <-closed:
				if current != 2 {
					t.Error("Veralteter Abschluss nach doppeltem Öffnen")
				}
			case <-time.After(time.Second):
				t.Error("Abschluss fehlt")
			}
			a.Close()
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
func TestTrayLogoHasFrameAndBlackCorners(t *testing.T) {
	logo := TrayLogoImage(32)
	frame := logo.RGBAAt(1, 16)
	if frame.R < 40 || frame.R > 180 || logo.RGBAAt(0, 0).R != 0 || logo.RGBAAt(16, 16).R != 255 {
		t.Fatalf("Tray-Rahmen oder S unlesbar: %+v", frame)
	}
}
func TestListFailureKeepsStopActionForRunningService(t *testing.T) {
	s := Presentation(View{Phase: "running", Language: "de", Error: "tvdb_missing"})
	if s.Action != "stop" || s.Color != 0xff6262 || s.Text != "TVDB-Zuordnung fehlt" {
		t.Fatalf("Fehlerstatus verhindert Bedienung: %+v", s)
	}
}
