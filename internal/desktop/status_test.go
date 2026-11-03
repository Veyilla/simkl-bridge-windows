package desktop

import "testing"

func TestStatusExplainsStateAndAction(t *testing.T) {
	for _, v := range []struct {
		phase, code, text, action string
		color                     uint32
	}{{"running", "", "Bridge läuft", "stop", 0x43df87}, {"stopped", "", "Bridge gestoppt", "start", 0xff6262}, {"starting", "", "Bridge startet …", "wait", 0xffbe4a}, {"error", "authentication_required", "Simkl-Anmeldung erforderlich", "start", 0xff6262}, {"error", "restart_limit", "Nach wiederholten Fehlern gestoppt", "start", 0xff6262}} {
		s := Presentation(View{Phase: v.phase, Error: v.code, Language: "de"})
		if s.Text != v.text || s.Action != v.action || s.Color != v.color || s.Tooltip != "Simkl Bridge — "+v.text {
			t.Fatalf("Zustand %+v: %+v", v, s)
		}
	}
	en := Presentation(View{Phase: "error", Error: "authentication_required", Language: "en"})
	if en.Text != "Simkl sign-in required" {
		t.Fatal(en.Text)
	}
}
