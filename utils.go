package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/dc0d/goroutines"
	"github.com/hashicorp/hcl"
)

var (
	appCtx    context.Context
	appCancel context.CancelFunc
	appWG     = &sync.WaitGroup{}
)

//-----------------------------------------------------------------------------

// build flags
var (
	BuildTime  string
	CommitHash string
	GoVersion  string
	GitTag     string
)

//-----------------------------------------------------------------------------

func init() {
	appCtx, appCancel = context.WithCancel(context.Background())
	onSignal(func() { appCancel() })
}

func onSignal(f func(), sig ...os.Signal) {
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

// loadHCL loads hcl conf file. default conf file names (if filePath not provided)
// in the same directory are <appname>.conf and if not fount app.conf
func loadHCL(ptr interface{}, filePath ...string) error {
	var fp string
	if len(filePath) > 0 {
		fp = filePath[0]
	}
	if fp == "" {
		fp := fmt.Sprintf("%s.conf", filepath.Base(os.Args[0]))
		if _, err := os.Stat(fp); err != nil {
			fp = "app.conf"
		}
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

//-----------------------------------------------------------------------------
