package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDataSource saves the given data source to the storage without performing validation.
func (k Keeper) SetDataSource(
	ctx sdk.Context, id types.DataSourceID, dataSource types.DataSource,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSourceStoreKey(id), k.cdc.MustMarshalBinaryBare(dataSource))
}

// AddDataSource adds the given data source to the storage.
func (k Keeper) AddDataSource(
	ctx sdk.Context, owner sdk.AccAddress, name string, description string,
	fee sdk.Coins, executable []byte,
) (types.DataSourceID, sdk.Error) {
	if int64(len(executable)) > k.MaxDataSourceExecutableSize(ctx) {
		return 0, types.ErrBadDataValue(
			"AddDataSource: Executable size (%d) exceeds the maximum size (%d).",
			len(executable), int(k.MaxDataSourceExecutableSize(ctx)),
		)
	}
	if int64(len(name)) > k.MaxNameLength(ctx) {
		return 0, types.ErrBadDataValue(
			"AddDataSource: Name length (%d) exceeds the maximum length (%d).",
			len(name), int(k.MaxNameLength(ctx)),
		)
	}
	if int64(len(description)) > k.MaxDescriptionLength(ctx) {
		return 0, types.ErrBadDataValue(
			"AddDataSource: Description length (%d) exceeds the maximum length (%d).",
			len(description), int(k.MaxDescriptionLength(ctx)),
		)
	}

	newDataSourceID := k.GetNextDataSourceID(ctx)
	newDataSource := types.NewDataSource(owner, name, description, fee, executable)
	k.SetDataSource(ctx, newDataSourceID, newDataSource)
	return newDataSourceID, nil
}

// EditDataSource edits the given data source by given data source id to the storage.
func (k Keeper) EditDataSource(
	ctx sdk.Context, dataSourceID types.DataSourceID, owner sdk.AccAddress, name string,
	description string, fee sdk.Coins, executable []byte,
) sdk.Error {
	if !k.CheckDataSourceExists(ctx, dataSourceID) {
		return types.ErrItemNotFound("EditDataSource: Unknown data source ID %d.", dataSourceID)
	}

	if int64(len(executable)) > k.MaxDataSourceExecutableSize(ctx) {
		return types.ErrBadDataValue(
			"EditDataSource: Executable size (%d) exceeds the maximum size (%d).",
			len(executable), int(k.MaxDataSourceExecutableSize(ctx)),
		)
	}
	if int64(len(name)) > k.MaxNameLength(ctx) {
		return types.ErrBadDataValue(
			"EditDataSource: Name length (%d) exceeds the maximum length (%d).",
			len(name), int(k.MaxNameLength(ctx)),
		)
	}
	if int64(len(description)) > k.MaxDescriptionLength(ctx) {
		return types.ErrBadDataValue(
			"EditDataSource: Description length (%d) exceeds the maximum length (%d).",
			len(description), int(k.MaxDescriptionLength(ctx)),
		)
	}

	updatedDataSource := types.NewDataSource(owner, name, description, fee, executable)
	k.SetDataSource(ctx, dataSourceID, updatedDataSource)
	return nil
}

// GetDataSource returns the entire DataSource struct for the given ID.
func (k Keeper) GetDataSource(
	ctx sdk.Context, id types.DataSourceID,
) (types.DataSource, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckDataSourceExists(ctx, id) {
		return types.DataSource{}, types.ErrItemNotFound(
			"GetDataSource: Unknown data source ID %d.", id,
		)
	}

	var dataSource types.DataSource
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.DataSourceStoreKey(id)), &dataSource)
	return dataSource, nil
}

// CheckDataSourceExists checks if the data source of this ID exists in the storage.
func (k Keeper) CheckDataSourceExists(ctx sdk.Context, id types.DataSourceID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.DataSourceStoreKey(id))
}

// GetDataSourceIterator returns an iterator for all data sources in the store.
func (k Keeper) GetDataSourceIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.DataSourceStoreKeyPrefix)
}

// GetAllDataSources returns list of all data sources.
func (k Keeper) GetAllDataSources(ctx sdk.Context) []types.DataSource {
	var dataSource types.DataSource
	dataSources := []types.DataSource{}
	iterator := k.GetDataSourceIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &dataSource)
		dataSources = append(dataSources, dataSource)
	}
	return dataSources
}
