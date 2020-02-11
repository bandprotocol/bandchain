package types

import (
	"encoding/json"

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
	Value      json.RawMessage `json:"value"`
	ReportedAt uint64          `json:"reportedAt"`
	Validator  sdk.ValAddress  `json:"validator"`
}

// NewValidatorReport is a contructor of ValidatorReport
func NewValidatorReport(
	value json.RawMessage,
	reportedAt uint64,
	valAddress sdk.ValAddress,
) ValidatorReport {
	return ValidatorReport{
		Value:      value,
		ReportedAt: reportedAt,
		Validator:  valAddress,
	}
}

// ExternalData is a struct that capsulized external data id and raw data
type ExternalData struct {
	ExternalDataID int64  `json:"externalDataID"`
	Data           []byte `json:"data"`
}

// NewExternalData is a constructor of ExternalData
func NewExternalData(externalDataID int64, data []byte) ExternalData {
	return ExternalData{
		ExternalDataID: externalDataID,
		Data:           data,
	}
}
