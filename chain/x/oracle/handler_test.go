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
	executable := []byte("executable")
	sender := sdk.AccAddress([]byte("sender"))
	msg := types.NewMsgCreateDataSource(owner, name, description, executable, sender)
	return handleMsgCreateDataSource(ctx, keeper, msg)
}

func mockOracleScript(ctx sdk.Context, keeper Keeper) (*sdk.Result, error) {
	owner := sdk.AccAddress([]byte("owner"))
	name := "oracle_script_1"
	description := "description"
	code := []byte("code")
	sender := sdk.AccAddress([]byte("sender"))
	schema := "schema"
	sourceCodeURL := "sourceCodeURL"
	msg := types.NewMsgCreateOracleScript(owner, name, description, code, schema, sourceCodeURL, sender)
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
	require.Equal(t, []byte("executable"), dataSource.Executable)

	events := ctx.EventManager().Events()
	require.Equal(t, 1, len(events))
	require.Equal(t, sdk.Event{
		Type:       EventTypeCreateDataSource,
		Attributes: []tmkv.Pair{{Key: []byte(AttributeKeyID), Value: []byte("1")}},
	}, events[0])
}

func TestEditDataSourceSuccess(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	_, err := mockDataSource(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("anotherowner"))
	newName := "data_source_2"
	newDescription := "new_description"
	newExecutable := []byte("executable_2")
	sender := sdk.AccAddress([]byte("owner"))

	msg := types.NewMsgEditDataSource(1, newOwner, newName, newDescription, newExecutable, sender)
	_, err = handleMsgEditDataSource(ctx, keeper, msg)
	require.Nil(t, err)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, dataSource.Owner)
	require.Equal(t, newName, dataSource.Name)
	require.Equal(t, newDescription, dataSource.Description)
	require.Equal(t, newExecutable, dataSource.Executable)

	events := ctx.EventManager().Events()
	require.Equal(t, 2, len(events))
	require.Equal(t, sdk.Event{
		Type:       EventTypeCreateDataSource,
		Attributes: []tmkv.Pair{{Key: []byte(AttributeKeyID), Value: []byte("1")}},
	}, events[0])
	require.Equal(t, sdk.Event{
		Type:       EventTypeEditDataSource,
		Attributes: []tmkv.Pair{{Key: []byte(AttributeKeyID), Value: []byte("1")}},
	}, events[1])
}

func TestEditDataSourceByNotOwner(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	mockDataSource(ctx, keeper)

	newOwner := sdk.AccAddress([]byte("anotherowner"))
	newName := "data_source_2"
	newDescription := "new_description"
	newExecutable := []byte("executable_2")
	sender := sdk.AccAddress([]byte("sender"))

	msg := types.NewMsgEditDataSource(1, newOwner, newName, newDescription, newExecutable, sender)
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
		Attributes: []tmkv.Pair{{Key: []byte(AttributeKeyID), Value: []byte("1")}},
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
	schema := "schema"
	sourceCodeURL := "sourceCodeURL"

	msg := types.NewMsgEditOracleScript(1, newOwner, newName, newDescription, newCode, schema, sourceCodeURL, sender)
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
		Attributes: []tmkv.Pair{{Key: []byte(AttributeKeyID), Value: []byte("1")}},
	}, events[0])
	require.Equal(t, sdk.Event{
		Type:       EventTypeEditOracleScript,
		Attributes: []tmkv.Pair{{Key: []byte(AttributeKeyID), Value: []byte("1")}},
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
	schema := "schema"
	soureCodeURL := "sourceCodeURL"

	msg := types.NewMsgEditOracleScript(1, newOwner, newName, newDescription, newCode, schema, soureCodeURL, sender)
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

	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
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
		[]byte("executable2"),
	)
	keeper.SetDataSource(ctx, 2, dataSource2)

	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

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
		2, 1581589790, "clientID", nil,
	)
	require.Equal(t, expectRequest, actualRequest)

	require.Equal(t, int64(2), keeper.GetRawRequestCount(ctx, 1))

	// rawRequests := []types.RawDataRequest{
	// 	types.NewRawRequest(1, []byte("band-protocol")), types.NewRawRequest(2, []byte("band-chain")),
	// }
	// require.Equal(t, rawRequests, keeper.GetRawRequests(ctx, 1))
	// check consumed gas must more than 100000
	// TODO: Write a better test than just checking number comparison
	require.GreaterOrEqual(t, afterGas-beforeGas, uint64(100000))
}

func TestRequestIBCSuccess(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))
	_, err := keeper.CoinKeeper.AddCoins(ctx, sender, keep.NewUBandCoins(410))
	require.Nil(t, err)

	sourcePort := "sourcePort"
	sourceChannel := "sourceChannel"

	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
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
		[]byte("executable2"),
	)
	keeper.SetDataSource(ctx, 2, dataSource2)

	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

	// Test here
	beforeGas := ctx.GasMeter().GasConsumed()
	_, err = handleMsgRequestDataIBC(ctx, keeper, msg, sourcePort, sourceChannel)
	afterGas := ctx.GasMeter().GasConsumed()
	require.Nil(t, err)

	// Check global request count
	require.Equal(t, int64(1), keeper.GetRequestCount(ctx))
	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	expectRequest := types.NewRequest(1, calldata,
		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
		2, 1581589790, "clientID", &types.RequestIBC{sourcePort, sourceChannel},
	)
	require.Equal(t, expectRequest, actualRequest)

	require.Equal(t, int64(2), keeper.GetRawRequestCount(ctx, 1))

	// rawRequests := []types.RawDataRequest{
	// 	types.NewRawRequest(1, []byte("band-protocol")), types.NewRawRequest(2, []byte("band-chain")),
	// }
	// require.Equal(t, rawRequests, keeper.GetRawRequests(ctx, 1))
	// check consumed gas must more than 100000
	// TODO: Write a better test than just checking number comparison
	require.GreaterOrEqual(t, afterGas-beforeGas, uint64(100000))
}

func TestRequestInvalidDataSource(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

	_, err := handleMsgRequestData(ctx, keeper, msg)
	require.NotNil(t, err)

	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
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

func TestRequestIBCInvalidDataSource(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))
	sourcePort := "sourcePort"
	sourceChannel := "sourceChannel"

	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

	_, err := handleMsgRequestDataIBC(ctx, keeper, msg, sourcePort, sourceChannel)
	require.NotNil(t, err)

	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	_, err = handleMsgRequestDataIBC(ctx, keeper, msg, sourcePort, sourceChannel)
	require.NotNil(t, err)
}

func TestRequestWithPrepareGasExceed(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))

	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

	_, err := handleMsgRequestData(ctx, keeper, msg)
	require.NotNil(t, err)
}

func TestRequestIBCWithPrepareGasExceed(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")
	sender := sdk.AccAddress([]byte("sender"))
	sourcePort := "sourcePort"
	sourceChannel := "sourceChannel"
	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
	keeper.SetOracleScript(ctx, 1, script)

	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}

	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

	dataSource := keep.GetTestDataSource()
	keeper.SetDataSource(ctx, 1, dataSource)

	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

	_, err := handleMsgRequestDataIBC(ctx, keeper, msg, sourcePort, sourceChannel)
	require.NotNil(t, err)
}

func TestReportSuccess(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")

	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
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
		2, 1581589790, "clientID", nil,
	)
	keeper.SetRequest(ctx, 1, request)
	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(42, 1, []byte("calldata1")))

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	msg := types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(42, 0, []byte("data1")),
	}, validatorAddress1, reporterAddress1)

	_, err = handleMsgReportData(ctx, keeper, msg)
	require.Nil(t, err)
	list := keeper.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{}, list)

	msg = types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(42, 0, []byte("data2")),
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

	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
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
		2, 1581589790, "clientID", nil,
	)
	keeper.SetRequest(ctx, 1, request)
	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(42, 1, []byte("calldata1")))

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	msg := types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(41, 0, []byte("data1")),
	}, validatorAddress1, reporterAddress1)

	// Test only 1 failed case, other case tested in keeper/report_test.go
	_, err = handleMsgReportData(ctx, keeper, msg)
	require.NotNil(t, err)
}

// func TestEndBlock(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)

// 	ctx = ctx.WithBlockHeight(2)
// 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
// 	calldata := []byte("calldata")
// 	sender := sdk.AccAddress([]byte("sender"))

// 	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
// 	keeper.SetOracleScript(ctx, 1, script)

// 	pubStr := []string{
// 		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
// 		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
// 	}

// 	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
// 	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

// 	dataSource := keep.GetTestDataSource()
// 	keeper.SetDataSource(ctx, 1, dataSource)

// 	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

// 	handleMsgRequestData(ctx, keeper, msg)

// 	keeper.SetReport(ctx, 1, 1, validatorAddress1, types.NewRawDataReport(0, []byte("answer1")))
// 	keeper.SetReport(ctx, 1, 1, validatorAddress2, types.NewRawDataReport(0, []byte("answer2")))

// 	keeper.SetPendingResolveList(ctx, []types.RequestID{1})

// 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589999), 0))
// 	handleEndBlock(ctx, keeper)

// 	require.Equal(t, []types.RequestID{}, keeper.GetPendingResolveList(ctx))

// 	result, err := keeper.GetResult(ctx, 1, 1, calldata)
// 	require.Nil(t, err)
// 	require.Equal(t,
// 		types.Result{
// 			RequestTime:              1581589790,
// 			AggregationTime:          1581589999,
// 			RequestedValidatorsCount: 2,
// 			MinCount: 2,
// 			ReportedValidatorsCount:  0,
// 			Data:                     []byte("answer2"),
// 		},
// 		result,
// 	)

// 	actualRequest, err := keeper.GetRequest(ctx, 1)
// 	require.Nil(t, err)
// 	require.Equal(t, types.ResolveStatus_Success, actualRequest.ResolveStatus)
// }

func TestEndBlockExecuteFailedIfExecuteGasLessThanGasUsed(t *testing.T) {
	// TODO: Write this test properly. Pending on having owasm that can easily control gas usage.
}

func TestAddAndRemoveOracleAddress(t *testing.T) {
	// Setup test environment
	ctx, keeper := keep.CreateTestInput(t, false)

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")

	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
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
		2, 1581589790, "clientID", nil,
	)
	keeper.SetRequest(ctx, 1, request)
	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(42, 1, []byte("calldata1")))

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	keeper.AddReporter(ctx, validatorAddress1, reporterAddress2)
	err = keeper.AddReporter(ctx, validatorAddress1, reporterAddress2)

	require.NotNil(t, err)

	msg := types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(42, 0, []byte("data1")),
	}, validatorAddress1, reporterAddress2)

	_, err = handleMsgReportData(ctx, keeper, msg)
	require.Nil(t, err)
	list := keeper.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{}, list)

	keeper.RemoveReporter(ctx, validatorAddress1, reporterAddress2)
	err = keeper.RemoveReporter(ctx, validatorAddress1, reporterAddress2)
	require.NotNil(t, err)

	msg = types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(42, 0, []byte("data2")),
	}, validatorAddress1, reporterAddress2)

	_, err = handleMsgReportData(ctx, keeper, msg)
	require.NotNil(t, err)
}
