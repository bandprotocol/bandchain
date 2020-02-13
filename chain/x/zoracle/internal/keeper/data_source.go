package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddDataSource adds the given data source to the storage.
func (k Keeper) AddDataSource(ctx sdk.Context, owner sdk.AccAddress, name string, fee sdk.Coins, executable []byte) sdk.Error {
	newDataSourceID := k.GetNextDataSourceID(ctx)

	// TODO: check executable size.
	newDataSource := types.NewDataSource(owner, name, fee, executable)

	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSourceStoreKey(newDataSourceID), k.cdc.MustMarshalBinaryBare(newDataSource))
	return nil
}

// EditDataSource edits the given data source by given data source id to the storage.
func (k Keeper) EditDataSource(ctx sdk.Context, owner sdk.AccAddress, dataSourceID int64, name string, fee sdk.Coins, executable []byte) sdk.Error {
	// TODO: check executable size.
	updatedDataSource := types.NewDataSource(owner, name, fee, executable)

	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSourceStoreKey(dataSourceID), k.cdc.MustMarshalBinaryBare(updatedDataSource))
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
