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

	keeper.SetDataSource(ctx, 1, types.NewDataSource(owner, name, fee, executable))

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
			1,
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
			1,
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

func TestQueryRequests(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	// query before set new request
	acsBytes, err := querier(
		ctx,
		[]string{"requests", "1", "3"},
		abci.RequestQuery{},
	)
	// It must return empty array of requests
	require.Nil(t, err)
	require.Equal(t, []byte("[]"), acsBytes)

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
	keeper.GetNextRequestID(ctx)

	// request 2
	keeper.SetRequest(ctx, 2, request)

	keeper.SetRawDataRequest(ctx, 2, 100, types.NewRawDataRequest(1, []byte("only calldata")))
	keeper.GetNextRequestID(ctx)

	// create query
	acsBytes, err = querier(
		ctx,
		[]string{"requests", "1", "3"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	// Use bytes format for comparison
	acs, errJSON := codec.MarshalJSONIndent(
		keeper.cdc,
		[]types.RequestQuerierInfo{
			types.NewRequestQuerierInfo(
				1,
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
			types.NewRequestQuerierInfo(
				2,
				request,
				[]types.RawDataRequestWithExternalID{
					types.NewRawDataRequestWithExternalID(
						100,
						types.NewRawDataRequest(1, []byte("only calldata")),
					),
				},
				[]types.ReportWithValidator{},
				nil,
			),
		},
	)
	require.Nil(t, errJSON)
	require.Equal(t, acs, acsBytes)
}

func TestQueryOracleScriptById(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxOracleScriptCodeSize(ctx, 20)
	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	// query before add a oracle script
	_, err := querier(
		ctx,
		[]string{"oracle_script", "1"},
		abci.RequestQuery{},
	)
	// Should return error oracle script not found
	require.NotNil(t, err)

	owner := sdk.AccAddress([]byte("owner"))
	name := "oracle_script"
	code := []byte("code")
	expectedResult := types.NewOracleScriptQuerierInfo(1, owner, name, code)

	keeper.SetOracleScript(ctx, 1, types.NewOracleScript(owner, name, code))

	// This time querier should be able to find a oracle script
	oracleScript, err := querier(
		ctx,
		[]string{"oracle_script", "1"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON := codec.MarshalJSONIndent(keeper.cdc, expectedResult)
	require.Nil(t, errJSON)

	require.Equal(t, expectedResultBytes, oracleScript)
}

func TestQueryOracleScriptsByStartIdAndNumberOfOracleScripts(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxOracleScriptCodeSize(ctx, 20)
	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	expectedResult := []types.OracleScriptQuerierInfo{}

	// Add a new 10 oracle scripts
	for i := 1; i <= 10; i++ {
		owner := sdk.AccAddress([]byte("owner" + strconv.Itoa(i)))
		name := "oracle_script_" + strconv.Itoa(i)
		code := []byte("code" + strconv.Itoa(i))
		eachOracleScript := types.NewOracleScriptQuerierInfo(int64(i), owner, name, code)

		err := keeper.AddOracleScript(ctx, eachOracleScript.Owner, eachOracleScript.Name, eachOracleScript.Code)
		require.Nil(t, err)

		expectedResult = append(expectedResult, eachOracleScript)
	}

	// Query first 5 oracle scripts
	oracleScripts, err := querier(
		ctx,
		[]string{"oracle_scripts", "1", "5"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON := codec.MarshalJSONIndent(keeper.cdc, expectedResult[0:5])
	require.Nil(t, errJSON)
	require.Equal(t, expectedResultBytes, oracleScripts)

	// Query last 5 oracle scripts
	oracleScripts, err = querier(
		ctx,
		[]string{"oracle_scripts", "6", "5"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON = codec.MarshalJSONIndent(keeper.cdc, expectedResult[5:])
	require.Nil(t, errJSON)
	require.Equal(t, expectedResultBytes, oracleScripts)

	// Query first 15 oracle scripts which exceed number of all oracle script right now
	// This should return all exist oracle scripts (10 oracle scripts)
	oracleScripts, err = querier(
		ctx,
		[]string{"oracle_scripts", "1", "15"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON = codec.MarshalJSONIndent(keeper.cdc, expectedResult)
	require.Nil(t, errJSON)
	require.Equal(t, expectedResultBytes, oracleScripts)

	// Query oracle scripts from id=8 to id=17
	// But we only have id=1 to id=10
	// So the result should be [id=8, id=9, id=10]
	oracleScripts, err = querier(
		ctx,
		[]string{"oracle_scripts", "8", "10"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON = codec.MarshalJSONIndent(keeper.cdc, expectedResult[7:])
	require.Nil(t, errJSON)
	require.Equal(t, expectedResultBytes, oracleScripts)
}

func TestQueryOracleScriptsGotEmptyArrayBecauseNoOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxOracleScriptCodeSize(ctx, 20)
	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	oracleScripts, err := querier(
		ctx,
		[]string{"oracle_scripts", "1", "5"},
		abci.RequestQuery{},
	)
	require.Nil(t, err)

	expectedResultBytes, errJSON := codec.MarshalJSONIndent(keeper.cdc, []types.OracleScriptQuerierInfo{})
	require.Nil(t, errJSON)

	require.Equal(t, expectedResultBytes, oracleScripts)
}

func TestQueryOracleScriptsFailBecauseInvalidNumberOfOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxOracleScriptCodeSize(ctx, 20)
	// Create variable "querier" which is a function
	querier := NewQuerier(keeper)

	// Number of oracle scripts should <= 100
	_, err := querier(
		ctx,
		[]string{"oracle_scripts", "1", "101"},
		abci.RequestQuery{},
	)
	require.NotNil(t, err)

	// Number of oracle scripts should >= 1
	_, err = querier(
		ctx,
		[]string{"oracle_scripts", "1", "0"},
		abci.RequestQuery{},
	)
	require.NotNil(t, err)
}
