package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"syscall"
	"unsafe"
)

func embedIcon(path string, icon []byte) error {
	if len(icon) < 6 {
		return errors.New("Ungültige Symboldatei")
	}
	count := int(binary.LittleEndian.Uint16(icon[4:6]))
	if count < 1 || len(icon) < 6+count*16 {
		return errors.New("Ungültige Symboldatei")
	}
	library := syscall.NewLazyDLL("kernel32.dll")
	file, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	handle, _, _ := library.NewProc("BeginUpdateResourceW").Call(uintptr(unsafe.Pointer(file)), 0)
	if handle == 0 {
		return errors.New("Programmsymbol konnte nicht vorbereitet werden")
	}
	committed := false
	defer func() {
		if !committed {
			library.NewProc("EndUpdateResourceW").Call(handle, 1)
		}
	}()
	var group bytes.Buffer
	binary.Write(&group, binary.LittleEndian, []uint16{0, 1, uint16(count)})
	for i := 0; i < count; i++ {
		entry := icon[6+i*16 : 22+i*16]
		size := int(binary.LittleEndian.Uint32(entry[8:12]))
		offset := int(binary.LittleEndian.Uint32(entry[12:16]))
		if offset < 0 || size < 1 || offset+size > len(icon) {
			return errors.New("Ungültige Symboldaten")
		}
		data := icon[offset : offset+size]
		ok, _, _ := library.NewProc("UpdateResourceW").Call(handle, 3, uintptr(i+1), 0, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
		if ok == 0 {
			return errors.New("Programmsymbol konnte nicht eingebettet werden")
		}
		group.Write(entry[:12])
		binary.Write(&group, binary.LittleEndian, uint16(i+1))
	}
	data := group.Bytes()
	ok, _, _ := library.NewProc("UpdateResourceW").Call(handle, 14, 1, 0, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	if ok == 0 {
		return errors.New("Symbolgruppe konnte nicht eingebettet werden")
	}
	ok, _, _ = library.NewProc("EndUpdateResourceW").Call(handle, 0)
	if ok == 0 {
		return errors.New("Programmsymbol konnte nicht gespeichert werden")
	}
	committed = true
	return nil
}
