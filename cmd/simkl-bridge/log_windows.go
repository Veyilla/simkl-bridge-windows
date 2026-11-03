package main

import (
	"io"
	"os"
)

func readLogTail(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return "Noch kein Protokoll verfügbar."
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return "Protokoll konnte nicht gelesen werden."
	}
	if info.Size() > 32768 {
		file.Seek(info.Size()-32768, io.SeekStart)
	}
	data, err := io.ReadAll(io.LimitReader(file, 32768))
	if err != nil {
		return "Protokoll konnte nicht gelesen werden."
	}
	return string(data)
}
