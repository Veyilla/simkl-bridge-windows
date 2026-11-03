//go:build windows

package desktop

import (
	"context"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

var user32 = syscall.NewLazyDLL("user32.dll")
var gdi32 = syscall.NewLazyDLL("gdi32.dll")
var shell32 = syscall.NewLazyDLL("shell32.dll")
var kernel32 = syscall.NewLazyDLL("kernel32.dll")

func u(s string) *uint16 { p, _ := syscall.UTF16PtrFromString(s); return p }

//go:uintptrescapes
func call(d *syscall.LazyDLL, n string, a ...uintptr) uintptr {
	r, _, _ := d.NewProc(n).Call(a...)
	return r
}

type winRect struct{ Left, Top, Right, Bottom int32 }
type point struct{ X, Y int32 }
type message struct {
	Hwnd    uintptr
	Msg     uint32
	W, L    uintptr
	Time    uint32
	P       point
	Private uint32
}
type windowClass struct {
	Style                              uint32
	Proc                               uintptr
	ClassExtra, WindowExtra            int32
	Instance, Icon, Cursor, Background uintptr
	Menu, Name                         *uint16
}
type paintInfo struct {
	DC              uintptr
	Erase           int32
	Rect            winRect
	Restore, Update int32
	Reserved        [32]byte
}
type notifyData struct {
	Size                uint32
	Hwnd                uintptr
	ID, Flags, Callback uint32
	Icon                uintptr
	Tip                 [128]uint16
	State, Mask         uint32
	Info                [256]uint16
	Timeout             uint32
	Title               [64]uint16
	InfoFlags           uint32
	GUID                [16]byte
	Balloon             uintptr
}

type Actions struct {
	Start, Stop, Restart, OpenSettings, OpenLog, OpenAbout, Quit func()
	SetLanguage                                                  func(string)
	SetStartup                                                   func(bool)
}
type View struct {
	Phase, Language, Error string
	Startup                bool
}
type App struct {
	hwnd                     uintptr
	view                     View
	actions                  Actions
	font, bold, icon         uintptr
	small, title, headerIcon uintptr
	scale                    float64
	buffer                   *surface
	anchor                   point
	anchored                 bool
	hover                    int
	visible                  bool
	notify                   notifyData
	queue                    chan func()
	done                     chan struct{}
	settings                 *settingsWindow
}

func rgb(c uint32) uintptr { return uintptr(c&0xFF)<<16 | uintptr(c&0xFF00) | uintptr((c>>16)&0xFF) }
func (a *App) fill(dc uintptr, r winRect, c uint32) {
	b := call(gdi32, "CreateSolidBrush", rgb(c))
	call(user32, "FillRect", dc, uintptr(unsafe.Pointer(&r)), b)
	call(gdi32, "DeleteObject", b)
}
func (a *App) text(dc uintptr, s string, x, y int, bold bool, c uint32) {
	f := a.font
	if bold {
		f = a.bold
	}
	old := call(gdi32, "SelectObject", dc, f)
	call(gdi32, "SetTextColor", dc, rgb(c))
	call(gdi32, "SetBkMode", dc, 1)
	v := syscall.StringToUTF16(s)
	call(gdi32, "TextOutW", dc, uintptr(x), uintptr(y), uintptr(unsafe.Pointer(&v[0])), uintptr(len(v)-1))
	call(gdi32, "SelectObject", dc, old)
}
func (a *App) label(de, en string) string {
	if a.view.Language == "en" {
		return en
	}
	return de
}
func (a *App) status() (string, uint32) {
	s := Presentation(a.view)
	return s.Text, s.Color
}

func invoke(f func()) {
	if f != nil {
		go f()
	}
}
func (a *App) click(x, y int) {
	if y >= 510 && y < 542 {
		lang := "de"
		if x >= 100 {
			lang = "en"
		}
		a.view.Language = lang
		if a.actions.SetLanguage != nil {
			a.actions.SetLanguage(lang)
		}
		a.redraw()
		return
	}
	i := rowAt(y)
	if i < 0 {
		return
	}
	switch i {
	case 0:
		if a.view.Phase == "starting" {
			return
		}
		if a.view.Phase == "running" {
			invoke(a.actions.Stop)
		} else {
			invoke(a.actions.Start)
		}
	case 1:
		if a.view.Phase == "starting" {
			return
		}
		invoke(a.actions.Restart)
	case 2:
		if a.actions.SetStartup != nil {
			a.actions.SetStartup(!a.view.Startup)
		}
	case 3:
		invoke(a.actions.OpenSettings)
	case 4:
		invoke(a.actions.OpenLog)
	case 5:
		invoke(a.actions.OpenAbout)
	case 6:
		invoke(a.actions.Quit)
	}
	a.redraw()
}
func (a *App) redraw() { call(user32, "InvalidateRect", a.hwnd, 0, 0) }
func (a *App) show() {
	if !a.visible {
		a.anchored = false
	}
	size := Size{368, 552}
	if a.settings != nil {
		size = Size{420, 650}
	}
	a.position(size)
}
func (a *App) hide() {
	if a.settings != nil {
		return
	}
	call(user32, "ShowWindow", a.hwnd, 0)
	a.visible = false
}
func (a *App) Update(v View) { a.Dispatch(func() { a.update(v) }) }
func (a *App) update(v View) {
	a.view = v
	status, _ := a.status()
	tip := syscall.StringToUTF16("Simkl Bridge — " + status)
	clear(a.notify.Tip[:])
	copy(a.notify.Tip[:], tip)
	call(shell32, "Shell_NotifyIconW", 1, uintptr(unsafe.Pointer(&a.notify)))
	a.redraw()
}
func (a *App) Close() { call(user32, "PostMessageW", a.hwnd, 0x10, 0, 0) }

func (a *App) ShowMenu() { a.Dispatch(a.show) }

// Dispatch hält sämtliche Fenster- und Darstellungsänderungen auf dem Fensterthread.
func (a *App) Dispatch(f func()) {
	select {
	case <-a.done:
		return
	case a.queue <- f:
		call(user32, "PostMessageW", a.hwnd, 0x8002, 0, 0)
	}
}

func Run(ctx context.Context, actions Actions) error {
	return RunView(ctx, View{Phase: "stopped", Language: "de"}, actions, nil)
}
func RunView(ctx context.Context, view View, actions Actions, ready func(*App)) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if p := user32.NewProc("SetThreadDpiAwarenessContext"); p.Find() == nil {
		old, _, _ := p.Call(^uintptr(3))
		if old != 0 {
			defer p.Call(old)
		}
	}
	a := &App{view: view, actions: actions, hover: -1, queue: make(chan func(), 128), done: make(chan struct{})}
	defer close(a.done)
	defer func() {
		if a.settings != nil {
			call(user32, "UnregisterClassW", uintptr(unsafe.Pointer(u("SimklBridgeSettingsWindow"))), call(kernel32, "GetModuleHandleW", 0))
		}
	}()
	proc := syscall.NewCallback(func(hwnd uintptr, msg uint32, w uintptr, pointer unsafe.Pointer) uintptr {
		l := uintptr(pointer)
		switch msg {
		case 0x14:
			return 1
		case 0x111:
			if a.settings != nil && w>>16 == 0 {
				a.settings.command(int(w & 0xffff))
			}
			return 0
		case 0x138, 0x133, 0x134:
			if a.settings != nil {
				call(gdi32, "SetTextColor", w, rgb(0xffffff))
				call(gdi32, "SetBkColor", w, 0)
				return a.settings.brush
			}
		case 0x2b:
			if a.settings != nil {
				a.settings.drawButton(pointer)
				return 1
			}
		case 0x8002:
			select {
			case f := <-a.queue:
				f()
			default:
			}
			return 0
		case 0xf:
			var p paintInfo
			dc := call(user32, "BeginPaint", hwnd, uintptr(unsafe.Pointer(&p)))
			a.paint(dc)
			call(user32, "EndPaint", hwnd, uintptr(unsafe.Pointer(&p)))
			return 0
		case 0x200:
			if a.settings != nil {
				return 0
			}
			_, y := a.logical(l)
			if a.trackHover(y) {
				a.redraw()
			}
			return 0
		case 0x202:
			if a.settings != nil {
				return 0
			}
			x, y := a.logical(l)
			a.click(x, y)
			return 0
		case 0x100:
			if w == 27 {
				if a.settings != nil {
					a.settings.hide()
					return 0
				}
				a.hide()
				return 0
			}
			if w == 38 {
				a.hover = (a.hover + 8) % 9
				a.redraw()
				return 0
			}
			if w == 40 {
				a.hover = (a.hover + 1) % 9
				a.redraw()
				return 0
			}
			if w == 13 && a.hover >= 0 {
				if a.hover >= 7 {
					x := 24
					if a.hover == 8 {
						x = 110
					}
					a.click(x, 522)
					return 0
				}
				y := HomeLayout().Rows[a.hover].Y + 20
				a.click(60, y)
				return 0
			}
		case 0x1c:
			if w == 0 {
				a.hide()
			}
			return 0
		case 0x8001:
			if uint32(l) == 0x205 || uint32(l) == 0x202 {
				a.show()
			}
			return 0
		case 0x10:
			if a.settings != nil {
				a.settings.shutdown()
			}
			call(user32, "DestroyWindow", hwnd)
			return 0
		case 2:
			call(shell32, "Shell_NotifyIconW", 2, uintptr(unsafe.Pointer(&a.notify)))
			call(user32, "PostQuitMessage", 0)
			return 0
		}
		return call(user32, "DefWindowProcW", hwnd, uintptr(msg), w, l)
	})
	instance := call(kernel32, "GetModuleHandleW", 0)
	name := u("SimklBridgeTrayWindow")
	wc := windowClass{Proc: proc, Instance: instance, Cursor: call(user32, "LoadCursorW", 0, 32512), Name: name}
	if call(user32, "RegisterClassW", uintptr(unsafe.Pointer(&wc))) == 0 {
		return fmt.Errorf("Windows-Fensterklasse konnte nicht erstellt werden")
	}
	defer call(user32, "UnregisterClassW", uintptr(unsafe.Pointer(name)), instance)
	a.hwnd = call(user32, "CreateWindowExW", 0x88, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(u("Simkl Bridge"))), 0x82000000, 0, 0, 368, 552, 0, 0, instance, 0)
	if a.hwnd == 0 {
		return fmt.Errorf("Windows-Menü konnte nicht erstellt werden")
	}
	a.fonts(1)
	defer a.releaseDrawing()
	var iconErr error
	a.icon, iconErr = nativeIcon(TrayLogoImage(64))
	if iconErr != nil {
		return iconErr
	}
	defer call(user32, "DestroyIcon", a.icon)
	a.notify = notifyData{Hwnd: a.hwnd, ID: 1, Flags: 7, Callback: 0x8001, Icon: a.icon}
	a.notify.Size = uint32(unsafe.Sizeof(a.notify))
	a.update(view)
	if call(shell32, "Shell_NotifyIconW", 0, uintptr(unsafe.Pointer(&a.notify))) == 0 {
		call(user32, "DestroyWindow", a.hwnd)
		return fmt.Errorf("Tray-Symbol konnte nicht erstellt werden")
	}
	if ready != nil {
		ready(a)
	}
	go func() {
		select {
		case <-ctx.Done():
			a.Close()
		case <-a.done:
		}
	}()
	var m message
	for {
		r := call(user32, "GetMessageW", uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if r == 0 || r == ^uintptr(0) {
			break
		}
		if a.settings != nil && call(user32, "IsWindowVisible", a.settings.hwnd) != 0 && call(user32, "IsDialogMessageW", a.settings.hwnd, uintptr(unsafe.Pointer(&m))) != 0 {
			continue
		}
		call(user32, "TranslateMessage", uintptr(unsafe.Pointer(&m)))
		call(user32, "DispatchMessageW", uintptr(unsafe.Pointer(&m)))
	}
	return nil
}

func WakeExisting() {
	hwnd := call(user32, "FindWindowW", uintptr(unsafe.Pointer(u("SimklBridgeTrayWindow"))), 0)
	if hwnd != 0 {
		call(user32, "PostMessageW", hwnd, 0x8001, 0, 0x205)
	}
}
