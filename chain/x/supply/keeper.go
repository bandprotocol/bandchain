package supply

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// WrappedSupplyKeeper encapsulates the underlying supply keeper and overrides
// its BurnCoins function to send the coins to the community pool instead of
// just destroying them.
//
// Note that distrKeeper keeps the reference to the distr module keeper.
// Due to the circular dependency between supply-distr-staking, distrKeeper
// cannot be initialized when the struct is created. Rather, SetDistrKeeper
// is expected to be called to set `distrKeeper`.
type WrappedSupplyKeeper struct {
	supply.Keeper

	distrKeeper *distr.Keeper
}

// WrapSupplyKeeperBurnToCommunityPool creates a new instance of WrappedSupplyKeeper
// with its distrKeeper member set to nil.
func WrapSupplyKeeperBurnToCommunityPool(supplyKeeper supply.Keeper) WrappedSupplyKeeper {
	return WrappedSupplyKeeper{supplyKeeper, nil}
}

// SetDistrKeeper sets distr module keeper for this WrappedSupplyKeeper instance.
func (k *WrappedSupplyKeeper) SetDistrKeeper(distrKeeper *distr.Keeper) {
	k.distrKeeper = distrKeeper
}

// BurnCoins moves the specified amount of coins from the given module name to
// the community pool. The total supply of the coins will not change.
func (k WrappedSupplyKeeper) BurnCoins(
	ctx sdk.Context, moduleName string, amt sdk.Coins,
) error {
	// If distrKeeper is not set OR we want to burn coins in distr itself, we will
	// just use the original BurnCoins function.
	if k.distrKeeper == nil || moduleName == distr.ModuleName {
		return k.Keeper.BurnCoins(ctx, moduleName, amt)
	}

	// Create the account if it doesn't yet exist.
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(sdkerrors.Wrapf(
			sdkerrors.ErrUnknownAddress,
			"module account %s does not exist", moduleName,
		))
	}

	if !acc.HasPermission(supply.Burner) {
		panic(sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"module account %s does not have permissions to burn tokens", moduleName,
		))
	}

	// Instead of burning coins, we send them to the community pool.
	k.SendCoinsFromModuleToModule(ctx, moduleName, distr.ModuleName, amt)
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(amt...)...)
	k.distrKeeper.SetFeePool(ctx, feePool)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf(
		"sent %s from %s module account to community pool", amt.String(), moduleName,
	))
	return nil
}
