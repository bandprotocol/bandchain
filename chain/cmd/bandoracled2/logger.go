package main

import (
	"os"

	"github.com/kyokomi/emoji"
	"github.com/tendermint/tendermint/libs/log"
)

type Logger struct {
	logger log.Logger
}

func NewLogger() *Logger {
	return &Logger{logger: log.NewTMLogger(os.Stdout)}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.logger.Debug(emoji.Sprintf(format, args...))
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.logger.Info(emoji.Sprintf(format, args...))
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.logger.Error(emoji.Sprintf(format, args...))
}

func (l *Logger) With(keyvals ...interface{}) *Logger {
	return &Logger{logger: l.logger.With(keyvals...)}
}
