package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Report represent detail of validator report
type Report struct {
	Data       []ExternalData `json:"data"`
	ReportedAt int64          `json:"reportedAt"`
}

// NewReport is a contructor of Report
func NewReport(data []ExternalData, reportedAt int64) Report {
	return Report{
		Data:       data,
		ReportedAt: reportedAt,
	}
}

// ExternalData encapsulates a raw data report for an external data source from a block validator.
type ExternalData struct {
	ExternalDataID int64  `json:"externalDataID"`
	Data           []byte `json:"data"`
}

// NewExternalData creates a new ExternalData instance.
func NewExternalData(externalDataID int64, data []byte) ExternalData {
	return ExternalData{
		ExternalDataID: externalDataID,
		Data:           data,
	}
}

// ValidatorReport is a report that contain operator address in struct
type ValidatorReport struct {
	Data       []ExternalData `json:"data"`
	ReportedAt int64          `json:"reportedAt"`
	Validator  sdk.ValAddress `json:"validator"`
}

// NewValidatorReport is a contructor of ValidatorReport
func NewValidatorReport(
	data []ExternalData,
	reportedAt int64,
	valAddress sdk.ValAddress,
) ValidatorReport {
	return ValidatorReport{
		Data:       data,
		ReportedAt: reportedAt,
		Validator:  valAddress,
	}
}
