package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	CoinKeeper    bank.Keeper
	StakingKeeper staking.Keeper
}

// NewKeeper creates a new zoracle Keeper instance.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, coinKeeper bank.Keeper, stakingKeeper staking.Keeper) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		CoinKeeper:    coinKeeper,
		StakingKeeper: stakingKeeper,
	}
}

// GetRequestCount returns the current number of all requests ever exist.
func (k Keeper) GetRequestCount(ctx sdk.Context) uint64 {
	var requestNumber uint64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RequestsCountStoreKey)
	if bz == nil {
		return 0
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &requestNumber)
	if err != nil {
		panic(err)
	}
	return requestNumber
}

// GetNextRequestID increments and returns the current number of requests.
// If the global request count is not set, it initializes it with value 0.
func (k Keeper) GetNextRequestID(ctx sdk.Context) uint64 {
	requestNumber := k.GetRequestCount(ctx)
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestNumber + 1)
	store.Set(types.RequestsCountStoreKey, bz)
	return requestNumber + 1
}

// GetDataSourceCount returns the current number of all data sources ever exist.
func (k Keeper) GetDataSourceCount(ctx sdk.Context) int64 {
	var dataSourceCount int64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DataSourceCountStoreKey)
	if bz == nil {
		return 0
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &dataSourceCount)
	if err != nil {
		panic(err)
	}
	return dataSourceCount
}

// GetNextDataSourceID increments and returns the current number of data source.
// If the global data source count is not set, it initializes the value and returns 1.
func (k Keeper) GetNextDataSourceID(ctx sdk.Context) int64 {
	dataSourceCount := k.GetDataSourceCount(ctx)
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(dataSourceCount + 1)
	store.Set(types.DataSourceCountStoreKey, bz)
	return dataSourceCount + 1
}
