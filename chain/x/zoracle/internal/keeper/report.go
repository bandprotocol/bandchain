package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/d3n/chain/wasm"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

// SetReport is a function that saves a data report to storeage.
func (k Keeper) SetReport(ctx sdk.Context, requestID uint64, validatorAddress sdk.ValAddress, data []byte) {
	key := types.ReportStoreKey(requestID, validatorAddress)
	report := types.NewReport(data, uint64(ctx.BlockHeight()))
	store := ctx.KVStore(k.storeKey)
	store.Set(key, k.cdc.MustMarshalBinaryBare(report))
}

// GetReportsIterator returns an iterator for all reports for a specific request ID.
func (k Keeper) GetReportsIterator(ctx sdk.Context, requestID uint64) sdk.Iterator {
	prefix := types.GetIteratorPrefix(types.ReportKeyPrefix, requestID)
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, prefix)
}

// GetDataReports returns all the reports for a specific request ID.
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

// GetDataReports returns all the reports (each including its reporter) for a specific request ID.
func (k Keeper) GetValidatorReports(ctx sdk.Context, requestID uint64) ([]types.ValidatorReport, sdk.Error) {
	iterator := k.GetReportsIterator(ctx, requestID)
	data := make([]types.ValidatorReport, 0)

	// Get code and params from request id
	request, err := k.GetRequest(ctx, requestID)
	if err != nil {
		return nil, err
	}

	code, err := k.GetCode(ctx, request.CodeHash)
	if err != nil {
		return nil, err
	}
	for ; iterator.Valid(); iterator.Next() {
		var report types.Report
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &report)

		rawValue, err := wasm.ParseRawData(code.Code, request.Params, report.Data)
		if err != nil {
			rawValue = []byte("null")
		}
		vReport := types.NewValidatorReport(
			rawValue,
			report.ReportedAt,
			types.GetValidatorAddress(iterator.Key(), types.ReportKeyPrefix, requestID),
		)
		data = append(data, vReport)
	}
	return data, nil
}
