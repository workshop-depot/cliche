package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dc0d/goroutines"
	"github.com/dc0d/public-club/appclub"
)

var (
	appCtx    context.Context
	appCancel context.CancelFunc
	appWG     = &sync.WaitGroup{}
)

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
	appclub.CallOnSignal(func() { appCancel() })
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
