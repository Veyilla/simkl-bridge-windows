package main

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"testing"
)

func TestIconContainsSharpSizesThrough256Pixels(t *testing.T) {
	data, err := makeIcon()
	if err != nil {
		t.Fatal(err)
	}
	if binary.LittleEndian.Uint16(data[4:6]) != 6 {
		t.Fatal("Symbolgrößen fehlen")
	}
	for i, want := range []int{16, 32, 48, 64, 128, 256} {
		entry := data[6+i*16 : 22+i*16]
		offset := binary.LittleEndian.Uint32(entry[12:16])
		length := binary.LittleEndian.Uint32(entry[8:12])
		img, err := png.Decode(bytes.NewReader(data[offset : offset+length]))
		if err != nil || img.Bounds().Dx() != want {
			t.Fatalf("Symbolgröße %d fehlt: %v", want, err)
		}
	}
}
