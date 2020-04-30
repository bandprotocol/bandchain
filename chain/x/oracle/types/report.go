package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RawReport encapsulates a raw data report for an external data source from a block validator.
type RawReport struct {
	ExternalID ExternalID `json:"externalID"`
	ExitCode   uint8      `json:"exitCode"`
	Data       []byte     `json:"data"`
}

// RawReport creates a new RawDataReport instance.
func NewRawReport(externalID ExternalID, exitCode uint8, data []byte) RawReport {
	return RawReport{
		ExternalID: externalID,
		ExitCode:   exitCode,
		Data:       data,
	}
}

// Report is a report that contains operator address in struct
type Report struct {
	RawReports []RawReport    `json:"detail"`
	Validator  sdk.ValAddress `json:"validator"`
}

// NewReport is a contructor of Report
func NewReport(validator sdk.ValAddress, reports []RawReport) Report {
	return Report{
		RawReports: reports,
		Validator:  validator,
	}
}
