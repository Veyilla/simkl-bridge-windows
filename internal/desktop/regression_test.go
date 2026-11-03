package desktop

import (
	"image"
	"image/color"
	"testing"
)

func TestLogoTransparencyPreservesSoftEdges(t *testing.T) {
	if got := compositeChannel(0, 32768, 255); got < 126 || got > 128 {
		t.Fatalf("Halbtransparenter schwarzer Rand muss grau bleiben, ist %d", got)
	}
}
func TestHomeRowsNeverCrossSeparatorsOrFooter(t *testing.T) {
	layout := HomeLayout()
	if len(layout.Rows) != 7 {
		t.Fatal("Menüaktionen fehlen")
	}
	for _, row := range layout.Rows {
		for _, separator := range layout.Separators {
			if separator >= row.Y && separator < row.Y+row.H {
				t.Fatalf("Trennlinie durch Menüzeile: %+v bei %d", row, separator)
			}
		}
		if row.Y+row.H > layout.Footer.Y {
			t.Fatalf("Zeile überschneidet Sprachauswahl: %+v", row)
		}
	}
	if layout.Rows[6].Y < layout.Rows[5].Y+layout.Rows[5].H+12 {
		t.Fatal("Beenden nicht sauber getrennt")
	}
}
func TestNativeIconKeepsTopAndBottomOrientation(t *testing.T) {
	source := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			c := color.RGBA{R: 255, A: 255}
			if y >= 16 {
				c = color.RGBA{B: 255, A: 255}
			}
			source.SetRGBA(x, y, c)
		}
	}
	icon, err := nativeIcon(source)
	if err != nil {
		t.Fatal(err)
	}
	defer call(user32, "DestroyIcon", icon)
	pixels, err := iconPixels(icon, 32)
	if err != nil {
		t.Fatal(err)
	}
	if pixels.RGBAAt(16, 3).R < 240 || pixels.RGBAAt(16, 28).B < 240 {
		t.Fatal("Tray-Symbol wurde vertikal gespiegelt")
	}
}
func TestHoverRedrawsOnlyWhenRowChanges(t *testing.T) {
	a := &App{hover: -1}
	redraws := 0
	for _, y := range []int{150, 151, 152, 153, 195, 196} {
		if a.trackHover(y) {
			redraws++
		}
	}
	if redraws != 2 || a.hover != 1 {
		t.Fatalf("%d Neuzeichnungen für zwei betretene Zeilen", redraws)
	}
}
func TestStatusDotHasSoftEdges(t *testing.T) {
	s, err := newSurface(32, 32)
	if err != nil {
		t.Fatal(err)
	}
	defer s.close()
	a := &App{scale: 1}
	a.dot(s.dc, 16, 16, 5, 0xffffff)
	img := s.rgba()
	partial := 0
	for y := 8; y < 24; y++ {
		for x := 8; x < 24; x++ {
			v := img.RGBAAt(x, y).R
			if v > 0 && v < 255 {
				partial++
			}
		}
	}
	if partial == 0 {
		t.Fatal("Statuspunkt hat harte, ausgefranste Kanten")
	}
}
func TestOfficeAndUltraWidePanelsStayInsideEveryWorkArea(t *testing.T) {
	for _, area := range []Rect{{0, 0, 1366, 728}, {0, 0, 1920, 1040}, {0, 0, 3840, 2120}, {0, 0, 5120, 1400}, {-3840, 40, 3840, 2120}, {40, 0, 1880, 1080}} {
		for _, scale := range []float64{1, 1.25, 1.5, 2} {
			for _, size := range []Size{{368, 552}, {420, 650}} {
				r, actual := FitPanel(Rect{area.X + area.W - 2, area.Y + area.H - 2, 1, 1}, area, size, scale)
				if r.X < area.X || r.Y < area.Y || r.X+r.W > area.X+area.W || r.Y+r.H > area.Y+area.H || actual <= 0 || actual > scale {
					t.Fatalf("Panel außerhalb des Bildschirms: %+v bei %+v", r, area)
				}
			}
		}
	}
}
