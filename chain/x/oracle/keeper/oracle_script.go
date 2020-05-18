package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasOracleScript checks if the oracle script of this ID exists in the storage.
func (k Keeper) HasOracleScript(ctx sdk.Context, id types.OracleScriptID) bool {
	return ctx.KVStore(k.storeKey).Has(types.OracleScriptStoreKey(id))
}

// GetOracleScript returns the oracle script struct for the given ID or error if not exists.
func (k Keeper) GetOracleScript(ctx sdk.Context, id types.OracleScriptID) (types.OracleScript, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.OracleScriptStoreKey(id))
	if bz == nil {
		return types.OracleScript{}, sdkerrors.Wrapf(types.ErrOracleScriptNotFound, "id: %d", id)
	}
	var oracleScript types.OracleScript
	k.cdc.MustUnmarshalBinaryBare(bz, &oracleScript)
	return oracleScript, nil
}

// MustGetOracleScript returns the oracle script struct for the given ID. Panic if not exists.
func (k Keeper) MustGetOracleScript(ctx sdk.Context, id types.OracleScriptID) types.OracleScript {
	oracleScript, err := k.GetOracleScript(ctx, id)
	if err != nil {
		panic(err)
	}
	return oracleScript
}

// SetOracleScript saves the given oracle script to the storage without performing validation.
func (k Keeper) SetOracleScript(ctx sdk.Context, id types.OracleScriptID, oracleScript types.OracleScript) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OracleScriptStoreKey(id), k.cdc.MustMarshalBinaryBare(oracleScript))
}

// AddOracleScript adds the given oracle script to the storage. Returns error if validation fails.
func (k Keeper) AddOracleScript(ctx sdk.Context, oracleScript types.OracleScript) (types.OracleScriptID, error) {
	id := k.GetNextOracleScriptID(ctx)
	k.SetOracleScript(ctx, id, oracleScript)
	return id, nil
}

// EditOracleScript edits the given oracle script by id and flushes it to the storage.
func (k Keeper) EditOracleScript(ctx sdk.Context, id types.OracleScriptID, new types.OracleScript) error {
	oracleScript, err := k.GetOracleScript(ctx, id)
	if err != nil {
		return err
	}
	oracleScript.Owner = new.Owner
	oracleScript.Name = modify(oracleScript.Name, new.Name)
	oracleScript.Description = modify(oracleScript.Description, new.Description)
	oracleScript.Code = new.Code // TODO: Revisit this after file cache is done.
	oracleScript.Schema = modify(oracleScript.Schema, new.Schema)
	oracleScript.SourceCodeURL = modify(oracleScript.SourceCodeURL, new.SourceCodeURL)
	k.SetOracleScript(ctx, id, oracleScript)
	return nil
}

// GetAllOracleScripts returns the list of all oracle scripts in the store, or nil if there is none.
func (k Keeper) GetAllOracleScripts(ctx sdk.Context) (oracleScripts []types.OracleScript) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleScriptStoreKeyPrefix)
	for ; iterator.Valid(); iterator.Next() {
		var oracleScript types.OracleScript
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &oracleScript)
		oracleScripts = append(oracleScripts, oracleScript)
	}
	return oracleScripts
}
