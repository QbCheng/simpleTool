package simpleLogger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	defaultLogger *simpleLogger
	once          sync.Once
)

const (
	defaultFlag = log.Lshortfile | log.LstdFlags
)

func DefaultLogger(options ...option) *simpleLogger {
	if defaultLogger == nil {
		once.Do(func() {
			defaultLogger = &simpleLogger{}
			for i := range options {
				options[i](defaultLogger)
			}
			if defaultLogger.flag == 0 {
				defaultLogger.logger = log.New(os.Stdout, "", defaultFlag)
			} else {
				defaultLogger.logger = log.New(os.Stdout, "", defaultLogger.flag)
			}
		})
	}
	return defaultLogger
}

type option func(log *simpleLogger)

func WithCallPath(callPath int) option {
	return func(log *simpleLogger) {
		log.callPath = callPath
	}
}

func WithFlag(flag int) option {
	return func(log *simpleLogger) {
		log.flag = flag
	}
}

type simpleLogger struct {
	logger   *log.Logger
	callPath int
	flag     int
}

func NewLogger(options ...option) *simpleLogger {
	ret := &simpleLogger{}
	for i := range options {
		options[i](ret)
	}
	if ret.flag == 0 {
		ret.logger = log.New(os.Stdout, "", defaultFlag)
	} else {
		ret.logger = log.New(os.Stdout, "", ret.flag)
	}
	return ret
}

func (sl simpleLogger) Logf(layout string, s ...interface{}) {
	_ = sl.logger.Output(sl.callPath, fmt.Sprintf(layout, s...))
}
