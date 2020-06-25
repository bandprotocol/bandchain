package keeper

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	fileCache     filecache.Cache
	ParamSpace    params.Subspace
	StakingKeeper types.StakingKeeper
}

// NewKeeper creates a new oracle Keeper instance.
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, fileDir string, paramSpace params.Subspace,
	stakingKeeper types.StakingKeeper,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(ParamKeyTable())
	}
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		fileCache:     filecache.New(fileDir),
		ParamSpace:    paramSpace,
		StakingKeeper: stakingKeeper,
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

// SetRollingSeed sets the rolling seed value to be provided value.
func (k Keeper) SetRollingSeed(ctx sdk.Context, rollingSeed []byte) {
	ctx.KVStore(k.storeKey).Set(types.RollingSeedStoreKey, rollingSeed)
}

// GetRollingSeed returns the current rolling seed value.
func (k Keeper) GetRollingSeed(ctx sdk.Context) []byte {
	bz := ctx.KVStore(k.storeKey).Get(types.RollingSeedStoreKey)
	if bz == nil {
		// If RollingSeedStoreKey is not yet set, we initialize it to a zero 32-byte slice.
		return make([]byte, 32)
	}
	return bz
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
	if bz == nil {
		return 0
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &requestNumber)
	return requestNumber
}

// SetRequestLastExpired sets the ID of the last expired request.
func (k Keeper) SetRequestLastExpired(ctx sdk.Context, id int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RequestLastExpiredStoreKey, k.cdc.MustMarshalBinaryLengthPrefixed(id))
}

// GetRequestLastExpired returns the ID of the last expired request.
func (k Keeper) GetRequestLastExpired(ctx sdk.Context) int64 {
	var requestNumber int64
	bz := ctx.KVStore(k.storeKey).Get(types.RequestLastExpiredStoreKey)
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
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestNumber + 1)
	ctx.KVStore(k.storeKey).Set(types.RequestCountStoreKey, bz)
	return types.RequestID(requestNumber + 1)
}

// GetDataSourceCount returns the current number of all data sources ever exist.
func (k Keeper) GetDataSourceCount(ctx sdk.Context) int64 {
	var dataSourceCount int64
	bz := ctx.KVStore(k.storeKey).Get(types.DataSourceCountStoreKey)
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
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(dataSourceCount + 1)
	ctx.KVStore(k.storeKey).Set(types.DataSourceCountStoreKey, bz)
	return types.DataSourceID(dataSourceCount + 1)
}

// GetOracleScriptCount returns the current number of all oracle scripts ever exist.
func (k Keeper) GetOracleScriptCount(ctx sdk.Context) int64 {
	var oracleScriptCount int64
	bz := ctx.KVStore(k.storeKey).Get(types.OracleScriptCountStoreKey)
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
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(oracleScriptCount + 1)
	ctx.KVStore(k.storeKey).Set(types.OracleScriptCountStoreKey, bz)
	return types.OracleScriptID(oracleScriptCount + 1)
}

// AddFile saves the given data to a file in HOME/files directory using sha256 sum as filename.
// Returns do-not-modify symbol is the given input file is do-not-modify.
func (k Keeper) AddFile(file []byte) string {
	if bytes.Equal(file, types.DoNotModifyBytes) {
		return types.DoNotModify
	}
	return k.fileCache.AddFile(file)
}

// GetFile loads the file from the file storage. Panics if the file does not exist.
func (k Keeper) GetFile(name string) []byte {
	return k.fileCache.MustGetFile(name)
}
