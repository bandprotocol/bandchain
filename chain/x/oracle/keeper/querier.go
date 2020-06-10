package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// NewQuerier is the module level router for state queries.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryParams:
			return queryParameters(ctx, keeper)
		case types.QueryCounts:
			return queryCounts(ctx, keeper)
		case types.QueryData:
			return queryData(ctx, path[1:], keeper)
		case types.QueryDataSources:
			return queryDataSourceByID(ctx, path[1:], keeper)
		case types.QueryOracleScripts:
			return queryOracleScriptByID(ctx, path[1:], keeper)
		case types.QueryRequests:
			return queryRequestByID(ctx, path[1:], keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown oracle query endpoint")
		}
	}
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, error) {
	return codec.MarshalJSONIndent(types.ModuleCdc, k.GetParams(ctx))
}

func queryCounts(ctx sdk.Context, k Keeper) ([]byte, error) {
	return codec.MarshalJSONIndent(types.ModuleCdc, types.QueryCountsResult{
		DataSourceCount:   k.GetDataSourceCount(ctx),
		OracleScriptCount: k.GetOracleScriptCount(ctx),
		RequestCount:      k.GetRequestCount(ctx),
	})
}

func queryData(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "data hash not specified")
	}
	return k.fileCache.GetFile(path[0])
}

func queryDataSourceByID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "data source not specified")
	}
	id, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, err.Error())
	}
	dataSource, err := k.GetDataSource(ctx, types.DataSourceID(id))
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSONIndent(types.ModuleCdc, dataSource)
}

func queryOracleScriptByID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "oracle script not specified")
	}
	id, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, err.Error())
	}
	oracleScript, err := k.GetOracleScript(ctx, types.OracleScriptID(id))
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSONIndent(types.ModuleCdc, oracleScript)
}

func queryRequestByID(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	if len(path) != 1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "request not specified")
	}
	id, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, err.Error())
	}
	request, err := k.GetRequest(ctx, types.RequestID(id))
	if err != nil {
		return nil, err
	}
	reports := k.GetReports(ctx, types.RequestID(id))

	if !k.HasResult(ctx, types.RequestID(id)) {
		return codec.MarshalJSONIndent(types.ModuleCdc, types.QueryRequestResult{
			Request: request,
			Reports: reports,
			Result:  nil,
		})
	}

	result := k.MustGetResult(ctx, types.RequestID(id))
	return codec.MarshalJSONIndent(types.ModuleCdc, types.QueryRequestResult{
		Request: request,
		Reports: reports,
		Result:  &result,
	})
}
