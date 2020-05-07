package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewValidatorReportInfo creates a new ValidatorReportInfo instance
func NewValidatorReportInfo(
	validator sdk.ValAddress, isFullTime bool, indexOffset, missedBlocksCounter uint64,
) ValidatorReportInfo {
	return ValidatorReportInfo{
		Validator:            validator,
		IsFullTime:           isFullTime,
		IndexOffset:          indexOffset,
		MissedReportsCounter: missedBlocksCounter,
	}
}
