package desktop

import (
	"context"
	"simklbridge/internal/config"
	"testing"
	"time"
)

func TestBackWaitsForCancelledSignInBeforeResumingWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	started := make(chan struct{})
	finish := make(chan struct{})
	closed := make(chan struct{}, 1)
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) {
		a.ShowSettings(config.Default(), SettingsActions{Connect: func(ctx context.Context, _ string, _ func(string, string)) error {
			close(started)
			<-ctx.Done()
			<-finish
			return ctx.Err()
		}, Closed: func(bool) { closed <- struct{}{} }})
		a.Dispatch(func() {
			a.settings.command(settingConnect)
			<-started
			a.settings.hide()
			select {
			case <-closed:
				t.Error("Worker vor Ende der abgebrochenen Anmeldung fortgesetzt")
			case <-time.After(30 * time.Millisecond):
			}
			close(finish)
			select {
			case <-closed:
			case <-time.After(time.Second):
				t.Error("Abschluss nach abgebrochener Anmeldung fehlt")
			}
			a.Close()
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
