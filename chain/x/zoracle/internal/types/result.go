package types

import (
	"encoding/binary"
	"fmt"
)

// Result is a data structure that stores the detail of the result of a specific request.
type Result struct {
	RequestTime              int64  `json:"requestTime"`
	AggregationTime          int64  `json:"aggregationTime"`
	RequestedValidatorsCount int64  `json:"requestedValidatorsCount"`
	SufficientValidatorCount int64  `json:"sufficientValidatorCount"`
	ReportedValidatorsCount  int64  `json:"reportedValidatorsCount"`
	Data                     []byte `json:"data"`
}

// NewResult creates a new Result instance.
func NewResult(
	requestTime int64,
	aggregationTime int64,
	requestedValidatorsCount int64,
	sufficientValidatorCount int64,
	reportedValidatorsCount int64,
	data []byte,
) Result {
	return Result{
		RequestTime:              requestTime,
		AggregationTime:          aggregationTime,
		RequestedValidatorsCount: requestedValidatorsCount,
		SufficientValidatorCount: sufficientValidatorCount,
		ReportedValidatorsCount:  reportedValidatorsCount,
		Data:                     data,
	}
}

// DecodeResult is a helper function for decoding bytes to Result.
func DecodeResult(b []byte) (Result, error) {
	if len(b) < 40 {
		return Result{}, fmt.Errorf(fmt.Sprintf("Expect size of input to be at least 40 bytes but got %d bytes", len(b)))
	}
	return Result{
		RequestTime:              int64(binary.BigEndian.Uint64(b[0:8])),
		AggregationTime:          int64(binary.BigEndian.Uint64(b[8:16])),
		RequestedValidatorsCount: int64(binary.BigEndian.Uint64(b[16:24])),
		SufficientValidatorCount: int64(binary.BigEndian.Uint64(b[24:32])),
		ReportedValidatorsCount:  int64(binary.BigEndian.Uint64(b[32:40])),
		Data:                     b[40:],
	}, nil
}

// Bytes is a helper function for encoding Result to bytes.
func (result Result) EncodeResult() []byte {
	bs := make([]byte, 40)

	binary.BigEndian.PutUint64(bs[0:8], uint64(result.RequestTime))
	binary.BigEndian.PutUint64(bs[8:16], uint64(result.AggregationTime))
	binary.BigEndian.PutUint64(bs[16:24], uint64(result.RequestedValidatorsCount))
	binary.BigEndian.PutUint64(bs[24:32], uint64(result.SufficientValidatorCount))
	binary.BigEndian.PutUint64(bs[32:40], uint64(result.ReportedValidatorsCount))

	return append(bs, result.Data...)
}
