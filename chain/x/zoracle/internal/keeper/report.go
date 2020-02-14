package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

func (k Keeper) AddReport(
	ctx sdk.Context, requestID int64, dataSet []types.RawDataReport, validator sdk.ValAddress,
) sdk.Error {
	request, err := k.GetRequest(ctx, requestID)
	if err != nil {
		return err
	}

	if request.IsResolved || request.ExpirationHeight < ctx.BlockHeight() {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	found := false
	for _, validValidator := range request.RequestedValidators {
		if validator.Equals(validValidator) {
			found = true
			break
		}
	}
	if !found {
		return types.ErrInvalidValidator(types.DefaultCodespace)
	}

	for _, submittedValidator := range request.ReceivedValidators {
		if validator.Equals(submittedValidator) {
			// TODO: fix error later
			return types.ErrInvalidValidator(types.DefaultCodespace)
		}
	}

	if int64(len(dataSet)) != k.GetRawDataRequestCount(ctx, requestID) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	lastExternalID := int64(0)
	for idx, rawReport := range dataSet {
		if idx != 0 && lastExternalID >= rawReport.ExternalDataID {
			// TODO: fix error later
			return types.ErrRequestNotFound(types.DefaultCodespace)
		}
		if !k.CheckRawDataRequestExists(ctx, requestID, rawReport.ExternalDataID) {
			// TODO: fix error later
			return types.ErrRequestNotFound(types.DefaultCodespace)
		}
		k.SetRawDataReport(ctx, requestID, rawReport.ExternalDataID, validator, rawReport.Data)
		lastExternalID = rawReport.ExternalDataID
	}

	request.ReceivedValidators = append(request.ReceivedValidators, validator)
	k.SetRequest(ctx, requestID, request)
	if k.ShouldBecomePendingResolve(ctx, requestID) {
		k.AddPendingRequest(ctx, requestID)
	}

	return nil
}

// SetRawDataReport is a function that saves a raw data report to store.
func (k Keeper) SetRawDataReport(
	ctx sdk.Context,
	requestID, externalID int64,
	validatorAddress sdk.ValAddress,
	data []byte,
) {
	key := types.RawDataReportStoreKey(requestID, externalID, validatorAddress)
	store := ctx.KVStore(k.storeKey)
	store.Set(key, data)
}

// GetRawDataReport is a function that gets a raw data report from store.
func (k Keeper) GetRawDataReport(
	ctx sdk.Context,
	requestID, externalID int64,
	validatorAddress sdk.ValAddress,
) ([]byte, sdk.Error) {
	key := types.RawDataReportStoreKey(requestID, externalID, validatorAddress)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return []byte{}, types.ErrReportNotFound(types.DefaultCodespace)
	}
	return store.Get(key), nil
}

// // GetReportsIterator returns an iterator for all reports for a specific request ID.
// func (k Keeper) GetReportsIterator(ctx sdk.Context, requestID int64) sdk.Iterator {
// 	prefix := types.GetIteratorPrefix(types.ReportKeyPrefix, requestID)
// 	store := ctx.KVStore(k.storeKey)
// 	return sdk.KVStorePrefixIterator(store, prefix)
// }

// // GetDataReports returns all the reports for a specific request ID.
// func (k Keeper) GetDataReports(ctx sdk.Context, requestID int64) []types.Report {
// 	iterator := k.GetReportsIterator(ctx, requestID)
// 	var data []types.Report
// 	for ; iterator.Valid(); iterator.Next() {
// 		var report types.Report
// 		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &report)
// 		data = append(data, report)
// 	}
// 	return data
// }

// // GetValidatorReports returns all the reports (each including its reporter) for a specific request ID.
// func (k Keeper) GetValidatorReports(ctx sdk.Context, requestID int64) ([]types.ReportWithValidator, sdk.Error) {
// 	iterator := k.GetReportsIterator(ctx, requestID)
// 	data := make([]types.ReportWithValidator, 0)

// 	// Check request is existed
// 	if !k.CheckRequestExists(ctx, requestID) {
// 		return []types.ReportWithValidator{}, types.ErrRequestNotFound(types.DefaultCodespace)
// 	}

// 	for ; iterator.Valid(); iterator.Next() {
// 		var report types.Report
// 		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &report)

// 		vReport := types.NewReportWithValidator(
// 			report.Data,
// 			report.ReportedAt,
// 			types.GetValidatorAddress(iterator.Key(), types.ReportKeyPrefix, requestID),
// 		)
// 		data = append(data, vReport)
// 	}
// 	return data, nil
// }
