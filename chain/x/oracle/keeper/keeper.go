package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	owasm "github.com/bandprotocol/go-owasm/api"
)

const (
	RollingSeedSizeInBytes = 32
)

type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	fileCache        filecache.Cache
	feeCollectorName string
	paramSpace       params.Subspace
	supplyKeeper     types.SupplyKeeper
	stakingKeeper    types.StakingKeeper
	distrKeeper      types.DistrKeeper
	owasmVM          *owasm.Vm
}

// NewKeeper creates a new oracle Keeper instance.
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, fileDir string, feeCollectorName string,
	paramSpace params.Subspace, supplyKeeper types.SupplyKeeper,
	stakingKeeper types.StakingKeeper, distrKeeper types.DistrKeeper,
	owasmVM *owasm.Vm,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(ParamKeyTable())
	}
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		fileCache:        filecache.New(fileDir),
		feeCollectorName: feeCollectorName,
		paramSpace:       paramSpace,
		supplyKeeper:     supplyKeeper,
		stakingKeeper:    stakingKeeper,
		distrKeeper:      distrKeeper,
		owasmVM:          owasmVM,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ParamKeyTable returns the parameter key table for oracle module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// GetParam returns the parameter as specified by key as an uint64.
func (k Keeper) GetParam(ctx sdk.Context, key []byte) (res uint64) {
	k.paramSpace.Get(ctx, key, &res)
	return res
}

// SetParam saves the given key-value parameter to the store.
func (k Keeper) SetParam(ctx sdk.Context, key []byte, value uint64) {
	k.paramSpace.Set(ctx, key, value)
}

// GetParams returns all current parameters as a types.Params instance.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetRollingSeed sets the rolling seed value to be provided value.
func (k Keeper) SetRollingSeed(ctx sdk.Context, rollingSeed []byte) {
	ctx.KVStore(k.storeKey).Set(types.RollingSeedStoreKey, rollingSeed)
}

// GetRollingSeed returns the current rolling seed value.
func (k Keeper) GetRollingSeed(ctx sdk.Context) []byte {
	return ctx.KVStore(k.storeKey).Get(types.RollingSeedStoreKey)
}

// SetRequestCount sets the number of request count to the given value. Useful for genesis state.
func (k Keeper) SetRequestCount(ctx sdk.Context, count int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RequestCountStoreKey, k.cdc.MustMarshalBinaryLengthPrefixed(count))
}

// GetRequestCount returns the current number of all requests ever exist.
func (k Keeper) GetRequestCount(ctx sdk.Context) int64 {
	var requestNumber int64
	bz := ctx.KVStore(k.storeKey).Get(types.RequestCountStoreKey)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &requestNumber)
	return requestNumber
}

// SetRequestLastExpired sets the ID of the last expired request.
func (k Keeper) SetRequestLastExpired(ctx sdk.Context, id types.RequestID) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RequestLastExpiredStoreKey, k.cdc.MustMarshalBinaryLengthPrefixed(id))
}

// GetRequestLastExpired returns the ID of the last expired request.
func (k Keeper) GetRequestLastExpired(ctx sdk.Context) types.RequestID {
	var requestNumber types.RequestID
	bz := ctx.KVStore(k.storeKey).Get(types.RequestLastExpiredStoreKey)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &requestNumber)
	return requestNumber
}

// GetNextRequestID increments and returns the current number of requests.
func (k Keeper) GetNextRequestID(ctx sdk.Context) types.RequestID {
	requestNumber := k.GetRequestCount(ctx)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestNumber + 1)
	ctx.KVStore(k.storeKey).Set(types.RequestCountStoreKey, bz)
	return types.RequestID(requestNumber + 1)
}

// SetDataSourceCount sets the number of data source count to the given value.
func (k Keeper) SetDataSourceCount(ctx sdk.Context, count int64) {
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(count)
	ctx.KVStore(k.storeKey).Set(types.DataSourceCountStoreKey, bz)
}

// GetDataSourceCount returns the current number of all data sources ever exist.
func (k Keeper) GetDataSourceCount(ctx sdk.Context) int64 {
	var dataSourceCount int64
	bz := ctx.KVStore(k.storeKey).Get(types.DataSourceCountStoreKey)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &dataSourceCount)
	return dataSourceCount
}

// GetNextDataSourceID increments and returns the current number of data sources.
func (k Keeper) GetNextDataSourceID(ctx sdk.Context) types.DataSourceID {
	dataSourceCount := k.GetDataSourceCount(ctx)
	k.SetDataSourceCount(ctx, dataSourceCount+1)
	return types.DataSourceID(dataSourceCount + 1)
}

// SetOracleScriptCount sets the number of oracle script count to the given value.
func (k Keeper) SetOracleScriptCount(ctx sdk.Context, count int64) {
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(count)
	ctx.KVStore(k.storeKey).Set(types.OracleScriptCountStoreKey, bz)
}

// GetOracleScriptCount returns the current number of all oracle scripts ever exist.
func (k Keeper) GetOracleScriptCount(ctx sdk.Context) int64 {
	var oracleScriptCount int64
	bz := ctx.KVStore(k.storeKey).Get(types.OracleScriptCountStoreKey)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &oracleScriptCount)
	return oracleScriptCount
}

// GetNextOracleScriptID increments and returns the current number of oracle scripts.
func (k Keeper) GetNextOracleScriptID(ctx sdk.Context) types.OracleScriptID {
	oracleScriptCount := k.GetOracleScriptCount(ctx)
	k.SetOracleScriptCount(ctx, oracleScriptCount+1)
	return types.OracleScriptID(oracleScriptCount + 1)
}

// GetFile loads the file from the file storage. Panics if the file does not exist.
func (k Keeper) GetFile(name string) []byte {
	return k.fileCache.MustGetFile(name)
}
