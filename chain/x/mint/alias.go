package mint

// nolint

import (
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/keeper"
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/types"
)

const (
	DefaultParamspace = types.DefaultParamspace
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	ParamKeyTable       = types.ParamKeyTable
	NewParams           = types.NewParams
	DefaultParams       = types.DefaultParams

	// variable aliases
	KeyMintAir = types.KeyMintAir
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
)
