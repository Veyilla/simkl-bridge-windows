package main

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"os"
	"simklbridge/internal/desktop"
)

func makeIcon() ([]byte, error) {
	sizes := []int{16, 32, 48, 64, 128, 256}
	images := make([][]byte, len(sizes))
	for i, size := range sizes {
		var buffer bytes.Buffer
		if err := png.Encode(&buffer, desktop.TrayLogoImage(size)); err != nil {
			return nil, err
		}
		images[i] = buffer.Bytes()
	}
	var output bytes.Buffer
	binary.Write(&output, binary.LittleEndian, []uint16{0, 1, uint16(len(sizes))})
	offset := 6 + 16*len(sizes)
	for i, size := range sizes {
		dimension := byte(size)
		if size == 256 {
			dimension = 0
		}
		output.Write([]byte{dimension, dimension, 0, 0})
		binary.Write(&output, binary.LittleEndian, []uint16{1, 32})
		binary.Write(&output, binary.LittleEndian, []uint32{uint32(len(images[i])), uint32(offset)})
		offset += len(images[i])
	}
	for _, data := range images {
		output.Write(data)
	}
	return output.Bytes(), nil
}
func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		os.Exit(2)
	}
	data, err := makeIcon()
	if err != nil {
		os.Exit(1)
	}
	if os.WriteFile(os.Args[1], data, 0644) != nil {
		os.Exit(1)
	}
	if len(os.Args) == 3 {
		if embedIcon(os.Args[2], data) != nil {
			os.Exit(1)
		}
	}
}
