package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Report represent detail of validator report
type Report struct {
	Data     []byte `json:"data"`
	ReportAt uint64 `json:"reportAt"`
}

// NewReport is a contructor of Report
func NewReport(data []byte, reportAt uint64) Report {
	return Report{
		Data:     data,
		ReportAt: reportAt,
	}
}

// ValidatorReport is a report that contain operator address in struct
type ValidatorReport struct {
	Value     RawJson        `json:"value"`
	ReportAt  uint64         `json:"reportAt"`
	Validator sdk.ValAddress `json:"validator"`
}

// NewValidatorReport is a contructor of ValidatorReport
func NewValidatorReport(value RawJson, reportAt uint64, valAddress sdk.ValAddress) ValidatorReport {
	return ValidatorReport{
		Value:     value,
		ReportAt:  reportAt,
		Validator: valAddress,
	}
}
