package executor

import (
	"errors"
	"time"
)

var (
	ErrExecutionimeout = errors.New("execution timeout")
)

type ExecResult struct {
	Output []byte
	Code   uint32
}

type Executor interface {
	Exec(timeout time.Duration, exec []byte, args ...string) (ExecResult, error)
}
