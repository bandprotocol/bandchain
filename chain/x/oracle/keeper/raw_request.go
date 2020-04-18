package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasRawDataRequest checks if the raw data request of this ID exists in the storage.
func (k Keeper) HasRawDataRequest(ctx sdk.Context, requestID types.RID, externalID types.EID) bool {
	return ctx.KVStore(k.storeKey).Has(types.RawDataRequestStoreKey(requestID, externalID))
}

// SetRawDataRequest saves the raw data request to the store without performing validation.
func (k Keeper) SetRawDataRequest(ctx sdk.Context, requestID types.RID, externalID types.EID, rawDataRequest types.RawDataRequest) {
	ctx.KVStore(k.storeKey).Set(
		types.RawDataRequestStoreKey(requestID, externalID),
		k.cdc.MustMarshalBinaryBare(rawDataRequest),
	)
}

// GetRawDataRequest returns the raw data request struct or error if not exists.
func (k Keeper) GetRawDataRequest(ctx sdk.Context, requestID types.RID, externalID types.EID) (types.RawDataRequest, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.RawDataRequestStoreKey(requestID, externalID))
	if bz == nil {
		return types.RawDataRequest{}, sdkerrors.Wrapf(
			types.ErrRawRequestNotFound, "reqID: %d, extID: %d", requestID, externalID,
		)
	}
	var rawRequest types.RawDataRequest
	k.cdc.MustUnmarshalBinaryBare(bz, &rawRequest)
	return rawRequest, nil
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

	if k.HasRawDataRequest(ctx, requestID, externalID) {
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
