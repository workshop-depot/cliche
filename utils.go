package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/comail/colog"
	"github.com/dc0d/goroutines"
	"github.com/hashicorp/hcl"
)

// value error
type valErr string

func (v valErr) Error() string { return string(v) }

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func getBuffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	buff.Reset()
	return buff
}

func putBuffer(buff *bytes.Buffer) {
	bufferPool.Put(buff)
}

// collection of errors
type colErr []error

func (x colErr) String() string {
	return x.Error()
}

func (x colErr) Error() string {
	if x == nil {
		return ``
	}

	buff := getBuffer()
	defer putBuffer(buff)

	for _, ve := range x {
		if ve == nil {
			continue
		}
		buff.WriteString(` [` + ve.Error() + `]`)
	}
	res := strings.TrimSpace(buff.String())

	return res
}

var (
	noCaller error = valErr("N/A")
)

func caller() (funcName, fileName string, fileLine int, callerErr error) {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "", "", -1, noCaller
	}
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "", "", -1, noCaller
	}
	name := fun.Name()
	fileName, fileLine = fun.FileLine(fun.Entry())
	fileName = filepath.Base(fileName)
	ix := strings.LastIndex(name, ".")
	if ix > 0 && (ix+1) < len(name) {
		name = name[ix+1:]
	}
	funcName = name
	return
}

func timerScope(name string, opCount ...int) func() {
	if name == "" {
		funcName, fileName, fileLine, err := caller()
		if err != nil {
			name = "N/A"
		} else {
			name = fmt.Sprintf("%s(?) @ %s-L%v", funcName, fileName, fileLine)
		}
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
	appCtx    context.Context
	appCancel context.CancelFunc
	appWG     = &sync.WaitGroup{}
)

func init() {
	colog.Register()
	colog.SetFlags(log.Lshortfile)

	appCtx, appCancel = context.WithCancel(context.Background())
	callOnSignal(func() { appCancel() })
}

func finit(timeout time.Duration, cancelApp ...bool) {
	if len(cancelApp) > 0 && cancelApp[0] {
		appCancel()
	}
	<-appCtx.Done()
	werr := goroutines.New().
		EnsureStarted().
		Timeout(timeout).
		Go(func() {
			appWG.Wait()
		})
	if werr != nil {
		log.Println("error:", werr)
	}
}

// build flags
var (
	BuildTime  string
	CommitHash string
	GoVersion  string
	GitTag     string
)

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

func loadHCL(ptr interface{}, filePath ...string) error {
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
