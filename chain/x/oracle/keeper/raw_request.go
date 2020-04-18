package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasRawRequest checks if the raw request of this ID tuple exists in the storage.
func (k Keeper) HasRawRequest(ctx sdk.Context, rid types.RID, eid types.EID) bool {
	return ctx.KVStore(k.storeKey).Has(types.RawRequestStoreKey(rid, eid))
}

// GetRawRequest returns the raw request struct for the given ID tuple or error if not exists.
func (k Keeper) GetRawRequest(ctx sdk.Context, rid types.RID, eid types.EID) (types.RawRequest, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.RawRequestStoreKey(rid, eid))
	if bz == nil {
		return types.RawRequest{}, sdkerrors.Wrapf(
			types.ErrRawRequestNotFound, "reqID: %d, extID: %d", rid, eid,
		)
	}
	var rawRequest types.RawRequest
	k.cdc.MustUnmarshalBinaryBare(bz, &rawRequest)
	return rawRequest, nil
}

// SetRawRequest saves the raw request to the storage without performing validation.
func (k Keeper) SetRawRequest(ctx sdk.Context, rid types.RID, eid types.EID, data types.RawRequest) {
	bz := k.cdc.MustMarshalBinaryBare(data)
	ctx.KVStore(k.storeKey).Set(types.RawRequestStoreKey(rid, eid), bz)
}

// AddRawRequest performs all sanity checks and adds a new raw request to the store.
func (k Keeper) AddRawRequest(ctx sdk.Context, rid types.RID, eid types.EID, did types.DID, calldata []byte) error {
	if err := k.EnsureMaxValue(ctx, types.KeyMaxCalldataSize, uint64(len(calldata))); err != nil {
		return err
	}
	if !k.HasRequest(ctx, rid) {
		return sdkerrors.Wrapf(types.ErrRequestNotFound, "id: %d", rid)
	}
	if !k.HasDataSource(ctx, did) {
		return sdkerrors.Wrapf(types.ErrDataSourceNotFound, "id: %d", did)
	}
	if k.HasRawRequest(ctx, rid, eid) {
		return sdkerrors.Wrapf(types.ErrRawRequestAlreadyExists, "reqID: %d, extID: %d", rid, eid)
	}
	// TODO: Make sure we consume gas for adding raw requests. That should be done in handler level.
	// TODO: Validate data source count. Also should be done in handler level.
	k.SetRawRequest(ctx, rid, eid, types.NewRawDataRequest(did, calldata))
	return nil
}

// GetRawRequestCount returns the number of raw requests for the given request ID.
func (k Keeper) GetRawRequestCount(ctx sdk.Context, rid types.RequestID) int64 {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store,
		types.GetIteratorPrefix(types.RawRequestStoreKeyPrefix, rid),
	)
	count := 0
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return int64(count)
}

// GetRawRequestsByRID returns all raw requests for the given request ID, or nil if there is none.
func (k Keeper) GetRawRequestsByRID(ctx sdk.Context, rid types.RID) (res []types.RawRequestWithExternalID) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store,
		types.GetIteratorPrefix(types.RawRequestStoreKeyPrefix, rid),
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
