package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RawDataReport struct {
	ExitCode uint8  `json:"exitCode"`
	Data     []byte `json:"data"`
}

func NewRawDataReport(exitCode uint8, data []byte) RawDataReport {
	return RawDataReport{
		ExitCode: exitCode,
		Data:     data,
	}
}

// RawDataReportWithID encapsulates a raw data report for an external data source from a block validator.
type RawDataReportWithID struct {
	ExternalID ExternalID `json:"externalID"`
	ExitCode   uint8      `json:"exitCode"`
	Data       []byte     `json:"data"`
}

// RawDataReportWithID creates a new RawDataReport instance.
func NewRawDataReportWithID(externalID ExternalID, exitCode uint8, data []byte) RawDataReportWithID {
	return RawDataReportWithID{
		ExternalID: externalID,
		ExitCode:   exitCode,
		Data:       data,
	}
}

var NewReport = NewRawDataReportWithID

// ReportWithValidator is a report that contains operator address in struct
type ReportWithValidator struct {
	RawDataReports []RawDataReportWithID `json:"detail"`
	Validator      sdk.ValAddress        `json:"validator"`
}

// NewReportWithValidator is a contructor of ReportWithValidator
func NewReportWithValidator(
	reports []RawDataReportWithID,
	valAddress sdk.ValAddress,

) ReportWithValidator {
	return ReportWithValidator{
		RawDataReports: reports,
		Validator:      valAddress,
	}
}

type Report = ReportWithValidator
type BatchReport = ReportWithValidator

var NewBatchReport = NewReportWithValidator
