package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

// AddResult is a function to validate result size and set to store.
func (k Keeper) AddResult(
	ctx sdk.Context, requestID int64, oracleScriptID int64, calldata []byte, result []byte,
) sdk.Error {
	if int64(len(result)) > k.MaxResultSize(ctx) {
		// TODO: better error later
		return types.ErrResultNotFound(types.DefaultCodespace)
	}
	k.SetResult(ctx, requestID, oracleScriptID, calldata, result)
	return nil
}

// SetResult is a function to save result of execute code to store.
func (k Keeper) SetResult(
	ctx sdk.Context, requestID int64, oracleScriptID int64, calldata []byte, result []byte,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.ResultStoreKey(requestID, oracleScriptID, calldata),
		result,
	)
}

// GetResult returns the result bytes in store.
func (k Keeper) GetResult(
	ctx sdk.Context, requestID int64, oracleScriptID int64, calldata []byte,
) ([]byte, sdk.Error) {
	if !k.HasResult(ctx, requestID, oracleScriptID, calldata) {
		return nil, types.ErrResultNotFound(types.DefaultCodespace)
	}
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.ResultStoreKey(requestID, oracleScriptID, calldata)), nil
}

// HasResult checks if the result at this request id is present in the store or not.
func (k Keeper) HasResult(ctx sdk.Context, requestID int64, oracleScriptID int64, calldata []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ResultStoreKey(requestID, oracleScriptID, calldata))
}
