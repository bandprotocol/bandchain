package executor

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	flagQueryTimeout = "timeout"
)

var (
	ErrExecutionimeout = errors.New("execution timeout")
	ErrRestNotOk       = errors.New("rest return non 2XX response")
)

type ExecResult struct {
	Output  []byte
	Code    uint32
	Version string
}

type Executor interface {
	Exec(exec []byte, arg string, env interface{}) (ExecResult, error)
}

var testProgram []byte = []byte("#!/usr/bin/env python3\nimport os\nimport sys\nprint(sys.argv[1], os.getenv('BAND_CHAIN_ID'))")

// NewExecutor returns executor by name and executor URL
func NewExecutor(executor string) (exec Executor, err error) {
	name, base, timeout, err := parseExecutor(executor)
	if err != nil {
		return nil, err
	}
	switch name {
	case "rest":
		exec = NewRestExec(base, timeout)
	case "docker":
		return nil, fmt.Errorf("Docker executor is currently not supported")
	default:
		return nil, fmt.Errorf("Invalid executor name: %s, base: %s", name, base)
	}

	// TODO: Remove hardcode in test execution
	res, err := exec.Exec(testProgram, "TEST_ARG", map[string]interface{}{
		"BAND_CHAIN_ID":    "test-chain-id",
		"BAND_VALIDATOR":   "test-validator",
		"BAND_REQUEST_ID":  "test-request-id",
		"BAND_EXTERNAL_ID": "test-external-id",
		"BAND_REPORTER":    "test-reporter",
		"BAND_SIGNATURE":   "test-signature",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to run test program: %s", err.Error())
	}
	if res.Code != 0 {
		return nil, fmt.Errorf("test program returned nonzero code: %d", res.Code)
	}
	if string(res.Output) != "TEST_ARG test-chain-id\n" {
		return nil, fmt.Errorf("test program returned wrong output: %s", res.Output)
	}
	return exec, nil
}

// parseExecutor splits the executor string in the form of "name:base?timeout=" into parts.
func parseExecutor(executorStr string) (name string, base string, timeout time.Duration, err error) {
	executor := strings.SplitN(executorStr, ":", 2)
	if len(executor) != 2 {
		return "", "", 0, fmt.Errorf("Invalid executor, cannot parse executor: %s", executorStr)
	}
	u, err := url.Parse(executor[1])
	if err != nil {
		return "", "", 0, fmt.Errorf("Invalid url, cannot parse %s to url with error: %s", executor[1], err.Error())
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
		return "", "", 0, fmt.Errorf("Invalid timeout, cannot parse duration with error: %s", err.Error())
	}
	return executor[0], u.String(), timeout, nil
}
