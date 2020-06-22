package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
)

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator stakingexported.ValidatorI) (stop bool))
	Validator(ctx sdk.Context, address sdk.ValAddress) stakingexported.ValidatorI
	Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec)
	Jail(ctx sdk.Context, consAddr sdk.ConsAddress)
}
