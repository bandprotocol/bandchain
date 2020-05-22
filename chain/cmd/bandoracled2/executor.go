package main

import (
	"time"

	"github.com/bandprotocol/bandchain/chain/pkg/byteexec"
)

type executor interface {
	Execute(l *Logger, exec []byte, timeout time.Duration, arg string) ([]byte, uint32)
}

type lambdaExecutor struct {
	URL string
}

func (e *lambdaExecutor) Execute(
	l *Logger, exec []byte, timeout time.Duration, arg string,
) ([]byte, uint32) {
	// TODO: Make URL configurable
	result, err := byteexec.RunOnAWSLambda(
		exec, timeout, arg,
		e.URL,
	)
	if err != nil {
		l.Error(":skull: LambdaExecutor failed with error: %s", err.Error())
		return []byte("EXECUTION_ERROR"), 255
	}

	return result, 0
}

// NewExecutor returns executor by name and executer URL
func NewExecutor(name string, url string) executor {
	switch name {
	case "lambda":
		return &lambdaExecutor{URL: url}
	default:
		return &lambdaExecutor{URL: url}
	}
}
