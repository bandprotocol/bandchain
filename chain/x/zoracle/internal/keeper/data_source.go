package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
) (types.DataSourceID, error) {
	if uint64(len(executable)) > k.GetParam(ctx, types.KeyMaxDataSourceExecutableSize) {
		return 0, sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddDataSource: Executable size (%d) exceeds the maximum size (%d).",
			len(executable), k.GetParam(ctx, types.KeyMaxDataSourceExecutableSize),
		)
	}
	if uint64(len(name)) > k.GetParam(ctx, types.KeyMaxNameLength) {
		return 0, sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddDataSource: Name length (%d) exceeds the maximum length (%d).",
			len(name), k.GetParam(ctx, types.KeyMaxNameLength),
		)
	}
	if uint64(len(description)) > k.GetParam(ctx, types.KeyMaxDescriptionLength) {
		return 0, sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddDataSource: Description length (%d) exceeds the maximum length (%d).",
			len(description), k.GetParam(ctx, types.KeyMaxDescriptionLength),
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
) error {
	if !k.CheckDataSourceExists(ctx, dataSourceID) {
		return sdkerrors.Wrapf(types.ErrItemNotFound, "EditDataSource: Unknown data source ID %d.", dataSourceID)
	}

	if uint64(len(executable)) > k.GetParam(ctx, types.KeyMaxDataSourceExecutableSize) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"EditDataSource: Executable size (%d) exceeds the maximum size (%d).",
			len(executable), k.GetParam(ctx, types.KeyMaxDataSourceExecutableSize),
		)
	}
	if uint64(len(name)) > k.GetParam(ctx, types.KeyMaxNameLength) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"EditDataSource: Name length (%d) exceeds the maximum length (%d).",
			len(name), k.GetParam(ctx, types.KeyMaxNameLength),
		)
	}
	if uint64(len(description)) > k.GetParam(ctx, types.KeyMaxDescriptionLength) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"EditDataSource: Description length (%d) exceeds the maximum length (%d).",
			len(description), k.GetParam(ctx, types.KeyMaxDescriptionLength),
		)
	}

	updatedDataSource := types.NewDataSource(owner, name, description, fee, executable)
	k.SetDataSource(ctx, dataSourceID, updatedDataSource)
	return nil
}

// GetDataSource returns the entire DataSource struct for the given ID.
func (k Keeper) GetDataSource(
	ctx sdk.Context, id types.DataSourceID,
) (types.DataSource, error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckDataSourceExists(ctx, id) {
		return types.DataSource{}, sdkerrors.Wrapf(types.ErrItemNotFound,
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
