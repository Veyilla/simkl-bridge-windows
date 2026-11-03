package desktop

import "math"

func FitPanel(anchor, work Rect, size Size, requestedScale float64) (Rect, float64) {
	scale := math.Min(requestedScale, math.Min(float64(work.W-12)/float64(size.W), float64(work.H-12)/float64(size.H)))
	return Place(anchor, work, Size{int(float64(size.W) * scale), int(float64(size.H) * scale)}), scale
}

func LogicalPoint(x, y, width, height int) (int, int) {
	if width <= 0 || height <= 0 {
		return 0, 0
	}
	layout := HomeLayout()
	return x * layout.Width / width, y * layout.Height / height
}

type Rect struct{ X, Y, W, H int }
type Size struct{ W, H int }

func Place(anchor Rect, area Rect, desired Size) Rect {
	w, h := min(desired.W, area.W), min(desired.H, area.H)
	return Rect{max(area.X, min(anchor.X+anchor.W-w, area.X+area.W-w)), max(area.Y, min(anchor.Y-h, area.Y+area.H-h)), w, h}
}

type Palette struct{ Background, Text, Outline uint32 }

func PopupPalette() Palette { return Palette{0, 0xFFFFFF, 0x777777} }
