package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasDataSource checks if the data source of this ID exists in the storage.
func (k Keeper) HasDataSource(ctx sdk.Context, id types.DID) bool {
	return ctx.KVStore(k.storeKey).Has(types.DataSourceStoreKey(id))
}

// GetDataSource returns the data source struct for the given ID or error if not exists.
func (k Keeper) GetDataSource(ctx sdk.Context, id types.DID) (types.DataSource, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.DataSourceStoreKey(id))
	if bz == nil {
		return types.DataSource{}, sdkerrors.Wrapf(types.ErrDataSourceNotFound, "id: %d", id)
	}
	var dataSource types.DataSource
	k.cdc.MustUnmarshalBinaryBare(bz, &dataSource)
	return dataSource, nil
}

// MustGetDataSource returns the data source struct for the given ID. Panic if not exists.
func (k Keeper) MustGetDataSource(ctx sdk.Context, id types.DID) types.DataSource {
	dataSource, err := k.GetDataSource(ctx, id)
	if err != nil {
		panic(err)
	}
	return dataSource
}

// SetDataSource saves the given data source to the storage without performing validation.
func (k Keeper) SetDataSource(ctx sdk.Context, id types.DID, dataSource types.DataSource) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSourceStoreKey(id), k.cdc.MustMarshalBinaryBare(dataSource))
}

// AddDataSource adds the given data source to the storage. Returns error if validation fails.
func (k Keeper) AddDataSource(ctx sdk.Context, dataSource types.DataSource) (types.DID, error) {
	if err := AnyError(
		k.EnsureLength(ctx, types.KeyMaxNameLength, len(dataSource.Name)),
		k.EnsureLength(ctx, types.KeyMaxDescriptionLength, len(dataSource.Description)),
		k.EnsureLength(ctx, types.KeyMaxExecutableSize, len(dataSource.Executable)),
	); err != nil {
		return 0, err
	}
	id := k.GetNextDataSourceID(ctx)
	k.SetDataSource(ctx, id, dataSource)
	return id, nil
}

// EditDataSource edits the given data source by id and flushes it to the storage.
func (k Keeper) EditDataSource(ctx sdk.Context, id types.DID, new types.DataSource) error {
	dataSource, err := k.GetDataSource(ctx, id)
	if err != nil {
		return err
	}
	dataSource.Owner = new.Owner // TODO: Allow NOT_MODIFY or nil in these fields.
	dataSource.Name = new.Name
	dataSource.Description = new.Description
	dataSource.Fee = new.Fee
	dataSource.Executable = new.Executable
	if err := AnyError(
		k.EnsureLength(ctx, types.KeyMaxNameLength, len(dataSource.Name)),
		k.EnsureLength(ctx, types.KeyMaxDescriptionLength, len(dataSource.Description)),
		k.EnsureLength(ctx, types.KeyMaxExecutableSize, len(dataSource.Executable)),
	); err != nil {
		return err
	}
	k.SetDataSource(ctx, id, dataSource)
	return nil
}

// PayDataSourceFee sends fee from sender to data source owner. Returns error on failure.
func (k Keeper) PayDataSourceFee(ctx sdk.Context, id types.DID, sender sdk.AccAddress) error {
	dataSource, err := k.GetDataSource(ctx, id)
	if err != nil {
		return err
	}
	if dataSource.Owner.Equals(sender) || dataSource.Fee.IsZero() {
		// If you are the owner or it's free, no payment action is needed.
		return nil
	}
	return k.CoinKeeper.SendCoins(ctx, sender, dataSource.Owner, dataSource.Fee)
}

// GetAllDataSources returns the list of all data sources in the store, or nil if there is none.
func (k Keeper) GetAllDataSources(ctx sdk.Context) (dataSources []types.DataSource) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DataSourceStoreKeyPrefix)
	for ; iterator.Valid(); iterator.Next() {
		var dataSource types.DataSource
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &dataSource)
		dataSources = append(dataSources, dataSource)
	}
	return dataSources
}
