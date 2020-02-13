package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRequest is a function to save request to the given ID.
func (k Keeper) SetRequest(ctx sdk.Context, id int64, request types.Request) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RequestStoreKey(id), k.cdc.MustMarshalBinaryBare(request))
}

// GetRequest returns the entire Request metadata struct.
func (k Keeper) GetRequest(ctx sdk.Context, id int64) (types.Request, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckRequestExists(ctx, id) {
		return types.Request{}, types.ErrRequestNotFound(types.DefaultCodespace)
	}

	bz := store.Get(types.RequestStoreKey(id))
	var request types.Request
	k.cdc.MustUnmarshalBinaryBare(bz, &request)
	return request, nil
}

// AddNewReceiveValidator checks that new validator is a valid validator and not in received list yet then add new
// validator to list.
func (k Keeper) AddNewReceiveValidator(ctx sdk.Context, id int64, validator sdk.ValAddress) sdk.Error {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		return err
	}
	for _, submittedValidator := range request.ReceivedValidators {
		if validator.Equals(submittedValidator) {
			return types.ErrDuplicateValidator(types.DefaultCodespace)
		}
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
	request.ReceivedValidators = append(request.ReceivedValidators, validator)
	k.SetRequest(ctx, id, request)
	return nil
}

// SetResolve set resolve status and save to context.
func (k Keeper) SetResolve(ctx sdk.Context, id int64, isResolved bool) sdk.Error {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		return err
	}

	request.IsResolved = isResolved
	k.SetRequest(ctx, id, request)
	return nil
}

// CheckRequestExists checks if the request at this id is present in the store or not.
func (k Keeper) CheckRequestExists(ctx sdk.Context, id int64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.RequestStoreKey(id))
}

// AddPendingRequest checks and append new request id to list if id already existed in list, it will return error.
func (k Keeper) AddPendingRequest(ctx sdk.Context, requestID int64) sdk.Error {
	pendingList := k.GetPendingRequests(ctx)
	for _, entry := range pendingList {
		if requestID == entry {
			return types.ErrDuplicateRequest(types.DefaultCodespace)
		}
	}
	pendingList = append(pendingList, requestID)
	k.SetPendingRequests(ctx, pendingList)
	return nil
}

// SetPendingRequests saves the list of pending request that will be resolved at end block.
func (k Keeper) SetPendingRequests(ctx sdk.Context, reqIDs []int64) {
	store := ctx.KVStore(k.storeKey)
	encoded := k.cdc.MustMarshalBinaryBare(reqIDs)
	if encoded == nil {
		encoded = []byte{}
	}
	store.Set(types.UnresolvedRequestListStoreKey, encoded)
}

// GetPendingRequests returns the list of pending request.
func (k Keeper) GetPendingRequests(ctx sdk.Context) []int64 {
	store := ctx.KVStore(k.storeKey)
	reqIDsBytes := store.Get(types.UnresolvedRequestListStoreKey)

	// If the state is empty
	if len(reqIDsBytes) == 0 {
		return []int64{}
	}

	var reqIDs []int64
	k.cdc.MustUnmarshalBinaryBare(reqIDsBytes, &reqIDs)

	return reqIDs
}
