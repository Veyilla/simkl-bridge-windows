package main

import (
	"context"
	"errors"

	"golang.org/x/sys/windows"
	"log"
	"os"
	"path/filepath"
	"simklbridge/internal/auth"
	"simklbridge/internal/config"
	"simklbridge/internal/control"
	"simklbridge/internal/desktop"
	"simklbridge/internal/diagnostics"
	"simklbridge/internal/feed"
	"simklbridge/internal/listinput"
	"simklbridge/internal/simkl"
	"simklbridge/internal/startup"
	"simklbridge/internal/supervisor"
	"simklbridge/internal/worker"
	"strings"
	"sync"
	"time"
)

var version = "0.1.0-dev"
var projectURL = ""

func main() {
	dir, err := dataDirectory()
	if err != nil {
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "--bridge-worker" {
		err = worker.Run(context.Background(), dir)
		switch {
		case errors.Is(err, supervisor.ErrAuthentication):
			os.Exit(41)
		case errors.Is(err, supervisor.ErrConfiguration):
			os.Exit(42)
		case err != nil:
			os.Exit(1)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] != "--autostart" {
		return
	}
	name, _ := windows.UTF16PtrFromString("Local\\SimklBridge")
	mutex, err := windows.CreateMutex(nil, false, name)
	if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
		if mutex != 0 {
			windows.CloseHandle(mutex)
		}
		desktop.WakeExisting()
		return
	}
	if err != nil {
		return
	}
	defer windows.CloseHandle(mutex)
	if err = os.MkdirAll(dir, 0700); err != nil {
		return
	}
	file, err := os.OpenFile(filepath.Join(dir, "bridge.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer file.Close()
	logger := log.New(file, "", log.Ldate|log.Ltime)
	store := config.Store{Path: filepath.Join(dir, "settings.json")}
	current, loadErr := store.Load()
	if loadErr != nil {
		current = config.Default()
		logger.Print("Einstellungen konnten nicht gelesen werden; bestehende Datei bleibt erhalten.")
	}
	exe, err := os.Executable()
	if err != nil {
		return
	}
	startupManager := startup.Manager{Entry: startup.RegistryEntry{}, Executable: exe}
	enabled, startupErr := startupManager.Enabled()
	if startupErr == nil {
		current.Startup = enabled
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	controller := supervisor.New(supervisor.NewOwnedWorker, nil)
	lifecycle := control.New(controller)
	var mu sync.Mutex
	var app *desktop.App
	snapshot := func() config.Config {
		mu.Lock()
		defer mu.Unlock()
		c := current
		c.Lists = append([]config.Target{}, c.Lists...)
		return c
	}
	var viewMu sync.Mutex
	var lastView desktop.View
	viewSet := false
	update := func() {
		c := snapshot()
		s := controller.Snapshot()
		code := s.ErrorCode
		if s.WorkerPID > 0 && (s.Phase == "running" || (s.Phase == "error" && s.ErrorCode != "restart_limit")) {
			if report, err := diagnostics.Load(dir, s.WorkerPID); err == nil && report.Code != "" {
				code = report.Code
			}
		}
		v := desktop.View{Phase: s.Phase, Error: code, Language: c.Language, Startup: c.Startup}
		viewMu.Lock()
		defer viewMu.Unlock()
		if !viewSet || v != lastView {
			app.Update(v)
			lastView = v
			viewSet = true
		}
	}
	stop := func() {
		stopCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()
		if lifecycle.Stop(stopCtx) != nil {
			logger.Print("Bridge konnte nicht rechtzeitig beendet werden.")
		}
	}
	var settingsActions desktop.SettingsActions
	clientFor := func(id string) (*auth.Manager, *simkl.Client) {
		manager := auth.New(id, "https://api.simkl.com", nil, auth.ProtectedStore{Path: auth.StorePath(dir, id), ClientID: id})
		return manager, simkl.New(id, "", nil, manager)
	}
	settingsActions = desktop.SettingsActions{Version: version, ProjectURL: validProjectURL(projectURL),
		Connect: func(ctx context.Context, id string, display func(string, string)) error {
			release, err := lifecycle.EnterSetupOperation(ctx)
			if err != nil {
				return err
			}
			defer release()
			manager, _ := clientFor(id)
			ch, err := manager.Begin(ctx)
			if err != nil {
				return err
			}
			display(ch.UserCode, ch.VerificationURI)
			return manager.Complete(ctx, ch)
		},
		Discover: func(ctx context.Context, id string) ([]desktop.ListChoice, error) {
			release, err := lifecycle.EnterSetupOperation(ctx)
			if err != nil {
				return nil, err
			}
			defer release()
			_, api := clientFor(id)
			account, err := api.AccountID(ctx)
			if err != nil {
				return nil, err
			}
			lists, err := api.UserLists(ctx, account)
			if err != nil {
				return nil, err
			}
			choices := []desktop.ListChoice{}
			for _, l := range lists {
				count := 0
				if l.Counts != nil {
					count = l.Counts.Items
				}
				choices = append(choices, desktop.ListChoice{ID: l.ID, Name: l.Name, MediaType: l.MediaType, Count: count})
			}
			return choices, nil
		},
		Inspect: func(ctx context.Context, clientID, input, destination string) (config.Target, error) {
			release, err := lifecycle.EnterSetupOperation(ctx)
			if err != nil {
				return config.Target{}, err
			}
			defer release()
			id, err := listinput.Parse(input)
			if err != nil {
				return config.Target{}, err
			}
			_, api := clientFor(clientID)
			cache := simkl.NewCache(api)
			list, err := cache.Get(ctx, id, simkl.Requested)
			if err != nil {
				return config.Target{}, err
			}
			if _, err = feed.Format(ctx, cache, list, destination); err != nil {
				return config.Target{}, err
			}
			return config.Target{ID: id, Name: list.Name, Destination: destination}, nil
		},
		Save: func(c config.Config) error {
			if loadErr != nil {
				return errors.New("Bestehende Einstellungsdatei zuerst prüfen; sie wurde nicht überschrieben")
			}
			c.Startup = snapshot().Startup
			if err := store.Save(c); err != nil {
				return errors.New("Einstellungen konnten nicht gespeichert werden")
			}
			mu.Lock()
			current = c
			mu.Unlock()
			update()
			return nil
		},
	}
	openSettings := func() {
		c := snapshot()
		ticket, err := lifecycle.BeginSettings(ctx, c.ClientID != "" && len(c.Lists) > 0)
		if err != nil {
			logger.Print("Einrichtung konnte nicht sicher geöffnet werden.")
			return
		}
		actions := settingsActions
		actions.Closed = func(saved bool) {
			current := snapshot()
			lifecycle.EndSettings(ctx, ticket, saved, current.ClientID != "" && len(current.Lists) > 0)
		}
		app.ShowSettings(c, actions)
	}
	actions := desktop.Actions{
		Start: func() { lifecycle.Start(ctx) }, Stop: stop, Restart: func() { lifecycle.Restart(ctx) }, Quit: cancel,
		OpenSettings: openSettings,
		OpenAbout: func() {
			about := desktop.SettingsActions{Version: version, ProjectURL: validProjectURL(projectURL), ReadOnly: true}
			app.ShowAbout(snapshot(), about)
		},
		OpenLog: func() {
			app.ShowLog(snapshot(), desktop.SettingsActions{ReadOnly: true, LogText: readLogTail(filepath.Join(dir, "bridge.log"))})
		},
		SetLanguage: func(lang string) {
			c := snapshot()
			c.Language = lang
			if loadErr == nil && store.Save(c) == nil {
				mu.Lock()
				current = c
				mu.Unlock()
			}
			update()
		},
		SetStartup: func(want bool) {
			if loadErr != nil {
				return
			}
			c := snapshot()
			old := c.Startup
			if err := startupManager.Set(want); err != nil {
				logger.Print("Autostart konnte nicht geändert werden.")
				update()
				return
			}
			c.Startup = want
			if err := store.Save(c); err != nil {
				startupManager.Set(old)
				logger.Print("Autostart wurde zurückgesetzt: Einstellungen nicht gespeichert.")
			} else {
				mu.Lock()
				current = c
				mu.Unlock()
			}
			update()
		},
	}
	view := desktop.View{Phase: "stopped", Language: current.Language, Startup: current.Startup}
	err = desktop.RunView(ctx, view, actions, func(a *desktop.App) {
		app = a
		changes := controller.Subscribe()
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case state := <-changes:
					logger.Printf("Bridge: %s (%s), Neustartversuche: %d", state.Phase, state.ErrorCode, state.Retries)
					update()
				case <-ticker.C:
					update()
				case <-ctx.Done():
					return
				}
			}
		}()
		if current.ClientID == "" || len(current.Lists) == 0 {
			go openSettings()
		} else {
			go lifecycle.Start(ctx)
		}
	})
	stop()
	if err != nil {
		logger.Print("Windows-Oberfläche konnte nicht gestartet werden.")
	}
}
func dataDirectory() (string, error) {
	if dir := os.Getenv("SIMKL_BRIDGE_DATA_DIR"); dir != "" {
		if !filepath.IsAbs(dir) {
			return "", errors.New("Ungültiger Einstellungsordner")
		}
		return dir, nil
	}
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "SimklBridge"), nil
}
func validProjectURL(value string) string {
	if strings.HasPrefix(value, "https://github.com/") && !strings.ContainsAny(value, "\r\n") {
		return value
	}
	return ""
}
