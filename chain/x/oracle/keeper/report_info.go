package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/gogo/protobuf/types"

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

// TODO: If it isn't used in handler, it will be removed
// IterateValidatorReportInfos iterates over the stored ValidatorReportInfo
// func (k Keeper) IterateValidatorReportInfos(ctx sdk.Context,
// 	handler func(address sdk.ConsAddress, info types.ValidatorReportInfo) (stop bool)) {

// 	store := ctx.KVStore(k.storeKey)
// 	iter := sdk.KVStorePrefixIterator(store, types.ValidatorReportInfoKey)
// 	defer iter.Close()
// 	for ; iter.Valid(); iter.Next() {
// 		address := types.GetValidatorReportInfoAddress(iter.Key())
// 		var info types.ValidatorReportInfo
// 		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &info)
// 		if handler(address, info) {
// 			break
// 		}
// 	}
// }

// GetValidatorMissedReportBitArray gets the bit for the missed report array
func (k Keeper) GetValidatorMissedReportBitArray(ctx sdk.Context, address sdk.ValAddress, index uint64) bool {
	bz := ctx.KVStore(k.storeKey).Get(types.ValidatorMissedReportBitArrayKey(address, index))
	var missed gogotypes.BoolValue
	if bz == nil {
		// lazy: treat empty key as not missed
		return false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &missed)

	return missed.Value
}

// TODO: If it isn't used in handler, it will be removed
// IterateValidatorMissedBlockBitArray iterates over the signed blocks window
// and performs a callback function
// func (k Keeper) IterateValidatorMissedBlockBitArray(ctx sdk.Context,
// 	address sdk.ConsAddress, handler func(index int64, missed bool) (stop bool)) {

// 	store := ctx.KVStore(k.storeKey)
// 	index := int64(0)
// 	// Array may be sparse
// 	for ; index < k.SignedBlocksWindow(ctx); index++ {
// 		var missed gogotypes.BoolValue
// 		bz := store.Get(types.GetValidatorMissedReportBitArrayKey(address, index))
// 		if bz == nil {
// 			continue
// 		}

// 		k.cdc.MustUnmarshalBinaryBare(bz, &missed)
// 		if handler(index, missed.Value) {
// 			break
// 		}
// 	}
// }

// SetValidatorMissedReportBitArray sets the bit that checks if the validator has
// missed a report in the current window
func (k Keeper) SetValidatorMissedReportBitArray(ctx sdk.Context, address sdk.ValAddress, index uint64, missed bool) {
	ctx.KVStore(k.storeKey).Set(
		types.ValidatorMissedReportBitArrayKey(address, index),
		k.cdc.MustMarshalBinaryBare(&gogotypes.BoolValue{Value: missed}),
	)
}

// clearValidatorMissedReportBitArray deletes every instance of ValidatorMissedReportBitArray in the store
func (k Keeper) clearValidatorMissedReportBitArray(ctx sdk.Context, address sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorMissedReportBitArrayPrefixKey(address))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
