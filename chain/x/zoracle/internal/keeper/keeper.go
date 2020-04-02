package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
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
	ParamSpace    params.Subspace
}

// NewKeeper creates a new zoracle Keeper instance.
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, coinKeeper bank.Keeper,
	stakingKeeper staking.Keeper, paramSpace params.Subspace,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		CoinKeeper:    coinKeeper,
		StakingKeeper: stakingKeeper,
		ParamSpace:    paramSpace.WithKeyTable(ParamKeyTable()),
	}
}

// ParamKeyTable returns the parameter key table for zoracle module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// GetParam returns the parameter as specified by key as an uint64.
func (k Keeper) GetParam(ctx sdk.Context, key []byte) (res uint64) {
	k.ParamSpace.Get(ctx, key, &res)
	return
}

// SetParam saves the given key-value parameter to the store.
func (k Keeper) SetParam(ctx sdk.Context, key []byte, value uint64) {
	k.ParamSpace.Set(ctx, key, value)
}

// GetParams returns all current parameters as a types.Params instance.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.GetParam(ctx, types.KeyMaxDataSourceExecutableSize),
		k.GetParam(ctx, types.KeyMaxOracleScriptCodeSize),
		k.GetParam(ctx, types.KeyMaxCalldataSize),
		k.GetParam(ctx, types.KeyMaxDataSourceCountPerRequest),
		k.GetParam(ctx, types.KeyMaxRawDataReportSize),
		k.GetParam(ctx, types.KeyMaxResultSize),
		k.GetParam(ctx, types.KeyEndBlockExecuteGasLimit),
		k.GetParam(ctx, types.KeyMaxNameLength),
		k.GetParam(ctx, types.KeyMaxDescriptionLength),
		k.GetParam(ctx, types.KeyGasPerRawDataRequestPerValidator),
	)
}

// GetRequestCount returns the current number of all requests ever exist.
func (k Keeper) GetRequestCount(ctx sdk.Context) int64 {
	var requestNumber int64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RequestsCountStoreKey)
	if bz == nil {
		return 0
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &requestNumber)
	return requestNumber
}

// GetNextRequestID increments and returns the current number of requests.
// If the global request count is not set, it initializes it with value 0.
func (k Keeper) GetNextRequestID(ctx sdk.Context) types.RequestID {
	requestNumber := k.GetRequestCount(ctx)
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestNumber + 1)
	store.Set(types.RequestsCountStoreKey, bz)
	return types.RequestID(requestNumber + 1)
}

// GetDataSourceCount returns the current number of all data sources ever exist.
func (k Keeper) GetDataSourceCount(ctx sdk.Context) int64 {
	var dataSourceCount int64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DataSourceCountStoreKey)
	if bz == nil {
		return 0
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &dataSourceCount)
	return dataSourceCount
}

// GetNextDataSourceID increments and returns the current number of data source.
// If the global data source count is not set, it initializes the value and returns 1.
func (k Keeper) GetNextDataSourceID(ctx sdk.Context) types.DataSourceID {
	dataSourceCount := k.GetDataSourceCount(ctx)
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(dataSourceCount + 1)
	store.Set(types.DataSourceCountStoreKey, bz)
	return types.DataSourceID(dataSourceCount + 1)
}

// GetOracleScriptCount returns the current number of all oracle scripts ever exist.
func (k Keeper) GetOracleScriptCount(ctx sdk.Context) int64 {
	var oracleScriptCount int64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.OracleScriptCountStoreKey)
	if bz == nil {
		return 0
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &oracleScriptCount)
	return oracleScriptCount
}

// GetNextOracleScriptID increments and returns the current number of oracle script.
// If the global oracle script count is not set, it initializes the value and returns 1.
func (k Keeper) GetNextOracleScriptID(ctx sdk.Context) types.OracleScriptID {
	oracleScriptCount := k.GetOracleScriptCount(ctx)
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(oracleScriptCount + 1)
	store.Set(types.OracleScriptCountStoreKey, bz)
	return types.OracleScriptID(oracleScriptCount + 1)
}
