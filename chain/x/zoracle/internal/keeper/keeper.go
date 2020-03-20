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

func (keeper Keeper) MaxDataSourceExecutableSize(ctx sdk.Context) (res int64) {
	keeper.ParamSpace.Get(ctx, types.KeyMaxDataSourceExecutableSize, &res)
	return
}

func (keeper Keeper) SetMaxDataSourceExecutableSize(ctx sdk.Context, value int64) {
	keeper.ParamSpace.Set(ctx, types.KeyMaxDataSourceExecutableSize, value)
}

func (keeper Keeper) MaxOracleScriptCodeSize(ctx sdk.Context) (res int64) {
	keeper.ParamSpace.Get(ctx, types.KeyMaxOracleScriptCodeSize, &res)
	return
}

func (keeper Keeper) SetMaxOracleScriptCodeSize(ctx sdk.Context, value int64) {
	keeper.ParamSpace.Set(ctx, types.KeyMaxOracleScriptCodeSize, value)
}

func (keeper Keeper) MaxCalldataSize(ctx sdk.Context) (res int64) {
	keeper.ParamSpace.Get(ctx, types.KeyMaxCalldataSize, &res)
	return
}

func (keeper Keeper) SetMaxCalldataSize(ctx sdk.Context, value int64) {
	keeper.ParamSpace.Set(ctx, types.KeyMaxCalldataSize, value)
}

func (keeper Keeper) MaxDataSourceCountPerRequest(ctx sdk.Context) (res int64) {
	keeper.ParamSpace.Get(ctx, types.KeyMaxDataSourceCountPerRequest, &res)
	return
}

func (keeper Keeper) SetMaxDataSourceCountPerRequest(ctx sdk.Context, value int64) {
	keeper.ParamSpace.Set(ctx, types.KeyMaxDataSourceCountPerRequest, value)
}

func (keeper Keeper) MaxRawDataReportSize(ctx sdk.Context) (res int64) {
	keeper.ParamSpace.Get(ctx, types.KeyMaxRawDataReportSize, &res)
	return
}

func (keeper Keeper) SetMaxRawDataReportSize(ctx sdk.Context, value int64) {
	keeper.ParamSpace.Set(ctx, types.KeyMaxRawDataReportSize, value)
}

func (keeper Keeper) MaxResultSize(ctx sdk.Context) (res int64) {
	keeper.ParamSpace.Get(ctx, types.KeyMaxResultSize, &res)
	return
}

func (keeper Keeper) SetMaxResultSize(ctx sdk.Context, value int64) {
	keeper.ParamSpace.Set(ctx, types.KeyMaxResultSize, value)
}

func (keeper Keeper) EndBlockExecuteGasLimit(ctx sdk.Context) (res uint64) {
	keeper.ParamSpace.Get(ctx, types.KeyEndBlockExecuteGasLimit, &res)
	return
}

func (keeper Keeper) SetEndBlockExecuteGasLimit(ctx sdk.Context, value uint64) {
	keeper.ParamSpace.Set(ctx, types.KeyEndBlockExecuteGasLimit, value)
}

func (keeper Keeper) SetMaxNameLength(ctx sdk.Context, value int64) {
	keeper.ParamSpace.Set(ctx, types.KeyMaxNameLength, value)
}

func (keeper Keeper) MaxNameLength(ctx sdk.Context) (res int64) {
	keeper.ParamSpace.Get(ctx, types.KeyMaxNameLength, &res)
	return
}
func (keeper Keeper) MaxDescriptionLength(ctx sdk.Context) (res int64) {
	keeper.ParamSpace.Get(ctx, types.KeyMaxDescriptionLength, &res)
	return
}

func (keeper Keeper) SetMaxDescriptionLength(ctx sdk.Context, value int64) {
	keeper.ParamSpace.Set(ctx, types.KeyMaxDescriptionLength, value)
}

func (keeper Keeper) GasPerRawDataRequestPerValidator(ctx sdk.Context) (res uint64) {
	keeper.ParamSpace.Get(ctx, types.KeyGasPerRawDataRequestPerValidator, &res)
	return
}

func (keeper Keeper) SetGasPerRawDataRequestPerValidator(ctx sdk.Context, value uint64) {
	keeper.ParamSpace.Set(ctx, types.KeyGasPerRawDataRequestPerValidator, value)
}

// GetParams returns all current parameters as a types.Params instance.
func (keeper Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		keeper.MaxDataSourceExecutableSize(ctx),
		keeper.MaxOracleScriptCodeSize(ctx),
		keeper.MaxCalldataSize(ctx),
		keeper.MaxDataSourceCountPerRequest(ctx),
		keeper.MaxRawDataReportSize(ctx),
		keeper.MaxResultSize(ctx),
		keeper.EndBlockExecuteGasLimit(ctx),
		keeper.MaxNameLength(ctx),
		keeper.MaxDescriptionLength(ctx),
		keeper.GasPerRawDataRequestPerValidator(ctx),
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
