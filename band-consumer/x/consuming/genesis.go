package consuming

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState is the band-consumer state that must be provided at genesis.
type GenesisState struct{}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}

// DefaultGenesisState returns the default genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{}
}
