// Package tad - comment/delete whatever you don't need or copy to utils.go
package tad

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/hcl"
)

//-----------------------------------------------------------------------------

// OnSignal runs function on receiving the OS signal
func OnSignal(f func(), sig ...os.Signal) {
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

//-----------------------------------------------------------------------------

// Supervise runs in sync, use as "go Supervise(...)",
// takes care of restarts (in case of panic or error)
func Supervise(action func() error, intensity int, period ...time.Duration) {
	dt := time.Second * 3
	if len(period) > 0 && period[0] > 0 {
		dt = period[0]
	}
	for intensity != 0 {
		if intensity > 0 {
			intensity--
		}
		if err := runOnce(action); err != nil {
			log.Error(err)
			time.Sleep(dt)
		} else {
			break
		}
	}
}

func runOnce(action func() error) (errrun error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				errrun = err
				return
			}
		}
	}()
	return action()
}

//-----------------------------------------------------------------------------

// Application-wide execution context (whatever)
var (
	AppCtx    context.Context
	AppCancel context.CancelFunc
	AppWG     = &sync.WaitGroup{}
)

func init() {
	AppCtx, AppCancel = context.WithCancel(context.Background())
	OnSignal(func() { AppCancel() })
}

//-----------------------------------------------------------------------------

// Finit .
func Finit(timeout time.Duration, cancelApp ...bool) {
	if len(cancelApp) > 0 && cancelApp[0] {
		AppCancel()
	}
	<-AppCtx.Done()

	done := make(chan struct{})
	go func() {
		defer close(done)
		AppWG.Wait()
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		log.Error(fmt.Errorf("TIMEOUT"))
	}
}

//-----------------------------------------------------------------------------

func defaultAppNameHandler() string {
	return filepath.Base(os.Args[0])
}

func defaultConfNameHandler() string {
	fp := fmt.Sprintf("%s.conf", defaultAppNameHandler())
	if _, err := os.Stat(fp); err != nil {
		fp = "app.conf"
	}
	return fp
}

// LoadHCL loads hcl conf file. default conf file names (if filePath not provided)
// in the same directory are <appname>.conf and if not fount app.conf
func LoadHCL(ptr interface{}, filePath ...string) error {
	var fp string
	if len(filePath) > 0 {
		fp = filePath[0]
	}
	if fp == "" {
		fp = defaultConfNameHandler()
	}
	cn, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	err = hcl.Unmarshal(cn, ptr)
	if err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------------------------------------

type _err string

func (v _err) Error() string { return string(v) }

// Errorf value type (string) error
func Errorf(format string, a ...interface{}) error {
	return _err(fmt.Sprintf(format, a...))
}

//-----------------------------------------------------------------------------

var (
	errNotAvailable = Errorf("N/A")
)

// Here .
func Here(skip ...int) (funcName, fileName string, fileLine int, callerErr error) {
	sk := 1
	if len(skip) > 0 && skip[0] > 1 {
		sk = skip[0]
	}
	var pc uintptr
	var ok bool
	pc, fileName, fileLine, ok = runtime.Caller(sk)
	if !ok {
		callerErr = errNotAvailable
		return
	}
	fn := runtime.FuncForPC(pc)
	name := fn.Name()
	ix := strings.LastIndex(name, ".")
	if ix > 0 && (ix+1) < len(name) {
		name = name[ix+1:]
	}
	funcName = name
	nd, nf := filepath.Split(fileName)
	fileName = filepath.Join(filepath.Base(nd), nf)
	return
}

//-----------------------------------------------------------------------------

// TimerScope .
func TimerScope(name string, opCount ...int) func() {
	if name == "" {
		funcName, fileName, fileLine, err := Here(2)
		if err != nil {
			name = "N/A"
		} else {
			name = fmt.Sprintf("%s:%02d %s()", fileName, fileLine, funcName)
		}
	}
	log.Info(name, `started`)
	start := time.Now()
	return func() {
		elapsed := time.Now().Sub(start)
		log.Infof("%s took %v", name, elapsed)
		if len(opCount) == 0 {
			return
		}

		N := opCount[0]
		if N <= 0 {
			return
		}

		E := float64(elapsed)
		FRC := E / float64(N)

		log.Infof("op/sec %.2f", float64(N)/(E/float64(time.Second)))

		switch {
		case FRC > float64(time.Second):
			log.Infof("sec/op %.2f", (E/float64(time.Second))/float64(N))
		case FRC > float64(time.Millisecond):
			log.Infof("milli-sec/op %.2f", (E/float64(time.Millisecond))/float64(N))
		case FRC > float64(time.Microsecond):
			log.Infof("micro-sec/op %.2f", (E/float64(time.Microsecond))/float64(N))
		default:
			log.Infof("nano-sec/op %.2f", (E/float64(time.Nanosecond))/float64(N))
		}
	}
}

//-----------------------------------------------------------------------------

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// GetBuffer .
func GetBuffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	return buff
}

// PutBuffer .
func PutBuffer(buff *bytes.Buffer) {
	bufferPool.Put(buff)
	buff.Reset()
}

//-----------------------------------------------------------------------------
