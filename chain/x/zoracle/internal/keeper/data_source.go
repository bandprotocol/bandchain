package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetDataSource saves the given data source with the given ID to the storage.
func (k Keeper) SetDataSource(ctx sdk.Context, id int64, dataSource types.DataSource) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSourceStoreKey(id), k.cdc.MustMarshalBinaryBare(dataSource))
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
