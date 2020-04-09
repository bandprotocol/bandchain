package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SetOracleScript saves the given oracle script to the storage without performing validation.
func (k Keeper) SetOracleScript(
	ctx sdk.Context, id types.OracleScriptID, oracleScript types.OracleScript,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OracleScriptStoreKey(id), k.cdc.MustMarshalBinaryBare(oracleScript))
}

// AddOracleScript adds the given oracle script to the storage.
func (k Keeper) AddOracleScript(
	ctx sdk.Context, owner sdk.AccAddress, name string, description string, code []byte,
) (types.OracleScriptID, error) {
	if uint64(len(code)) > k.GetParam(ctx, types.KeyMaxOracleScriptCodeSize) {
		return 0, sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddOracleScript: Code size (%d) exceeds the maximum size (%d).",
			len(code), k.GetParam(ctx, types.KeyMaxOracleScriptCodeSize),
		)
	}
	if uint64(len(name)) > k.GetParam(ctx, types.KeyMaxNameLength) {
		return 0, sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddOracleScript: Name length (%d) exceeds the maximum length (%d). 211",
			len(name), k.GetParam(ctx, types.KeyMaxNameLength),
		)
	}
	if uint64(len(description)) > k.GetParam(ctx, types.KeyMaxDescriptionLength) {
		return 0, sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddOracleScript: Name length (%d) exceeds the maximum length (%d).",
			len(name), k.GetParam(ctx, types.KeyMaxNameLength),
		)
	}

	newOracleScriptID := k.GetNextOracleScriptID(ctx)
	newOracleScript := types.NewOracleScript(owner, name, description, code)
	k.SetOracleScript(ctx, newOracleScriptID, newOracleScript)
	return newOracleScriptID, nil
}

// EditOracleScript edits the given oracle script by given oracle script id to the storage.
func (k Keeper) EditOracleScript(
	ctx sdk.Context, oracleScriptID types.OracleScriptID, owner sdk.AccAddress,
	name string, description string, code []byte,
) error {
	if !k.CheckOracleScriptExists(ctx, oracleScriptID) {
		return sdkerrors.Wrapf(types.ErrItemNotFound,
			"EditOracleScript: Unknown oracle script ID %d.", oracleScriptID,
		)
	}
	if uint64(len(code)) > k.GetParam(ctx, types.KeyMaxOracleScriptCodeSize) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"EditDataSource: Code size (%d) exceeds the maximum size (%d).",
			len(code), k.GetParam(ctx, types.KeyMaxOracleScriptCodeSize),
		)
	}
	if uint64(len(name)) > k.GetParam(ctx, types.KeyMaxNameLength) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"EditOracleScript: Name length (%d) exceeds the maximum length (%d).",
			len(name), k.GetParam(ctx, types.KeyMaxNameLength),
		)
	}
	if uint64(len(description)) > k.GetParam(ctx, types.KeyMaxDescriptionLength) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"EditDataSource: Description length (%d) exceeds the maximum length (%d).",
			len(description), k.GetParam(ctx, types.KeyMaxDescriptionLength),
		)
	}

	updatedOracleScript := types.NewOracleScript(owner, name, description, code)
	k.SetOracleScript(ctx, oracleScriptID, updatedOracleScript)
	return nil
}

// GetOracleScript returns the entire OracleScript struct for the given ID.
func (k Keeper) GetOracleScript(
	ctx sdk.Context, id types.OracleScriptID,
) (types.OracleScript, error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckOracleScriptExists(ctx, id) {
		return types.OracleScript{}, sdkerrors.Wrapf(types.ErrItemNotFound,
			"GetOracleScript: Unknown oracle script ID %d.", id,
		)
	}

	var oracleScript types.OracleScript
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.OracleScriptStoreKey(id)), &oracleScript)
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
