package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/levigross/grequests"
)

type executor interface {
	Execute(l *Logger, exec []byte, timeout time.Duration, arg string) ([]byte, uint32)
}

type lambdaExecutor struct {
	URL string
}

type externalExecutionResponse struct {
	Returncode uint32 `json:"returncode"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
}

func (e *lambdaExecutor) Execute(
	l *Logger, exec []byte, timeout time.Duration, arg string,
) ([]byte, uint32) {
	resp, err := grequests.Post(
		e.URL,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			JSON: map[string]string{
				"executable": string(exec),
				"calldata":   arg,
			},
		},
	)

	if err != nil {
		l.Error(":skull: LambdaExecutor failed with error: %s", err.Error())
		return []byte("EXECUTION_ERROR"), 255
	}

	if resp.Ok != true {
		l.Error(":skull: LambdaExecutor failed with error: %s", resp.Error)
		return []byte("EXECUTION_ERROR"), 255
	}

	r := externalExecutionResponse{}
	err = resp.JSON(&r)

	if err != nil {
		l.Error(":skull: LambdaExecutor failed with error: %s", err.Error())
		return []byte("EXECUTION_ERROR"), 255
	}

	if r.Returncode == 0 {
		return []byte(r.Stdout), r.Returncode
	} else {
		return []byte{}, r.Returncode
	}
}

// NewExecutor returns executor by name and executor URL
func NewExecutor(executor string) (executor, error) {
	name, url, err := parseExecutor(executor)
	if err != nil {
		return nil, err
	}
	switch name {
	case "lambda":
		return &lambdaExecutor{URL: url}, nil
	default:
		return nil, fmt.Errorf("Invalid executor name: %s, url: %s", name, url)
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
