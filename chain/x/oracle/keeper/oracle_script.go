package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasOracleScript checks if the oracle script of this ID exists in the storage.
func (k Keeper) HasOracleScript(ctx sdk.Context, id types.OID) bool {
	return ctx.KVStore(k.storeKey).Has(types.OracleScriptStoreKey(id))
}

// GetOracleScript returns the oracle script struct for the given ID or error if not exists.
func (k Keeper) GetOracleScript(ctx sdk.Context, id types.OID) (types.OracleScript, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.OracleScriptStoreKey(id))
	if bz == nil {
		return types.OracleScript{}, sdkerrors.Wrapf(types.ErrOracleScriptNotFound, "id: %d", id)
	}
	var oracleScript types.OracleScript
	k.cdc.MustUnmarshalBinaryBare(bz, &oracleScript)
	return oracleScript, nil
}

// MustGetOracleScript returns the oracle script struct for the given ID. Panic if not exists.
func (k Keeper) MustGetOracleScript(ctx sdk.Context, id types.OID) types.OracleScript {
	oracleScript, err := k.GetOracleScript(ctx, id)
	if err != nil {
		panic(err)
	}
	return oracleScript
}

// SetOracleScript saves the given oracle script to the storage without performing validation.
func (k Keeper) SetOracleScript(ctx sdk.Context, id types.OID, oracleScript types.OracleScript) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OracleScriptStoreKey(id), k.cdc.MustMarshalBinaryBare(oracleScript))
}

// AddOracleScript adds the given oracle script to the storage. Returns error if validation fails.
func (k Keeper) AddOracleScript(ctx sdk.Context, oracleScript types.OracleScript) (types.OID, error) {
	if err := AnyError(
		k.EnsureLength(ctx, types.KeyMaxNameLength, len(oracleScript.Name)),
		k.EnsureLength(ctx, types.KeyMaxDescriptionLength, len(oracleScript.Description)),
		k.EnsureLength(ctx, types.KeyMaxOracleScriptCodeSize, len(oracleScript.Code)),
	); err != nil {
		return 0, err
	}
	id := k.GetNextOracleScriptID(ctx)
	k.SetOracleScript(ctx, id, oracleScript)
	return id, nil
}

// EditOracleScript edits the given oracle script by id and flushes it to the storage.
func (k Keeper) EditOracleScript(ctx sdk.Context, id types.OID, new types.OracleScript) error {
	oracleScript, err := k.GetOracleScript(ctx, id)
	if err != nil {
		return err
	}
	oracleScript.Owner = new.Owner // TODO: Allow NOT_MODIFY or nil in these fields.
	oracleScript.Name = new.Name
	oracleScript.Description = new.Description
	oracleScript.Code = new.Code
	oracleScript.Schema = new.Schema
	oracleScript.SourceCodeURL = new.SourceCodeURL
	if err := AnyError(
		k.EnsureLength(ctx, types.KeyMaxNameLength, len(oracleScript.Name)),
		k.EnsureLength(ctx, types.KeyMaxDescriptionLength, len(oracleScript.Description)),
		k.EnsureLength(ctx, types.KeyMaxOracleScriptCodeSize, len(oracleScript.Code)),
	); err != nil {
		return err
	}
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
