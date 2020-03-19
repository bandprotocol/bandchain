package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
)

func (k Keeper) AddReport(
	ctx sdk.Context,
	requestID types.RequestID,
	dataSet []types.RawDataReportWithID,
	validator sdk.ValAddress,
	reporter sdk.AccAddress,
) sdk.Error {
	request, err := k.GetRequest(ctx, requestID)
	if err != nil {
		return err
	}

	if request.ResolveStatus != types.Open {
		return types.ErrInvalidState(
			"AddReport: Request ID %d: Expect resolve status to be %d, but actual value is %d",
			requestID,
			types.Open,
			request.ResolveStatus,
		)
	}

	if request.ExpirationHeight < ctx.BlockHeight() {
		return types.ErrInvalidState(
			"AddReport: Request ID %d: Current block height is %d, but request expired at height %d",
			requestID,
			ctx.BlockHeight(),
			request.ExpirationHeight,
		)
	}

	found := false
	for _, validValidator := range request.RequestedValidators {
		if validator.Equals(validValidator) {
			found = true
			break
		}
	}

	if k.CheckReporter(ctx, validator, reporter) {
		agent := sdk.ValAddress(reporter)
		for _, validValidator := range request.RequestedValidators {
			if agent.Equals(validValidator) {
				found = true
				break
			}
		}
	}

	if !found {
		return types.ErrUnauthorizedPermission(
			"AddReport: Reporter (%s) and (%s) is not on the reporter list.",
			validator.String(),
			reporter.String(),
		)
	}

	for _, submittedValidator := range request.ReceivedValidators {
		if validator.Equals(submittedValidator) {
			return types.ErrItemDuplication(
				"AddReport: Duplicate report to request ID %d from reporter %s.",
				requestID,
				submittedValidator.String(),
			)
		}
	}

	rawDataRequestCount := k.GetRawDataRequestCount(ctx, requestID)
	if int64(len(dataSet)) != rawDataRequestCount {
		return types.ErrBadDataValue(
			"AddReport: Request ID %d: Expects %d raw data reports, but received %d raw data reports.",
			requestID,
			rawDataRequestCount,
			len(dataSet),
		)
	}

	lastExternalID := types.ExternalID(0)
	for idx, rawReport := range dataSet {
		if idx != 0 && lastExternalID >= rawReport.ExternalDataID {
			return types.ErrBadDataValue("AddReport: Raw data reports are not in an incresaing order.")
		}
		if !k.CheckRawDataRequestExists(ctx, requestID, rawReport.ExternalDataID) {
			return types.ErrBadDataValue(
				"AddReport: RequestID %d: Unknown external ID %d",
				requestID,
				rawReport.ExternalDataID,
			)
		}
		if int64(len(rawReport.Data)) > k.MaxRawDataReportSize(ctx) {
			return types.ErrBadDataValue(
				"AddReport: Raw report data size (%d) exceeds the maximum limit (%d).",
				len(rawReport.Data),
				k.MaxRawDataReportSize(ctx),
			)
		}
		k.SetRawDataReport(
			ctx,
			requestID,
			rawReport.ExternalDataID,
			validator,
			types.RawDataReport{
				ExitCode: rawReport.ExitCode,
				Data:     rawReport.Data,
			},
		)
		lastExternalID = rawReport.ExternalDataID
	}

	request.ReceivedValidators = append(request.ReceivedValidators, validator)
	k.SetRequest(ctx, requestID, request)
	if k.ShouldBecomePendingResolve(ctx, requestID) {
		err := k.AddPendingRequest(ctx, requestID)
		if err != nil {
			// This should never happen, but we detect it anyway just in case.
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

// GetRawDataReport is a function that gets a raw data report from store.
func (k Keeper) GetRawDataReport(
	ctx sdk.Context,
	requestID types.RequestID, externalID types.ExternalID,
	validatorAddress sdk.ValAddress,
) (types.RawDataReport, sdk.Error) {
	key := types.RawDataReportStoreKey(requestID, externalID, validatorAddress)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return types.RawDataReport{}, types.ErrItemNotFound(
			"GetRawDataReport: Unable to find raw data report with request ID %d external ID %d from %s",
			requestID,
			externalID,
			validatorAddress.String(),
		)
	}
	bz := store.Get(key)
	var rawDataReport types.RawDataReport
	k.cdc.MustUnmarshalBinaryBare(bz, &rawDataReport)
	return rawDataReport, nil
}

// GetRawDataReportsIterator returns an iterator for all reports for a specific request ID.
func (k Keeper) GetRawDataReportsIterator(ctx sdk.Context, requestID types.RequestID) sdk.Iterator {
	prefix := types.GetIteratorPrefix(types.RawDataReportStoreKeyPrefix, requestID)
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, prefix)
}
