package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

// SetCode is a function to save codeHash as key and code as value
func (k Keeper) SetCode(ctx sdk.Context, code []byte) []byte {
	store := ctx.KVStore(k.storeKey)
	codeHash := crypto.Keccak256(code)
	key := types.CodeHashStoreKey(codeHash)
	store.Set(key, code)
	return codeHash
}

// GetCode is a function to get code by using codeHash
func (k Keeper) GetCode(ctx sdk.Context, codeHash []byte) ([]byte, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	key := types.CodeHashStoreKey(codeHash)
	if !k.CheckCodeHashExists(ctx, codeHash) {
		return []byte{}, types.ErrCodeHashNotFound(types.DefaultCodespace)
	}
	code := store.Get(key)
	return code, nil
}

// CheckCodeHashExists checks if the code at this codeHash is valid in the store or not
func (k Keeper) CheckCodeHashExists(ctx sdk.Context, codeHash []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.CodeHashStoreKey(codeHash))
}
