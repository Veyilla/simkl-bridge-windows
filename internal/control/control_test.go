package control

import (
	"context"
	"testing"
)

type testWorker struct {
	starts, stops int
	running       bool
}

func (w *testWorker) Start(context.Context) error { w.starts++; w.running = true; return nil }
func (w *testWorker) Stop(context.Context) error  { w.stops++; w.running = false; return nil }
func (w *testWorker) Restart(context.Context) error {
	w.stops++
	w.starts++
	w.running = true
	return nil
}
func TestBackPreservesManualStop(t *testing.T) {
	w := &testWorker{}
	c := New(w)
	ctx := context.Background()
	c.Start(ctx)
	c.Stop(ctx)
	ticket, _ := c.BeginSettings(ctx, true)
	c.EndSettings(ctx, ticket, false, true)
	if w.running || w.starts != 1 {
		t.Fatal("Zurück hob den manuellen Stopp auf")
	}
}
func TestStaleCloseCannotStartWorkerInsideNewSettings(t *testing.T) {
	w := &testWorker{}
	c := New(w)
	ctx := context.Background()
	c.Start(ctx)
	old, _ := c.BeginSettings(ctx, true)
	current, _ := c.BeginSettings(ctx, true)
	c.EndSettings(ctx, old, false, true)
	if w.running {
		t.Fatal("Veralteter Abschluss startete während der neuen Einrichtung")
	}
	c.EndSettings(ctx, current, false, true)
	if !w.running || w.starts != 2 {
		t.Fatal("Gewünschter Betriebszustand wurde nicht wiederhergestellt")
	}
}
func TestFirstSavedSetupStartsBridge(t *testing.T) {
	w := &testWorker{}
	c := New(w)
	ctx := context.Background()
	ticket, _ := c.BeginSettings(ctx, false)
	c.EndSettings(ctx, ticket, true, true)
	if !w.running {
		t.Fatal("Erste fertige Einrichtung startet nicht")
	}
}
