package desktop

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"simklbridge/assets"
	"sync"
)

var iconOnce sync.Once
var iconImage image.Image

// IconPixel composites the supplied transparent S on white, with black corners.
func IconPixel(x, y int) uint32 {
	iconOnce.Do(func() { iconImage, _, _ = image.Decode(bytes.NewReader(assets.SimklPNG)) })
	if x < 0 || x >= 32 || y < 0 || y >= 32 || iconImage == nil {
		return 0
	}
	r, g, b, a := iconImage.At(x*iconImage.Bounds().Dx()/32, y*iconImage.Bounds().Dy()/32).RGBA()
	backing := uint32(0)
	if x >= 4 && x < 28 && y >= 5 && y < 27 {
		backing = 255
	}
	channel := func(v uint32) uint32 { return compositeChannel(v, a, backing) }
	return channel(r)<<16 | channel(g)<<8 | channel(b)
}

func compositeChannel(value, alpha, backing uint32) uint32 {
	return min(255, (value+backing*257*(65535-alpha)/65535)/257)
}

// LogoImage verkleinert die Originalauflösung mit Flächenmittelung.
func LogoImage(size int) *image.RGBA {
	iconOnce.Do(func() { iconImage, _, _ = image.Decode(bytes.NewReader(assets.SimklPNG)) })
	out := image.NewRGBA(image.Rect(0, 0, size, size))
	if iconImage == nil {
		return out
	}
	bounds := iconImage.Bounds()
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			var sum [3]uint64
			count := uint64(0)
			for sy := y * bounds.Dy() / size; sy < max(y*bounds.Dy()/size+1, (y+1)*bounds.Dy()/size); sy++ {
				for sx := x * bounds.Dx() / size; sx < max(x*bounds.Dx()/size+1, (x+1)*bounds.Dx()/size); sx++ {
					r, g, b, a := iconImage.At(sx+bounds.Min.X, sy+bounds.Min.Y).RGBA()
					back := uint32(0)
					if sx*32 >= 4*bounds.Dx() && sx*32 < 28*bounds.Dx() && sy*32 >= 5*bounds.Dy() && sy*32 < 27*bounds.Dy() {
						back = 255
					}
					sum[0] += uint64(compositeChannel(r, a, back))
					sum[1] += uint64(compositeChannel(g, a, back))
					sum[2] += uint64(compositeChannel(b, a, back))
					count++
				}
			}
			out.SetRGBA(x, y, color.RGBA{uint8(sum[0] / count), uint8(sum[1] / count), uint8(sum[2] / count), 255})
		}
	}
	return out
}
