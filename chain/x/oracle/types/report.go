package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewRawReport creates a new RawDataReport instance.
func NewRawReport(externalID ExternalID, exitCode uint32, data []byte) RawReport {
	return RawReport{
		ExternalID: externalID,
		ExitCode:   exitCode,
		Data:       data,
	}
}

type RawReports []RawReport

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
