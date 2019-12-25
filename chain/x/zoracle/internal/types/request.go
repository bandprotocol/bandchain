package types

import (
	"fmt"
	"strings"
)

// Request is a type to store detail of request
type Request struct {
	CodeHash    []byte `json:"codeHash"`
	Params      []byte `json:"params"`
	ReportEndAt uint64 `json:"reportEnd"`
}

// NewRequest - contructor of Request struct
func NewRequest(codeHash, params []byte, reportEndAt uint64) Request {
	return Request{
		CodeHash:    codeHash,
		Params:      params,
		ReportEndAt: reportEndAt,
	}
}

func (req Request) String() string {
	return strings.TrimSpace(fmt.Sprintf(`CodeHash: %x
Params: %x
ReportEndAt: %d`,
		req.CodeHash,
		req.Params,
		req.ReportEndAt,
	))
}
