package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

// HasDataSource checks if the data source of this ID exists in the storage.
func (k Keeper) HasDataSource(ctx sdk.Context, id types.DataSourceID) bool {
	return ctx.KVStore(k.storeKey).Has(types.DataSourceStoreKey(id))
}

// GetDataSource returns the data source struct for the given ID or error if not exists.
func (k Keeper) GetDataSource(ctx sdk.Context, id types.DataSourceID) (types.DataSource, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.DataSourceStoreKey(id))
	if bz == nil {
		return types.DataSource{}, sdkerrors.Wrapf(types.ErrDataSourceNotFound, "id: %d", id)
	}
	var dataSource types.DataSource
	k.cdc.MustUnmarshalBinaryBare(bz, &dataSource)
	return dataSource, nil
}

// MustGetDataSource returns the data source struct for the given ID. Panic if not exists.
func (k Keeper) MustGetDataSource(ctx sdk.Context, id types.DataSourceID) types.DataSource {
	dataSource, err := k.GetDataSource(ctx, id)
	if err != nil {
		panic(err)
	}
	return dataSource
}

// SetDataSource saves the given data source to the storage without performing validation.
func (k Keeper) SetDataSource(ctx sdk.Context, id types.DataSourceID, dataSource types.DataSource) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSourceStoreKey(id), k.cdc.MustMarshalBinaryBare(dataSource))
}

// AddDataSource adds the given data source to the storage.
func (k Keeper) AddDataSource(ctx sdk.Context, dataSource types.DataSource) types.DataSourceID {
	id := k.GetNextDataSourceID(ctx)
	k.SetDataSource(ctx, id, dataSource)
	return id
}

// MustEditDataSource edits the given data source by id and flushes it to the storage.
func (k Keeper) MustEditDataSource(ctx sdk.Context, id types.DataSourceID, new types.DataSource) {
	dataSource := k.MustGetDataSource(ctx, id)
	dataSource.Owner = new.Owner
	dataSource.Name = modify(dataSource.Name, new.Name)
	dataSource.Description = modify(dataSource.Description, new.Description)
	dataSource.Filename = modify(dataSource.Filename, new.Filename)
	k.SetDataSource(ctx, id, dataSource)
}

// GetAllDataSources returns the list of all data sources in the store, or nil if there is none.
func (k Keeper) GetAllDataSources(ctx sdk.Context) (dataSources []types.DataSource) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DataSourceStoreKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dataSource types.DataSource
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &dataSource)
		dataSources = append(dataSources, dataSource)
	}
	return dataSources
}

// AddExecutableFile saves the given executable file to a file to filecahe storage and returns
// its sha256sum reference name. Returns do-not-modify symbol if the input is do-not-modify.
func (k Keeper) AddExecutableFile(file []byte) string {
	if bytes.Equal(file, types.DoNotModifyBytes) {
		return types.DoNotModify
	}
	return k.fileCache.AddFile(file)
}
