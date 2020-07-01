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
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
}

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingexported.ValidatorI
	IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator stakingexported.ValidatorI) (stop bool))
	Validator(ctx sdk.Context, address sdk.ValAddress) stakingexported.ValidatorI
}

// DistrKeeper defines the expected distribution keeper.
type DistrKeeper interface {
	GetCommunityTax(ctx sdk.Context) (percent sdk.Dec)
	GetFeePool(ctx sdk.Context) (feePool distr.FeePool)
	SetFeePool(ctx sdk.Context, feePool distr.FeePool)
	AllocateTokensToValidator(ctx sdk.Context, val stakingexported.ValidatorI, tokens sdk.DecCoins)
}
