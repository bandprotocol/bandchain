package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
)

// AddResult is a function to validate result size and set to store.
func (k Keeper) AddResult(
	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID, calldata []byte, result []byte,
) sdk.Error {
	if int64(len(result)) > k.MaxResultSize(ctx) {
		return types.ErrBadDataValue(
			"AddResult: Result size (%d) exceeds the maximum size (%d).",
			len(result),
			int(k.MaxResultSize(ctx)),
		)
	}

	request, err := k.GetRequest(ctx, requestID)
	if err != nil {
		return err
	}

	k.SetResult(ctx, requestID, oracleScriptID, calldata, types.NewResult(
		request.RequestTime,
		ctx.BlockTime().Unix(),
		int64(len(request.RequestedValidators)),
		request.SufficientValidatorCount,
		int64(len(request.ReceivedValidators)),
		result,
	))

	return nil
}

// SetResult is a function to save result of execute code to store.
func (k Keeper) SetResult(
	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID, calldata []byte, result types.Result,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.ResultStoreKey(requestID, oracleScriptID, calldata),
		result.EncodeResult(),
	)
}

// GetResult returns the result bytes in store.
func (k Keeper) GetResult(
	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID, calldata []byte,
) (types.Result, sdk.Error) {
	if !k.HasResult(ctx, requestID, oracleScriptID, calldata) {
		return types.Result{}, types.ErrItemNotFound(
			"GetResult: Result for request ID %d is not available.",
			requestID,
		)
	}
	store := ctx.KVStore(k.storeKey)
	return types.MustDecodeResult(
		store.Get(types.ResultStoreKey(requestID, oracleScriptID, calldata)),
	), nil
}

// HasResult checks if the result at this request id is present in the store or not.
func (k Keeper) HasResult(ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID, calldata []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ResultStoreKey(requestID, oracleScriptID, calldata))
}
