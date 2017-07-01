package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/comail/colog"
	"github.com/dc0d/goroutines"
)

// simpleSupervisor is a simple supervisor in the context of this app (uses wg),
// if intensity < 0 runs forever, default period is 1 second.
func simpleSupervisor(
	intensity int,
	fn func(),
	period ...time.Duration) {
	if intensity == 0 {
		return
	}
	if intensity > 0 {
		intensity--
	}
	dt := time.Second
	if len(period) > 0 {
		dt = period[0]
	}
	goroutines.New().
		AddToGroup(wg).
		Recover(func(e interface{}) {
			time.Sleep(dt)
			go simpleSupervisor(intensity, fn, dt)
		}).
		Go(fn)
}

const (
	noCaller = "N/A"
)

func caller() string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return noCaller
	}
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return noCaller
	}
	name := fun.Name()
	ix := strings.LastIndex(name, ".")
	if ix > 0 && (ix+1) < len(name) {
		name = name[ix+1:]
	}
	return name
}

func timerScope(name string, opCount ...int) func() {
	if name == "" {
		name = caller()
	}
	log.Println(name, `started`)
	start := time.Now()
	return func() {
		elapsed := time.Now().Sub(start)
		log.Printf("%s took %v", name, elapsed)
		if len(opCount) == 0 {
			return
		}

		N := opCount[0]
		if N <= 0 {
			return
		}

		E := float64(elapsed)
		FRC := E / float64(N)

		log.Printf("op/sec %.2f", float64(N)/(E/float64(time.Second)))

		switch {
		case FRC > float64(time.Second):
			log.Printf("sec/op %.2f", (E/float64(time.Second))/float64(N))
		case FRC > float64(time.Millisecond):
			log.Printf("milli-sec/op %.2f", (E/float64(time.Millisecond))/float64(N))
		case FRC > float64(time.Microsecond):
			log.Printf("micro-sec/op %.2f", (E/float64(time.Microsecond))/float64(N))
		default:
			log.Printf("nano-sec/op %.2f", (E/float64(time.Nanosecond))/float64(N))
		}
	}
}

func callOnSignal(f func(), sig ...os.Signal) {
	if f == nil {
		return
	}
	sigc := make(chan os.Signal, 1)
	if len(sig) > 0 {
		signal.Notify(sigc, sig...)
	} else {
		signal.Notify(sigc,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGSTOP,
			syscall.SIGABRT,
			syscall.SIGTSTP,
			syscall.SIGKILL)
	}
	go func() {
		<-sigc
		f()
	}()
}

var (
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
)

func init() {
	wg = &sync.WaitGroup{}
	colog.Register()
	colog.SetFlags(log.Lshortfile)

	ctx, cancel = context.WithCancel(context.Background())
	callOnSignal(func() { cancel() })
}

func finit(timeout time.Duration, cancelApp ...bool) {
	if len(cancelApp) > 0 && cancelApp[0] {
		cancel()
	}
	<-ctx.Done()
	werr := goroutines.New().
		EnsureStarted().
		Timeout(timeout).
		Go(func() {
			wg.Wait()
		})
	if werr != nil {
		log.Println("error:", werr)
	}
}
