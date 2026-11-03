package desktop

import (
	"context"
	"simklbridge/internal/config"
	"strconv"
	"syscall"
	"unsafe"
)

func (a *App) openSettings(c config.Config, actions SettingsActions, page int) {
	if a.settings != nil {
		if !a.settings.actions.ReadOnly {
			if page == 0 && !actions.ReadOnly {
				a.settings.actions.Closed = actions.Closed
			}
			if page == 4 {
				setText(a.settings.controls[207], actions.LogText)
			}
			a.settings.showPage(page)
			return
		}
		a.settings.shutdown()
		a.settings = nil
	}
	a.position(Size{420, 650})
	s := &settingsWindow{app: a, hwnd: a.hwnd, draft: c, actions: actions, controls: map[int]uintptr{}, pages: map[int]int{}, scale: a.scale}
	s.draft.Lists = append([]config.Target{}, c.Lists...)
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.brush = call(gdi32, "CreateSolidBrush", 0)
	s.font = call(gdi32, "CreateFontW", uintptr(-int32(15*s.scale)), 0, 0, 0, 400, 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(u("Segoe UI"))))
	a.settings = s
	s.button(119, -1, "← Zurück", "← Back", 20, 16, 105)
	s.control(206, -1, "STATIC", "Simkl Bridge", 145, 22, 250, 30, 0)
	s.button(10, -1, "Anmeldung", "Sign-in", 20, 66, 120)
	s.button(11, -1, "Listen", "Lists", 150, 66, 120)
	s.button(12, -1, "Allgemein", "General", 280, 66, 120)
	s.control(200, 0, "STATIC", s.label("Client-ID deiner AUTH-V2-App", "Your AUTH V2 app client ID"), 20, 124, 380, 24, 0)
	s.control(settingClient, 0, "EDIT", c.ClientID, 20, 156, 380, 32, 0x10000|0x800000|0x80)
	s.button(settingConnect, 0, "Mit Simkl verbinden", "Connect to Simkl", 20, 204, 180)
	s.button(settingBrowser, 0, "Browser öffnen", "Open browser", 220, 204, 180)
	call(user32, "EnableWindow", s.controls[settingBrowser], 0)
	s.control(201, 0, "STATIC", s.label("Eine neue AUTH-V2-App als „TV, devices & command line“ erstellen. Hier die Client-ID einfügen und den Code im Browser bestätigen.\r\n\r\nFür eigene Listen brauchst du Simkl PRO oder VIP.", "Create an AUTH V2 app using ‘TV, devices & command line’. Paste its client ID here and approve the code in your browser.\r\n\r\nCustom Lists require Simkl PRO or VIP."), 20, 265, 380, 150, 0)
	s.button(settingGuide, 0, "Simkl-App registrieren", "Register a Simkl app", 20, 440, 380)
	s.button(settingDiscover, 1, "Meine Listen laden", "Load my lists", 20, 116, 380)
	s.control(settingAvailable, 1, "LISTBOX", "", 20, 160, 380, 82, 0x10000|0x800000|0x200000|1)
	s.button(settingSelectSonarr, 1, "Auswahl → Sonarr", "Selection → Sonarr", 20, 254, 180)
	s.button(settingSelectRadarr, 1, "Auswahl → Radarr", "Selection → Radarr", 220, 254, 180)
	s.control(202, 1, "STATIC", s.label("Oder: Listenlink / Listen-ID", "Or: list link / ID"), 20, 300, 380, 24, 0)
	s.control(settingLink, 1, "EDIT", "", 20, 330, 380, 30, 0x10000|0x800000|0x80)
	s.button(settingAddSonarr, 1, "Link → Sonarr", "Link → Sonarr", 20, 372, 180)
	s.button(settingAddRadarr, 1, "Link → Radarr", "Link → Radarr", 220, 372, 180)
	s.control(settingConfigured, 1, "LISTBOX", "", 20, 416, 380, 78, 0x10000|0x800000|0x200000|1)
	s.button(settingCopy, 1, "Adresse kopieren", "Copy address", 20, 505, 180)
	s.button(settingRemove, 1, "Zuordnung entfernen", "Remove mapping", 220, 505, 180)
	s.control(203, 2, "STATIC", s.label("Lokaler Bridge-Port", "Local bridge port"), 20, 130, 230, 30, 0)
	s.control(settingPort, 2, "EDIT", strconv.Itoa(c.Port), 270, 128, 130, 32, 0x10000|0x800000|0x2000)
	s.control(204, 2, "STATIC", s.label("Die Bridge ist nur auf diesem Rechner erreichbar.\r\n\r\nSonarr und Radarr erhalten vollständige Listen. Unveränderte Inhalte werden zwischengespeichert.\r\n\r\nAutostart und Sprache kannst du direkt im Hauptmenü ändern.", "The bridge is accessible only on this computer.\r\n\r\nSonarr and Radarr receive complete lists. Unchanged contents are cached.\r\n\r\nStartup and language are available in the main menu."), 20, 210, 380, 210, 0)
	version := actions.Version
	if version == "" {
		version = "Vorabversion"
	}
	s.control(205, 3, "STATIC", s.label("Über Simkl Bridge", "About Simkl Bridge")+"\r\n"+version+"\r\n\r\n© 2026 Veyilla Nightwhisper\r\n\r\n"+s.label("Verbindet deine Simkl-Listen mit Sonarr und Radarr. Unabhängige Anwendung, kein offizielles Simkl-Produkt.\r\n\r\nQuellen und Lizenzen kannst du unten nachlesen.", "Connects your Simkl lists to Sonarr and Radarr. Independent app, not an official Simkl product.\r\n\r\nRead the sources and licenses below."), 20, 130, 380, 185, 0)
	s.button(settingLicense, 3, "Lizenz lesen", "Read license", 20, 330, 180)
	s.button(settingNotices, 3, "Quellenhinweise", "Third-party notices", 220, 330, 180)
	s.button(settingLicenseFolder, 3, "Weitere Lizenztexte öffnen", "Open other license texts", 20, 385, 380)
	s.button(settingSimkl, 3, "Simkl öffnen", "Open Simkl", 20, 440, 180)
	s.button(settingProject, 3, "GitHub-Projektseite", "GitHub project page", 220, 440, 180)
	if actions.ProjectURL == "" {
		call(user32, "EnableWindow", s.controls[settingProject], 0)
		s.control(208, 3, "STATIC", s.label("Noch keine öffentliche Projektseite verknüpft.", "No public project page linked yet."), 20, 490, 380, 45, 0)
	}
	s.control(207, 4, "EDIT", actions.LogText, 20, 120, 380, 470, 0x10000|0x800000|0x200000|0x4|0x800|0x40)
	s.control(settingStatus, -1, "STATIC", "", 20, 554, 380, 42, 0)
	s.button(settingSave, -1, "Speichern und zurück", "Save and back", 20, 606, 380)
	s.fillConfigured()
	s.showPage(page)
	a.hover = -1
	a.redraw()
	call(user32, "SetFocus", s.controls[settingClient])
}
func (s *settingsWindow) drawButton(l unsafe.Pointer) {
	d := (*drawItem)(l)
	s.app.fill(d.DC, d.Rect, 0x151515)
	call(gdi32, "SelectObject", d.DC, s.font)
	call(gdi32, "SetBkMode", d.DC, 1)
	c := uint32(0xffffff)
	if d.State&4 != 0 {
		c = 0x777777
	}
	call(gdi32, "SetTextColor", d.DC, rgb(c))
	text := windowText(d.Window)
	v := syscall.StringToUTF16(text)
	rect := d.Rect
	call(user32, "DrawTextW", d.DC, uintptr(unsafe.Pointer(&v[0])), uintptr(len(v)-1), uintptr(unsafe.Pointer(&rect)), 0x25|0x800)
	if d.State&0x10 != 0 {
		call(user32, "DrawFocusRect", d.DC, uintptr(unsafe.Pointer(&rect)))
	}
}
