package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

// SetResult is a function to save result of execute code to store.
func (k Keeper) SetResult(ctx sdk.Context, requestID int64, result []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.ResultStoreKey(requestID),
		result,
	)
}

// GetResult returns the result bytes in store.
func (k Keeper) GetResult(ctx sdk.Context, requestID int64) ([]byte, sdk.Error) {
	if !k.HasResult(ctx, requestID) {
		return nil, types.ErrResultNotFound(types.DefaultCodespace)
	}
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.ResultStoreKey(requestID)), nil
}

// HasResult checks if the result at this request id is present in the store or not.
func (k Keeper) HasResult(ctx sdk.Context, requestID int64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ResultStoreKey(requestID))
}
