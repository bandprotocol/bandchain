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

// SetRawRequest saves the raw request to the storage without performing validation.
func (k Keeper) SetRawRequest(ctx sdk.Context, rid types.RID, data types.RawRequest) {
	bz := k.cdc.MustMarshalBinaryBare(data)
	ctx.KVStore(k.storeKey).Set(types.RawRequestStoreKey(rid, data.ExternalID), bz)
}

// AddRawRequest performs all sanity checks and adds a new raw request to the store.
func (k Keeper) AddRawRequest(ctx sdk.Context, rid types.RID, data types.RawRequest) error {
	err := k.EnsureMax(ctx, types.KeyMaxCalldataSize, uint64(len(data.Calldata)))
	if err != nil {
		return err
	}
	if !k.HasRequest(ctx, rid) {
		return sdkerrors.Wrapf(types.ErrRequestNotFound, "id: %d", rid)
	}
	if !k.HasDataSource(ctx, data.DataSourceID) {
		return sdkerrors.Wrapf(types.ErrDataSourceNotFound, "id: %d", data.DataSourceID)
	}
	if k.HasRawRequest(ctx, rid, data.ExternalID) {
		return sdkerrors.Wrapf(types.ErrRawRequestAlreadyExists,
			"reqID: %d, extID: %d", rid, data.ExternalID,
		)
	}
	// TODO: Make sure we consume gas for adding raw requests. That should be done in handler level.
	// TODO: Validate data source count. Also should be done in handler level.
	k.SetRawRequest(ctx, rid, data)
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
func (k Keeper) GetRawRequestsByRID(ctx sdk.Context, rid types.RID) (res []types.RawRequest) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store,
		types.GetIteratorPrefix(types.RawRequestStoreKeyPrefix, rid),
	)
	for ; iterator.Valid(); iterator.Next() {
		var rawRequest types.RawRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &rawRequest)
		res = append(res, rawRequest)
	}
	return res
}
