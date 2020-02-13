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

// AddSubmitValidator checks that new validator is a valid validator and not in submitted list yet then add new
// validator to list.
func (k Keeper) AddSubmitValidator(ctx sdk.Context, id int64, validator sdk.ValAddress) sdk.Error {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		return err
	}
	for _, submittedValidator := range request.SubmittedValidatorList {
		if validator.Equals(submittedValidator) {
			return types.ErrDuplicateValidator(types.DefaultCodespace)
		}
	}
	found := false
	for _, validValidator := range request.Validators {
		if validator.Equals(validValidator) {
			found = true
			break
		}
	}

	if !found {
		return types.ErrInvalidValidator(types.DefaultCodespace)
	}
	request.SubmittedValidatorList = append(request.SubmittedValidatorList, validator)
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

// uniqueReqIDs is used to create array with all elements being unique (deduplicated).
func uniqueReqIDs(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// SetUnresolvedRequests saves the list of request in pending period.
func (k Keeper) SetUnresolvedRequests(ctx sdk.Context, reqIDs []int64) {
	store := ctx.KVStore(k.storeKey)
	urIDs := uniqueReqIDs(reqIDs)
	encoded := k.cdc.MustMarshalBinaryBare(urIDs)
	if encoded == nil {
		encoded = []byte{}
	}
	store.Set(types.UnresolvedRequestListStoreKey, encoded)
}

// GetUnresolvedRequests returns the list of request IDs in pending period.
func (k Keeper) GetUnresolvedRequests(ctx sdk.Context) []int64 {
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
