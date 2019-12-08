package keeper

import (
	"github.com/bandprotocol/bandx/oracle/x/oracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetReport is a function that save report to storeage
func (k Keeper) SetReport(ctx sdk.Context, requestID uint64, validatorAddress sdk.ValAddress, data []byte) {
	key := types.ReportStoreKey(requestID, validatorAddress)
	report := types.NewReport(data, uint64(ctx.BlockHeight()))
	store := ctx.KVStore(k.storeKey)
	store.Set(key, k.cdc.MustMarshalBinaryBare(report))
}

func (k Keeper) GetReportsIterator(ctx sdk.Context, requestID uint64) sdk.Iterator {
	prefix := types.GetIteratorPrefix(types.ReportKeyPrefix, requestID)
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, prefix)
}

func (k Keeper) GetDataReports(ctx sdk.Context, requestID uint64) []types.Report {
	iterator := k.GetReportsIterator(ctx, requestID)
	var data []types.Report
	for ; iterator.Valid(); iterator.Next() {
		var report types.Report
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &report)
		data = append(data, report)
	}
	return data
}

func (k Keeper) GetValidatorReports(ctx sdk.Context, requestID uint64) []types.ValidatorReport {
	iterator := k.GetReportsIterator(ctx, requestID)
	var data []types.ValidatorReport
	for ; iterator.Valid(); iterator.Next() {
		var report types.Report
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &report)
		vReport := types.NewValidatorReport(
			report,
			types.GetValidatorAddress(iterator.Key(), types.ReportKeyPrefix, requestID),
		)
		data = append(data, vReport)
	}
	return data
}
