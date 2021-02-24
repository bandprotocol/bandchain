package keeper

import (
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the odinmint store
type Keeper struct {
	mint.Keeper
	paramSpace params.Subspace
	cdc        *codec.Codec
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(cdc *codec.Codec, keeper mint.Keeper, paramSpace params.Subspace) Keeper {

	return Keeper{
		Keeper:     keeper,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
		cdc:        cdc,
	}
}

//______________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/odinmint")
}

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
