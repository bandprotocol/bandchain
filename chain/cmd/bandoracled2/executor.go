package main

import (
	"time"

	"github.com/bandprotocol/bandchain/chain/pkg/byteexec"
)

type executor interface {
	Execute(l *Logger, exec []byte, timeout time.Duration, arg string) ([]byte, uint32)
}

type lambdaExecutor struct{}

func (e *lambdaExecutor) Execute(
	l *Logger, exec []byte, timeout time.Duration, arg string,
) ([]byte, uint32) {
	// TODO: Make URL configurable
	result, err := byteexec.RunOnAWSLambda(
		exec, timeout, arg,
		"https://dmptasv4j8.execute-api.ap-southeast-1.amazonaws.com/bash-execute",
	)
	if err != nil {
		l.Error(":skull: LambdaExecutor failed with error: %s", err.Error())
		return []byte("EXECUTION_ERROR"), 255
	}

	return result, 0
}
