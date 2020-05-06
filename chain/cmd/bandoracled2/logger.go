package main

import (
	"fmt"
	"os"

	"github.com/tendermint/tendermint/libs/log"
)

type Logger struct {
	logger log.Logger
}

var logger Logger

func init() {
	logger = Logger{logger: log.NewTMLogger(os.Stdout)}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}
