package workercontext

import (
	"context"
	"sync"
)

// WaitGroup JobDone is the standard Done() in WaitGroup, renamed because of duplicate
type WaitGroup interface {
	Add(delta int)
	JobDone()
	Wait()
}

type waitGroup struct {
	sync.WaitGroup
}

func (wg *waitGroup) Add(delta int) { wg.WaitGroup.Add(delta) }
func (wg *waitGroup) JobDone()      { wg.WaitGroup.Done() }
func (wg *waitGroup) Wait()         { wg.WaitGroup.Wait() }

// WorkerContext combination of context.Context & WaitGroup
type WorkerContext interface {
	context.Context
	WaitGroup
}

type workerContext struct {
	context.Context
	WaitGroup
}

// New .
func New(ctx context.Context) WorkerContext {
	// TODO: should panic on ctx == nil
	return &workerContext{
		Context:   ctx,
		WaitGroup: &waitGroup{},
	}
}
