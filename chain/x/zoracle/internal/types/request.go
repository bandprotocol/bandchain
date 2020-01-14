package types

import (
	"fmt"
	"strings"
)

// Request is a type to store detail of request
type Request struct {
	Name        string `json:"name"`
	CodeHash    []byte `json:"codeHash"`
	Params      []byte `json:"params"`
	ReportEndAt uint64 `json:"reportEnd"`
}

// NewRequest - contructor of Request struct
func NewRequest(name string, codeHash, params []byte, reportEndAt uint64) Request {
	return Request{
		Name:        name,
		CodeHash:    codeHash,
		Params:      params,
		ReportEndAt: reportEndAt,
	}
}

func (req Request) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Name: %s
CodeHash: %x
Params: %x
ReportEndAt: %d`,
		req.Name,
		req.CodeHash,
		req.Params,
		req.ReportEndAt,
	))
}
