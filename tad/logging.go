package tad

import (
	"fmt"
	internalog "log"
)

// Logger generic logger
type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
}

var log Logger

type loggerFunc func(args ...interface{})

func (lf loggerFunc) Error(args ...interface{})                 { lf(args...) }
func (lf loggerFunc) Info(args ...interface{})                  { lf(args...) }
func (lf loggerFunc) Errorf(format string, args ...interface{}) { lf(fmt.Sprintf(format, args...)) }
func (lf loggerFunc) Infof(format string, args ...interface{})  { lf(fmt.Sprintf(format, args...)) }

func init() {
	if log == nil {
		log = loggerFunc(internalog.Print)
	}
}

// SetLogger replaces logger for this package
func SetLogger(logger Logger) { log = logger }

// GetLogger .
func GetLogger() Logger { return log }
