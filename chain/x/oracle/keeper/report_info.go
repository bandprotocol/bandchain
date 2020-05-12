package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// HasValidatorReportInfo checks if a given validator has report information persited.
func (k Keeper) HasValidatorReportInfo(ctx sdk.Context, address sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.ValidatorReportInfoStoreKey(address))
}

// GetValidatorReportInfo returns the ValidatorReportInfo for the giver validator address.
func (k Keeper) GetValidatorReportInfo(
	ctx sdk.Context, address sdk.ValAddress,
) (types.ValidatorReportInfo, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.ValidatorReportInfoStoreKey(address))
	if bz == nil {
		return types.ValidatorReportInfo{}, sdkerrors.Wrapf(
			types.ErrValidatorReportInfoNotFound, "address: %s", address.String(),
		)
	}
	var info types.ValidatorReportInfo
	k.cdc.MustUnmarshalBinaryBare(bz, &info)
	return info, nil
}

// MustGetValidatorReportInfo returns the ValidatorReportInfo for the giver validator address.
// Panics error if not exists.
func (k Keeper) MustGetValidatorReportInfo(
	ctx sdk.Context, address sdk.ValAddress,
) types.ValidatorReportInfo {
	info, err := k.GetValidatorReportInfo(ctx, address)
	if err != nil {
		panic(err)
	}
	return info
}

// SetValidatorReportInfo sets the validator report info to a validator address key
func (k Keeper) SetValidatorReportInfo(ctx sdk.Context, address sdk.ValAddress, info types.ValidatorReportInfo) {
	ctx.KVStore(k.storeKey).Set(types.ValidatorReportInfoStoreKey(address), k.cdc.MustMarshalBinaryBare(info))
}
