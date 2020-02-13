package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

// SetResult is a function to save result of execute code to store
func (k Keeper) SetResult(ctx sdk.Context, requestID int64, codeHash []byte, params []byte, result []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.ResultStoreKey(requestID, codeHash, params),
		result,
	)
}

func (k Keeper) GetResult(ctx sdk.Context, requestID int64, codeHash []byte, params []byte) ([]byte, sdk.Error) {
	if !k.HasResult(ctx, requestID, codeHash, params) {
		return nil, types.ErrResultNotFound(types.DefaultCodespace)
	}
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.ResultStoreKey(requestID, codeHash, params)), nil
}

func (k Keeper) HasResult(ctx sdk.Context, requestID int64, codeHash []byte, params []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ResultStoreKey(requestID, codeHash, params))
}
