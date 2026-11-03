package supervisor

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type testProcess struct {
	exit chan error
	once sync.Once
}

func (p *testProcess) Wait() error { return <-p.exit }
func (p *testProcess) Stop() error { p.once.Do(func() { p.exit <- nil }); return nil }
func (p *testProcess) crash()      { p.once.Do(func() { p.exit <- errors.New("Absturz") }) }

type testClock struct {
	mu    sync.Mutex
	now   time.Time
	waits chan time.Duration
	tick  chan time.Time
}

func newClock() *testClock {
	return &testClock{now: time.Unix(1000, 0), waits: make(chan time.Duration, 8), tick: make(chan time.Time, 8)}
}
func (c *testClock) Now() time.Time                         { c.mu.Lock(); defer c.mu.Unlock(); return c.now }
func (c *testClock) After(d time.Duration) <-chan time.Time { c.waits <- d; return c.tick }
func (c *testClock) advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	n := c.now
	c.mu.Unlock()
	c.tick <- n
}
func receive[T any](t *testing.T, c <-chan T) T {
	t.Helper()
	select {
	case v := <-c:
		return v
	case <-time.After(3 * time.Second):
		t.Fatal("Zeitüberschreitung")
		var zero T
		return zero
	}
}
func fixture(t *testing.T) (*Controller, *testClock, chan *testProcess) {
	t.Helper()
	clock := newClock()
	launched := make(chan *testProcess, 16)
	c := New(func(context.Context) (Process, error) {
		p := &testProcess{exit: make(chan error, 1)}
		launched <- p
		return p, nil
	}, clock)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := c.Stop(ctx); err != nil {
			t.Error(err)
		}
	})
	return c, clock, launched
}
func TestManualStopDoesNotRestart(t *testing.T) {
	c, clock, launched := fixture(t)
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	receive(t, launched)
	if err := c.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if s := c.Snapshot(); s.Phase != "stopped" {
		t.Fatalf("Zustand: %+v", s)
	}
	select {
	case <-clock.waits:
		t.Fatal("Manuelles Stoppen löst Neustart aus")
	default:
	}
	select {
	case <-launched:
		t.Fatal("Zusätzlicher Prozess")
	default:
	}
}
func TestQuitStopsOwnedWorker(t *testing.T) {
	c, _, launched := fixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	c.Start(ctx)
	p := receive(t, launched)
	cancel()
	deadline := time.After(3 * time.Second)
	for c.Snapshot().Phase != "stopped" {
		select {
		case <-deadline:
			t.Fatal("Prozess wurde nicht beendet")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	p.once.Do(func() { t.Error("Eigener Prozess blieb aktiv") })
}
func TestOnlyOneWorkerDuringConcurrentStart(t *testing.T) {
	c, _, launched := fixture(t)
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Go(func() {
			if err := c.Start(context.Background()); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if len(launched) != 1 {
		t.Fatalf("%d Prozesse statt einem", len(launched))
	}
}
func TestThreeCrashRetriesThenStop(t *testing.T) {
	c, clock, launched := fixture(t)
	c.Start(context.Background())
	for _, want := range []time.Duration{2 * time.Second, 5 * time.Second, 15 * time.Second} {
		receive(t, launched).crash()
		if got := receive(t, clock.waits); got != want {
			t.Fatalf("Wartezeit %s statt %s", got, want)
		}
		clock.advance(want)
	}
	receive(t, launched).crash()
	deadline := time.After(3 * time.Second)
	for c.Snapshot().Phase != "error" {
		select {
		case <-deadline:
			t.Fatal("Fehlerzustand fehlt")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	if s := c.Snapshot(); s.Retries != 3 || s.ErrorCode != "restart_limit" {
		t.Fatalf("Zustand: %+v", s)
	}
	select {
	case <-clock.waits:
		t.Fatal("Vierter Neustart wurde geplant")
	default:
	}
}
func TestStableMinuteResetsRetryCount(t *testing.T) {
	c, clock, launched := fixture(t)
	c.Start(context.Background())
	receive(t, launched).crash()
	receive(t, clock.waits)
	clock.advance(2 * time.Second)
	p := receive(t, launched)
	clock.mu.Lock()
	clock.now = clock.now.Add(60 * time.Second)
	clock.mu.Unlock()
	p.crash()
	if d := receive(t, clock.waits); d != 2*time.Second {
		t.Fatalf("Zähler nicht zurückgesetzt: %s", d)
	}
}
func TestActionRequiredDoesNotRestart(t *testing.T) {
	clock := newClock()
	c := New(func(context.Context) (Process, error) { return nil, ErrAuthentication }, clock)
	if err := c.Start(context.Background()); !errors.Is(err, ErrAuthentication) {
		t.Fatalf("Fehler: %v", err)
	}
	if s := c.Snapshot(); s.ErrorCode != "authentication_required" {
		t.Fatalf("Zustand: %+v", s)
	}
	select {
	case <-clock.waits:
		t.Fatal("Anmeldefehler löst Neustart aus")
	default:
	}
}
func TestRestartStopsBeforeStartingNext(t *testing.T) {
	c, _, launched := fixture(t)
	c.Start(context.Background())
	p := receive(t, launched)
	if err := c.Restart(context.Background()); err != nil {
		t.Fatal(err)
	}
	receive(t, launched)
	p.once.Do(func() { t.Error("Alter Prozess blieb aktiv") })
}
