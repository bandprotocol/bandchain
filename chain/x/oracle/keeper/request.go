package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

// HasRequest checks if the request of this ID exists in the storage.
func (k Keeper) HasRequest(ctx sdk.Context, id types.RequestID) bool {
	return ctx.KVStore(k.storeKey).Has(types.RequestStoreKey(id))
}

// GetRequest returns the request struct for the given ID or error if not exists.
func (k Keeper) GetRequest(ctx sdk.Context, id types.RequestID) (types.Request, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.RequestStoreKey(id))
	if bz == nil {
		return types.Request{}, sdkerrors.Wrapf(types.ErrRequestNotFound, "id: %d", id)
	}
	var request types.Request
	k.cdc.MustUnmarshalBinaryBare(bz, &request)
	return request, nil
}

// MustGetRequest returns the request struct for the given ID. Panics error if not exists.
func (k Keeper) MustGetRequest(ctx sdk.Context, id types.RequestID) types.Request {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		panic(err)
	}
	return request
}

// SetRequest saves the given data request to the store without performing any validation.
func (k Keeper) SetRequest(ctx sdk.Context, id types.RequestID, request types.Request) {
	ctx.KVStore(k.storeKey).Set(types.RequestStoreKey(id), k.cdc.MustMarshalBinaryBare(request))
}

// DeleteRequest removes the given data request from the store.
func (k Keeper) DeleteRequest(ctx sdk.Context, id types.RequestID) {
	ctx.KVStore(k.storeKey).Delete(types.RequestStoreKey(id))
}

// AddRequest attempts to create and save a new request.
func (k Keeper) AddRequest(ctx sdk.Context, req types.Request) types.RequestID {
	id := k.GetNextRequestID(ctx)
	k.SetRequest(ctx, id, req)
	return id
}

// ProcessExpiredRequests resolves all expired requests and deactivates missed validators.
func (k Keeper) ProcessExpiredRequests(ctx sdk.Context) {
	currentReqID := k.GetRequestLastExpired(ctx) + 1
	lastReqID := types.RequestID(k.GetRequestCount(ctx))
	expirationBlockCount := int64(k.GetParam(ctx, types.KeyExpirationBlockCount))
	// Loop through all data requests in chronological order. If a request reaches its
	// expiration height, we will deactivate validators that didn't report data on the
	// request. We also resolve requests to status EXPIRED if they are not yet resolved.
	for ; currentReqID <= lastReqID; currentReqID++ {
		req := k.MustGetRequest(ctx, currentReqID)
		// This request is not yet expired, so there's nothing to do here. Ditto for
		// all other requests that come after this. Thus we can just break the loop.
		if req.RequestHeight+expirationBlockCount > ctx.BlockHeight() {
			break
		}
		// If the request still does not have result, we resolve it as EXPIRED.
		if !k.HasResult(ctx, currentReqID) {
			k.ResolveExpired(ctx, currentReqID)
		}
		// Deactivate all validators that do not report to this request.
		for _, val := range req.RequestedValidators {
			if !k.HasReport(ctx, currentReqID, val) {
				k.MissReport(ctx, val, req.RequestTime)
			}
		}
		// Set last expired request ID to be this current request.
		k.SetRequestLastExpired(ctx, currentReqID)
	}
}

// AddPendingRequest adds the request to the pending list. DO NOT add same request more than once.
func (k Keeper) AddPendingRequest(ctx sdk.Context, id types.RequestID) {
	pendingList := k.GetPendingResolveList(ctx)
	pendingList = append(pendingList, id)
	k.SetPendingResolveList(ctx, pendingList)
}

// SetPendingResolveList saves the list of pending request that will be resolved at end block.
func (k Keeper) SetPendingResolveList(ctx sdk.Context, ids []types.RequestID) {
	bz := k.cdc.MustMarshalBinaryBare(ids)
	if bz == nil {
		bz = []byte{}
	}
	ctx.KVStore(k.storeKey).Set(types.PendingResolveListStoreKey, bz)
}

// GetPendingResolveList returns the list of pending requests to be executed during EndBlock.
func (k Keeper) GetPendingResolveList(ctx sdk.Context) (ids []types.RequestID) {
	bz := ctx.KVStore(k.storeKey).Get(types.PendingResolveListStoreKey)
	if len(bz) == 0 { // Return an empty list if the key does not exist in the store.
		return []types.RequestID{}
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &ids)
	return ids
}
