package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	ParamSpace    params.Subspace
	CoinKeeper    bank.Keeper
	StakingKeeper staking.Keeper
	ChannelKeeper types.ChannelKeeper
	scopedKeeper  capability.ScopedKeeper
}

// NewKeeper creates a new oracle Keeper instance.
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	coinKeeper bank.Keeper, stakingKeeper staking.Keeper,
	channelKeeper types.ChannelKeeper, scopedKeeper capability.ScopedKeeper,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(ParamKeyTable())
	}
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		ParamSpace:    paramSpace,
		CoinKeeper:    coinKeeper,
		StakingKeeper: stakingKeeper,
		ChannelKeeper: channelKeeper,
		scopedKeeper:  scopedKeeper,
	}
}

// ParamKeyTable returns the parameter key table for oracle module.
func ParamKeyTable() params.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&types.Params{})
}

// GetParam returns the parameter as specified by key as an uint64.
func (k Keeper) GetParam(ctx sdk.Context, key []byte) (res uint64) {
	k.ParamSpace.Get(ctx, key, &res)
	return res
}

// SetParam saves the given key-value parameter to the store.
func (k Keeper) SetParam(ctx sdk.Context, key []byte, value uint64) {
	k.ParamSpace.Set(ctx, key, value)
}

// GetParams returns all current parameters as a types.Params instance.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.ParamSpace.GetParamSet(ctx, &params)
	return params
}

// GetRequestBeginID TODO
func (k Keeper) GetRequestBeginID(ctx sdk.Context) types.RequestID {
	var requestBeginID int64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RequestBeginStoreKey)
	if bz == nil {
		return types.RequestID(1)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &requestBeginID)
	return types.RequestID(requestBeginID)
}

// SetRequestBeginID TODO
func (k Keeper) SetRequestBeginID(ctx sdk.Context, id types.RequestID) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(id)
	store.Set(types.RequestBeginStoreKey, bz)
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
