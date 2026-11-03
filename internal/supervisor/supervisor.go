// Package supervisor verwaltet ausschließlich den eigenen Bridge-Prozess.
package supervisor

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrAuthentication = errors.New("Simkl-Anmeldung erforderlich")
var ErrConfiguration = errors.New("Bridge-Einstellungen prüfen")

type State struct {
	Phase     string
	Retries   int
	ErrorCode string
	WorkerPID int
}
type Process interface {
	Wait() error
	Stop() error
}
type Launcher func(context.Context) (Process, error)
type Clock interface {
	Now() time.Time
	After(time.Duration) <-chan time.Time
}
type realClock struct{}

func (realClock) Now() time.Time                         { return time.Now() }
func (realClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

type Controller struct {
	op          sync.Mutex
	mu          sync.Mutex
	state       State
	subscribers []chan State
	launch      Launcher
	clock       Clock
	cancel      context.CancelFunc
	done        chan struct{}
}

func New(launch Launcher, clock Clock) *Controller {
	if clock == nil {
		clock = realClock{}
	}
	return &Controller{launch: launch, clock: clock, state: State{Phase: "stopped"}}
}
func (c *Controller) Snapshot() State { c.mu.Lock(); defer c.mu.Unlock(); return c.state }
func (c *Controller) Subscribe() <-chan State {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch := make(chan State, 1)
	ch <- c.state
	c.subscribers = append(c.subscribers, ch)
	return ch
}
func (c *Controller) set(s State) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = s
	for _, ch := range c.subscribers {
		select {
		case <-ch:
		default:
		}
		ch <- s
	}
}
func (c *Controller) Start(ctx context.Context) error {
	c.op.Lock()
	defer c.op.Unlock()
	return c.start(ctx)
}
func (c *Controller) start(ctx context.Context) error {
	if c.done != nil {
		select {
		case <-c.done:
		default:
			return nil
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	c.set(State{Phase: "starting"})
	lifetime, cancel := context.WithCancel(ctx)
	p, err := c.launch(lifetime)
	if err != nil {
		cancel()
		c.set(State{Phase: "error", ErrorCode: errorCode(err)})
		return err
	}
	c.cancel = cancel
	c.done = make(chan struct{})
	c.set(State{Phase: "running", WorkerPID: processID(p)})
	go c.monitor(lifetime, p, c.done)
	return nil
}
func (c *Controller) Stop(ctx context.Context) error {
	c.op.Lock()
	defer c.op.Unlock()
	return c.stop(ctx)
}
func (c *Controller) stop(ctx context.Context) error {
	if c.done == nil {
		c.set(State{Phase: "stopped"})
		return nil
	}
	c.cancel()
	select {
	case <-c.done:
		c.set(State{Phase: "stopped"})
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (c *Controller) Restart(ctx context.Context) error {
	c.op.Lock()
	defer c.op.Unlock()
	if err := c.stop(ctx); err != nil {
		return err
	}
	return c.start(ctx)
}
func errorCode(err error) string {
	if errors.Is(err, ErrAuthentication) {
		return "authentication_required"
	}
	if errors.Is(err, ErrConfiguration) {
		return "configuration_required"
	}
	return "worker_failed"
}
func (c *Controller) monitor(ctx context.Context, p Process, done chan struct{}) {
	defer close(done)
	retries := 0
	delays := []time.Duration{2 * time.Second, 5 * time.Second, 15 * time.Second}
	for {
		started := c.clock.Now()
		exited := make(chan error, 1)
		go func(current Process) { exited <- current.Wait() }(p)
		var err error
		select {
		case <-ctx.Done():
			p.Stop()
			<-exited
			c.set(State{Phase: "stopped"})
			return
		case err = <-exited:
		}
		if ctx.Err() != nil {
			c.set(State{Phase: "stopped"})
			return
		}
		if errors.Is(err, ErrAuthentication) || errors.Is(err, ErrConfiguration) {
			c.set(State{Phase: "error", Retries: retries, ErrorCode: errorCode(err), WorkerPID: processID(p)})
			return
		}
		if c.clock.Now().Sub(started) >= time.Minute {
			retries = 0
		}
		for {
			if retries == len(delays) {
				c.set(State{Phase: "error", Retries: retries, ErrorCode: "restart_limit", WorkerPID: processID(p)})
				return
			}
			delay := delays[retries]
			retries++
			c.set(State{Phase: "starting", Retries: retries, ErrorCode: "worker_failed"})
			select {
			case <-ctx.Done():
				c.set(State{Phase: "stopped"})
				return
			case <-c.clock.After(delay):
			}
			p, err = c.launch(ctx)
			if err == nil {
				c.set(State{Phase: "running", Retries: retries, WorkerPID: processID(p)})
				break
			}
			if errors.Is(err, ErrAuthentication) || errors.Is(err, ErrConfiguration) {
				c.set(State{Phase: "error", Retries: retries, ErrorCode: errorCode(err)})
				return
			}
		}
	}
}
func processID(p Process) int {
	if owned, ok := p.(interface{ PID() int }); ok {
		return owned.PID()
	}
	return 0
}
