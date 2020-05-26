package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// HasValidatorReportInfo checks if a given validator has report information exists.
func (k Keeper) HasValidatorReportInfo(ctx sdk.Context, address sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.ValidatorReportInfoStoreKey(address))
}

// GetValidatorReportInfoWithDefault returns the ValidatorReportInfo for the given validator.
// Returns the default value with 0 misses if the info does not exist in the store.
func (k Keeper) GetValidatorReportInfoWithDefault(ctx sdk.Context, val sdk.ValAddress) types.ValidatorReportInfo {
	bz := ctx.KVStore(k.storeKey).Get(types.ValidatorReportInfoStoreKey(val))
	if bz == nil {
		return types.NewValidatorReportInfo(val, 0)
	}
	var info types.ValidatorReportInfo
	k.cdc.MustUnmarshalBinaryBare(bz, &info)
	return info
}

// SetValidatorReportInfo sets the validator report info to for the given validator.
func (k Keeper) SetValidatorReportInfo(ctx sdk.Context, val sdk.ValAddress, info types.ValidatorReportInfo) {
	ctx.KVStore(k.storeKey).Set(types.ValidatorReportInfoStoreKey(val), k.cdc.MustMarshalBinaryBare(info))
}

// GetAllValidatorReportInfos returns the list of all validator report info in the store, or nil if there is none.
func (k Keeper) GetAllValidatorReportInfos(ctx sdk.Context) (infos []types.ValidatorReportInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorReportInfoKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var info types.ValidatorReportInfo
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &info)
		infos = append(infos, info)
	}
	return infos
}
