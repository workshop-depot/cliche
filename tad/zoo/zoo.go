package zoo

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dc0d/cliche/tad"
)

var log = tad.GetLogger()

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

var (
	errNotAvailable = tad.Errorf("N/A")
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