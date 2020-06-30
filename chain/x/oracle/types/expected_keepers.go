package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
}

// SupplyKeeper defines the expected supply Keeper.
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)

	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingexported.ValidatorI
	IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator stakingexported.ValidatorI) (stop bool))
	Validator(ctx sdk.Context, address sdk.ValAddress) stakingexported.ValidatorI
	Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec)
	Jail(ctx sdk.Context, consAddr sdk.ConsAddress)
}

// DistrKeeper defines the expected distribution keeper.
type DistrKeeper interface {
	GetCommunityTax(ctx sdk.Context) (percent sdk.Dec)
	GetFeePool(ctx sdk.Context) (feePool distr.FeePool)
	SetFeePool(ctx sdk.Context, feePool distr.FeePool)
	AllocateTokensToValidator(ctx sdk.Context, val stakingexported.ValidatorI, tokens sdk.DecCoins)
}
