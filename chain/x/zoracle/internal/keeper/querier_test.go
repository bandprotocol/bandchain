package keeper

// import (
// 	"encoding/hex"
// 	"path/filepath"
// 	"testing"

// 	"github.com/bandprotocol/d3n/chain/wasm"
// 	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
// 	"github.com/cosmos/cosmos-sdk/codec"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/require"
// 	abci "github.com/tendermint/tendermint/abci/types"
// )

// func TestQueryRequestById(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	// Create variable "querier" which is a function
// 	querier := NewQuerier(keeper)

// 	// query before set new request
// 	acsBytes, err := querier(
// 		ctx,
// 		[]string{"request", "1"},
// 		abci.RequestQuery{},
// 	)
// 	// It must return error request not found
// 	require.Equal(t, types.CodeRequestNotFound, err.Code())

// 	// set code
// 	absPath, _ := filepath.Abs("../../../../wasm/res/result.wasm")
// 	code, _ := wasm.ReadBytes(absPath)
// 	owner := sdk.AccAddress([]byte("owner"))
// 	name := "Crypto Price"
// 	codeHash := keeper.SetCode(ctx, code, name, owner)
// 	params, _ := hex.DecodeString("00000000")

// 	// set request

// 	request := types.NewRequest(codeHash, params, 3)
// 	keeper.SetRequest(ctx, 1, request)
// 	result, _ := hex.DecodeString("0000000000002710")
// 	keeper.SetResult(ctx, 1, codeHash, params, result)

// 	// create query
// 	acsBytes, err = querier(
// 		ctx,
// 		[]string{"request", "1"},
// 		abci.RequestQuery{},
// 	)
// 	require.Nil(t, err)

// 	paramsMap := []byte(`{"symbol":"BTC"}`)
// 	parsedResult := []byte(`{"price_in_usd": 10000}`)

// 	// Use bytes format for comparison
// 	request = types.NewRequest(codeHash, params, 3)
// 	acs, errJSON := codec.MarshalJSONIndent(
// 		keeper.cdc,
// 		types.NewRequestInfo(
// 			request.CodeHash,
// 			paramsMap,
// 			params,
// 			request.ReportEndAt,
// 			[]types.ReportWithValidator{},
// 			parsedResult,
// 			result,
// 		),
// 	)
// 	require.Nil(t, errJSON)
// 	require.Equal(t, acs, acsBytes)
// }

// func TestQueryPendingRequest(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)
// 	// Create variable "querier" which is a function
// 	querier := NewQuerier(keeper)

// 	// Read the state. The state should contain 0 pending request.
// 	acsBytes, err := querier(
// 		ctx,
// 		[]string{"pending_request"},
// 		abci.RequestQuery{},
// 	)
// 	require.Nil(t, err)
// 	// Use bytes format for comparison
// 	acs, errJSON := codec.MarshalJSONIndent(keeper.cdc, []uint64{})
// 	require.Nil(t, errJSON)
// 	require.Equal(t, acs, acsBytes)

// 	// set request
// 	name := "Crypto price"
// 	owner := sdk.AccAddress([]byte("owner"))
// 	code := []byte("code")
// 	codeHash := keeper.SetCode(ctx, code, name, owner)
// 	request := types.NewRequest(codeHash, []byte("params"), 3)
// 	keeper.SetRequest(ctx, 2, request)

// 	// set pending
// 	pendingRequests := keeper.GetUnresolvedRequests(ctx)
// 	pendingRequests = append(pendingRequests, 2)
// 	keeper.SetUnresolvedRequests(ctx, pendingRequests)

// 	// Read the state agian. The state should contain 1 pending request. That is reqID = 2.
// 	acsBytes, err = querier(
// 		ctx,
// 		[]string{"pending_request"},
// 		abci.RequestQuery{},
// 	)
// 	require.Nil(t, err)
// 	// Use bytes format for comparison
// 	acs, errJSON = codec.MarshalJSONIndent(keeper.cdc, []uint64{2})
// 	require.Nil(t, errJSON)
// 	require.Equal(t, acs, acsBytes)
// }

// func TestQueryScript(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	absPath, _ := filepath.Abs("../../../../wasm/res/result.wasm")
// 	code, _ := wasm.ReadBytes(absPath)
// 	owner := sdk.AccAddress([]byte("owner"))
// 	name := "Crypto Price"
// 	codeHash := keeper.SetCode(ctx, code, name, owner)

// 	// Create variable "querier" which is a function
// 	querier := NewQuerier(keeper)

// 	rawQueryBytes, err := querier(
// 		ctx,
// 		[]string{"script", hex.EncodeToString(codeHash)},
// 		abci.RequestQuery{},
// 	)
// 	require.Nil(t, err)

// 	expectJson, errJSON := codec.MarshalJSONIndent(
// 		keeper.cdc,
// 		types.NewScriptInfo(
// 			name,
// 			codeHash,
// 			[]types.Field{
// 				types.Field{Name: "symbol", Type: "coins::Coins"},
// 			},
// 			[]types.Field{
// 				types.Field{Name: "coin_gecko", Type: "f32"},
// 				types.Field{Name: "crypto_compare", Type: "f32"},
// 			},
// 			[]types.Field{types.Field{Name: "price_in_usd", Type: "u64"}},
// 			owner,
// 		),
// 	)
// 	require.Nil(t, errJSON)
// 	require.Equal(t, expectJson, rawQueryBytes)
// }

// func TestSerializeParams(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	absPath, _ := filepath.Abs("../../../../wasm/res/serialized_params.wasm")
// 	code, _ := wasm.ReadBytes(absPath)
// 	owner := sdk.AccAddress([]byte("owner"))
// 	name := "Crypto Price"
// 	codeHash := keeper.SetCode(ctx, code, name, owner)

// 	// Create variable "querier" which is a function
// 	querier := NewQuerier(keeper)

// 	rawQueryBytes, err := querier(
// 		ctx,
// 		[]string{"serialize_params", hex.EncodeToString(codeHash), `{"crypto_symbol":"ETH", "stock_symbol":"GOOG","alphavantage_api_key":"WVKPOO76169EX950"}`},
// 		abci.RequestQuery{},
// 	)
// 	require.Nil(t, err)

// 	expectBytes, _ := hex.DecodeString("000000010000000000000004474f4f47000000000000001057564b504f4f37363136394558393530")
// 	expectJson, errJSON := codec.MarshalJSONIndent(keeper.cdc, expectBytes)

// 	require.Nil(t, errJSON)
// 	require.Equal(t, expectJson, rawQueryBytes)
// }
