package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

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
		// case types.QueryScript:
		// 	return queryScript(ctx, path[1:], req, keeper)
		// case types.QueryAllScripts:
		// 	return queryAllScripts(ctx, path[1:], req, keeper)
		// case types.SerializeParams:
		// 	return serializeParams(ctx, path[1:], req, keeper)
		case types.QueryRequestNumber:
			return queryRequestNumber(ctx, req, keeper)
		case types.QueryDataSource:
			return queryDataSource(ctx, path[1:], req, keeper)
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
	id, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("wrong format for requestid %s", err.Error()))
	}
	request, sdkErr := keeper.GetRequest(ctx, id)
	if sdkErr != nil {
		return nil, sdkErr
	}

	rawRequests := keeper.GetRawDataRequestWithExternalIDs(ctx, id)

	iterator := keeper.GetRawDataReportsIterator(ctx, id)
	reportMap := make(map[string]([]types.RawDataReport))
	for ; iterator.Valid(); iterator.Next() {
		validator, externalID := types.GetValidatorAddressAndExternalID(iterator.Key(), id)
		if _, ok := reportMap[string(validator)]; !ok {
			reportMap[string(validator)] = make([]types.RawDataReport, 0)
		}

		reportMap[string(validator)] = append(
			reportMap[string(validator)],
			types.NewRawDataReport(externalID, iterator.Value()),
		)
	}

	reports := make([]types.ReportWithValidator, 0)

	for _, validator := range request.RequestedValidators {
		valReport, ok := reportMap[string(validator)]
		if ok {
			reports = append(reports, types.NewReportWithValidator(
				valReport, validator,
			))
		}
	}
	var result []byte
	if keeper.HasResult(ctx, id, request.OracleScriptID, request.Calldata) {
		var sdkErr sdk.Error
		result, sdkErr = keeper.GetResult(ctx, id, request.OracleScriptID, request.Calldata)
		if sdkErr != nil {
			return nil, sdkErr
		}
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewRequestQuerierInfo(
		request,
		rawRequests,
		reports,
		result,
	))
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return res, nil
}

// queryPending is a query function to get the list of request IDs that are still on pending status.
func queryPending(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	res, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetPendingResolveList(ctx))
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return res, nil
}

// // queryScript is a query function to get infomation of stored wasm code
// func queryScript(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
// 	if len(path) == 0 {
// 		return nil, sdk.ErrInternal("must specify the codehash")
// 	}
// 	codeHash, err := hex.DecodeString(path[0])
// 	if err != nil {
// 		return nil, sdk.ErrUnknownRequest("cannot decode hexstr")
// 	}
// 	if !keeper.CheckCodeHashExists(ctx, codeHash) {
// 		return nil, sdk.ErrUnknownRequest("codehash not found")
// 	}
// 	code, err := keeper.GetCode(ctx, codeHash)
// 	if err != nil {
// 		panic("cannot get codehash")
// 	}

// 	// Get raw params info
// 	rawParamsInfo, err := wasm.ParamsInfo(code.Code)
// 	var paramsInfo []types.Field
// 	if err != nil {
// 		paramsInfo = nil
// 	} else {
// 		paramsInfo, err = types.ParseFields(rawParamsInfo)
// 		if err != nil {
// 			paramsInfo = nil
// 		}
// 	}

// 	// Get raw data sources info
// 	rawDataInfo, err := wasm.RawDataInfo(code.Code)
// 	var dataInfo []types.Field
// 	if err != nil {
// 		dataInfo = nil
// 	} else {
// 		dataInfo, err = types.ParseFields(rawDataInfo)
// 		if err != nil {
// 			dataInfo = nil
// 		}
// 	}

// 	// Get result info
// 	rawResultInfo, err := wasm.ResultInfo(code.Code)
// 	var resultInfo []types.Field
// 	if err != nil {
// 		resultInfo = nil
// 	} else {
// 		resultInfo, err = types.ParseFields(rawResultInfo)
// 		if err != nil {
// 			resultInfo = nil
// 		}
// 	}

// 	return codec.MustMarshalJSONIndent(keeper.cdc, types.NewScriptInfo(
// 		code.Name,
// 		codeHash,
// 		paramsInfo,
// 		dataInfo,
// 		resultInfo,
// 		code.Owner,
// 	)), nil
// }

// func queryAllScripts(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
// 	if len(path) != 2 {
// 		return nil, sdk.ErrInternal("must specify the page and limit")
// 	}
// 	page, err := strconv.ParseInt(path[0], 10, 64)
// 	if err != nil {
// 		return nil, sdk.ErrInternal("page must be a number")
// 	}

// 	limit, err := strconv.ParseInt(path[1], 10, 64)
// 	if err != nil {
// 		return nil, sdk.ErrInternal("limit must be a number")
// 	}

// 	start := int((page - 1) * limit)
// 	end := int(page * limit)

// 	// TODO: Current fetch all scripts and return only script in list
// 	iterator := keeper.GetCodesIterator(ctx)
// 	results := make([]types.ScriptInfo, 0)
// 	for i := 0; i < end && iterator.Valid(); iterator.Next() {
// 		if i >= start {
// 			var script types.ScriptInfo
// 			rawInfo, sdkErr := queryScript(ctx, []string{hex.EncodeToString(iterator.Key()[1:])}, req, keeper)
// 			if sdkErr != nil {
// 				continue
// 			}
// 			err := keeper.cdc.UnmarshalJSON(rawInfo, &script)
// 			if err != nil {
// 				continue
// 			}
// 			results = append(results, script)
// 		}
// 		i++
// 	}
// 	return codec.MustMarshalJSONIndent(keeper.cdc, results), nil
// }

// // serializeParams is a function that receive codeHash and params and return a serialized format of params
// func serializeParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
// 	if len(path) != 2 {
// 		return nil, sdk.ErrInternal(fmt.Sprintf("number of arguments must be 2, but got %d", len(path)))
// 	}

// 	codeHash, err := hex.DecodeString(path[0])
// 	if err != nil {
// 		return nil, sdk.ErrUnknownRequest("cannot decode hexstr")
// 	}
// 	if !keeper.CheckCodeHashExists(ctx, codeHash) {
// 		return nil, sdk.ErrUnknownRequest("codehash not found")
// 	}
// 	code, err := keeper.GetCode(ctx, codeHash)
// 	if err != nil {
// 		panic("cannot get codehash")
// 	}

// 	rawParams, err := wasm.SerializeParams(code.Code, []byte(path[1]))

// 	return codec.MustMarshalJSONIndent(keeper.cdc, rawParams), nil
// }

func queryRequestNumber(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	return codec.MustMarshalJSONIndent(keeper.cdc, keeper.GetRequestCount(ctx)), nil
}

func queryDataSource(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) == 0 {
		return nil, sdk.ErrInternal("must specify the data source id")
	}
	id, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("wrong format for data source id %s", err.Error()))
	}

	dataSource, sdkErr := keeper.GetDataSource(ctx, id)
	if sdkErr != nil {
		return nil, sdkErr
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewDataSourceQuerierInfo(
		id,
		dataSource.Owner,
		dataSource.Name,
		dataSource.Fee,
		dataSource.Executable,
	))
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
