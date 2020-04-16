package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SetRawDataRequest saves the raw data request to the store without performing validation.
func (k Keeper) SetRawDataRequest(
	ctx sdk.Context, requestID types.RequestID, externalID types.ExternalID,
	rawDataRequest types.RawDataRequest,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.RawDataRequestStoreKey(requestID, externalID),
		k.cdc.MustMarshalBinaryBare(rawDataRequest),
	)
}

// GetRawDataRequest returns the raw data request detail by the given request ID and external ID.
func (k Keeper) GetRawDataRequest(
	ctx sdk.Context, requestID types.RequestID, externalID types.ExternalID,
) (types.RawDataRequest, error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckRawDataRequestExists(ctx, requestID, externalID) {
		return types.RawDataRequest{}, sdkerrors.Wrapf(types.ErrItemNotFound,
			"GetRawDataRequest: Unknown raw data request for request ID %d external ID %d.",
			requestID, externalID,
		)
	}

	bz := store.Get(types.RawDataRequestStoreKey(requestID, externalID))
	var requestDetail types.RawDataRequest
	k.cdc.MustUnmarshalBinaryBare(bz, &requestDetail)
	return requestDetail, nil
}

// CheckRawDataRequestExists checks if the raw data request at this request ID and external ID
// exists in the store or not.
func (k Keeper) CheckRawDataRequestExists(
	ctx sdk.Context, requestID types.RequestID, externalID types.ExternalID,
) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.RawDataRequestStoreKey(requestID, externalID))
}

// AddNewRawDataRequest performs all sanity checks and adds a new raw data request to the store.
func (k Keeper) AddNewRawDataRequest(
	ctx sdk.Context, requestID types.RequestID, externalID types.ExternalID,
	dataSourceID types.DataSourceID, calldata []byte,
) error {
	if uint64(len(calldata)) > k.GetParam(ctx, types.KeyMaxCalldataSize) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddNewRawDataRequest: Calldata size (%d) exceeds the maximum size (%d).",
			len(calldata), k.GetParam(ctx, types.KeyMaxCalldataSize),
		)
	}

	request, err := k.GetRequest(ctx, requestID)
	if err != nil {
		return err
	}

	if !k.HasDataSource(ctx, dataSourceID) {
		return sdkerrors.Wrapf(types.ErrDataSourceNotFound, "id: %d", dataSourceID)
	}

	if k.CheckRawDataRequestExists(ctx, requestID, externalID) {
		return sdkerrors.Wrapf(types.ErrItemDuplication,
			"AddNewRawDataRequest: Request ID %d: Raw data with external ID %d already exists.",
			requestID, externalID,
		)
	}

	ctx.GasMeter().ConsumeGas(
		k.GetParam(ctx, types.KeyGasPerRawDataRequestPerValidator)*uint64(len(request.RequestedValidators)),
		"RawDataRequest",
	)
	k.SetRawDataRequest(
		ctx, requestID, externalID,
		types.NewRawDataRequest(dataSourceID, calldata),
	)
	return k.ValidateDataSourceCount(ctx, requestID)
}

// GetRawDataRequestIterator is a function to get iterator on all raw data request that belong to
// given request id
func (k Keeper) GetRawDataRequestIterator(ctx sdk.Context, requestID types.RequestID) sdk.Iterator {
	prefix := types.GetIteratorPrefix(types.RawDataRequestStoreKeyPrefix, requestID)
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, prefix)
}

// GetRawDataRequestCount returns the number of raw data requests for the given request.
func (k Keeper) GetRawDataRequestCount(ctx sdk.Context, requestID types.RequestID) int64 {
	iterator := k.GetRawDataRequestIterator(ctx, requestID)
	count := 0
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return int64(count)
}

// GetRawDataRequests returns a list of raw data requests for the given request.
func (k Keeper) GetRawDataRequests(
	ctx sdk.Context, requestID types.RequestID,
) []types.RawDataRequest {
	iterator := k.GetRawDataRequestIterator(ctx, requestID)
	rawRequests := make([]types.RawDataRequest, 0)
	for ; iterator.Valid(); iterator.Next() {
		var rawRequest types.RawDataRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &rawRequest)
		rawRequests = append(rawRequests, rawRequest)
	}
	return rawRequests
}

// GetRawDataRequestWithExternalIDs returns a list of raw data requests bundled with external IDs
// for the given request.
func (k Keeper) GetRawDataRequestWithExternalIDs(
	ctx sdk.Context, requestID types.RequestID,
) []types.RawDataRequestWithExternalID {
	iterator := k.GetRawDataRequestIterator(ctx, requestID)
	rawRequests := make([]types.RawDataRequestWithExternalID, 0)
	for ; iterator.Valid(); iterator.Next() {
		var rawRequest types.RawDataRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &rawRequest)
		rawRequests = append(rawRequests,
			types.NewRawDataRequestWithExternalID(
				types.GetExternalIDFromRawDataRequestKey(iterator.Key()),
				rawRequest,
			),
		)
	}
	return rawRequests
}
