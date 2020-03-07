package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDataSource saves the given data source with the given ID to the storage.
// WARNING: This function doesn't perform any check on ID.
func (k Keeper) SetDataSource(ctx sdk.Context, id int64, dataSource types.DataSource) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSourceStoreKey(id), k.cdc.MustMarshalBinaryBare(dataSource))
}

// AddDataSource adds the given data source to the storage.
func (k Keeper) AddDataSource(ctx sdk.Context, owner sdk.AccAddress, name string, description string, fee sdk.Coins, executable []byte) sdk.Error {
	newDataSourceID := k.GetNextDataSourceID(ctx)

	if len(executable) > int(k.MaxDataSourceExecutableSize(ctx)) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	if len(name) > int(k.MaxNameLength(ctx)) {
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	if len(description) > int(k.MaxDescriptionLength(ctx)) {
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	newDataSource := types.NewDataSource(owner, name, description, fee, executable)
	k.SetDataSource(ctx, newDataSourceID, newDataSource)
	return nil
}

// EditDataSource edits the given data source by given data source id to the storage.
func (k Keeper) EditDataSource(ctx sdk.Context, dataSourceID int64, owner sdk.AccAddress, name string, description string, fee sdk.Coins, executable []byte) sdk.Error {
	if !k.CheckDataSourceExists(ctx, dataSourceID) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	if len(executable) > int(k.MaxDataSourceExecutableSize(ctx)) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}
	if len(name) > int(k.MaxNameLength(ctx)) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}
	if len(description) > int(k.MaxDescriptionLength(ctx)) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	updatedDataSource := types.NewDataSource(owner, name, description, fee, executable)
	k.SetDataSource(ctx, dataSourceID, updatedDataSource)
	return nil
}

// GetDataSource returns the entire DataSource struct for the given ID.
func (k Keeper) GetDataSource(ctx sdk.Context, id int64) (types.DataSource, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckDataSourceExists(ctx, id) {
		// TODO: fix error later
		return types.DataSource{}, types.ErrRequestNotFound(types.DefaultCodespace)
	}

	bz := store.Get(types.DataSourceStoreKey(id))
	var dataSource types.DataSource
	k.cdc.MustUnmarshalBinaryBare(bz, &dataSource)
	return dataSource, nil
}

// CheckDataSourceExists checks if the data source of this ID exists in the storage.
func (k Keeper) CheckDataSourceExists(ctx sdk.Context, id int64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.DataSourceStoreKey(id))
}
