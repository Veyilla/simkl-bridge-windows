package supervisor

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if len(os.Args) == 2 && os.Args[1] == "--bridge-worker" {
		if os.Getenv("SIMKL_TEST_WORKER_EXIT") == "41" {
			os.Exit(41)
		}
		for {
			time.Sleep(time.Hour)
		}
	}
	os.Exit(m.Run())
}
func TestOwnedWorkerRunsAndStops(t *testing.T) {
	p, err := NewOwnedWorker(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- p.Wait() }()
	select {
	case err := <-done:
		t.Fatalf("Worker vorzeitig beendet: %v", err)
	case <-time.After(30 * time.Millisecond):
	}
	if err := p.Stop(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Eigener Worker wurde nicht beendet")
	}
	if err := p.Stop(); err != nil {
		t.Fatal("Wiederholtes Stoppen:", err)
	}
}
func TestWorkerAuthenticationExitIsActionable(t *testing.T) {
	t.Setenv("SIMKL_TEST_WORKER_EXIT", "41")
	p, err := NewOwnedWorker(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Wait(); !errors.Is(err, ErrAuthentication) {
		t.Fatalf("Exit 41 nicht erkannt: %v", err)
	}
}
