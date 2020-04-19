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

// Report is a report that contains operator address in struct
type Report struct {
	RawDataReports []RawDataReportWithID `json:"detail"`
	Validator      sdk.ValAddress        `json:"validator"`
}

// NewReport is a contructor of Report
func NewReport(validator sdk.ValAddress, reports []RawDataReportWithID) Report {
	return Report{
		RawDataReports: reports,
		Validator:      validator,
	}
}
