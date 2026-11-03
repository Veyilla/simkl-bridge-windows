package desktop

import (
	"math"
	"syscall"
	"unsafe"
)

func (a *App) px(v int) uintptr { return uintptr(int(math.Round(float64(v) * a.scale))) }
func (a *App) rect(r Rect) winRect {
	return winRect{int32(a.px(r.X)), int32(a.px(r.Y)), int32(a.px(r.X + r.W)), int32(a.px(r.Y + r.H))}
}
func (a *App) drawText(dc uintptr, s string, x, y int, font uintptr, c uint32) {
	old := call(gdi32, "SelectObject", dc, font)
	call(gdi32, "SetTextColor", dc, rgb(c))
	call(gdi32, "SetBkMode", dc, 1)
	v := syscall.StringToUTF16(s)
	call(gdi32, "TextOutW", dc, a.px(x), a.px(y), uintptr(unsafe.Pointer(&v[0])), uintptr(len(v)-1))
	call(gdi32, "SelectObject", dc, old)
}
func (a *App) fonts(scale float64) {
	if a.font != 0 && math.Abs(a.scale-scale) < 0.001 {
		return
	}
	for _, h := range []uintptr{a.font, a.bold, a.small, a.title} {
		if h != 0 {
			call(gdi32, "DeleteObject", h)
		}
	}
	if a.headerIcon != 0 {
		call(user32, "DestroyIcon", a.headerIcon)
	}
	a.scale = scale
	create := func(height, weight int) uintptr {
		return call(gdi32, "CreateFontW", uintptr(-int32(float64(height)*scale)), 0, 0, 0, uintptr(weight), 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(u("Segoe UI"))))
	}
	a.font = create(16, 400)
	a.bold = create(16, 600)
	a.small = create(14, 400)
	a.title = create(18, 600)
	a.headerIcon, _ = nativeIcon(LogoImage(max(16, int(40*scale))))
}
func (a *App) releaseDrawing() {
	for _, h := range []uintptr{a.font, a.bold, a.small, a.title} {
		if h != 0 {
			call(gdi32, "DeleteObject", h)
		}
	}
	if a.headerIcon != 0 {
		call(user32, "DestroyIcon", a.headerIcon)
	}
	if a.buffer != nil {
		a.buffer.close()
	}
}
func (a *App) paint(dc uintptr) {
	var r winRect
	call(user32, "GetClientRect", a.hwnd, uintptr(unsafe.Pointer(&r)))
	w, h := int(r.Right), int(r.Bottom)
	if a.buffer == nil || a.buffer.width != w || a.buffer.height != h {
		if a.buffer != nil {
			a.buffer.close()
		}
		a.buffer, _ = newSurface(w, h)
	}
	if a.buffer == nil {
		return
	}
	a.draw(a.buffer.dc)
	call(gdi32, "BitBlt", dc, 0, 0, uintptr(w), uintptr(h), a.buffer.dc, 0, 0, 0x00cc0020)
}
func (a *App) draw(dc uintptr) {
	var client winRect
	call(user32, "GetClientRect", a.hwnd, uintptr(unsafe.Pointer(&client)))
	a.fill(dc, client, 0)
	v := vector(dc, a.scale)
	v.roundRect(0, 0, int(float64(client.Right)/a.scale)-1, int(float64(client.Bottom)/a.scale)-1, 12, 0x303030)
	v.close()
	if a.settings != nil {
		return
	}
	layout := HomeLayout()
	call(user32, "DrawIconEx", dc, a.px(24), a.px(24), a.headerIcon, a.px(40), a.px(40), 0, 0, 3)
	v = vector(dc, a.scale)
	v.roundRect(20, 20, 48, 48, 10, 0x777777)
	v.close()
	a.drawText(dc, "Simkl Bridge", 82, 24, a.title, 0xffffff)
	a.drawText(dc, a.label("Deine Listen. Direkt verbunden.", "Your lists. Directly connected."), 82, 50, a.small, 0xb8b8b8)
	status, col := a.status()
	a.dot(dc, 28, 99, 4, col)
	a.drawText(dc, status, 44, 88, a.font, 0xffffff)
	for _, y := range layout.Separators {
		a.fill(dc, a.rect(Rect{20, y, 328, 1}), 0x292929)
	}
	labels := []string{a.label("Bridge starten", "Start bridge"), a.label("Neu starten", "Restart bridge"), a.label("Mit Windows starten", "Start with Windows"), a.label("Einstellungen", "Settings"), a.label("Protokoll öffnen", "Open log"), a.label("Über", "About"), a.label("Beenden", "Quit")}
	if a.view.Phase == "running" {
		labels[0] = a.label("Bridge stoppen", "Stop bridge")
	}
	if a.view.Phase == "starting" {
		labels[0] = a.label("Bridge startet …", "Bridge starting …")
	}
	for i, row := range layout.Rows {
		if a.hover == i {
			a.fill(dc, a.rect(row), 0x151515)
		}
		a.symbol(dc, i, 28, row.Y+21)
		a.drawText(dc, labels[i], 56, row.Y+10, a.font, 0xffffff)
		if i == 2 && a.view.Startup {
			a.stroke(dc, []point{{326, int32(row.Y + 21)}, {330, int32(row.Y + 25)}, {338, int32(row.Y + 15)}}, 0xffffff)
		}
	}
	a.drawText(dc, "Deutsch", 24, 514, a.small, 0xffffff)
	a.drawText(dc, "English", 102, 514, a.small, 0xc4c4c4)
	x, width := 24, 50
	if a.view.Language == "en" {
		x, width = 102, 45
	}
	a.fill(dc, a.rect(Rect{x, 534, width, 1}), 0xffffff)
	if a.hover == 7 || a.hover == 8 {
		x = 22
		if a.hover == 8 {
			x = 100
		}
		r := a.rect(Rect{x, 510, 64, 30})
		call(user32, "DrawFocusRect", dc, uintptr(unsafe.Pointer(&r)))
	}
}
func (a *App) dot(dc uintptr, x, y, r int, c uint32) {
	v := vector(dc, a.scale)
	defer v.close()
	v.ellipse(x-r, y-r, 2*r, 2*r, c, true)
}
func (a *App) stroke(dc uintptr, points []point, c uint32) {
	v := vector(dc, a.scale)
	defer v.close()
	v.line(points, c)
}
func (a *App) symbol(dc uintptr, id, x, y int) {
	line := func(coords ...int) {
		p := []point{}
		for i := 0; i < len(coords); i += 2 {
			p = append(p, point{int32(x + coords[i]), int32(y + coords[i+1])})
		}
		a.stroke(dc, p, 0xeeeeee)
	}
	switch id {
	case 0:
		if a.view.Phase == "running" {
			line(-6, -6, 6, -6, 6, 6, -6, 6, -6, -6)
		} else {
			line(-5, -7, 7, 0, -5, 7, -5, -7)
		}
	case 1:
		v := vector(dc, a.scale)
		v.arc(x-7, y-7, 14, 14, 35, 300, 0xeeeeee)
		v.close()
		line(1, -7, 7, -7, 7, -1)
	case 2:
		line(0, -9, 0, -1)
		v := vector(dc, a.scale)
		v.arc(x-7, y-7, 14, 14, -50, 280, 0xeeeeee)
		v.close()
	case 3:
		line(-8, -4, -3, -4)
		line(1, -4, 8, -4)
		line(-8, 5, 3, 5)
		line(7, 5, 8, 5)
		line(-3, -7, 1, -7, 1, -1, -3, -1, -3, -7)
		line(3, 2, 7, 2, 7, 8, 3, 8, 3, 2)
	case 4:
		line(-6, -8, 3, -8, 7, -4, 7, 8, -6, 8, -6, -8)
		line(2, -8, 2, -3, 7, -3)
		line(-3, 1, 4, 1)
		line(-3, 4, 4, 4)
	case 5:
		v := vector(dc, a.scale)
		v.ellipse(x-7, y-7, 14, 14, 0xeeeeee, false)
		v.close()
		line(0, -1, 0, 4)
		a.dot(dc, x, y-4, 1, 0xeeeeee)
	case 6:
		line(-1, -7, -7, -7, -7, 7, -1, 7)
		line(-2, 0, 8, 0, 4, -4)
		line(8, 0, 4, 4)
	}
}
func rowAt(y int) int {
	for i, r := range HomeLayout().Rows {
		if y >= r.Y && y < r.Y+r.H {
			return i
		}
	}
	return -1
}
func (a *App) trackHover(y int) bool {
	next := rowAt(y)
	if y >= 510 && y < 542 {
		next = 7
	}
	if next == a.hover {
		return false
	}
	a.hover = next
	return true
}
func (a *App) position(size Size) {
	if !a.anchored {
		call(user32, "GetCursorPos", uintptr(unsafe.Pointer(&a.anchor)))
		a.anchored = true
	}
	work, scale := screenAt(a.anchor)
	r, scale := FitPanel(Rect{int(a.anchor.X), int(a.anchor.Y), 1, 1}, work, size, scale)
	a.fonts(scale)
	call(user32, "SetWindowPos", a.hwnd, ^uintptr(0), uintptr(r.X), uintptr(r.Y), uintptr(r.W), uintptr(r.H), 0x40)
	call(user32, "SetForegroundWindow", a.hwnd)
	a.visible = true
}
