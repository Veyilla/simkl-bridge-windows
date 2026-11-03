package desktop

import (
	"image"
	"unsafe"
)

func TrayLogoImage(size int) *image.RGBA {
	logo := LogoImage(size)
	s, err := newSurface(size, size)
	if err != nil {
		return logo
	}
	defer s.close()
	b := unsafe.Slice((*byte)(s.bits), size*size*4)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			c := logo.RGBAAt(x, y)
			i := (y*size + x) * 4
			b[i] = c.B
			b[i+1] = c.G
			b[i+2] = c.R
			b[i+3] = 255
		}
	}
	v := vector(s.dc, float64(size)/32)
	v.roundRect(1, 1, 30, 30, 6, 0x777777)
	v.close()
	return s.rgba()
}
