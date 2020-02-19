package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetOracleScript saves the given oracle script with the given ID to the storage.
// WARNING: This function doesn't perform any check on ID.
func (k Keeper) SetOracleScript(ctx sdk.Context, id int64, oracleScript types.OracleScript) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OracleScriptStoreKey(id), k.cdc.MustMarshalBinaryBare(oracleScript))
}

// AddOracleScript adds the given oracle script to the storage.
func (k Keeper) AddOracleScript(ctx sdk.Context, owner sdk.AccAddress, name string, code []byte) sdk.Error {
	newOracleScriptID := k.GetNextOracleScriptID(ctx)

	if len(code) > int(k.MaxOracleScriptCodeSize(ctx)) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	newOracleScript := types.NewOracleScript(owner, name, code)
	k.SetOracleScript(ctx, newOracleScriptID, newOracleScript)
	return nil
}

// EditOracleScript edits the given oracle script by given oracle script id to the storage.
func (k Keeper) EditOracleScript(ctx sdk.Context, oracleScriptID int64, owner sdk.AccAddress, name string, code []byte) sdk.Error {
	if !k.CheckOracleScriptExists(ctx, oracleScriptID) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	if len(code) > int(k.MaxOracleScriptCodeSize(ctx)) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	updatedOracleScript := types.NewOracleScript(owner, name, code)
	k.SetOracleScript(ctx, oracleScriptID, updatedOracleScript)
	return nil
}

// GetOracleScript returns the entire OracleScript struct for the given ID.
func (k Keeper) GetOracleScript(ctx sdk.Context, id int64) (types.OracleScript, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckOracleScriptExists(ctx, id) {
		// TODO: fix error later
		return types.OracleScript{}, types.ErrRequestNotFound(types.DefaultCodespace)
	}

	bz := store.Get(types.OracleScriptStoreKey(id))
	var oracleScript types.OracleScript
	k.cdc.MustUnmarshalBinaryBare(bz, &oracleScript)
	return oracleScript, nil
}

// CheckOracleScriptExists checks if the oracle script of this ID exists in the storage.
func (k Keeper) CheckOracleScriptExists(ctx sdk.Context, id int64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.OracleScriptStoreKey(id))
}
