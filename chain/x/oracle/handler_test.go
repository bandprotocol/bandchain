package oracle

import (
	"testing"
	"time"

	keep "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmkv "github.com/tendermint/tendermint/libs/kv"
)

func mockDataSource(ctx sdk.Context, keeper Keeper) (*sdk.Result, error) {
	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source_1"
	description := "description"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	executable := []byte("executable")
	sender := sdk.AccAddress([]byte("sender"))
	msg := types.NewMsgCreateDataSource(owner, name, description, fee, executable, sender)
	return handleMsgCreateDataSource(ctx, keeper, msg)
}

func mockOracleScript(ctx sdk.Context, keeper Keeper) (*sdk.Result, error) {
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

	_, err := mockDataSource(ctx, keeper)
	require.Nil(t, err)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, sdk.AccAddress([]byte("owner")), dataSource.Owner)
	require.Equal(t, "data_source_1", dataSource.Name)
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 10)), dataSource.Fee)
	require.Equal(t, []byte("executable"), dataSource.Executable)

	events := ctx.EventManager().Events()
	require.Equal(t, 1, len(events))
	require.Equal(t, sdk.Event{
		Type:       EventTypeCreateDataSource,
		Attributes: []tmkv.Pair{tmkv.Pair{Key: []byte(AttributeKeyID), Value: []byte("1")}},
	}, events[0])
}

func TestEditDataSourceSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	_, err := mockDataSource(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("anotherowner"))
	newName := "data_source_2"
	newDescription := "new_description"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 99))
	newExecutable := []byte("executable_2")
	sender := sdk.AccAddress([]byte("owner"))

	msg := types.NewMsgEditDataSource(1, newOwner, newName, newDescription, newFee, newExecutable, sender)
	_, err = handleMsgEditDataSource(ctx, keeper, msg)
	require.Nil(t, err)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, dataSource.Owner)
	require.Equal(t, newName, dataSource.Name)
	require.Equal(t, newDescription, dataSource.Description)
	require.Equal(t, newFee, dataSource.Fee)
	require.Equal(t, newExecutable, dataSource.Executable)

	events := ctx.EventManager().Events()
	require.Equal(t, 2, len(events))
	require.Equal(t, sdk.Event{
		Type:       EventTypeCreateDataSource,
		Attributes: []tmkv.Pair{tmkv.Pair{Key: []byte(AttributeKeyID), Value: []byte("1")}},
	}, events[0])
	require.Equal(t, sdk.Event{
		Type:       EventTypeEditDataSource,
		Attributes: []tmkv.Pair{tmkv.Pair{Key: []byte(AttributeKeyID), Value: []byte("1")}},
	}, events[1])
}

func TestEditDataSourceByNotOwner(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	mockDataSource(ctx, keeper)

	newOwner := sdk.AccAddress([]byte("anotherowner"))
	newName := "data_source_2"
	newDescription := "new_description"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 99))
	newExecutable := []byte("executable_2")
	sender := sdk.AccAddress([]byte("sender"))

	msg := types.NewMsgEditDataSource(1, newOwner, newName, newDescription, newFee, newExecutable, sender)
	_, err := handleMsgEditDataSource(ctx, keeper, msg)
	require.NotNil(t, err)
}

func TestCreateOracleScriptSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	_, err := mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	expect, err := keeper.GetOracleScript(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, sdk.AccAddress([]byte("owner")), expect.Owner)
	require.Equal(t, "oracle_script_1", expect.Name)
	require.Equal(t, []byte("code"), expect.Code)

	events := ctx.EventManager().Events()
	require.Equal(t, 1, len(events))
	require.Equal(t, sdk.Event{
		Type:       EventTypeCreateOracleScript,
		Attributes: []tmkv.Pair{tmkv.Pair{Key: []byte(AttributeKeyID), Value: []byte("1")}},
	}, events[0])
}

func TestEditOracleScriptSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	_, err := mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("anotherowner"))
	newName := "oracle_script_2"
	newDescription := "description_2"
	newCode := []byte("code_2")
	sender := sdk.AccAddress([]byte("owner"))

	msg := types.NewMsgEditOracleScript(1, newOwner, newName, newDescription, newCode, sender)
	_, err = handleMsgEditOracleScript(ctx, keeper, msg)
	require.Nil(t, err)

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, oracleScript.Owner)
	require.Equal(t, newName, oracleScript.Name)
	require.Equal(t, newCode, oracleScript.Code)

	events := ctx.EventManager().Events()
	require.Equal(t, 2, len(events))
	require.Equal(t, sdk.Event{
		Type:       EventTypeCreateOracleScript,
		Attributes: []tmkv.Pair{tmkv.Pair{Key: []byte(AttributeKeyID), Value: []byte("1")}},
	}, events[0])
	require.Equal(t, sdk.Event{
		Type:       EventTypeEditOracleScript,
		Attributes: []tmkv.Pair{tmkv.Pair{Key: []byte(AttributeKeyID), Value: []byte("1")}},
	}, events[1])
}

func TestEditOracleScriptByNotOwner(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	_, err := mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("anotherowner"))
	newName := "data_source_2"
	newDescription := "description_2"
	newCode := []byte("code_2")
	sender := sdk.AccAddress([]byte("not_owner"))

	msg := types.NewMsgEditOracleScript(1, newOwner, newName, newDescription, newCode, sender)
	_, err = handleMsgEditOracleScript(ctx, keeper, msg)
	require.NotNil(t, err)
}

func TestRequestSuccess(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))
	_, err := keeper.CoinKeeper.AddCoins(ctx, sender, keep.NewUBandCoins(410))
	require.Nil(t, err)

	owner := sdk.AccAddress([]byte("owner"))
	owner2 := sdk.AccAddress([]byte("anotherowner"))

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource1 := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource1)

	dataSource2 := types.NewDataSource(
		sdk.AccAddress([]byte("anotherowner")),
		"data_source2",
		"description2",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 400)),
		[]byte("executable2"),
	)
	keeper.SetDataSource(ctx, 2, dataSource2)

	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 1000000, 1000000, "clientID", sender)

	// Test here
	beforeGas := ctx.GasMeter().GasConsumed()
	_, err = handleMsgRequestData(ctx, keeper, msg)
	afterGas := ctx.GasMeter().GasConsumed()
	require.Nil(t, err)

	// Check global request count
	require.Equal(t, int64(1), keeper.GetRequestCount(ctx))
	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	expectRequest := types.NewRequest(1, calldata,
		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
		2, 1581589790, 102, 1000000, "clientID",
	)
	expectRequest.ExecuteGas = 1000000
	require.Equal(t, expectRequest, actualRequest)

	require.Equal(t, int64(2), keeper.GetRawDataRequestCount(ctx, 1))

	rawRequests := []types.RawDataRequest{
		types.NewRawDataRequest(1, []byte("band-protocol")), types.NewRawDataRequest(2, []byte("band-chain")),
	}
	require.Equal(t, rawRequests, keeper.GetRawDataRequests(ctx, 1))
	// check consumed gas must more than 2000000 (prepareGas + executeGas)
	require.True(t, afterGas-beforeGas > 2000000)

	senderBalance := keeper.CoinKeeper.GetAllBalances(ctx, sender)
	require.Equal(t, sdk.Coins(nil), senderBalance)

	ownerBalance := keeper.CoinKeeper.GetAllBalances(ctx, owner)
	require.Equal(t, keep.NewUBandCoins(10), ownerBalance)

	owner2Balance := keeper.CoinKeeper.GetAllBalances(ctx, owner2)
	require.Equal(t, keep.NewUBandCoins(400), owner2Balance)
}

func TestRequestInvalidDataSource(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 30, 20000, "clientID", sender)
	_, err := handleMsgRequestData(ctx, keeper, msg)
	require.NotNil(t, err)

	script := keep.GetTestOracleScript("../../owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	_, err = handleMsgRequestData(ctx, keeper, msg)
	require.NotNil(t, err)
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
	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 3, 1000000, "clientID", sender)

	_, err := handleMsgRequestData(ctx, keeper, msg)
	require.NotNil(t, err)
}

func TestRequestWithInsufficientFee(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))
	_, err := keeper.CoinKeeper.AddCoins(ctx, sender, keep.NewUBandCoins(50))
	require.Nil(t, err)

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

	dataSource2 := types.NewDataSource(
		sdk.AccAddress([]byte("anotherowner")),
		"data_source2",
		"description2",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 300)),
		[]byte("executable2"),
	)
	keeper.SetDataSource(ctx, 2, dataSource2)

	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 1000000, 1000000, "clientID", sender)

	_, err = handleMsgRequestData(ctx, keeper, msg)
	require.NotNil(t, err)
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
	reporterAddress1 := sdk.AccAddress(validatorAddress1)

	address1 := keep.GetAddressFromPub(pubStr[0])

	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)
	reporterAddress2 := sdk.AccAddress(validatorAddress2)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	_, err := keeper.CoinKeeper.AddCoins(ctx, address1, keep.NewUBandCoins(1000000))
	require.Nil(t, err)

	request := types.NewRequest(1, calldata,
		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
		2, 1581589790, 102, 1000000, "clientID",
	)
	keeper.SetRequest(ctx, 1, request)
	keeper.SetRawDataRequest(ctx, 1, 42, types.NewRawDataRequest(1, []byte("calldata1")))

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	msg := types.NewMsgReportData(1, []types.RawDataReportWithID{
		types.NewRawDataReportWithID(42, 0, []byte("data1")),
	}, validatorAddress1, reporterAddress1)

	_, err = handleMsgReportData(ctx, keeper, msg)
	require.Nil(t, err)
	list := keeper.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{}, list)

	msg = types.NewMsgReportData(1, []types.RawDataReportWithID{
		types.NewRawDataReportWithID(42, 0, []byte("data2")),
	}, validatorAddress2, reporterAddress2)

	_, err = handleMsgReportData(ctx, keeper, msg)
	require.Nil(t, err)

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

	reporterAddress1 := sdk.AccAddress(validatorAddress1)

	address1 := keep.GetAddressFromPub(pubStr[0])
	_, err := keeper.CoinKeeper.AddCoins(ctx, address1, keep.NewUBandCoins(1000000))
	require.Nil(t, err)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	request := types.NewRequest(1, calldata,
		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
		2, 1581589790, 102, 1000000, "clientID",
	)
	keeper.SetRequest(ctx, 1, request)
	keeper.SetRawDataRequest(ctx, 1, 42, types.NewRawDataRequest(1, []byte("calldata1")))

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	msg := types.NewMsgReportData(1, []types.RawDataReportWithID{
		types.NewRawDataReportWithID(41, 0, []byte("data1")),
	}, validatorAddress1, reporterAddress1)

	// Test only 1 failed case, other case tested in keeper/report_test.go
	_, err = handleMsgReportData(ctx, keeper, msg)
	require.NotNil(t, err)
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

	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 30, 2500, "clientID", sender)

	handleMsgRequestData(ctx, keeper, msg)

	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress1, types.NewRawDataReport(0, []byte("answer1")))
	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress2, types.NewRawDataReport(0, []byte("answer2")))

	keeper.SetPendingResolveList(ctx, []types.RequestID{1})

	ctx = ctx.WithBlockTime(time.Unix(int64(1581589999), 0))
	handleEndBlock(ctx, keeper)

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
	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 500, 500, "clientID", sender)

	handleMsgRequestData(ctx, keeper, msg)

	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress1, types.NewRawDataReport(0, []byte("answer1")))
	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress2, types.NewRawDataReport(0, []byte("answer2")))

	keeper.SetPendingResolveList(ctx, []types.RequestID{1})

	handleEndBlock(ctx, keeper)

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
	msg := types.NewMsgRequestData(1, calldata, 2, 2, 100, 1000000, 100000, "clientID", sender)
	handleMsgRequestData(ctx, keeper, msg)

	msg = types.NewMsgRequestData(1, calldata, 2, 2, 100, 1000000, 50000, "clientID", sender)
	handleMsgRequestData(ctx, keeper, msg)

	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress1, types.NewRawDataReport(0, []byte("answer1")))
	keeper.SetRawDataReport(ctx, 1, 1, validatorAddress2, types.NewRawDataReport(0, []byte("answer2")))

	keeper.SetRawDataReport(ctx, 2, 1, validatorAddress1, types.NewRawDataReport(0, []byte("answer1")))
	keeper.SetRawDataReport(ctx, 2, 1, validatorAddress2, types.NewRawDataReport(0, []byte("answer2")))

	keeper.SetParam(ctx, KeyEndBlockExecuteGasLimit, 75000)

	keeper.SetPendingResolveList(ctx, []types.RequestID{1, 2})
	handleEndBlock(ctx, keeper)
	require.Equal(t, []types.RequestID{}, keeper.GetPendingResolveList(ctx))

	_, err := keeper.GetResult(ctx, 1, 1, calldata)
	require.NotNil(t, err)

	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, types.Failure, actualRequest.ResolveStatus)

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
			types.NewMsgRequestData(scriptID, calldata, 2, 2, 100, 2000, 2500, "clientID", sender),
		)

		keeper.SetRawDataReport(ctx, i, 1, validatorAddress1, types.NewRawDataReport(0, []byte("answer1")))
		keeper.SetRawDataReport(ctx, i, 1, validatorAddress2, types.NewRawDataReport(0, []byte("answer2")))
		pendingList = append(pendingList, i)
	}

	// Each execute use 2270 gas, so it can resolve 3 requests per block
	keeper.SetParam(ctx, KeyEndBlockExecuteGasLimit, 7500)
	keeper.SetPendingResolveList(ctx, pendingList)

	handleEndBlock(ctx, keeper)
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

	handleEndBlock(ctx, keeper)
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
		types.NewMsgRequestData(scriptID, calldata, 2, 2, 100, 2000, 2500, "clientID", sender),
	)

	handleEndBlock(ctx, keeper)
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

	keeper.SetRawDataReport(ctx, 11, 1, validatorAddress1, types.NewRawDataReport(0, []byte("answer1")))
	keeper.SetRawDataReport(ctx, 11, 1, validatorAddress2, types.NewRawDataReport(0, []byte("answer2")))
	keeper.SetPendingResolveList(ctx, []types.RequestID{10, 11})

	handleEndBlock(ctx, keeper)
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
	executeGasList := []uint64{2500, 50, 3000, 2500}

	for i := types.RequestID(1); i <= types.RequestID(4); i++ {
		handleMsgRequestData(
			ctx, keeper,
			types.NewMsgRequestData(scriptID, calldata, 2, 2, 100, 2000, executeGasList[i-1], "clientID", sender),
		)

		keeper.SetRawDataReport(ctx, i, 1, validatorAddress1, types.NewRawDataReport(0, []byte("answer1")))
		keeper.SetRawDataReport(ctx, i, 1, validatorAddress2, types.NewRawDataReport(0, []byte("answer2")))
		pendingList = append(pendingList, i)
	}

	keeper.SetParam(ctx, KeyEndBlockExecuteGasLimit, 2600)
	keeper.SetPendingResolveList(ctx, pendingList)

	handleEndBlock(ctx, keeper)
	require.Equal(t, []types.RequestID{4}, keeper.GetPendingResolveList(ctx))

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
	require.Equal(t, types.Failure, actualRequest.ResolveStatus)

	_, err = keeper.GetResult(ctx, 4, scriptID, calldata)
	require.NotNil(t, err)

	actualRequest, err = keeper.GetRequest(ctx, 4)
	require.Nil(t, err)
	require.Equal(t, types.Open, actualRequest.ResolveStatus)

}

func TestAddAndRemoveOracleAddress(t *testing.T) {
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

	address1 := keep.GetAddressFromPub(pubStr[0])

	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)
	reporterAddress2 := sdk.AccAddress(validatorAddress2)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	_, err := keeper.CoinKeeper.AddCoins(ctx, address1, keep.NewUBandCoins(1000000))
	require.Nil(t, err)

	request := types.NewRequest(1, calldata,
		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
		2, 1581589790, 102, 1000000, "clientID",
	)
	keeper.SetRequest(ctx, 1, request)
	keeper.SetRawDataRequest(ctx, 1, 42, types.NewRawDataRequest(1, []byte("calldata1")))

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	keeper.AddReporter(ctx, validatorAddress1, reporterAddress2)
	err = keeper.AddReporter(ctx, validatorAddress1, reporterAddress2)

	require.NotNil(t, err)

	msg := types.NewMsgReportData(1, []types.RawDataReportWithID{
		types.NewRawDataReportWithID(42, 0, []byte("data1")),
	}, validatorAddress1, reporterAddress2)

	_, err = handleMsgReportData(ctx, keeper, msg)
	require.Nil(t, err)
	list := keeper.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{}, list)

	keeper.RemoveReporter(ctx, validatorAddress1, reporterAddress2)
	err = keeper.RemoveReporter(ctx, validatorAddress1, reporterAddress2)
	require.NotNil(t, err)

	msg = types.NewMsgReportData(1, []types.RawDataReportWithID{
		types.NewRawDataReportWithID(42, 0, []byte("data2")),
	}, validatorAddress1, reporterAddress2)

	_, err = handleMsgReportData(ctx, keeper, msg)
	require.NotNil(t, err)
}
