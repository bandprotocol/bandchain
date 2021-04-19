package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Reporter represents address of reporter and its associated validator
type Reporter struct {
	Reporter  sdk.AccAddress `json:"reporter" yaml:"reporter"`
	Validator sdk.ValAddress `json:"validator" yaml:"validator"`
}

// NewReporter creates new instance of Reporter
func NewReporter(reporter sdk.AccAddress, validator sdk.ValAddress) Reporter {
	return Reporter{
		Reporter:  reporter,
		Validator: validator,
	}
}
