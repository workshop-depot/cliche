package tad

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test01(t *testing.T) {
	SetLogger(loggerFunc(func(...interface{}) {}))
	var sum int64
	Supervise(func() error {
		atomic.AddInt64(&sum, 1)
		return Errorf("DUMMY")
	}, 3, time.Millisecond*50)
	assert.Equal(t, int64(3), sum)
}

func Test02(t *testing.T) {
	SetLogger(loggerFunc(func(...interface{}) {}))
	var sum int64
	Supervise(func() error {
		atomic.AddInt64(&sum, 1)
		panic(Errorf("DUMMY"))
	}, 3, time.Millisecond*50)
	assert.Equal(t, int64(3), sum)
}

func Test03(t *testing.T) {
	SetLogger(loggerFunc(func(...interface{}) {}))
	var sum int64
	Supervise(func() error {
		atomic.AddInt64(&sum, 1)
		return nil
	}, 3, time.Millisecond*50)
	assert.Equal(t, int64(1), sum)
}
