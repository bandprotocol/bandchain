package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasReport checks if the report of this ID triple exists in the storage.
func (k Keeper) HasReport(ctx sdk.Context, rid types.RID, eid types.EID, val sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.RawDataReportStoreKey(rid, eid, val))
}

// SetDataReport saves the report to the storage without performing validation.
func (k Keeper) SetReport(ctx sdk.Context, rid types.RID, val sdk.ValAddress, report types.Report) {
	key := types.RawDataReportStoreKey(rid, report.ExternalID, val)
	ctx.KVStore(k.storeKey).Set(key, k.cdc.MustMarshalBinaryBare(report))
}

// AddReports performs sanity checks and adds a new batch report from one validator to one request
// to the store. Note that we expect each validator to report to all raw data requests at once.
func (k Keeper) AddBatchReport(ctx sdk.Context, rid types.RID, reps types.BatchReport) error {
	req, err := k.GetRequest(ctx, rid)
	if err != nil {
		return err
	}

	// TODO: Make sure we handle this correctly. We basically want to allow reporters to report
	// even after a request is resolved so they get to keep their stats.
	// if req.ResolveStatus != types.Open {
	// 	return sdkerrors.Wrapf(types.ErrInvalidState,
	// 		"AddReport: Request ID %d: Resolve status (%d) is not Open (%d).",
	// 		rid, req.ResolveStatus, types.Open,
	// 	)
	// }
	// if req.ExpirationHeight < ctx.BlockHeight() {
	// 	return sdkerrors.Wrapf(types.ErrInvalidState,
	// 		"AddReport: Request ID %d: Request already expired at height %d. Current height is %d.",
	// 		rid, req.ExpirationHeight, ctx.BlockHeight(),
	// 	)
	// }
	if !ContainsVal(req.RequestedValidators, reps.Validator) {
		return sdkerrors.Wrapf(types.ErrUnauthorizedPermission,
			"AddReport: Request ID %d: %s is not one of the requested validators.",
			rid, reps.Validator.String(),
		)
	}
	if ContainsVal(req.ReceivedValidators, reps.Validator) {
		return sdkerrors.Wrapf(types.ErrItemDuplication,
			"AddReport: Duplicate report to request ID %d from validator %s.",
			rid, reps.Validator.String(),
		)
	}

	rawDataRequestCount := k.GetRawRequestCount(ctx, rid)
	if int64(len(reps.RawDataReports)) != rawDataRequestCount {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddReport: Request ID %d: Incorrect number (%d) of raw data reports. Expect %d.",
			rid, len(reps.RawDataReports), rawDataRequestCount,
		)
	}

	for _, report := range reps.RawDataReports {
		// Here we can safely assume that external IDs are unique, as this has already been
		// checked by ValidateBasic performed in baseapp's runTx function.
		if !k.HasRawRequest(ctx, rid, report.ExternalID) {
			return sdkerrors.Wrapf(types.ErrBadDataValue,
				"AddReport: RequestID %d: Unknown external data ID %d",
				rid, report.ExternalID,
			)
		}
		if uint64(len(report.Data)) > k.GetParam(ctx, types.KeyMaxRawDataReportSize) {
			return sdkerrors.Wrapf(types.ErrBadDataValue,
				"AddReport: RequestID %d: Raw report data size (%d) exceeds the limit (%d).",
				rid, len(report.Data), k.GetParam(ctx, types.KeyMaxRawDataReportSize),
			)
		}
		k.SetReport(ctx, rid, reps.Validator, report)
	}

	req.ReceivedValidators = append(req.ReceivedValidators, reps.Validator)
	k.SetRequest(ctx, rid, req)
	if k.ShouldBecomePendingResolve(ctx, rid) {
		err := k.AddPendingRequest(ctx, rid)
		if err != nil {
			// Should never happen because we already perform ShouldBecomePendingResolve check.
			return err
		}
	}
	return nil
}

// GetRawDataReport returns the raw data report from the store for a specific combination of
// request ID, external ID, and validator address.
func (k Keeper) GetRawDataReport(
	ctx sdk.Context,
	rid types.RID, externalID types.ExternalID,
	validatorAddress sdk.ValAddress,
) (types.RawDataReport, error) {
	key := types.RawDataReportStoreKey(rid, externalID, validatorAddress)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return types.RawDataReport{}, sdkerrors.Wrapf(types.ErrItemNotFound,
			"GetRawDataReport: No raw data report with request ID %d external ID %d from %s",
			rid,
			externalID,
			validatorAddress.String(),
		)
	}
	var rawDataReport types.RawDataReport
	k.cdc.MustUnmarshalBinaryBare(store.Get(key), &rawDataReport)
	return rawDataReport, nil
}

// GetRawDataReportsIterator returns an iterator for all reports for a specific request ID.
func (k Keeper) GetRawDataReportsIterator(
	ctx sdk.Context, rid types.RID,
) sdk.Iterator {
	prefix := types.GetIteratorPrefix(types.RawDataReportStoreKeyPrefix, rid)
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, prefix)
}
