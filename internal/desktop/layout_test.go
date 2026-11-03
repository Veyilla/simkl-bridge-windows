package desktop

import "testing"

func TestPlaceStaysInsideWorkArea(t *testing.T) {
	tests := []struct {
		name         string
		anchor, area Rect
		size         Size
		want         Rect
	}{
		{"top left", Rect{0, 0, 20, 20}, Rect{0, 0, 1440, 900}, Size{320, 400}, Rect{0, 0, 320, 400}},
		{"bottom right", Rect{1420, 880, 20, 20}, Rect{0, 0, 1440, 900}, Size{320, 400}, Rect{1120, 480, 320, 400}},
		{"left monitor", Rect{-1920, 1060, 20, 20}, Rect{-1920, 0, 1920, 1080}, Size{320, 400}, Rect{-1920, 660, 320, 400}},
		{"scale 150", Rect{1900, 1060, 20, 20}, Rect{0, 0, 1920, 1080}, Size{480, 600}, Rect{1440, 460, 480, 600}},
		{"scale 200", Rect{2540, 1420, 20, 20}, Rect{0, 0, 2560, 1440}, Size{640, 800}, Rect{1920, 620, 640, 800}},
		{"small screen", Rect{300, 200, 20, 20}, Rect{0, 0, 320, 240}, Size{400, 500}, Rect{0, 0, 320, 240}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Place(tt.anchor, tt.area, tt.size); got != tt.want {
				t.Fatalf("got %+v; want %+v", got, tt.want)
			}
		})
	}
}

func TestPopupPalette(t *testing.T) {
	if p := PopupPalette(); p.Background != 0x000000 || p.Text != 0xFFFFFF || p.Outline != 0x777777 {
		t.Fatalf("unreadable popup palette: %+v", p)
	}
	if IconPixel(0, 0) != 0x000000 || IconPixel(16, 16) != 0xFFFFFF {
		t.Fatal("icon corners must be black and the S must remain white")
	}
}

func TestScaledPointerKeepsMenuActionsAligned(t *testing.T) {
	for _, v := range []struct{ x, y, w, h, wantX, wantY int }{{40, 240, 736, 1104, 20, 120}, {30, 180, 552, 828, 20, 120}, {20, 120, 368, 552, 20, 120}} {
		x, y := LogicalPoint(v.x, v.y, v.w, v.h)
		if x != v.wantX || y != v.wantY {
			t.Fatalf("Falsche Aktion bei Skalierung: %d,%d", x, y)
		}
	}
}
