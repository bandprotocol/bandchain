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

var from = sdk.AccAddress([]byte("BandBandBandBandBand"))

func hexToBytes(hexstr string) []byte {
	b, err := hex.DecodeString(hexstr)
	if err != nil {
		panic(err)
	}
	return b
}

var ScriptInfoMap = map[string]types.ScriptInfo{
	"06077e219b9bceb3ca90240d5f9d383e418e9916a9da02fce7aa441d279af2d4": types.NewScriptInfo("cryptocompare_price", 
		hexToBytes("06077e219b9bceb3ca90240d5f9d383e418e9916a9da02fce7aa441d279af2d4"), 
		[]types.Field{
			types.NewField("crypto_symbol", "coins::Coins"),
		}, 
		[]types.Field{
			types.NewField("crypto_compare_price", "f32"), 
			types.NewField("timestamp", "u64"),
		}, 
		[]types.Field{
			types.NewField("crypto_price_in_usd", "u64"),
			types.NewField("timestamp", "u64"),
		},
		from),
		"1464c1ad00a7209b7ba693603b56ee2da851d94aa7339b2de0bce5316f465c22": types.NewScriptInfo("coingecko_volume", 
		hexToBytes("1464c1ad00a7209b7ba693603b56ee2da851d94aa7339b2de0bce5316f465c22"), 
		[]types.Field{
			types.NewField("crypto_symbol", "coins::Coins"),
		}, 
		[]types.Field{
			types.NewField("coin_gecko_vol24h", "f64"), 
			types.NewField("timestamp", "u64"),
		}, 
		[]types.Field{
			types.NewField("vol24h_in_usd", "u64"),
			types.NewField("timestamp", "u64"),
		},
		from),
		"1aa0bb8d7584921357cd2c6847f7a0cd896e3ae99a6b95908f22a693d9e6c5c1": types.NewScriptInfo("weather_info", 
		hexToBytes("1aa0bb8d7584921357cd2c6847f7a0cd896e3ae99a6b95908f22a693d9e6c5c1"), 
		[]types.Field{
			types.NewField("city", "String"),
			types.NewField("key", "String"),
			types.NewField("sub_key", "String"),
		}, 
		[]types.Field{
			types.NewField("weather_value", "f64"),
			types.NewField("timestamp", "u64"),
		}, 
		[]types.Field{
			types.NewField("weather_value", "u64"),
			types.NewField("timestamp", "u64"),
		},
		from),
		"31eea00625c5fe2ea80a7fde255469a2f57534dfadc715ed30228e7392a6529c": types.NewScriptInfo("ethereum_gas_price", 
		hexToBytes("31eea00625c5fe2ea80a7fde255469a2f57534dfadc715ed30228e7392a6529c"), 
		[]types.Field{
			types.NewField("gas_option", "String"),
		}, 
		[]types.Field{
			types.NewField("gas_price", "f32"), 
			types.NewField("timestamp", "u64"),
		}, 
		[]types.Field{
			types.NewField("gas_price_in_gwei", "u64"),
			types.NewField("timestamp", "u64"),
		},
		from),
		"5bc4ba4f391f83ebe066387c43cc1576bcf80602671bd4831e008973dfea8bb9": types.NewScriptInfo("flight_verification", 
		hexToBytes("5bc4ba4f391f83ebe066387c43cc1576bcf80602671bd4831e008973dfea8bb9"), 
		[]types.Field{
			types.NewField("icao24", "String"),
			types.NewField("flight_option", "flight_option::FlightOption"),
			types.NewField("airport", "String"),
			types.NewField("should_happen_before", "u64"),
			types.NewField("should_happen_after", "u64"),
		}, 
		[]types.Field{
			types.NewField("has_flight_found", "bool"), 
		}, 
		[]types.Field{
			types.NewField("has_flight_found", "bool"),
		},
		from),
		"6b7be61b150aec5eb853afb3b53e41438959554580d31259a1095e51645bcd28": types.NewScriptInfo("binance_price", 
		hexToBytes("6b7be61b150aec5eb853afb3b53e41438959554580d31259a1095e51645bcd28"), 
		[]types.Field{
			types.NewField("crypto_symbol", "coins::Coins"),
		}, 
		[]types.Field{
			types.NewField("binance_price", "f32"), 
			types.NewField("timestamp", "u64"), 
		}, 
		[]types.Field{
			types.NewField("crypto_price_in_usd", "u64"),
			types.NewField("timestamp", "u64"),
		},
		from),
		"6fd319bd571dcb0828dc2b5317df806bda1f70bdf08e3c3fb2046d5fd9a8982d": types.NewScriptInfo("alphavantage", 
		hexToBytes("6fd319bd571dcb0828dc2b5317df806bda1f70bdf08e3c3fb2046d5fd9a8982d"), 
		[]types.Field{
			types.NewField("stock_symbol", "String"),
		}, 
		[]types.Field{
			types.NewField("alphavantage_price", "f32"), 
			types.NewField("timestamp", "u64"), 
		}, 
		[]types.Field{
			types.NewField("stock_price_in_usd", "u64"),
			types.NewField("timestamp", "u64"),
		},
		from),
		"d10a439e08c31902cfed8d337e68f297c84665c36c1025caa60ce6bbd2767715": types.NewScriptInfo("coingecko_price", 
		hexToBytes("d10a439e08c31902cfed8d337e68f297c84665c36c1025caa60ce6bbd2767715"), 
		[]types.Field{
			types.NewField("crypto_symbol", "coins::Coins"),
		}, 
		[]types.Field{
			types.NewField("coin_gecko_price", "f32"), 
			types.NewField("timestamp", "u64"), 
		}, 
		[]types.Field{
			types.NewField("crypto_price_in_usd", "u64"),
			types.NewField("timestamp", "u64"),
		},
		from),
		"d3497a03c318476612524fb6d8f6c761e411fdf67966b08d157b28a83980cf7c": types.NewScriptInfo("cryptocompare_volume", 
		hexToBytes("d3497a03c318476612524fb6d8f6c761e411fdf67966b08d157b28a83980cf7c"), 
		[]types.Field{
			types.NewField("crypto_symbol", "coins::Coins"),
		}, 
		[]types.Field{
			types.NewField("crypto_compare_vol24h", "f64"), 
			types.NewField("timestamp", "u64"),
		}, 
		[]types.Field{
			types.NewField("vol24h_in_usd", "u64"),
			types.NewField("timestamp", "u64"),
		},
		from),
		"e7944e5e24dc856dcb6d9926460926ec10b9b66cf44b664f9971b5a5e9255989": types.NewScriptInfo("random_u64", 
		hexToBytes("e7944e5e24dc856dcb6d9926460926ec10b9b66cf44b664f9971b5a5e9255989"), 
		[]types.Field{
			types.NewField("max_range", "u64"),
		}, 
		[]types.Field{
			types.NewField("random_bytes8", "Vec<u8>"), 
			types.NewField("timestamp", "u64"),
		}, 
		[]types.Field{
			types.NewField("random_u64", "u64"),
			types.NewField("timestamp", "u64"),
		},
		from),
		"ebfefd555d5b7b3e5bd7d9c6a9f396cbf22697edd1687a54d3425f8b4e4eb480": types.NewScriptInfo("bitcoin_info", 
		hexToBytes("ebfefd555d5b7b3e5bd7d9c6a9f396cbf22697edd1687a54d3425f8b4e4eb480"), 
		[]types.Field{
			types.NewField("block_height", "u64"),
		}, 
		[]types.Field{
			types.NewField("block_hash", "[u64; 4]"), 
			types.NewField("block_count", "u64"),
		}, 
		[]types.Field{
			types.NewField("block_hash", "[u64; 4]"),
			types.NewField("confirmation", "u64"),
		},
		from),
}


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
		case types.SerializeParams:
			return serializeParams(ctx, path[1:], req, keeper)
		case types.QueryRequestNumber:
			return queryRequestNumber(ctx, req, keeper)
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

	code, sdkErr := keeper.GetCode(ctx, request.CodeHash)
	if sdkErr != nil {
		return nil, sdkErr
	}

	reports, sdkErr := keeper.GetValidatorReports(ctx, reqID)
	if sdkErr != nil {
		return nil, sdkErr
	}
	result, sdkErr := keeper.GetResult(ctx, reqID, request.CodeHash, request.Params)
	var parsedResult []byte
	if sdkErr != nil {
		result = []byte{}
		parsedResult = []byte("{}")
	} else {
		parsedResult, err = wasm.ParseResult(code.Code, result)
		if err != nil {
			parsedResult = []byte("{}")
		}
	}

	parsedParams, err := wasm.ParseParams(code.Code, request.Params)
	if err != nil {
		parsedParams = []byte("{}")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewRequestInfo(
		request.CodeHash,
		parsedParams,
		request.Params,
		request.ReportEndAt,
		reports,
		parsedResult,
		result,
	))
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

	// TODO: Hard-coded all scripts.
	scriptInfo, ok := ScriptInfoMap[path[0]]
	if !ok {
		return nil, sdk.ErrUnknownRequest("codehash not found")
	}

	// codeHash, err := hex.DecodeString(path[0])
	// if err != nil {
	// 	return nil, sdk.ErrUnknownRequest("cannot decode hexstr")
	// }
	// if !keeper.CheckCodeHashExists(ctx, codeHash) {
	// 	return nil, sdk.ErrUnknownRequest("codehash not found")
	// }
	// code, err := keeper.GetCode(ctx, codeHash)
	// if err != nil {
	// 	panic("cannot get codehash")
	// }

	// Get raw params info
	// rawParamsInfo, err := wasm.ParamsInfo(code.Code)
	// var paramsInfo []types.Field
	// if err != nil {
	// 	paramsInfo = nil
	// } else {
	// 	paramsInfo, err = types.ParseFields(rawParamsInfo)
	// 	if err != nil {
	// 		paramsInfo = nil
	// 	}
	// }

	// Get raw data sources info
	// rawDataInfo, err := wasm.RawDataInfo(code.Code)
	// var dataInfo []types.Field
	// if err != nil {
	// 	dataInfo = nil
	// } else {
	// 	dataInfo, err = types.ParseFields(rawDataInfo)
	// 	if err != nil {
	// 		dataInfo = nil
	// 	}
	// }

	// Get result info
	// rawResultInfo, err := wasm.ResultInfo(code.Code)
	// var resultInfo []types.Field
	// if err != nil {
	// 	resultInfo = nil
	// } else {
	// 	resultInfo, err = types.ParseFields(rawResultInfo)
	// 	if err != nil {
	// 		resultInfo = nil
	// 	}
	// }

	return codec.MustMarshalJSONIndent(keeper.cdc, scriptInfo), nil
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

// serializeParams is a function that receive codeHash and params and return a serialized format of params
func serializeParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) != 2 {
		return nil, sdk.ErrInternal(fmt.Sprintf("number of arguments must be 2, but got %d", len(path)))
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

	rawParams, err := wasm.SerializeParams(code.Code, []byte(path[1]))

	return codec.MustMarshalJSONIndent(keeper.cdc, rawParams), nil
}

func queryRequestNumber(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	return codec.MustMarshalJSONIndent(keeper.cdc, keeper.GetRequestCount(ctx)), nil
}
