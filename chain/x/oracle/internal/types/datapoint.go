package types

import (
	"fmt"
	"strings"
)

// DataPoint is a type to store result of request
type DataPoint struct {
	RequestID   uint64 `json:"requestID"`
	CodeHash    []byte `json:"codeHash"`
	ReportEndAt uint64 `json:"reportEnd"`
	Result      []byte `json:"result"`
}

// NewDataPoint - contructor of DataPoint struct
func NewDataPoint(requestID uint64, codeHash []byte, reportEndAt uint64) DataPoint {
	return DataPoint{
		RequestID:   requestID,
		CodeHash:    codeHash,
		ReportEndAt: reportEndAt,
	}
}

func (req DataPoint) String() string {
	return strings.TrimSpace(fmt.Sprintf(`requestID: %d
CodeHash: %x
ReportEnd: %d
Result: %x`,
		req.RequestID,
		req.CodeHash,
		req.ReportEndAt,
		req.Result))
}
