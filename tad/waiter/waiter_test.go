package waiter

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dc0d/cliche/tad/workercontext"
	"github.com/stretchr/testify/assert"
)

func Test01(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wctx := workercontext.New(ctx)

	var sum int64

	for i := 0; i < 10; i++ {
		i := i + 1
		wctx.Add(1)
		go func() {
			defer wctx.JobDone()
			<-wctx.Done()
			atomic.AddInt64(&sum, int64(i))
		}()
	}

	New(wctx).
		Timeout(time.Millisecond * 100).
		Cancel(cancel).
		Wait()
	assert.Equal(t, int64(55), sum)
}

func Test02(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wctx := workercontext.New(ctx)

	var sum int64

	for i := 0; i < 10; i++ {
		i := i + 1
		wctx.Add(1)
		go func() {
			// defer wctx.JobDone()
			<-wctx.Done()
			atomic.AddInt64(&sum, int64(i))
		}()
	}

	New(wctx).
		Timeout(time.Millisecond * 100).
		Cancel(cancel).
		Wait()
	assert.Equal(t, int64(55), sum)
}
