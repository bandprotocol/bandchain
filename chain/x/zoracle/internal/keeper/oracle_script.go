package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetOracleScript saves the given oracle script with the given ID to the storage.
// WARNING: This function doesn't perform any check on ID.
func (k Keeper) SetOracleScript(ctx sdk.Context, id types.OracleScriptID, oracleScript types.OracleScript) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OracleScriptStoreKey(id), k.cdc.MustMarshalBinaryBare(oracleScript))
}

// AddOracleScript adds the given oracle script to the storage.
func (k Keeper) AddOracleScript(ctx sdk.Context, owner sdk.AccAddress, name string, description string, code []byte) sdk.Error {
	newOracleScriptID := k.GetNextOracleScriptID(ctx)

	if len(code) > int(k.MaxOracleScriptCodeSize(ctx)) {
		return types.ErrBadDataValue(
			"AddOracleScript: Code size (%d) exceeds the maximum size (%d).",
			len(code),
			int(k.MaxOracleScriptCodeSize(ctx)),
		)
	}
	if len(name) > int(k.MaxNameLength(ctx)) {
		return types.ErrBadDataValue(
			"AddOracleScript: Name length (%d) exceeds the maximum length (%d).",
			len(name),
			int(k.MaxNameLength(ctx)),
		)
	}
	if len(description) > int(k.MaxDescriptionLength(ctx)) {
		return types.ErrBadDataValue(
			"AddOracleScript: Name length (%d) exceeds the maximum length (%d).",
			len(name),
			int(k.MaxNameLength(ctx)),
		)
	}

	newOracleScript := types.NewOracleScript(owner, name, description, code)
	k.SetOracleScript(ctx, newOracleScriptID, newOracleScript)
	return nil
}

// EditOracleScript edits the given oracle script by given oracle script id to the storage.
func (k Keeper) EditOracleScript(ctx sdk.Context, oracleScriptID types.OracleScriptID, owner sdk.AccAddress, name string, description string, code []byte) sdk.Error {
	if !k.CheckOracleScriptExists(ctx, oracleScriptID) {
		return types.ErrItemNotFound(
			"EditOracleScript: Unknown oracle script ID %d.",
			oracleScriptID,
		)
	}

	if len(code) > int(k.MaxOracleScriptCodeSize(ctx)) {
		return types.ErrBadDataValue(
			"EditDataSource: Code size (%d) exceeds the maximum size (%d).",
			len(code),
			int(k.MaxOracleScriptCodeSize(ctx)),
		)
	}
	if len(name) > int(k.MaxNameLength(ctx)) {
		return types.ErrBadDataValue(
			"EditOracleScript: Name length (%d) exceeds the maximum length (%d).",
			len(name),
			int(k.MaxNameLength(ctx)),
		)
	}
	if len(description) > int(k.MaxDescriptionLength(ctx)) {
		return types.ErrBadDataValue(
			"EditDataSource: Description length (%d) exceeds the maximum length (%d).",
			len(description),
			int(k.MaxDescriptionLength(ctx)),
		)
	}

	updatedOracleScript := types.NewOracleScript(owner, name, description, code)
	k.SetOracleScript(ctx, oracleScriptID, updatedOracleScript)
	return nil
}

// GetOracleScript returns the entire OracleScript struct for the given ID.
func (k Keeper) GetOracleScript(ctx sdk.Context, id types.OracleScriptID) (types.OracleScript, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckOracleScriptExists(ctx, id) {
		return types.OracleScript{}, types.ErrItemNotFound(
			"GetOracleScript: Unknown oracle script ID %d.",
			id,
		)
	}

	bz := store.Get(types.OracleScriptStoreKey(id))
	var oracleScript types.OracleScript
	k.cdc.MustUnmarshalBinaryBare(bz, &oracleScript)
	return oracleScript, nil
}

// CheckOracleScriptExists checks if the oracle script of this ID exists in the storage.
func (k Keeper) CheckOracleScriptExists(ctx sdk.Context, id types.OracleScriptID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.OracleScriptStoreKey(id))
}

// GetOracleScriptIterator returns an iterator for all oracle scripts in the store.
func (k Keeper) GetOracleScriptIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.OracleScriptStoreKeyPrefix)
}

// GetAllOracleScripts returns list of all oracle scripts.
func (k Keeper) GetAllOracleScripts(ctx sdk.Context) []types.OracleScript {
	var oracleScript types.OracleScript
	oracleScripts := []types.OracleScript{}
	iterator := k.GetOracleScriptIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &oracleScript)
		oracleScripts = append(oracleScripts, oracleScript)
	}
	return oracleScripts
}
