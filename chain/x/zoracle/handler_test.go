package zoracle

import (
	"testing"
	"time"

	keep "github.com/bandprotocol/d3n/chain/x/zoracle/internal/keeper"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func mockDataSource(ctx sdk.Context, keeper Keeper) sdk.Result {
	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source_1"
	description := "description"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	executable := []byte("executable")
	sender := sdk.AccAddress([]byte("sender"))
	msg := types.NewMsgCreateDataSource(owner, name, description, fee, executable, sender)
	return handleMsgCreateDataSource(ctx, keeper, msg)
}

func mockOracleScript(ctx sdk.Context, keeper Keeper) sdk.Result {
	owner := sdk.AccAddress([]byte("owner"))
	name := "oracle_script_1"
	description := "description"
	code := []byte("code")
	sender := sdk.AccAddress([]byte("sender"))
	msg := types.NewMsgCreateOracleScript(owner, name, description, code, sender)
	return handleMsgCreateOracleScript(ctx, keeper, msg)
}

func TestCreateDataSourceSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	got := mockDataSource(ctx, keeper)
	require.True(t, got.IsOK(), "expected set data source to be ok, got %v", got)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, sdk.AccAddress([]byte("owner")), dataSource.Owner)
	require.Equal(t, "data_source_1", dataSource.Name)
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 10)), dataSource.Fee)
	require.Equal(t, []byte("executable"), dataSource.Executable)
}

func TestEditDataSourceSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	mockDataSource(ctx, keeper)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "data_source_2"
	newDescription := "new_description"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 99))
	newExecutable := []byte("executable_2")
	sender := sdk.AccAddress([]byte("owner"))

	msg := types.NewMsgEditDataSource(1, newOwner, newName, newDescription, newFee, newExecutable, sender)
	got := handleMsgEditDataSource(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected edit data source to be ok, got %v", got)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, dataSource.Owner)
	require.Equal(t, newName, dataSource.Name)
	require.Equal(t, newDescription, dataSource.Description)
	require.Equal(t, newFee, dataSource.Fee)
	require.Equal(t, newExecutable, dataSource.Executable)
}

func TestEditDataSourceByNotOwner(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	mockDataSource(ctx, keeper)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "data_source_2"
	newDescription := "new_description"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 99))
	newExecutable := []byte("executable_2")
	sender := sdk.AccAddress([]byte("sender"))

	msg := types.NewMsgEditDataSource(1, newOwner, newName, newDescription, newFee, newExecutable, sender)
	got := handleMsgEditDataSource(ctx, keeper, msg)
	require.False(t, got.IsOK())
	require.Equal(t, types.CodeInvalidOwner, got.Code)
}

func TestCreateOracleScriptSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	got := mockOracleScript(ctx, keeper)
	require.True(t, got.IsOK(), "expected set oracle script to be ok, got %v", got)

	expect, err := keeper.GetOracleScript(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, sdk.AccAddress([]byte("owner")), expect.Owner)
	require.Equal(t, "oracle_script_1", expect.Name)
	require.Equal(t, []byte("code"), expect.Code)
}

func TestEditOracleScriptSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	mockOracleScript(ctx, keeper)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "oracle_script_2"
	newDescription := "description_2"
	newCode := []byte("code_2")
	sender := sdk.AccAddress([]byte("owner"))

	msg := types.NewMsgEditOracleScript(1, newOwner, newName, newDescription, newCode, sender)
	got := handleMsgEditOracleScript(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected edit oracle script to be ok, got %v", got)

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, oracleScript.Owner)
	require.Equal(t, newName, oracleScript.Name)
	require.Equal(t, newCode, oracleScript.Code)
}

func TestEditOracleScriptByNotOwner(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	mockOracleScript(ctx, keeper)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "data_source_2"
	newDescription := "description_2"
	newCode := []byte("code_2")
	sender := sdk.AccAddress([]byte("not_owner"))

	msg := types.NewMsgEditOracleScript(1, newOwner, newName, newDescription, newCode, sender)
	got := handleMsgEditOracleScript(ctx, keeper, msg)
	require.False(t, got.IsOK())
	require.Equal(t, types.CodeInvalidOwner, got.Code)
}

func TestRequestSuccess(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 1000000, 1000000, sender)

	// Test here
	beforeGas := ctx.GasMeter().GasConsumed()
	got := handleMsgRequestData(ctx, keeper, msg)
	afterGas := ctx.GasMeter().GasConsumed()
	require.True(t, got.IsOK(), "expected request to be ok, got %v", got)

	// Check global request count
	require.Equal(t, int64(1), keeper.GetRequestCount(ctx))
	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	expectRequest := types.NewRequest(1, calldata,
		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
		2, 1581589790, 102, 1000000,
	)
	expectRequest.ExecuteGas = 1000000
	require.Equal(t, expectRequest, actualRequest)

	require.Equal(t, int64(1), keeper.GetRawDataRequestCount(ctx, 1))

	rawRequests := []types.RawDataRequest{
		types.NewRawDataRequest(1, []byte("band-protocol")),
	}
	require.Equal(t, rawRequests, keeper.GetRawDataRequests(ctx, 1))
	// check consumed gas must more than 2000000 (prepareGas + executeGas)
	require.True(t, afterGas-beforeGas > 2000000)
}

func TestRequestInvalidDataSource(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 30, 20000, sender)
	got := handleMsgRequestData(ctx, keeper, msg)
	require.False(t, got.IsOK())

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	got = handleMsgRequestData(ctx, keeper, msg)
	require.False(t, got.IsOK())
}

func TestRequestWithPrepareGasExceed(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	// set prepare gas to 3 (not enough for using) then it occurs error.
	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 3, 1000000, sender)

	got := handleMsgRequestData(ctx, keeper, msg)
	require.False(t, got.IsOK())
}

func TestReportSuccess(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	request := types.NewRequest(1, calldata,
		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
		2, 1581589790, 102, 1000000,
	)
	keeper.SetRequest(ctx, 1, request)
	keeper.SetRawDataRequest(ctx, 1, 42, types.NewRawDataRequest(1, []byte("calldata1")))

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	msg := types.NewMsgReportData(1, sdk.NewDecCoins(sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(0)))), []types.RawDataReport{
		types.NewRawDataReport(42, []byte("data1")),
	}, validatorAddress1)

	got := handleMsgReportData(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected report to be ok, got %v", got)
	list := keeper.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{}, list)

	msg = types.NewMsgReportData(1, sdk.NewDecCoins(sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(0)))), []types.RawDataReport{
		types.NewRawDataReport(42, []byte("data2")),
	}, validatorAddress2)

	got = handleMsgReportData(ctx, keeper, msg)
	require.True(t, got.IsOK(), "expected report to be ok, got %v", got)

	list = keeper.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{1}, list)
}

func TestReportFailed(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)
	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	request := types.NewRequest(1, calldata,
		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
		2, 1581589790, 102, 1000000,
	)
	keeper.SetRequest(ctx, 1, request)
	keeper.SetRawDataRequest(ctx, 1, 42, types.NewRawDataRequest(1, []byte("calldata1")))

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	msg := types.NewMsgReportData(1, sdk.NewDecCoins(sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(0)))), []types.RawDataReport{
		types.NewRawDataReport(41, []byte("data1")),
	}, validatorAddress1)

	// Test only 1 failed case, other case tested in keeper/report_test.go
	got := handleMsgReportData(ctx, keeper, msg)
	require.False(t, got.IsOK())
}

func TestEndBlock(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 30, 2500, sender)

	handleMsgRequestData(ctx, keeper, msg)

	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress1, []byte("answer1"))
	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress2, []byte("answer2"))

	keeper.SetPendingResolveList(ctx, []types.RequestID{1})

	ctx = ctx.WithBlockTime(time.Unix(int64(1581589999), 0))
	got := handleEndBlock(ctx, keeper)
	require.True(t, got.IsOK(), "expected set request to be ok, got %v", got)

	require.Equal(t, []types.RequestID{}, keeper.GetPendingResolveList(ctx))

	result, err := keeper.GetResult(ctx, 1, 1, calldata)
	require.Nil(t, err)
	require.Equal(t,
		types.Result{
			RequestTime:              1581589790,
			AggregationTime:          1581589999,
			RequestedValidatorsCount: 2,
			SufficientValidatorCount: 2,
			ReportedValidatorsCount:  0,
			Data:                     []byte("answer2"),
		},
		result,
	)

	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, types.Success, actualRequest.ResolveStatus)
}

func TestEndBlockExecuteFailedIfExecuteGasLessThanGasUsed(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	// Set gas for execution to 500
	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 500, 500, sender)

	handleMsgRequestData(ctx, keeper, msg)

	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress1, []byte("answer1"))
	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress2, []byte("answer2"))

	keeper.SetPendingResolveList(ctx, []types.RequestID{1})

	got := handleEndBlock(ctx, keeper)
	require.True(t, got.IsOK(), "expected set request to be ok, got %v", got)

	require.Equal(t, []types.RequestID{}, keeper.GetPendingResolveList(ctx))

	_, err := keeper.GetResult(ctx, 1, 1, calldata)
	require.NotNil(t, err)

	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, types.Failure, actualRequest.ResolveStatus)
}

func TestSkipInvalidExecuteGas(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	// Set gas for execution to 100000
	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 1000000, 100000, sender)
	handleMsgRequestData(ctx, keeper, msg)

	msg = types.NewMsgRequestData(1, calldata, 2, 2, 100, 1000000, 50000, sender)
	handleMsgRequestData(ctx, keeper, msg)

	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress1, []byte("answer1"))
	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress2, []byte("answer2"))

	keeper.SetRawDataReport(ctx, 2, 1, validatorAddress1, []byte("answer1"))
	keeper.SetRawDataReport(ctx, 2, 1, validatorAddress2, []byte("answer2"))

	keeper.SetEndBlockExecuteGasLimit(ctx, 75000)

	keeper.SetPendingResolveList(ctx, []types.RequestID{1, 2})
	got := handleEndBlock(ctx, keeper)
	require.True(t, got.IsOK(), "expected set request to be ok, got %v", got)
	require.Equal(t, []types.RequestID{}, keeper.GetPendingResolveList(ctx))

	_, err := keeper.GetResult(ctx, 1, 1, calldata)
	require.NotNil(t, err)

	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, types.Open, actualRequest.ResolveStatus)

	_, err = keeper.GetResult(ctx, 2, 1, calldata)
	require.Nil(t, err)

	actualRequest, err = keeper.GetRequest(ctx, 2)
	require.Nil(t, err)
	require.Equal(t, types.Success, actualRequest.ResolveStatus)
}

func TestStopResolveWhenOutOfGas(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	scriptID := types.OracleScriptID(1)
	keeper.SetOracleScript(ctx, scriptID, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	pendingList := []types.RequestID{}
	for i := types.RequestID(1); i <= types.RequestID(10); i++ {
		handleMsgRequestData(
			ctx, keeper,
			types.NewMsgRequestData(scriptID, calldata, 2, 2, 100, 2000, 2500, sender),
		)

		keeper.SetRawDataReport(ctx, i, 1, validatorAddress1, []byte("answer1"))
		keeper.SetRawDataReport(ctx, i, 1, validatorAddress2, []byte("answer2"))
		pendingList = append(pendingList, i)
	}

	// Each execute use 2270 gas, so it can resolve 3 requests per block
	keeper.SetEndBlockExecuteGasLimit(ctx, 7500)
	keeper.SetPendingResolveList(ctx, pendingList)

	got := handleEndBlock(ctx, keeper)
	require.True(t, got.IsOK(), "expected set request to be ok, got %v", got)
	require.Equal(t, []types.RequestID{4, 5, 6, 7, 8, 9, 10}, keeper.GetPendingResolveList(ctx))

	for i := types.RequestID(1); i <= types.RequestID(3); i++ {
		_, err := keeper.GetResult(ctx, i, scriptID, calldata)
		require.Nil(t, err)

		actualRequest, err := keeper.GetRequest(ctx, i)
		require.Nil(t, err)
		require.Equal(t, types.Success, actualRequest.ResolveStatus)
	}

	for i := types.RequestID(4); i <= types.RequestID(10); i++ {
		_, err := keeper.GetResult(ctx, i, scriptID, calldata)
		require.NotNil(t, err)

		actualRequest, err := keeper.GetRequest(ctx, i)
		require.Nil(t, err)
		require.Equal(t, types.Open, actualRequest.ResolveStatus)
	}

	got = handleEndBlock(ctx, keeper)
	require.True(t, got.IsOK(), "expected set request to be ok, got %v", got)
	require.Equal(t, []types.RequestID{7, 8, 9, 10}, keeper.GetPendingResolveList(ctx))

	for i := types.RequestID(1); i <= types.RequestID(6); i++ {
		_, err := keeper.GetResult(ctx, i, scriptID, calldata)
		require.Nil(t, err)

		actualRequest, err := keeper.GetRequest(ctx, i)
		require.Nil(t, err)
		require.Equal(t, types.Success, actualRequest.ResolveStatus)
	}

	for i := types.RequestID(7); i <= types.RequestID(10); i++ {
		_, err := keeper.GetResult(ctx, i, scriptID, calldata)
		require.NotNil(t, err)

		actualRequest, err := keeper.GetRequest(ctx, i)
		require.Nil(t, err)
		require.Equal(t, types.Open, actualRequest.ResolveStatus)
	}

	// New request
	handleMsgRequestData(
		ctx, keeper,
		types.NewMsgRequestData(scriptID, calldata, 2, 2, 100, 2000, 2500, sender),
	)

	got = handleEndBlock(ctx, keeper)
	require.True(t, got.IsOK(), "expected set request to be ok, got %v", got)
	require.Equal(t, []types.RequestID{10}, keeper.GetPendingResolveList(ctx))

	for i := types.RequestID(1); i <= types.RequestID(9); i++ {
		_, err := keeper.GetResult(ctx, i, scriptID, calldata)
		require.Nil(t, err)

		actualRequest, err := keeper.GetRequest(ctx, i)
		require.Nil(t, err)
		require.Equal(t, types.Success, actualRequest.ResolveStatus)
	}

	for i := types.RequestID(10); i <= types.RequestID(10); i++ {
		_, err := keeper.GetResult(ctx, i, scriptID, calldata)
		require.NotNil(t, err)

		actualRequest, err := keeper.GetRequest(ctx, i)
		require.Nil(t, err)
		require.Equal(t, types.Open, actualRequest.ResolveStatus)
	}

	keeper.SetRawDataReport(ctx, 11, 1, validatorAddress1, []byte("answer1"))
	keeper.SetRawDataReport(ctx, 11, 1, validatorAddress2, []byte("answer2"))
	keeper.SetPendingResolveList(ctx, []types.RequestID{10, 11})

	got = handleEndBlock(ctx, keeper)
	require.True(t, got.IsOK(), "expected set request to be ok, got %v", got)
	require.Equal(t, []types.RequestID{}, keeper.GetPendingResolveList(ctx))

	for i := types.RequestID(1); i <= types.RequestID(11); i++ {
		_, err := keeper.GetResult(ctx, i, scriptID, calldata)
		require.Nil(t, err)

		actualRequest, err := keeper.GetRequest(ctx, i)
		require.Nil(t, err)
		require.Equal(t, types.Success, actualRequest.ResolveStatus)
	}
}

func TestEndBlockInsufficientExecutionConsumeEndBlockGas(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	scriptID := types.OracleScriptID(1)
	keeper.SetOracleScript(ctx, scriptID, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	pendingList := []types.RequestID{}
	executeGasList := []uint64{2500, 50, 3000}

	for i := types.RequestID(1); i <= types.RequestID(3); i++ {
		handleMsgRequestData(
			ctx, keeper,
			types.NewMsgRequestData(scriptID, calldata, 2, 2, 100, 2000, executeGasList[i-1], sender),
		)

		keeper.SetRawDataReport(ctx, i, 1, validatorAddress1, []byte("answer1"))
		keeper.SetRawDataReport(ctx, i, 1, validatorAddress2, []byte("answer2"))
		pendingList = append(pendingList, i)
	}

	keeper.SetEndBlockExecuteGasLimit(ctx, 2600)
	keeper.SetPendingResolveList(ctx, pendingList)

	got := handleEndBlock(ctx, keeper)
	require.True(t, got.IsOK(), "expected set request to be ok, got %v", got)
	require.Equal(t, []types.RequestID{}, keeper.GetPendingResolveList(ctx))

	_, err := keeper.GetResult(ctx, 1, scriptID, calldata)
	require.Nil(t, err)

	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, types.Success, actualRequest.ResolveStatus)

	_, err = keeper.GetResult(ctx, 2, scriptID, calldata)
	require.NotNil(t, err)

	actualRequest, err = keeper.GetRequest(ctx, 2)
	require.Nil(t, err)
	require.Equal(t, types.Failure, actualRequest.ResolveStatus)

	_, err = keeper.GetResult(ctx, 3, scriptID, calldata)
	require.NotNil(t, err)

	actualRequest, err = keeper.GetRequest(ctx, 3)
	require.Nil(t, err)
	require.Equal(t, types.Open, actualRequest.ResolveStatus)

}
