// Package control ordnet Bedienung und Einrichtung um denselben Worker.
package control

import (
	"context"
	"sync"
)

type Worker interface {
	Start(context.Context) error
	Stop(context.Context) error
	Restart(context.Context) error
}
type Controller struct {
	operations                            sync.Mutex
	mu                                    sync.Mutex
	worker                                Worker
	wanted, editing, previouslyConfigured bool
	generation                            uint64
}

func (c *Controller) EnterSetupOperation(ctx context.Context) (func(), error) {
	c.operations.Lock()
	if err := ctx.Err(); err != nil {
		c.operations.Unlock()
		return nil, err
	}
	return c.operations.Unlock, nil
}

func New(worker Worker) *Controller { return &Controller{worker: worker} }
func (c *Controller) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.wanted = true
	if c.editing {
		return nil
	}
	return c.worker.Start(ctx)
}
func (c *Controller) Stop(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.wanted = false
	return c.worker.Stop(ctx)
}
func (c *Controller) Restart(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.wanted = true
	if c.editing {
		return nil
	}
	return c.worker.Restart(ctx)
}
func (c *Controller) BeginSettings(ctx context.Context, configured bool) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.generation++
	c.editing = true
	c.previouslyConfigured = configured
	return c.generation, c.worker.Stop(ctx)
}
func (c *Controller) EndSettings(ctx context.Context, ticket uint64, saved, configured bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ticket != c.generation || !c.editing {
		return nil
	}
	c.editing = false
	if saved && configured && !c.previouslyConfigured {
		c.wanted = true
	}
	if c.wanted && configured {
		return c.worker.Start(ctx)
	}
	return nil
}
