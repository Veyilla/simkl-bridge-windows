package desktop

import (
	"syscall"
	"unsafe"
)

type monitorInfo struct {
	Size          uint32
	Monitor, Work winRect
	Flags         uint32
}

func screenAt(p point) (Rect, float64) {
	packed := uintptr(uint64(uint32(p.X)) | uint64(uint32(p.Y))<<32)
	monitor := call(user32, "MonitorFromPoint", packed, 2)
	info := monitorInfo{}
	info.Size = uint32(unsafe.Sizeof(info))
	if call(user32, "GetMonitorInfoW", monitor, uintptr(unsafe.Pointer(&info))) == 0 {
		call(user32, "SystemParametersInfoW", 0x30, 0, uintptr(unsafe.Pointer(&info.Work)), 0)
	}
	scale := 1.0
	shcore := syscall.NewLazyDLL("shcore.dll")
	proc := shcore.NewProc("GetDpiForMonitor")
	if proc.Find() == nil {
		var x, y uint32
		if hr, _, _ := proc.Call(monitor, 0, uintptr(unsafe.Pointer(&x)), uintptr(unsafe.Pointer(&y))); hr == 0 && x > 0 {
			scale = float64(x) / 96
		}
	}
	return Rect{int(info.Work.Left), int(info.Work.Top), int(info.Work.Right - info.Work.Left), int(info.Work.Bottom - info.Work.Top)}, scale
}
func (a *App) logical(l uintptr) (int, int) {
	var r winRect
	call(user32, "GetClientRect", a.hwnd, uintptr(unsafe.Pointer(&r)))
	return LogicalPoint(int(int16(l)), int(int16(l>>16)), int(r.Right), int(r.Bottom))
}
