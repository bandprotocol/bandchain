package executor

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrExecutionimeout = errors.New("execution timeout")
	ErrRestNotOk       = errors.New("rest return non 2XX response")
)

type ExecResult struct {
	Output []byte
	Code   uint32
}

type Executor interface {
	Exec(timeout time.Duration, exec []byte, arg string) (ExecResult, error)
}

// NewExecutor returns executor by name and executor URL
func NewExecutor(executor string) (Executor, error) {
	name, base, err := parseExecutor(executor)
	if err != nil {
		return nil, err
	}
	switch name {
	case "rest":
		return NewRestExec(base), nil
	case "docker":
		return NewDockerExec(base), nil
	default:
		return nil, fmt.Errorf("Invalid executor name: %s, base: %s", name, base)
	}
}

// parseExecutor splits the executor string in the form of "name:url" into parts.
func parseExecutor(executorStr string) (name string, url string, err error) {
	executor := strings.SplitN(executorStr, ":", 2)
	if len(executor) != 2 {
		return "", "", fmt.Errorf("Invalid executor, cannot parse executor: %s", executorStr)
	}
	return executor[0], executor[1], nil
}
