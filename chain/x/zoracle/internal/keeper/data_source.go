package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) setDataSource(ctx sdk.Context, id int64, dataSource types.DataSource) {
	// TODO: check executable size.
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSourceStoreKey(id), k.cdc.MustMarshalBinaryBare(dataSource))
}

// AddDataSource adds the given data source to the storage.
func (k Keeper) AddDataSource(ctx sdk.Context, owner sdk.AccAddress, name string, fee sdk.Coins, executable []byte) sdk.Error {
	newDataSourceID := k.GetNextDataSourceID(ctx)

	newDataSource := types.NewDataSource(owner, name, fee, executable)
	k.setDataSource(ctx, newDataSourceID, newDataSource)
	return nil
}

// EditDataSource edits the given data source by given data source id to the storage.
func (k Keeper) EditDataSource(ctx sdk.Context, dataSourceID int64, owner sdk.AccAddress, name string, fee sdk.Coins, executable []byte) sdk.Error {
	if !k.CheckDataSourceExists(ctx, dataSourceID) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}

	updatedDataSource := types.NewDataSource(owner, name, fee, executable)
	k.setDataSource(ctx, dataSourceID, updatedDataSource)
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
