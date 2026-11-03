package desktop

import (
	"context"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRenderMenuAtCommonScales(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	checked := false
	err := RunView(ctx, View{Phase: "running", Language: "de"}, Actions{}, func(a *App) {
		a.Dispatch(func() {
			for _, scale := range []float64{1, 1.25, 1.5, 2} {
				a.fonts(scale)
				w, h := int(368*scale), int(552*scale)
				call(user32, "SetWindowPos", a.hwnd, 0, 0, 0, uintptr(w), uintptr(h), 6)
				s, err := newSurface(w, h)
				if err != nil {
					t.Error(err)
					continue
				}
				a.draw(s.dc)
				img := s.rgba()
				s.close()
				if img.Bounds().Dx() != w || img.RGBAAt(w/2, h/2).A != 255 {
					t.Error("Unvollständige Zeichenfläche")
				}
				if dir := os.Getenv("SIMKL_BRIDGE_RENDER_DIR"); dir != "" {
					os.MkdirAll(dir, 0700)
					path := filepath.Join(dir, fmt.Sprintf("Menue-%d-Prozent.png", int(scale*100)))
					f, err := os.Create(path)
					if err != nil {
						t.Error(err)
					} else {
						png.Encode(f, img)
						f.Close()
					}
				}
			}
			t.Logf("Logo-Messpunkte: %06x %06x %06x", IconPixel(16, 10), IconPixel(16, 16), IconPixel(16, 8))
			checked = true
			a.Close()
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checked {
		t.Fatal("Darstellung wurde nicht geprüft")
	}
}
