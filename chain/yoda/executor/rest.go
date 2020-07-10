package executor

import (
	"encoding/base64"
	"time"

	"github.com/levigross/grequests"
)

type RestExec struct {
	url string
}

func NewRestExec(url string) *RestExec {
	return &RestExec{url: url}
}

type externalExecutionResponse struct {
	Returncode uint32 `json:"returncode"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
}

func (e *RestExec) Exec(timeout time.Duration, code []byte, arg string) (ExecResult, error) {
	executable := base64.StdEncoding.EncodeToString(code)
	resp, err := grequests.Post(
		e.url,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			JSON: map[string]interface{}{
				"executable": executable,
				"calldata":   arg,
				"timeout":    timeout.Milliseconds(),
			},
		},
	)

	if err != nil {
		return ExecResult{}, err
	}

	if resp.Ok != true {
		return ExecResult{}, ErrRestNotOk
	}

	r := externalExecutionResponse{}
	err = resp.JSON(&r)

	if err != nil {
		return ExecResult{}, err
	}

	if r.Returncode == 0 {
		return ExecResult{Output: []byte(r.Stdout), Code: 0}, nil
	} else {
		return ExecResult{Output: []byte(r.Stderr), Code: r.Returncode}, nil
	}
}
