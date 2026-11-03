package desktop

import (
	"context"
	"simklbridge/internal/config"
	"strings"
	"testing"
	"time"
)

func TestAboutShowsAuthorAndOpensLicenseDocuments(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var opened []string
	checked := false
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) {
		a.ShowAbout(config.Default(), SettingsActions{ReadOnly: true, OpenDocument: func(name string) error {
			opened = append(opened, name)
			return nil
		}})
		a.Dispatch(func() {
			s := a.settings
			text := windowText(s.controls[205])
			if !strings.Contains(text, "Veyilla Nightwhisper") || strings.Contains(text, "Prozessüberwachung") {
				t.Error("Über enthält nicht die verständlichen Autorenangaben")
			}
			for _, id := range []int{settingLicense, settingNotices, settingLicenseFolder} {
				if s.controls[id] == 0 {
					t.Error("Lizenzschaltfläche fehlt")
				}
				s.command(id)
			}
			checked = len(opened) == 3 && opened[0] == "LICENSE" && opened[1] == "THIRD_PARTY_NOTICES.md" && opened[2] == "licenses"
			a.Close()
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checked {
		t.Fatal("Lizenzschaltflächen öffnen nicht die vorgesehenen Unterlagen")
	}
}
