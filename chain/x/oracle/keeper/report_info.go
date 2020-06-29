package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// HasReportInfo checks if a given validator has report information exists.
func (k Keeper) HasReportInfo(ctx sdk.Context, address sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.ReportInfoStoreKey(address))
}

// GetReportInfoWithDefault returns the ReportInfo for the given validator.
// Returns the default value with 0 misses if the info does not exist in the store.
func (k Keeper) GetReportInfoWithDefault(ctx sdk.Context, val sdk.ValAddress) types.ReportInfo {
	bz := ctx.KVStore(k.storeKey).Get(types.ReportInfoStoreKey(val))
	if bz == nil {
		return types.NewReportInfo(val, 0)
	}
	var info types.ReportInfo
	k.cdc.MustUnmarshalBinaryBare(bz, &info)
	return info
}

// SetReportInfo sets the validator report info to for the given validator.
func (k Keeper) SetReportInfo(ctx sdk.Context, val sdk.ValAddress, info types.ReportInfo) {
	ctx.KVStore(k.storeKey).Set(types.ReportInfoStoreKey(val), k.cdc.MustMarshalBinaryBare(info))
}

// GetAllReportInfos returns the list of all validator report info in the store, or nil if there is none.
func (k Keeper) GetAllReportInfos(ctx sdk.Context) (infos []types.ReportInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ReportInfoKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var info types.ReportInfo
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &info)
		infos = append(infos, info)
	}
	return infos
}
