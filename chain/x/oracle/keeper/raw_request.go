package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasRawDataRequest checks if the raw data request of this ID exists in the storage.
func (k Keeper) HasRawDataRequest(ctx sdk.Context, rid types.RID, eid types.EID) bool {
	return ctx.KVStore(k.storeKey).Has(types.RawDataRequestStoreKey(rid, eid))
}

// SetRawDataRequest saves the raw data request to the store without performing validation.
func (k Keeper) SetRawDataRequest(ctx sdk.Context, rid types.RID, eid types.EID, data types.RawRequest) {
	ctx.KVStore(k.storeKey).Set(
		types.RawDataRequestStoreKey(rid, eid),
		k.cdc.MustMarshalBinaryBare(data),
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

// GetRawDataRequestCount returns the number of raw data requests for the given request.
func (k Keeper) GetRawDataRequestCount(ctx sdk.Context, rid types.RequestID) int64 {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store,
		types.GetIteratorPrefix(types.RawDataRequestStoreKeyPrefix, rid),
	)
	count := 0
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return int64(count)
}

// GetRawRequestsByRID returns all raw requests for the given request ID, or nil if there is none.
func (k Keeper) GetRawRequestsByRID(ctx sdk.Context, rid types.RID) (res []types.RawDataRequestWithExternalID) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store,
		types.GetIteratorPrefix(types.RawDataRequestStoreKeyPrefix, rid),
	)
	for ; iterator.Valid(); iterator.Next() {
		var rawRequest types.RawRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &rawRequest)
		res = append(res, types.NewRawDataRequestWithExternalID(
			types.GetExternalIDFromRawDataRequestKey(iterator.Key()),
			rawRequest,
		))
	}
	return res
}
