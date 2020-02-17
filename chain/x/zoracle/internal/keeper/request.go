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

// AddRequest attempts to create a new request. An error is returned if some conditions failed.
func (k Keeper) AddRequest(
	ctx sdk.Context, oracleScriptID int64, calldata []byte,
	requestedValidatorCount, sufficientValidatorCount, expiration int64,
) (int64, sdk.Error) {
	if !k.CheckOracleScriptExists(ctx, oracleScriptID) {
		// TODO: fix error later
		return 0, types.ErrRequestNotFound(types.DefaultCodespace)
	}

	if len(calldata) > int(k.MaxCalldataSize(ctx)) {
		// TODO: fix error later
		return 0, types.ErrRequestNotFound(types.DefaultCodespace)
	}

	validatorsByPower := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	if int64(len(validatorsByPower)) < requestedValidatorCount {
		// TODO: Fix error later
		return 0, types.ErrRequestNotFound(types.DefaultCodespace)
	}

	validators := make([]sdk.ValAddress, requestedValidatorCount)
	for i := int64(0); i < requestedValidatorCount; i++ {
		validators[i] = validatorsByPower[i].GetOperator()
	}

	requestID := k.GetNextRequestID(ctx)
	k.SetRequest(ctx, requestID, types.NewRequest(
		oracleScriptID,
		calldata,
		validators,
		sufficientValidatorCount,
		ctx.BlockHeight(),
		ctx.BlockTime().Unix(),
		ctx.BlockHeight()+expiration,
	))

	return requestID, nil
}

// ValidateDataSourceCount validates that the number of raw data requests is
// not greater than `MaxDataSourceCountPerRequest`
func (k Keeper) ValidateDataSourceCount(ctx sdk.Context, id int64) sdk.Error {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		return err
	}

	if request.DataSourceCount > k.MaxDataSourceCountPerRequest(ctx) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	return nil
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

// ShouldBecomePendingResolve checks and returns whether the given request should be moved to the
// pending resolve list, which will be resolved during the EndBlock call. The move will happen exactly when
// the request receives sufficient raw reports from the validators.
func (k Keeper) ShouldBecomePendingResolve(ctx sdk.Context, id int64) bool {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		return false
	}
	return int64(len(request.ReceivedValidators)) == request.SufficientValidatorCount
}

// AddPendingRequest checks and append new request id to list if id already existed in list, it will return error.
func (k Keeper) AddPendingRequest(ctx sdk.Context, requestID int64) sdk.Error {
	pendingList := k.GetPendingResolveList(ctx)
	for _, entry := range pendingList {
		if requestID == entry {
			return types.ErrDuplicateRequest(types.DefaultCodespace)
		}
	}
	pendingList = append(pendingList, requestID)
	k.SetPendingResolveList(ctx, pendingList)
	return nil
}

// SetPendingResolveList saves the list of pending request that will be resolved at end block.
func (k Keeper) SetPendingResolveList(ctx sdk.Context, reqIDs []int64) {
	store := ctx.KVStore(k.storeKey)
	encoded := k.cdc.MustMarshalBinaryBare(reqIDs)
	if encoded == nil {
		encoded = []byte{}
	}
	store.Set(types.PendingResolveListStoreKey, encoded)
}

// GetPendingResolveList returns the list of pending request.
func (k Keeper) GetPendingResolveList(ctx sdk.Context) []int64 {
	store := ctx.KVStore(k.storeKey)
	reqIDsBytes := store.Get(types.PendingResolveListStoreKey)

	// If the state is empty
	if len(reqIDsBytes) == 0 {
		return []int64{}
	}

	var reqIDs []int64
	k.cdc.MustUnmarshalBinaryBare(reqIDsBytes, &reqIDs)

	return reqIDs
}
