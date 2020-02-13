package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	CoinKeeper    bank.Keeper
	StakingKeeper staking.Keeper
	Paramspace    params.Subspace
}

// NewKeeper creates a new zoracle Keeper instance.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, coinKeeper bank.Keeper, stakingKeeper staking.Keeper, paramspace params.Subspace) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		CoinKeeper:    coinKeeper,
		StakingKeeper: stakingKeeper,
		Paramspace:    paramspace.WithKeyTable(ParamKeyTable()),
	}
}

// ParamKeyTable for zoracle module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// Get MaxDataSourceExecutableSize.
func (keeper Keeper) MaxDataSourceExecutableSize(ctx sdk.Context) (res int64) {
	keeper.Paramspace.Get(ctx, types.KeyMaxDataSourceExecutableSize, &res)
	return
}

// Set SetMaxDataSourceExecutableSize.
func (keeper Keeper) SetMaxDataSourceExecutableSize(ctx sdk.Context, value int64) {
	keeper.Paramspace.Set(ctx, types.KeyMaxDataSourceExecutableSize, value)
}

// Get MaxOracleScriptCodeSize.
func (keeper Keeper) MaxOracleScriptCodeSize(ctx sdk.Context) (res int64) {
	keeper.Paramspace.Get(ctx, types.KeyMaxOracleScriptCodeSize, &res)
	return
}

// Set SetMaxOracleScriptCodeSize.
func (keeper Keeper) SetMaxOracleScriptCodeSize(ctx sdk.Context, value int64) {
	keeper.Paramspace.Set(ctx, types.KeyMaxOracleScriptCodeSize, value)
}

// Get MaxCalldataSize.
func (keeper Keeper) MaxCalldataSize(ctx sdk.Context) (res int64) {
	keeper.Paramspace.Get(ctx, types.KeyMaxCalldataSize, &res)
	return
}

// Set SetMaxCalldataSize.
func (keeper Keeper) SetMaxCalldataSize(ctx sdk.Context, value int64) {
	keeper.Paramspace.Set(ctx, types.KeyMaxCalldataSize, value)
}

// Get MaxDataSourceCountPerRequest.
func (keeper Keeper) MaxDataSourceCountPerRequest(ctx sdk.Context) (res int64) {
	keeper.Paramspace.Get(ctx, types.KeyMaxDataSourceCountPerRequest, &res)
	return
}

// Set SetMaxDataSourceCountPerRequest.
func (keeper Keeper) SetMaxDataSourceCountPerRequest(ctx sdk.Context, value int64) {
	keeper.Paramspace.Set(ctx, types.KeyMaxDataSourceCountPerRequest, value)
}

// Get MaxRawDataReportSize.
func (keeper Keeper) MaxRawDataReportSize(ctx sdk.Context) (res int64) {
	keeper.Paramspace.Get(ctx, types.KeyMaxRawDataReportSize, &res)
	return
}

// Set SetMaxRawDataReportSize.
func (keeper Keeper) SetMaxRawDataReportSize(ctx sdk.Context, value int64) {
	keeper.Paramspace.Set(ctx, types.KeyMaxRawDataReportSize, value)
}

// Get all parameteras as types.Params.
func (keeper Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		keeper.MaxDataSourceExecutableSize(ctx),
		keeper.MaxOracleScriptCodeSize(ctx),
		keeper.MaxCalldataSize(ctx),
		keeper.MaxDataSourceCountPerRequest(ctx),
		keeper.MaxRawDataReportSize(ctx),
	)
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

// GetOracleScriptCount returns the current number of all oracle scripts ever exist.
func (k Keeper) GetOracleScriptCount(ctx sdk.Context) int64 {
	var oracleScriptCount int64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.OracleScriptCountStoreKey)
	if bz == nil {
		return 0
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &oracleScriptCount)
	if err != nil {
		panic(err)
	}
	return oracleScriptCount
}

// GetNextOracleScriptID increments and returns the current number of oracle script.
// If the global oracle script count is not set, it initializes the value and returns 1.
func (k Keeper) GetNextOracleScriptID(ctx sdk.Context) int64 {
	oracleScriptCount := k.GetOracleScriptCount(ctx)
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(oracleScriptCount + 1)
	store.Set(types.OracleScriptCountStoreKey, bz)
	return oracleScriptCount + 1
}
