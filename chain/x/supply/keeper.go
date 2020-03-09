package supply

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

type WrappedSupplyKeeper struct {
	supply.Keeper

	distrKeeper *distr.Keeper
}

func WrapSupplyKeeperBurnToCommunityPool(supplyKeeper supply.Keeper) WrappedSupplyKeeper {
	return WrappedSupplyKeeper{supplyKeeper, nil}
}

func (k *WrappedSupplyKeeper) SetDistrKeeper(distrKeeper *distr.Keeper) {
	k.distrKeeper = distrKeeper
}

func (k WrappedSupplyKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error {
	// if distrKeeper is not set or we want to burn coins in distr, we use the original BurnCoins function
	if k.distrKeeper == nil || moduleName == distr.ModuleName {
		return k.Keeper.BurnCoins(ctx, moduleName, amt)
	}

	// create the account if it doesn't yet exist
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(supply.Burner) {
		panic(fmt.Sprintf("module account %s does not have permissions to burn tokens", moduleName))
	}

	// instead of burning coins, send them to the community pool
	k.SendCoinsFromModuleToModule(ctx, moduleName, distr.ModuleName, amt)
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoins(amt))
	k.distrKeeper.SetFeePool(ctx, feePool)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("sent %s from %s module account to community pool", amt.String(), moduleName))
	return nil
}
