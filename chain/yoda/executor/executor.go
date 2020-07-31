package executor

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var (
	flagQueryTimeout   = "timeout"
	ErrExecutionimeout = errors.New("execution timeout")
	ErrRestNotOk       = errors.New("rest return non 2XX response")
)

type ExecResult struct {
	Output []byte
	Code   uint32
}

type Executor interface {
	Exec(exec []byte, arg string) (ExecResult, error)
}

// NewExecutor returns executor by name and executor URL
func NewExecutor(executor string) (Executor, error) {
	name, base, timeout, err := parseExecutor(executor)
	if err != nil {
		return nil, err
	}
	switch name {
	case "rest":
		return NewRestExec(base, timeout), nil
	case "docker":
		return NewDockerExec(base, timeout), nil
	default:
		return nil, fmt.Errorf("Invalid executor name: %s, base: %s", name, base)
	}
}

// parseExecutor splits the executor string in the form of "name:base?timeout=" into parts.
func parseExecutor(executorStr string) (name string, base string, timeout time.Duration, err error) {
	executor := strings.SplitN(executorStr, ":", 2)
	if len(executor) != 2 {
		return "", "", 0, fmt.Errorf("Invalid executor, cannot parse executor: %s", executorStr)
	}
	u, err := url.Parse(executor[1])
	if err != nil {
		return "", "", 0, fmt.Errorf("Invalid url, cannot parse url: %s", executor[1])
	}

	query := u.Query()
	timeoutStr := query.Get(flagQueryTimeout)
	if timeoutStr == "" {
		return "", "", 0, fmt.Errorf("Invalid timeout, executor requires query timeout")
	}
	// Remove timeout from query because we need to return `base`
	query.Del(flagQueryTimeout)
	u.RawQuery = query.Encode()

	timeout, err = time.ParseDuration(timeoutStr)
	if err != nil {
		return "", "", 0, fmt.Errorf("Invalid timeout, cannot parse timeout: %s", timeoutStr)
	}
	return executor[0], u.String(), timeout, nil
}
