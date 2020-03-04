package supply

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

type WrappedSupplyKeeper struct {
	supply.Keeper
}

func (k WrappedSupplyKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error {
	// create the account if it doesn't yet exist
	// acc := k.GetModuleAccount(ctx, moduleName)
	// if acc == nil {
	// 		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", moduleName))
	// }

	// if !acc.HasPermission(types.Burner) {
	// 		panic(fmt.Sprintf("module account %s does not have permissions to burn tokens", moduleName))
	// }

	// _, err := k.bk.SubtractCoins(ctx, acc.GetAddress(), amt)
	// if err != nil {
	// 		panic(err)
	// }

	// // update total supply
	// supply := k.GetSupply(ctx)
	// supply = supply.Deflate(amt)
	// k.SetSupply(ctx, supply)

	// logger := k.Logger(ctx)
	// logger.Info(fmt.Sprintf("burned %s from %s module account", amt.String(), moduleName))

	// return nil
	return nil
}

func WrapSupplyKeeperBurnToCommunityPool(supplyKeeper supply.Keeper) WrappedSupplyKeeper {
	return WrappedSupplyKeeper{supplyKeeper}
}
