package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/zoracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AddReport adds a new data report from the given reporter on behalf of the given validator
// for a specific data request to the store. The function performs validations to make sure
// that the report is valid and saves it to the store.
func (k Keeper) AddReport(
	ctx sdk.Context, requestID types.RequestID, dataSet []types.RawDataReportWithID,
	validator sdk.ValAddress, reporter sdk.AccAddress,
) error {
	request, err := k.GetRequest(ctx, requestID)
	if err != nil {
		return err
	}

	if request.ResolveStatus != types.Open {
		return sdkerrors.Wrapf(types.ErrInvalidState,
			"AddReport: Request ID %d: Resolve status (%d) is not Open (%d).",
			requestID, request.ResolveStatus, types.Open,
		)
	}
	if request.ExpirationHeight < ctx.BlockHeight() {
		return sdkerrors.Wrapf(types.ErrInvalidState,
			"AddReport: Request ID %d: Request already expired at height %d. Current height is %d.",
			requestID, request.ExpirationHeight, ctx.BlockHeight(),
		)
	}
	if !k.CheckReporter(ctx, validator, reporter) {
		return sdkerrors.Wrapf(types.ErrUnauthorizedPermission,
			"AddReport: Request ID %d: %s is not an authorized reporter of %s.",
			requestID, reporter.String(), validator.String(),
		)
	}

	found := false
	for _, validValidator := range request.RequestedValidators {
		if validator.Equals(validValidator) {
			found = true
			break
		}
	}
	if !found {
		return sdkerrors.Wrapf(types.ErrUnauthorizedPermission,
			"AddReport: Request ID %d: %s is not one of the requested validators.",
			requestID, validator.String(),
		)
	}

	for _, submittedValidator := range request.ReceivedValidators {
		if validator.Equals(submittedValidator) {
			return sdkerrors.Wrapf(types.ErrItemDuplication,
				"AddReport: Duplicate report to request ID %d from validator %s.",
				requestID, submittedValidator.String(),
			)
		}
	}

	rawDataRequestCount := k.GetRawDataRequestCount(ctx, requestID)
	if int64(len(dataSet)) != rawDataRequestCount {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddReport: Request ID %d: Incorrect number (%d) of raw data reports. Expect %d.",
			requestID, len(dataSet), rawDataRequestCount,
		)
	}

	for _, rawReport := range dataSet {
		// Here we can safely assume that external IDs are unique, as this has already been
		// checked by ValidateBasic performed in baseapp's runTx function.
		if !k.CheckRawDataRequestExists(ctx, requestID, rawReport.ExternalDataID) {
			return sdkerrors.Wrapf(types.ErrBadDataValue,
				"AddReport: RequestID %d: Unknown external data ID %d",
				requestID, rawReport.ExternalDataID,
			)
		}
		if uint64(len(rawReport.Data)) > k.GetParam(ctx, types.KeyMaxRawDataReportSize) {
			return sdkerrors.Wrapf(types.ErrBadDataValue,
				"AddReport: RequestID %d: Raw report data size (%d) exceeds the limit (%d).",
				requestID, len(rawReport.Data), k.GetParam(ctx, types.KeyMaxRawDataReportSize),
			)
		}
		k.SetRawDataReport(
			ctx, requestID, rawReport.ExternalDataID, validator,
			types.RawDataReport{
				ExitCode: rawReport.ExitCode,
				Data:     rawReport.Data,
			},
		)
	}

	request.ReceivedValidators = append(request.ReceivedValidators, validator)
	k.SetRequest(ctx, requestID, request)
	if k.ShouldBecomePendingResolve(ctx, requestID) {
		err := k.AddPendingRequest(ctx, requestID)
		if err != nil {
			// Should never happen because we already perform ShouldBecomePendingResolve check.
			return err
		}
	}
	return nil
}

// SetRawDataReport is a function that saves a raw data report to store.
func (k Keeper) SetRawDataReport(
	ctx sdk.Context,
	requestID types.RequestID, externalID types.ExternalID,
	validatorAddress sdk.ValAddress,
	rawDataReport types.RawDataReport,
) {
	key := types.RawDataReportStoreKey(requestID, externalID, validatorAddress)
	store := ctx.KVStore(k.storeKey)
	store.Set(key, k.cdc.MustMarshalBinaryBare(rawDataReport))
}

// GetRawDataReport returns the raw data report from the store for a specific combination of
// request ID, external ID, and validator address.
func (k Keeper) GetRawDataReport(
	ctx sdk.Context,
	requestID types.RequestID, externalID types.ExternalID,
	validatorAddress sdk.ValAddress,
) (types.RawDataReport, error) {
	key := types.RawDataReportStoreKey(requestID, externalID, validatorAddress)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return types.RawDataReport{}, sdkerrors.Wrapf(types.ErrItemNotFound,
			"GetRawDataReport: No raw data report with request ID %d external ID %d from %s",
			requestID,
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
	ctx sdk.Context, requestID types.RequestID,
) sdk.Iterator {
	prefix := types.GetIteratorPrefix(types.RawDataReportStoreKeyPrefix, requestID)
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, prefix)
}
