package keeper

import (
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryRequestById(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	// query before set new request
	acsBytes, err := querier(
		ctx,
		[]string{"request", "1"},
		abci.RequestQuery{},
	)
	// It must return error request not found
	require.Equal(t, types.CodeRequestNotFound, err.Code())

	request := newDefaultRequest()
	keeper.SetRequest(ctx, 1, request)

	keeper.SetRawDataRequest(ctx, 1, 1, types.NewRawDataRequest(0, []byte("calldata1")))
	keeper.SetRawDataRequest(ctx, 1, 2, types.NewRawDataRequest(1, []byte("calldata2")))

	keeper.SetRawDataReport(ctx, 1, 1, request.RequestedValidators[0], []byte("report1"))
	keeper.SetRawDataReport(ctx, 1, 2, request.RequestedValidators[0], []byte("report2"))

	keeper.SetRawDataReport(ctx, 1, 1, request.RequestedValidators[1], []byte("report1-2"))
	keeper.SetRawDataReport(ctx, 1, 2, request.RequestedValidators[1], []byte("report2-2"))

	result, _ := hex.DecodeString("0000000000002710")
	keeper.SetResult(ctx, 1, request.OracleScriptID, request.Calldata, result)

	// create query
	acsBytes, err = querier(
		ctx,
		[]string{"request", "1"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	// Use bytes format for comparison
	acs, errJSON := codec.MarshalJSONIndent(
		keeper.cdc,
		types.NewRequestQuerierInfo(
			request,
			[]types.RawDataRequestWithExternalID{
				types.NewRawDataRequestWithExternalID(
					1,
					types.NewRawDataRequest(0, []byte("calldata1")),
				),
				types.NewRawDataRequestWithExternalID(
					2,
					types.NewRawDataRequest(1, []byte("calldata2")),
				),
			},
			[]types.ReportWithValidator{
				types.NewReportWithValidator([]types.RawDataReport{
					types.NewRawDataReport(1, []byte("report1")),
					types.NewRawDataReport(2, []byte("report2")),
				}, request.RequestedValidators[0]),
				types.NewReportWithValidator([]types.RawDataReport{
					types.NewRawDataReport(1, []byte("report1-2")),
					types.NewRawDataReport(2, []byte("report2-2")),
				}, request.RequestedValidators[1]),
			},
			result,
		),
	)
	require.Nil(t, errJSON)
	require.Equal(t, acs, acsBytes)
}

func TestQueryRequestIncompleteValidator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	// query before set new request
	acsBytes, err := querier(
		ctx,
		[]string{"request", "1"},
		abci.RequestQuery{},
	)
	// It must return error request not found
	require.Equal(t, types.CodeRequestNotFound, err.Code())

	request := newDefaultRequest()
	keeper.SetRequest(ctx, 1, request)

	keeper.SetRawDataRequest(ctx, 1, 1, types.NewRawDataRequest(0, []byte("calldata1")))
	keeper.SetRawDataRequest(ctx, 1, 2, types.NewRawDataRequest(1, []byte("calldata2")))

	keeper.SetRawDataReport(ctx, 1, 1, request.RequestedValidators[1], []byte("report1-2"))
	keeper.SetRawDataReport(ctx, 1, 2, request.RequestedValidators[1], []byte("report2-2"))

	// create query
	acsBytes, err = querier(
		ctx,
		[]string{"request", "1"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	// Use bytes format for comparison
	acs, errJSON := codec.MarshalJSONIndent(
		keeper.cdc,
		types.NewRequestQuerierInfo(
			request,
			[]types.RawDataRequestWithExternalID{
				types.NewRawDataRequestWithExternalID(
					1,
					types.NewRawDataRequest(0, []byte("calldata1")),
				),
				types.NewRawDataRequestWithExternalID(
					2,
					types.NewRawDataRequest(1, []byte("calldata2")),
				),
			},
			[]types.ReportWithValidator{
				types.NewReportWithValidator([]types.RawDataReport{
					types.NewRawDataReport(1, []byte("report1-2")),
					types.NewRawDataReport(2, []byte("report2-2")),
				}, request.RequestedValidators[1]),
			},
			nil,
		),
	)
	require.Nil(t, errJSON)
	require.Equal(t, acs, acsBytes)
}

func TestQueryDataSourceById(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxDataSourceExecutableSize(ctx, 20)
	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	// query before add a data source
	_, err := querier(
		ctx,
		[]string{"data_source", "1"},
		abci.RequestQuery{},
	)
	// Should return error data source not found
	require.NotNil(t, err)

	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	executable := []byte("executable")
	expectedResult := types.NewDataSourceQuerierInfo(1, owner, name, fee, executable)

	// Add a new data source
	err = keeper.AddDataSource(ctx, expectedResult.Owner, expectedResult.Name, expectedResult.Fee, expectedResult.Executable)
	require.Nil(t, err)

	// This time querier should be able to find a data source
	dataSource, err := querier(
		ctx,
		[]string{"data_source", "1"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON := codec.MarshalJSONIndent(keeper.cdc, expectedResult)
	require.Nil(t, errJSON)

	require.Equal(t, expectedResultBytes, dataSource)
}

func TestQueryDataSourcesByStartIdAndNumberOfDataSources(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxDataSourceExecutableSize(ctx, 20)
	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	expectedResult := []types.DataSourceQuerierInfo{}

	// Add a new 10 data sources
	for i := 1; i <= 10; i++ {
		owner := sdk.AccAddress([]byte("owner" + strconv.Itoa(i)))
		name := "data_source_" + strconv.Itoa(i)
		fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
		executable := []byte("executable" + strconv.Itoa(i))
		eachDataSource := types.NewDataSourceQuerierInfo(int64(i), owner, name, fee, executable)

		err := keeper.AddDataSource(ctx, eachDataSource.Owner, eachDataSource.Name, eachDataSource.Fee, eachDataSource.Executable)
		require.Nil(t, err)

		expectedResult = append(expectedResult, eachDataSource)
	}

	// Query first 5 data sources
	dataSources, err := querier(
		ctx,
		[]string{"data_sources", "1", "5"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON := codec.MarshalJSONIndent(keeper.cdc, expectedResult[0:5])
	require.Nil(t, errJSON)
	require.Equal(t, expectedResultBytes, dataSources)

	// Query last 5 data sources
	dataSources, err = querier(
		ctx,
		[]string{"data_sources", "6", "5"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON = codec.MarshalJSONIndent(keeper.cdc, expectedResult[5:])
	require.Nil(t, errJSON)
	require.Equal(t, expectedResultBytes, dataSources)

	// Query first 15 data sources which exceed number of all data source right now
	// This should return all exist data sources (10 data sources)
	dataSources, err = querier(
		ctx,
		[]string{"data_sources", "1", "15"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON = codec.MarshalJSONIndent(keeper.cdc, expectedResult)
	require.Nil(t, errJSON)
	require.Equal(t, expectedResultBytes, dataSources)

	// Query data sources from id=8 to id=17
	// But we only have id=1 to id=10
	// So the result should be [id=8, id=9, id=10]
	dataSources, err = querier(
		ctx,
		[]string{"data_sources", "8", "10"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON = codec.MarshalJSONIndent(keeper.cdc, expectedResult[7:])
	require.Nil(t, errJSON)
	require.Equal(t, expectedResultBytes, dataSources)
}

func TestQueryDataSourcesGotEmptyArrayBecauseNoDataSource(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxDataSourceExecutableSize(ctx, 20)
	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	dataSources, err := querier(
		ctx,
		[]string{"data_sources", "1", "5"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON := codec.MarshalJSONIndent(keeper.cdc, []types.DataSourceQuerierInfo{})
	require.Nil(t, errJSON)

	require.Equal(t, expectedResultBytes, dataSources)
}

func TestQueryDataSourcesFailBecauseInvalidNumberOfDataSource(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxDataSourceExecutableSize(ctx, 20)
	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	// Number of data sources should <= 100
	_, err := querier(
		ctx,
		[]string{"data_sources", "1", "101"},
		abci.RequestQuery{},
	)
	require.NotNil(t, err)

	// Number of data sources should >= 1
	_, err = querier(
		ctx,
		[]string{"data_sources", "1", "0"},
		abci.RequestQuery{},
	)
	require.NotNil(t, err)
}

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
// 	pendingRequests := keeper.GetPendingResolveList(ctx)
// 	pendingRequests = append(pendingRequests, 2)
// 	keeper.SetPendingResolveList(ctx, pendingRequests)

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
