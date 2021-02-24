package types

import "github.com/cosmos/cosmos-sdk/x/mint"

// GenesisState - minter state
type GenesisState struct {
	Minter mint.Minter `json:"minter" yaml:"minter"` // minter object
	Params Params      `json:"params" yaml:"params"` // inflation params
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter mint.Minter, params Params) GenesisState {
	return GenesisState{
		Minter: minter,
		Params: params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Minter: mint.DefaultInitialMinter(),
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return mint.ValidateMinter(data.Minter)
}
