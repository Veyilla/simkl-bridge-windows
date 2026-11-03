package desktop

import (
	"context"
	"testing"
	"time"
)

// Ein Aufruf aus einem Hintergrundprozess muss auf dem Fensterthread ankommen.
func TestDispatchUsesWindowThread(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := make(chan bool, 1)
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) {
		owner := call(kernel32, "GetCurrentThreadId")
		go func() { a.Dispatch(func() { result <- call(kernel32, "GetCurrentThreadId") == owner; a.Close() }) }()
	})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case ok := <-result:
		if !ok {
			t.Fatal("Aktualisierung lief außerhalb des Fensterthreads")
		}
	default:
		t.Fatal("Aktualisierung wurde nicht ausgeführt")
	}
}

func TestContextClosesWindow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	err := RunView(ctx, View{Language: "de"}, Actions{}, func(a *App) { cancel() })
	if err != nil {
		t.Fatal(err)
	}
}
