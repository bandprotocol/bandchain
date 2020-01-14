package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

// SetCode is a function to save codeHash as key and code as value.
func (k Keeper) SetCode(ctx sdk.Context, code []byte, owner sdk.AccAddress) []byte {
	store := ctx.KVStore(k.storeKey)
	sc := types.NewStoredCode(code, owner)
	codeHash := sc.GetCodeHash()
	key := types.CodeHashStoreKey(codeHash)
	store.Set(key, k.cdc.MustMarshalBinaryBare(sc))
	return codeHash
}

// GetCode is a function to get strored code struct by using codeHash.
func (k Keeper) GetCode(ctx sdk.Context, codeHash []byte) (types.StoredCode, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	key := types.CodeHashStoreKey(codeHash)
	if !k.CheckCodeHashExists(ctx, codeHash) {
		return types.StoredCode{}, types.ErrCodeHashNotFound(types.DefaultCodespace)
	}
	bz := store.Get(key)
	var storedCode types.StoredCode
	k.cdc.MustUnmarshalBinaryBare(bz, &storedCode)
	return storedCode, nil
}

// CheckCodeHashExists checks if the code at this codeHash is valid in the store or not.
func (k Keeper) CheckCodeHashExists(ctx sdk.Context, codeHash []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.CodeHashStoreKey(codeHash))
}

// DeleteCode deletes code from storage
func (k Keeper) DeleteCode(ctx sdk.Context, codeHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.CodeHashStoreKey(codeHash))
}

// GetCodesIterator returns an iterator for all codes in chain.
func (k Keeper) GetCodesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.CodeHashKeyPrefix)
}
