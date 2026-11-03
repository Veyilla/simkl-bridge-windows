package control

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCancelledNewSetupCannotOverlapOldTokenOperation(t *testing.T) {
	c := New(&testWorker{})
	release, err := c.EnterSetupOperation(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		unlock, err := c.EnterSetupOperation(ctx)
		if err == nil {
			unlock()
		}
		result <- err
	}()
	select {
	case <-result:
		t.Fatal("Anmeldeoperationen überlappten sich")
	case <-time.After(20 * time.Millisecond):
	}
	cancel()
	release()
	if err := <-result; !errors.Is(err, context.Canceled) {
		t.Fatal("Abgebrochene Einrichtung wurde trotzdem fortgesetzt", err)
	}
}
