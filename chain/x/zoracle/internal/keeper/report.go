package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

// SetReport is a function that saves a data report to storeage.
func (k Keeper) SetReport(
	ctx sdk.Context,
	requestID int64,
	validatorAddress sdk.ValAddress,
	data []types.RawDataReport,
) {
	key := types.ReportStoreKey(requestID, validatorAddress)
	report := types.NewReport(data, ctx.BlockHeight())
	store := ctx.KVStore(k.storeKey)
	store.Set(key, k.cdc.MustMarshalBinaryBare(report))
}

// GetReportsIterator returns an iterator for all reports for a specific request ID.
func (k Keeper) GetReportsIterator(ctx sdk.Context, requestID int64) sdk.Iterator {
	prefix := types.GetIteratorPrefix(types.ReportKeyPrefix, requestID)
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, prefix)
}

// GetDataReports returns all the reports for a specific request ID.
func (k Keeper) GetDataReports(ctx sdk.Context, requestID int64) []types.Report {
	iterator := k.GetReportsIterator(ctx, requestID)
	var data []types.Report
	for ; iterator.Valid(); iterator.Next() {
		var report types.Report
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &report)
		data = append(data, report)
	}
	return data
}

// GetValidatorReports returns all the reports (each including its reporter) for a specific request ID.
func (k Keeper) GetValidatorReports(ctx sdk.Context, requestID int64) ([]types.ReportWithValidator, sdk.Error) {
	iterator := k.GetReportsIterator(ctx, requestID)
	data := make([]types.ReportWithValidator, 0)

	// Check request is existed
	if !k.CheckRequestExists(ctx, requestID) {
		return []types.ReportWithValidator{}, types.ErrRequestNotFound(types.DefaultCodespace)
	}

	for ; iterator.Valid(); iterator.Next() {
		var report types.Report
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &report)

		vReport := types.NewReportWithValidator(
			report.Data,
			report.ReportedAt,
			types.GetValidatorAddress(iterator.Key(), types.ReportKeyPrefix, requestID),
		)
		data = append(data, vReport)
	}
	return data, nil
}
