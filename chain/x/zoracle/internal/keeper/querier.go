package keeper

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/d3n/chain/wasm"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

// NewQuerier is the module level router for state queries.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryRequest:
			return queryRequest(ctx, path[1:], req, keeper)
		case types.QueryPending:
			return queryPending(ctx, path[1:], req, keeper)
		case types.QueryScript:
			return queryScript(ctx, path[1:], req, keeper)
		case types.QueryAllScripts:
			return queryAllScripts(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

// queryRequest is a query function to get request information by request ID.
func queryRequest(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) == 0 {
		return nil, sdk.ErrInternal("must specify the requestid")
	}
	reqID, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("wrong format for requestid %s", err.Error()))
	}
	request, sdkErr := keeper.GetRequest(ctx, reqID)
	if sdkErr != nil {
		return nil, sdkErr
	}
	reports, sdkErr := keeper.GetValidatorReports(ctx, reqID)
	if sdkErr != nil {
		return nil, sdkErr
	}
	result, sdkErr := keeper.GetResult(ctx, reqID, request.CodeHash, request.Params)
	if sdkErr != nil {
		result = []byte{}
	}

	code, sdkErr := keeper.GetCode(ctx, request.CodeHash)
	if sdkErr != nil {
		return nil, sdkErr
	}

	rawParams, err := wasm.ParseParams(code.Code, request.Params)
	if err != nil {
		rawParams = []byte{}
	}
	res, err := codec.MarshalJSONIndent(
		keeper.cdc,
		types.NewRequestInfo(request.CodeHash, request.Params, rawParams, request.ReportEndAt, reports, result))
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return res, nil
}

// queryPending is a query function to get the list of request IDs that are still on pending status.
func queryPending(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	res, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetPending(ctx))
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return res, nil
}

// queryScript is a query function to get infomation of stored wasm code
func queryScript(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) == 0 {
		return nil, sdk.ErrInternal("must specify the codehash")
	}
	codeHash, err := hex.DecodeString(path[0])
	if err != nil {
		return nil, sdk.ErrUnknownRequest("cannot decode hexstr")
	}
	if !keeper.CheckCodeHashExists(ctx, codeHash) {
		return nil, sdk.ErrUnknownRequest("codehash not found")
	}
	code, err := keeper.GetCode(ctx, codeHash)
	if err != nil {
		panic("cannot get codehash")
	}

	// Get raw params info
	rawParamsInfo, err := wasm.ParamsInfo(code.Code)
	var paramsInfo []types.Field
	if err != nil {
		paramsInfo = nil
	} else {
		paramsInfo, err = types.ParseFields(rawParamsInfo)
		if err != nil {
			paramsInfo = nil
		}
	}

	// Get raw data sources info
	rawDataInfo, err := wasm.RawDataInfo(code.Code)
	var dataInfo []types.Field
	if err != nil {
		dataInfo = nil
	} else {
		dataInfo, err = types.ParseFields(rawDataInfo)
		if err != nil {
			dataInfo = nil
		}
	}

	return codec.MustMarshalJSONIndent(keeper.cdc, types.NewScriptInfo(code.Name, codeHash, paramsInfo, dataInfo, code.Owner)), nil
}

func queryAllScripts(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) != 2 {
		return nil, sdk.ErrInternal("must specify the page and limit")
	}
	page, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return nil, sdk.ErrInternal("page must be a number")
	}

	limit, err := strconv.ParseInt(path[1], 10, 64)
	if err != nil {
		return nil, sdk.ErrInternal("limit must be a number")
	}

	start := int((page - 1) * limit)
	end := int(page * limit)

	// TODO: Current fetch all scripts and return only script in list
	iterator := keeper.GetCodesIterator(ctx)
	results := make([]types.ScriptInfo, 0)
	for i := 0; i < end && iterator.Valid(); iterator.Next() {
		if i >= start {
			var script types.ScriptInfo
			rawInfo, sdkErr := queryScript(ctx, []string{hex.EncodeToString(iterator.Key()[1:])}, req, keeper)
			if sdkErr != nil {
				continue
			}
			err := keeper.cdc.UnmarshalJSON(rawInfo, &script)
			if err != nil {
				continue
			}
			results = append(results, script)
		}
		i++
	}
	return codec.MustMarshalJSONIndent(keeper.cdc, results), nil
}
