package desktop

import (
	"context"
	"fmt"
	"simklbridge/internal/config"
	"strconv"
	"sync"
	"syscall"
	"unsafe"
)

type ListChoice struct {
	ID              int64
	Name, MediaType string
	Count           int
}
type SettingsActions struct {
	Connect             func(context.Context, string, func(string, string)) error
	Discover            func(context.Context, string) ([]ListChoice, error)
	Inspect             func(context.Context, string, string, string) (config.Target, error)
	Save                func(config.Config) error
	Closed              func(bool)
	OpenDocument        func(string) error
	ProjectURL, Version string
	ReadOnly            bool
	LogText             string
}

const (
	settingClient        = 100
	settingConnect       = 101
	settingBrowser       = 102
	settingStatus        = 103
	settingDiscover      = 104
	settingAvailable     = 105
	settingSelectSonarr  = 106
	settingSelectRadarr  = 107
	settingLink          = 108
	settingAddSonarr     = 109
	settingAddRadarr     = 110
	settingConfigured    = 111
	settingCopy          = 112
	settingRemove        = 113
	settingPort          = 114
	settingSave          = 115
	settingGuide         = 116
	settingProject       = 117
	settingSimkl         = 118
	settingLicense       = 120
	settingNotices       = 121
	settingLicenseFolder = 122
)

type settingsWindow struct {
	app               *App
	hwnd, font, brush uintptr
	controls          map[int]uintptr
	pages             map[int]int
	page              int
	scale             float64
	draft             config.Config
	actions           SettingsActions
	choices           []ListChoice
	busy              bool
	uri               string
	ctx               context.Context
	cancel            context.CancelFunc
	jobs              sync.WaitGroup
	saved             bool
}

func windowText(hwnd uintptr) string {
	n := call(user32, "GetWindowTextLengthW", hwnd)
	b := make([]uint16, n+1)
	call(user32, "GetWindowTextW", hwnd, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	return syscall.UTF16ToString(b)
}
func setText(hwnd uintptr, s string) {
	call(user32, "SetWindowTextW", hwnd, uintptr(unsafe.Pointer(u(s))))
}
func (s *settingsWindow) label(de, en string) string {
	if s.draft.Language == "en" {
		return en
	}
	return de
}
func (s *settingsWindow) px(v int) uintptr { return uintptr(float64(v) * s.scale) }
func (s *settingsWindow) control(id, page int, class, text string, x, y, w, h int, style uintptr) {
	if class == "STATIC" {
		style |= 0x80
	}
	hnd := call(user32, "CreateWindowExW", 0, uintptr(unsafe.Pointer(u(class))), uintptr(unsafe.Pointer(u(text))), 0x50000000|style, s.px(x), s.px(y), s.px(w), s.px(h), s.hwnd, uintptr(id), call(kernel32, "GetModuleHandleW", 0), 0)
	s.controls[id] = hnd
	s.pages[id] = page
	call(user32, "SendMessageW", hnd, 0x30, s.font, 1)
}
func (s *settingsWindow) button(id, page int, de, en string, x, y, w int) {
	s.control(id, page, "BUTTON", s.label(de, en), x, y, w, 32, 0x10000|0xb)
}
func (s *settingsWindow) status(text string) { setText(s.controls[settingStatus], text) }
func (s *settingsWindow) showPage(page int) {
	s.page = page
	for id, h := range s.controls {
		visible := s.pages[id] < 0 || s.pages[id] == page
		if page >= 3 && (id == 10 || id == 11 || id == 12 || id == settingSave || id == settingStatus) {
			visible = false
		}
		cmd := uintptr(0)
		if visible {
			cmd = 5
		}
		call(user32, "ShowWindow", h, cmd)
	}
	call(user32, "InvalidateRect", s.hwnd, 0, 1)
}
func (s *settingsWindow) fillConfigured() {
	h := s.controls[settingConfigured]
	call(user32, "SendMessageW", h, 0x184, 0, 0)
	for _, v := range s.draft.Lists {
		label := fmt.Sprintf("%s  ·  %s  ·  %d", v.Name, v.Destination, v.ID)
		call(user32, "SendMessageW", h, 0x180, 0, uintptr(unsafe.Pointer(u(label))))
	}
}
func (s *settingsWindow) complete(err error) {
	s.busy = false
	if err != nil {
		s.status(err.Error())
	} else {
		s.status(s.label("Bereit.", "Ready."))
	}
}
func (s *settingsWindow) run(f func(context.Context) error) {
	if s.busy {
		return
	}
	s.busy = true
	s.status(s.label("Bitte warten …", "Please wait …"))
	ctx := s.ctx
	s.jobs.Add(1)
	go func() {
		defer s.jobs.Done()
		err := f(ctx)
		s.app.Dispatch(func() {
			if ctx.Err() == nil {
				s.complete(err)
			}
		})
	}()
}
func (s *settingsWindow) add(input, dest string) {
	if s.actions.Inspect == nil {
		return
	}
	client := windowText(s.controls[settingClient])
	s.run(func(ctx context.Context) error {
		target, err := s.actions.Inspect(ctx, client, input, dest)
		if err != nil {
			return err
		}
		s.app.Dispatch(func() {
			if ctx.Err() != nil {
				return
			}
			for _, old := range s.draft.Lists {
				if old.ID == target.ID && old.Destination == target.Destination {
					s.status(s.label("Diese Liste ist bereits eingerichtet.", "This list is already configured."))
					return
				}
			}
			s.draft.Lists = append(s.draft.Lists, target)
			s.fillConfigured()
		})
		return nil
	})
}
func (s *settingsWindow) command(id int) {
	switch {
	case id == 2:
		s.hide()
	case id == 119:
		s.hide()
	case id >= 10 && id <= 13:
		s.showPage(id - 10)
	case id == settingConnect:
		if s.actions.Connect == nil {
			return
		}
		client := windowText(s.controls[settingClient])
		s.run(func(ctx context.Context) error {
			return s.actions.Connect(ctx, client, func(code, uri string) {
				s.app.Dispatch(func() {
					if ctx.Err() == nil {
						s.uri = uri
						call(user32, "EnableWindow", s.controls[settingBrowser], 1)
						s.status(s.label("Dein Anmeldecode: ", "Your sign-in code: ") + code + s.label("\r\nIm Browser bestätigen; das Menü bleibt offen.", "\r\nApprove in your browser; this menu stays open."))
					}
				})
			})
		})
	case id == settingBrowser:
		if s.uri != "" {
			OpenURL(s.uri)
		}
	case id == settingGuide:
		OpenURL("https://simkl.com/settings/developer/")
	case id == settingSimkl:
		OpenURL("https://simkl.com/")
	case id == settingProject:
		if s.actions.ProjectURL != "" {
			OpenURL(s.actions.ProjectURL)
		}
	case id == settingLicense || id == settingNotices || id == settingLicenseFolder:
		name := "LICENSE"
		if id == settingNotices {
			name = "THIRD_PARTY_NOTICES.md"
		}
		if id == settingLicenseFolder {
			name = "licenses"
		}
		open := s.actions.OpenDocument
		if open == nil {
			open = OpenLocalDocument
		}
		if err := open(name); err != nil {
			s.status(s.label("Lizenzunterlagen konnten nicht geöffnet werden.", "Could not open license documents."))
			call(user32, "ShowWindow", s.controls[settingStatus], 1)
		}
	case id == settingDiscover:
		if s.actions.Discover == nil {
			return
		}
		client := windowText(s.controls[settingClient])
		s.run(func(ctx context.Context) error {
			choices, err := s.actions.Discover(ctx, client)
			if err != nil {
				return err
			}
			s.app.Dispatch(func() {
				if ctx.Err() != nil {
					return
				}
				s.choices = choices
				h := s.controls[settingAvailable]
				call(user32, "SendMessageW", h, 0x184, 0, 0)
				for _, v := range choices {
					label := fmt.Sprintf("%s  ·  %s  ·  %d", v.Name, v.MediaType, v.Count)
					call(user32, "SendMessageW", h, 0x180, 0, uintptr(unsafe.Pointer(u(label))))
				}
			})
			return nil
		})
	case id == settingSelectSonarr || id == settingSelectRadarr:
		index := int(int32(call(user32, "SendMessageW", s.controls[settingAvailable], 0x188, 0, 0)))
		if index < 0 || index >= len(s.choices) {
			s.status(s.label("Bitte zuerst eine Liste auswählen.", "Select a list first."))
			return
		}
		dest := "sonarr"
		if id == settingSelectRadarr {
			dest = "radarr"
		}
		s.add(strconv.FormatInt(s.choices[index].ID, 10), dest)
	case id == settingAddSonarr || id == settingAddRadarr:
		dest := "sonarr"
		if id == settingAddRadarr {
			dest = "radarr"
		}
		s.add(windowText(s.controls[settingLink]), dest)
	case id == settingCopy || id == settingRemove:
		index := int(int32(call(user32, "SendMessageW", s.controls[settingConfigured], 0x188, 0, 0)))
		if index < 0 || index >= len(s.draft.Lists) {
			s.status(s.label("Bitte zuerst eine eingerichtete Liste auswählen.", "Select a configured list first."))
			return
		}
		if id == settingCopy {
			port, err := strconv.Atoi(windowText(s.controls[settingPort]))
			if err != nil || port < 1 || port > 65535 {
				s.status(s.label("Port prüfen.", "Check port."))
				return
			}
			target := s.draft.Lists[index]
			if err = CopyText(s.hwnd, fmt.Sprintf("http://127.0.0.1:%d/%s/%d", port, target.Destination, target.ID)); err != nil {
				s.status(s.label("Kopieren fehlgeschlagen.", "Copy failed."))
			} else {
				s.status(s.label("Listenadresse kopiert.", "List address copied."))
			}
		} else {
			s.draft.Lists = append(s.draft.Lists[:index], s.draft.Lists[index+1:]...)
			s.fillConfigured()
		}
	case id == settingSave:
		if s.actions.ReadOnly {
			s.hide()
			return
		}
		if s.busy {
			s.status(s.label("Bitte zuerst den laufenden Vorgang abschließen.", "Wait for the current operation to finish."))
			return
		}
		draft := s.draft
		draft.ClientID = windowText(s.controls[settingClient])
		port, err := strconv.Atoi(windowText(s.controls[settingPort]))
		if err != nil {
			s.status(s.label("Port prüfen.", "Check port."))
			return
		}
		draft.Port = port
		if err = draft.Validate(); err != nil {
			s.status(err.Error())
			return
		}
		if s.actions.Save != nil {
			if err = s.actions.Save(draft); err != nil {
				s.status(err.Error())
				return
			}
		}
		s.draft = draft
		s.saved = true
		s.hide()
	}
}
func (s *settingsWindow) hide() {
	s.shutdown()
	s.app.settings = nil
	s.app.position(Size{368, 552})
	s.app.redraw()
	if s.actions.Closed != nil {
		closed, saved := s.actions.Closed, s.saved
		go func() { s.jobs.Wait(); closed(saved) }()
	}
}
func (s *settingsWindow) shutdown() {
	s.cancel()
	for _, h := range s.controls {
		call(user32, "DestroyWindow", h)
	}
	call(gdi32, "DeleteObject", s.font)
	call(gdi32, "DeleteObject", s.brush)
}

type drawItem struct {
	Kind, ID, Item, Action, State uint32
	Window, DC                    uintptr
	Rect                          winRect
	Data                          uintptr
}

func (a *App) ShowSettings(c config.Config, actions SettingsActions) {
	a.Dispatch(func() { a.openSettings(c, actions, 0) })
}
func (a *App) ShowAbout(c config.Config, actions SettingsActions) {
	a.Dispatch(func() { a.openSettings(c, actions, 3) })
}
func (a *App) ShowLog(c config.Config, actions SettingsActions) {
	a.Dispatch(func() { a.openSettings(c, actions, 4) })
}
func OpenURL(url string) {
	call(shell32, "ShellExecuteW", 0, uintptr(unsafe.Pointer(u("open"))), uintptr(unsafe.Pointer(u(url))), 0, 0, 1)
}
func CopyText(hwnd uintptr, text string) error {
	if call(user32, "OpenClipboard", hwnd) == 0 {
		return fmt.Errorf("Zwischenablage nicht verfügbar")
	}
	defer call(user32, "CloseClipboard")
	v := syscall.StringToUTF16(text)
	h := call(kernel32, "GlobalAlloc", 2, uintptr(len(v)*2))
	if h == 0 {
		return fmt.Errorf("Zwischenablage nicht verfügbar")
	}
	p := call(kernel32, "GlobalLock", h)
	if p == 0 {
		call(kernel32, "GlobalFree", h)
		return fmt.Errorf("Zwischenablage nicht verfügbar")
	}
	call(kernel32, "RtlMoveMemory", p, uintptr(unsafe.Pointer(&v[0])), uintptr(len(v)*2))
	call(kernel32, "GlobalUnlock", h)
	if call(user32, "EmptyClipboard") == 0 || call(user32, "SetClipboardData", 13, h) == 0 {
		call(kernel32, "GlobalFree", h)
		return fmt.Errorf("Kopieren fehlgeschlagen")
	}
	return nil
}
