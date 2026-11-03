package desktop

import "simklbridge/internal/diagnostics"

type Status struct {
	Text, Tooltip, Action string
	Color                 uint32
}

func Presentation(v View) Status {
	if v.Phase == "running" && v.Error != "" {
		text := diagnostics.Describe(v.Error, v.Language)
		return Status{Text: text, Tooltip: "Simkl Bridge — " + text, Action: "stop", Color: 0xff6262}
	}
	text, action, color := "Bridge gestoppt", "start", uint32(0xff6262)
	en := "Bridge stopped"
	switch v.Phase {
	case "running":
		text, en, action, color = "Bridge läuft", "Bridge running", "stop", 0x43df87
	case "starting":
		text, en, action, color = "Bridge startet …", "Bridge starting …", "wait", 0xffbe4a
	case "error":
		switch v.Error {
		case "authentication_required":
			text, en = "Simkl-Anmeldung erforderlich", "Simkl sign-in required"
		case "configuration_required":
			text, en = "Einstellungen prüfen", "Check settings"
		case "restart_limit":
			text, en = "Nach wiederholten Fehlern gestoppt", "Stopped after repeated failures"
		default:
			text, en = diagnostics.Describe(v.Error, "de"), diagnostics.Describe(v.Error, "en")
		}
	}
	if v.Language == "en" {
		text = en
	}
	return Status{Text: text, Tooltip: "Simkl Bridge — " + text, Action: action, Color: color}
}
