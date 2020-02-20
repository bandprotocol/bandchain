package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RawDataReport encapsulates a raw data report for an external data source from a block validator.
type RawDataReport struct {
	ExternalDataID int64  `json:"externalDataID"`
	Data           []byte `json:"data"`
}

// NewRawDataReport creates a new RawDataReport instance.
func NewRawDataReport(externalDataID int64, data []byte) RawDataReport {
	return RawDataReport{
		ExternalDataID: externalDataID,
		Data:           data,
	}
}

// ReportWithValidator is a report that contain operator address in struct
type ReportWithValidator struct {
	RawDataReports []RawDataReport `json:"detail"`
	Validator      sdk.ValAddress  `json:"validator"`
}

// NewReportWithValidator is a contructor of ReportWithValidator
func NewReportWithValidator(
	reports []RawDataReport,
	valAddress sdk.ValAddress,

) ReportWithValidator {
	return ReportWithValidator{
		RawDataReports: reports,
		Validator:      valAddress,
	}
}
