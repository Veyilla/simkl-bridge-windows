package main

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"unsafe"
)

func TestProgramReceivesItsOwnIconResource(t *testing.T) {
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(self)
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "icon-test.exe")
	if err = os.WriteFile(target, data, 0600); err != nil {
		t.Fatal(err)
	}
	icon, err := makeIcon()
	if err != nil {
		t.Fatal(err)
	}
	if err = embedIcon(target, icon); err != nil {
		t.Fatal(err)
	}
	dll := syscall.NewLazyDLL("kernel32.dll")
	name, _ := syscall.UTF16PtrFromString(target)
	module, _, _ := dll.NewProc("LoadLibraryExW").Call(uintptr(unsafe.Pointer(name)), 0, 2)
	if module == 0 {
		t.Fatal("Programmressourcen nicht lesbar")
	}
	defer dll.NewProc("FreeLibrary").Call(module)
	resource, _, _ := dll.NewProc("FindResourceW").Call(module, 1, 14)
	if resource == 0 {
		t.Fatal("Eigenes Programmsymbol fehlt")
	}
}
