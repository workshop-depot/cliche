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
	wg sync.WaitGroup
}

func (wg *waitGroup) Add(delta int) { wg.wg.Add(delta) }
func (wg *waitGroup) JobDone()      { wg.wg.Done() }
func (wg *waitGroup) Wait()         { wg.wg.Wait() }

// causes stackoverflow in test - keep it & post it
// type waitGroup struct {
// 	sync.WaitGroup
// }
// func (wg *waitGroup) Add(delta int) { wg.Add(delta) }
// func (wg *waitGroup) JobDone()      { wg.Done() }
// func (wg *waitGroup) Wait()         { wg.Wait() }

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
	return &workerContext{
		Context:   ctx,
		WaitGroup: &waitGroup{},
	}
}
