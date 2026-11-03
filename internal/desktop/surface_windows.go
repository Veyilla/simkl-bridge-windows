package desktop

import (
	"errors"
	"image"
	"image/color"
	"unsafe"
)

type bitmapHeader struct {
	Size                   uint32
	Width, Height          int32
	Planes, BitCount       uint16
	Compression, ImageSize uint32
	XPels, YPels           int32
	Colors, Important      uint32
}
type surface struct {
	dc, bitmap, old uintptr
	bits            unsafe.Pointer
	width, height   int
}

func newSurface(width, height int) (*surface, error) {
	if width <= 0 || height <= 0 || width > 8192 || height > 8192 {
		return nil, errors.New("Ungültige Zeichenfläche")
	}
	s := &surface{width: width, height: height}
	s.dc = call(gdi32, "CreateCompatibleDC", 0)
	header := bitmapHeader{Size: 40, Width: int32(width), Height: -int32(height), Planes: 1, BitCount: 32}
	s.bitmap = call(gdi32, "CreateDIBSection", s.dc, uintptr(unsafe.Pointer(&header)), 0, uintptr(unsafe.Pointer(&s.bits)), 0, 0)
	if s.bitmap == 0 || s.bits == nil {
		call(gdi32, "DeleteDC", s.dc)
		return nil, errors.New("Windows-Zeichenfläche konnte nicht erstellt werden")
	}
	s.old = call(gdi32, "SelectObject", s.dc, s.bitmap)
	return s, nil
}
func (s *surface) close() {
	call(gdi32, "SelectObject", s.dc, s.old)
	call(gdi32, "DeleteObject", s.bitmap)
	call(gdi32, "DeleteDC", s.dc)
}
func (s *surface) rgba() *image.RGBA {
	call(gdi32, "GdiFlush")
	out := image.NewRGBA(image.Rect(0, 0, s.width, s.height))
	b := unsafe.Slice((*byte)(s.bits), s.width*s.height*4)
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			i := (y*s.width + x) * 4
			out.SetRGBA(x, y, color.RGBA{b[i+2], b[i+1], b[i], 255})
		}
	}
	return out
}
func nativeIcon(source image.Image) (uintptr, error) {
	size := source.Bounds().Dx()
	s, err := newSurface(size, size)
	if err != nil {
		return 0, err
	}
	defer s.close()
	pixels := unsafe.Slice((*byte)(s.bits), size*size*4)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			r, g, b, a := source.At(x, y).RGBA()
			i := (y*size + x) * 4
			pixels[i] = byte(b / 257)
			pixels[i+1] = byte(g / 257)
			pixels[i+2] = byte(r / 257)
			pixels[i+3] = byte(a / 257)
		}
	}
	mask := call(gdi32, "CreateBitmap", uintptr(size), uintptr(size), 1, 1, 0)
	defer call(gdi32, "DeleteObject", mask)
	info := struct {
		Icon        int32
		X, Y        uint32
		Mask, Color uintptr
	}{Icon: 1, Mask: mask, Color: s.bitmap}
	icon := call(user32, "CreateIconIndirect", uintptr(unsafe.Pointer(&info)))
	if icon == 0 {
		return 0, errors.New("Tray-Symbol konnte nicht erstellt werden")
	}
	return icon, nil
}
func iconPixels(icon uintptr, size int) (*image.RGBA, error) {
	s, err := newSurface(size, size)
	if err != nil {
		return nil, err
	}
	defer s.close()
	call(user32, "DrawIconEx", s.dc, 0, 0, icon, uintptr(size), uintptr(size), 0, 0, 3)
	return s.rgba(), nil
}
