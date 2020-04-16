package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasOracleScript checks if the oracle script of this ID exists in the storage.
func (k Keeper) HasOracleScript(ctx sdk.Context, id types.OID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.OracleScriptStoreKey(id))
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
func (k Keeper) AddOracleScript(
	ctx sdk.Context, owner sdk.AccAddress, name string, description string, code []byte,
) (types.OID, error) {
	if err := AnyError(
		k.EnsureMaxValue(ctx, types.KeyMaxNameLength, uint64(len(name))),
		k.EnsureMaxValue(ctx, types.KeyMaxDescriptionLength, uint64(len(description))),
		k.EnsureMaxValue(ctx, types.KeyMaxOracleScriptCodeSize, uint64(len(code))),
	); err != nil {
		return 0, err
	}
	id := k.GetNextOracleScriptID(ctx)
	k.SetOracleScript(ctx, id, types.NewOracleScript(owner, name, description, code))
	return id, nil
}

// EditOracleScript edits the given oracle script by id and flushes it to the storage.
func (k Keeper) EditOracleScript(
	ctx sdk.Context, id types.OID, owner sdk.AccAddress, name string,
	description string, code []byte,
) error {
	oracleScript, err := k.GetOracleScript(ctx, id)
	if err != nil {
		return err
	}
	oracleScript.Owner = owner // TODO: Allow NOT_MODIFY or nil in these fields.
	oracleScript.Name = name
	oracleScript.Description = description
	oracleScript.Code = code
	if err := AnyError(
		k.EnsureMaxValue(ctx, types.KeyMaxNameLength, uint64(len(oracleScript.Name))),
		k.EnsureMaxValue(ctx, types.KeyMaxDescriptionLength, uint64(len(oracleScript.Description))),
		k.EnsureMaxValue(ctx, types.KeyMaxOracleScriptCodeSize, uint64(len(oracleScript.Code))),
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
