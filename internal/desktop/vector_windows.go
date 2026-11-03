package desktop

import (
	"math"
	"sync"
	"syscall"
	"unsafe"
)

var gdiplus = syscall.NewLazyDLL("gdiplus.dll")
var vectorOnce sync.Once
var vectorToken uintptr

type vectorCanvas struct {
	graphics uintptr
	scale    float64
}

func vector(dc uintptr, scale float64) *vectorCanvas {
	vectorOnce.Do(func() {
		input := struct {
			Version            uint32
			Debug              uintptr
			Background, Codecs int32
		}{Version: 1}
		call(gdiplus, "GdiplusStartup", uintptr(unsafe.Pointer(&vectorToken)), uintptr(unsafe.Pointer(&input)), 0)
	})
	v := &vectorCanvas{scale: scale}
	call(gdiplus, "GdipCreateFromHDC", dc, uintptr(unsafe.Pointer(&v.graphics)))
	call(gdiplus, "GdipSetSmoothingMode", v.graphics, 4)
	return v
}
func (v *vectorCanvas) close()           { call(gdiplus, "GdipDeleteGraphics", v.graphics) }
func (v *vectorCanvas) px(n int) uintptr { return uintptr(int(math.Round(float64(n) * v.scale))) }
func (v *vectorCanvas) pen(color uint32) uintptr {
	var p uintptr
	call(gdiplus, "GdipCreatePen1", uintptr(0xff000000|color), uintptr(math.Float32bits(float32(1.3*v.scale))), 2, uintptr(unsafe.Pointer(&p)))
	return p
}
func (v *vectorCanvas) line(points []point, color uint32) {
	if len(points) < 2 {
		return
	}
	p := v.pen(color)
	defer call(gdiplus, "GdipDeletePen", p)
	for i := 1; i < len(points); i++ {
		call(gdiplus, "GdipDrawLineI", v.graphics, p, v.px(int(points[i-1].X)), v.px(int(points[i-1].Y)), v.px(int(points[i].X)), v.px(int(points[i].Y)))
	}
}
func (v *vectorCanvas) ellipse(x, y, width, height int, color uint32, fill bool) {
	var brush uintptr
	if fill {
		call(gdiplus, "GdipCreateSolidFill", uintptr(0xff000000|color), uintptr(unsafe.Pointer(&brush)))
		defer call(gdiplus, "GdipDeleteBrush", brush)
		call(gdiplus, "GdipFillEllipseI", v.graphics, brush, v.px(x), v.px(y), v.px(width), v.px(height))
	} else {
		p := v.pen(color)
		defer call(gdiplus, "GdipDeletePen", p)
		call(gdiplus, "GdipDrawEllipseI", v.graphics, p, v.px(x), v.px(y), v.px(width), v.px(height))
	}
}
func (v *vectorCanvas) arc(x, y, width, height int, start, sweep float32, color uint32) {
	p := v.pen(color)
	defer call(gdiplus, "GdipDeletePen", p)
	call(gdiplus, "GdipDrawArcI", v.graphics, p, v.px(x), v.px(y), v.px(width), v.px(height), uintptr(math.Float32bits(start)), uintptr(math.Float32bits(sweep)))
}
func (v *vectorCanvas) roundRect(x, y, width, height, radius int, color uint32) {
	var path uintptr
	call(gdiplus, "GdipCreatePath", 0, uintptr(unsafe.Pointer(&path)))
	defer call(gdiplus, "GdipDeletePath", path)
	d := radius * 2
	for _, c := range []struct {
		x, y  int
		start float32
	}{{x, y, 180}, {x + width - d, y, 270}, {x + width - d, y + height - d, 0}, {x, y + height - d, 90}} {
		call(gdiplus, "GdipAddPathArcI", path, v.px(c.x), v.px(c.y), v.px(d), v.px(d), uintptr(math.Float32bits(c.start)), uintptr(math.Float32bits(90)))
	}
	call(gdiplus, "GdipClosePathFigure", path)
	p := v.pen(color)
	defer call(gdiplus, "GdipDeletePen", p)
	call(gdiplus, "GdipDrawPath", v.graphics, p, path)
}
