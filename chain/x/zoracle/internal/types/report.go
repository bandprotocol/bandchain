package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Report represent detail of validator report
type Report struct {
	Data       []byte `json:"data"`
	ReportedAt uint64 `json:"reportedAt"`
}

// NewReport is a contructor of Report
func NewReport(data []byte, reportedAt uint64) Report {
	return Report{
		Data:       data,
		ReportedAt: reportedAt,
	}
}

// ValidatorReport is a report that contain operator address in struct
type ValidatorReport struct {
	Value      RawJson        `json:"value"`
	ReportedAt uint64         `json:"reportedAt"`
	Validator  sdk.ValAddress `json:"validator"`
}

// NewValidatorReport is a contructor of ValidatorReport
func NewValidatorReport(value RawJson, reportedAt uint64, valAddress sdk.ValAddress) ValidatorReport {
	return ValidatorReport{
		Value:      value,
		ReportedAt: reportedAt,
		Validator:  valAddress,
	}
}
