package supply

import (
	"fmt"
	odinmint "github.com/GeoDB-Limited/odincore/chain/x/mint"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"
)

// WrappedSupplyKeeper encapsulates the underlying supply keeper and overrides
// its BurnCoins function to send the coins to the community pool instead of
// just destroying them.
//
// Note that distrKeeper keeps the reference to the distr module keeper.
// Due to the circular dependency between supply-distr, distrKeeper
// cannot be initialized when the struct is created. Rather, SetDistrKeeper
// is expected to be called to set `distrKeeper`.
type WrappedSupplyKeeper struct {
	supply.Keeper

	mintKeeper  *odinmint.Keeper
	distrKeeper *distr.Keeper
}

// WrapSupplyKeeperBurnToCommunityPool creates a new instance of WrappedSupplyKeeper
// with its distrKeeper member set to nil.
func WrapSupplyKeeperBurnToCommunityPool(sk supply.Keeper) WrappedSupplyKeeper {
	return WrappedSupplyKeeper{
		Keeper: sk,
	}
}

// SetDistrKeeper sets distr module keeper for this WrappedSupplyKeeper instance.
func (k *WrappedSupplyKeeper) SetDistrKeeper(distrKeeper *distr.Keeper) {
	k.distrKeeper = distrKeeper
}

func (k *WrappedSupplyKeeper) SetMintKeeper(mintKeeper *odinmint.Keeper) {
	k.mintKeeper = mintKeeper
}

// Logger returns a module-specific logger.
func (k WrappedSupplyKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprint("x/wrappedSupply"))
}

// BurnCoins moves the specified amount of coins from the given module name to
// the community pool. The total supply of the coins will not change.
func (k WrappedSupplyKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
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
			"module account %s does not have permissions to burn tokens",
			moduleName,
		))
	}

	logger := k.Logger(ctx)
	// Instead of burning coins, we send them to the community pool.
	err := k.SendCoinsFromModuleToModule(ctx, moduleName, distr.ModuleName, amt)
	if err != nil {
		err = sdkerrors.Wrap(err, fmt.Sprintf("failed to mint %s from %s module account", amt.String(), moduleName))
		logger.Error(err.Error())
		return err
	}
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(amt...)...)
	k.distrKeeper.SetFeePool(ctx, feePool)

	logger.Info(fmt.Sprintf(
		"sent %s from %s module account to community pool", amt.String(), moduleName,
	))
	return nil
}

// MintCoins does not create any new coins, just gets them from the community pull
func (k WrappedSupplyKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	vanillaMinting := k.mintKeeper.GetParams(ctx).MintAir
	if vanillaMinting {
		return k.Keeper.MintCoins(ctx, moduleName, amt)
	}
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(supply.Minter) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to mint tokens", moduleName))
	}

	logger := k.Logger(ctx)
	err := k.SendCoinsFromModuleToModule(ctx, distr.ModuleName, moduleName, amt)
	if err != nil {
		err = sdkerrors.Wrap(err, fmt.Sprintf("failed to mint %s from %s module account", amt.String(), moduleName))
		logger.Error(err.Error())
		return err
	}

	logger.Info(fmt.Sprintf("minted %s from %s module account", amt.String(), moduleName))

	return nil
}
