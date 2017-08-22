package waiter

import (
	"context"
	"time"

	"github.com/dc0d/cliche/tad/workercontext"
)

const (
	defaultTimeout = time.Second * 3
)

type _err string

func (v _err) Error() string { return string(v) }

// variables
var (
	ErrTimeout error = _err("TIMEOUT")
)

// Waiter .
type Waiter struct {
	wctx    workercontext.WorkerContext
	timeout time.Duration
	cancel  context.CancelFunc
}

// New .
func New(wctx workercontext.WorkerContext) *Waiter {
	return &Waiter{
		wctx:    wctx,
		timeout: defaultTimeout,
	}
}

// Timeout sets the timeout if greater than zero, default is 3 seconds
func (w *Waiter) Timeout(timeout time.Duration) *Waiter {
	if timeout > 0 {
		w.timeout = timeout
	}
	return w
}

// Cancel .
func (w *Waiter) Cancel(cancel context.CancelFunc) *Waiter {
	w.cancel = cancel
	return w
}

// Wait returns ErrTimeout if the underlying WaitGroup did not finished
func (w *Waiter) Wait() error {
	if w.cancel != nil {
		w.cancel()
	}
	<-w.wctx.Done()

	done := make(chan struct{})
	go func() {
		defer close(done)
		w.wctx.Wait()
	}()
	select {
	case <-done:
	case <-time.After(w.timeout):
		return ErrTimeout
	}
	return nil
}
