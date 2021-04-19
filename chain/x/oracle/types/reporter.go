package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ReportersPerValidator represents list of reporter address and their associated validator
type ReportersPerValidator struct {
	Validator sdk.ValAddress   `json:"validator" yaml:"validator"`
	Reporters []sdk.AccAddress `json:"reporters" yaml:"reporters"`
}

// NewReportersPerValidator creates new instance of ReportersPerValidator
func NewReportersPerValidator(validator sdk.ValAddress, reporters []sdk.AccAddress) ReportersPerValidator {
	return ReportersPerValidator{
		Validator: validator,
		Reporters: reporters,
	}
}
