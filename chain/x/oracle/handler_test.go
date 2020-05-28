package oracle_test

import (
	"bytes"
	gz "compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// func mockOracleScript(ctx sdk.Context, keeper Keeper) (*sdk.Result, error) {
// 	owner := sdk.AccAddress([]byte("owner"))
// 	name := "oracle_script_1"
// 	description := "description"
// 	code := []byte("code")
// 	sender := sdk.AccAddress([]byte("sender"))
// 	schema := "schema"
// 	sourceCodeURL := "sourceCodeURL"
// 	msg := types.NewMsgCreateOracleScript(owner, name, description, code, schema, sourceCodeURL, sender)
// 	return handleMsgCreateOracleScript(ctx, keeper, msg)
// }

func TestCreateDataSourceSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()

	owner := Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	msg := types.NewMsgCreateDataSource(owner, name, description, executable, Alice.Address)
	res, err := oracle.NewHandler(keeper)(ctx, msg)
	defer deleteFile(executable)
	require.Nil(t, err)
	require.NotNil(t, res)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, Owner.Address, dataSource.Owner)
	require.Equal(t, name, dataSource.Name)
	executableHash := sha256.Sum256(executable)
	expectFilename := hex.EncodeToString(executableHash[:])
	require.Equal(t, expectFilename, dataSource.Filename)
}

func TestCreateGzippedExecutableDataSourceSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()

	owner := Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(executable)
	zw.Close()
	gzippedExecutable := buf.Bytes()

	sender := Alice.Address
	msg := types.NewMsgCreateDataSource(owner, name, description, gzippedExecutable, sender)
	res, err := oracle.NewHandler(keeper)(ctx, msg)
	defer deleteFile(executable)
	require.NoError(t, err)
	require.NotNil(t, res)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, Owner.Address, dataSource.Owner)
	require.Equal(t, name, dataSource.Name)
	executableHash := sha256.Sum256(executable)
	expectFilename := hex.EncodeToString(executableHash[:])
	require.Equal(t, expectFilename, dataSource.Filename)
}

func TestCreateGzippedExecutableDataSourceFail(t *testing.T) {
	_, ctx, keeper := createTestInput()

	owner := Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(executable)
	zw.Close()
	gzippedExecutable := buf.Bytes()[:5]

	sender := Alice.Address
	msg := types.NewMsgCreateDataSource(owner, name, description, gzippedExecutable, sender)
	res, err := oracle.NewHandler(keeper)(ctx, msg)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestEditDataSourceSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()

	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	msg := types.NewMsgCreateDataSource(Owner.Address, name, description, executable, Alice.Address)
	defer deleteFile(executable)
	oracle.NewHandler(keeper)(ctx, msg)

	newName := "beeb"
	newDescription := "new_description"
	newExecutable := []byte("executable2")
	msgEdit := types.NewMsgEditDataSource(1, Owner.Address, newName, newDescription, newExecutable, Owner.Address)
	defer deleteFile(newExecutable)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)
	require.NoError(t, err)
	require.NotNil(t, res)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, Owner.Address, dataSource.Owner)
	require.Equal(t, newName, dataSource.Name)
	executableHash := sha256.Sum256(newExecutable)
	expectFilename := hex.EncodeToString(executableHash[:])
	require.Equal(t, expectFilename, dataSource.Filename)
}

func TestEditDataSourceFail(t *testing.T) {
	_, ctx, keeper := createTestInput()

	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	msg := types.NewMsgCreateDataSource(Owner.Address, name, description, executable, Alice.Address)
	oracle.NewHandler(keeper)(ctx, msg)
	defer deleteFile(executable)

	newName := "beeb"
	newDescription := "new_description"
	newExecutable := []byte("executable2")
	wrongDID := types.DataSourceID(99999)

	msgEdit := types.NewMsgEditDataSource(wrongDID, Owner.Address, newName, newDescription, newExecutable, Owner.Address)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)
	require.Nil(t, res)

	wrongSender := Bob.Address
	msgEdit = types.NewMsgEditDataSource(1, Owner.Address, newName, newDescription, newExecutable, wrongSender)
	res, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)
	require.Nil(t, res)

	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(executable)
	zw.Close()
	wrongGzippedExecutable := buf.Bytes()[:5]
	msgEdit = types.NewMsgEditDataSource(1, Owner.Address, newName, newDescription, wrongGzippedExecutable, Owner.Address)
	res, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)
	require.Nil(t, res)
}

// func TestEditOracleScriptSuccess(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	_, err := mockOracleScript(ctx, keeper)
// 	require.Nil(t, err)

// 	newOwner := sdk.AccAddress([]byte("anotherowner"))
// 	newName := "oracle_script_2"
// 	newDescription := "description_2"
// 	newCode := []byte("code_2")
// 	sender := sdk.AccAddress([]byte("owner"))
// 	schema := "schema"
// 	sourceCodeURL := "sourceCodeURL"

// 	msg := types.NewMsgEditOracleScript(1, newOwner, newName, newDescription, newCode, schema, sourceCodeURL, sender)
// 	_, err = handleMsgEditOracleScript(ctx, keeper, msg)
// 	require.Nil(t, err)

// 	oracleScript, err := keeper.GetOracleScript(ctx, 1)
// 	require.Nil(t, err)
// 	require.Equal(t, newOwner, oracleScript.Owner)
// 	require.Equal(t, newName, oracleScript.Name)
// 	require.Equal(t, newCode, oracleScript.Code)

// 	events := ctx.EventManager().Events()
// 	require.Equal(t, 2, len(events))
// 	require.Equal(t, sdk.Event{
// 		Type:       EventTypeCreateOracleScript,
// 		Attributes: []tmkv.Pair{{Key: []byte(AttributeKeyID), Value: []byte("1")}},
// 	}, events[0])
// 	require.Equal(t, sdk.Event{
// 		Type:       EventTypeEditOracleScript,
// 		Attributes: []tmkv.Pair{{Key: []byte(AttributeKeyID), Value: []byte("1")}},
// 	}, events[1])
// }

// func TestEditOracleScriptByNotOwner(t *testing.T) {
// 	ctx, keeper := keep.CreateTestInput(t, false)
// 	_, err := mockOracleScript(ctx, keeper)
// 	require.Nil(t, err)

// 	newOwner := sdk.AccAddress([]byte("anotherowner"))
// 	newName := "data_source_2"
// 	newDescription := "description_2"
// 	newCode := []byte("code_2")
// 	sender := sdk.AccAddress([]byte("not_owner"))
// 	schema := "schema"
// 	soureCodeURL := "sourceCodeURL"

// 	msg := types.NewMsgEditOracleScript(1, newOwner, newName, newDescription, newCode, schema, soureCodeURL, sender)
// 	_, err = handleMsgEditOracleScript(ctx, keeper, msg)
// 	require.NotNil(t, err)
// }

// func TestRequestSuccess(t *testing.T) {
// 	// Setup test environment
// 	ctx, keeper := keep.CreateTestInput(t, false)

// 	ctx = ctx.WithBlockHeight(2)
// 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
// 	calldata := []byte("calldata")
// 	sender := sdk.AccAddress([]byte("sender"))
// 	_, err := keeper.CoinKeeper.AddCoins(ctx, sender, keep.NewUBandCoins(410))
// 	require.Nil(t, err)

// 	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
// 	keeper.SetOracleScript(ctx, 1, script)

// 	pubStr := []string{
// 		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
// 		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
// 	}

// 	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
// 	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

// 	dataSource1 := keep.GetTestDataSource()
// 	keeper.SetDataSource(ctx, 1, dataSource1)

// 	dataSource2 := types.NewDataSource(
// 		sdk.AccAddress([]byte("anotherowner")),
// 		"data_source2",
// 		"description2",
// 		[]byte("executable2"),
// 	)
// 	keeper.SetDataSource(ctx, 2, dataSource2)

// 	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

// 	// Test here
// 	beforeGas := ctx.GasMeter().GasConsumed()
// 	_, err = handleMsgRequestData(ctx, keeper, msg)
// 	afterGas := ctx.GasMeter().GasConsumed()
// 	require.Nil(t, err)

// 	// Check global request count
// 	require.Equal(t, int64(1), keeper.GetRequestCount(ctx))
// 	actualRequest, err := keeper.GetRequest(ctx, 1)
// 	require.Nil(t, err)
// 	expectRequest := types.NewRequest(1, calldata,
// 		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
// 		2, 1581589790, "clientID", nil, nil,
// 	)
// 	require.Equal(t, expectRequest, actualRequest)

// 	require.Equal(t, int64(2), keeper.GetRawRequestCount(ctx, 1))

// 	// rawRequests := []types.RawDataRequest{
// 	// 	types.NewRawRequest(1, []byte("band-protocol")), types.NewRawRequest(2, []byte("band-chain")),
// 	// }
// 	// require.Equal(t, rawRequests, keeper.GetRawRequests(ctx, 1))
// 	// check consumed gas must more than 100000
// 	// TODO: Write a better test than just checking number comparison
// 	require.GreaterOrEqual(t, afterGas-beforeGas, uint64(50000))
// }

// // func TestIBCInfoSuccess(t *testing.T) {
// // 	// Setup test environment
// // 	ctx, keeper := keep.CreateTestInput(t, false)

// // 	ctx = ctx.WithBlockHeight(2)
// // 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
// // 	calldata := []byte("calldata")
// // 	sender := sdk.AccAddress([]byte("sender"))
// // 	_, err := keeper.CoinKeeper.AddCoins(ctx, sender, keep.NewUBandCoins(410))
// // 	require.Nil(t, err)

// // 	sourcePort := "sourcePort"
// // 	sourceChannel := "sourceChannel"

// // 	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
// // 	keeper.SetOracleScript(ctx, 1, script)

// // 	pubStr := []string{
// // 		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
// // 		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
// // 	}

// // 	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
// // 	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

// // 	dataSource1 := keep.GetTestDataSource()
// // 	keeper.SetDataSource(ctx, 1, dataSource1)

// // 	dataSource2 := types.NewDataSource(
// // 		sdk.AccAddress([]byte("anotherowner")),
// // 		"data_source2",
// // 		"description2",
// // 		[]byte("executable2"),
// // 	)
// // 	keeper.SetDataSource(ctx, 2, dataSource2)

// // 	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

// // 	// Test here
// // 	beforeGas := ctx.GasMeter().GasConsumed()
// // 	_, err = handleMsgRequestDataIBC(ctx, keeper, msg, sourcePort, sourceChannel)
// // 	afterGas := ctx.GasMeter().GasConsumed()
// // 	require.Nil(t, err)

// // 	// Check global request count
// // 	require.Equal(t, int64(1), keeper.GetRequestCount(ctx))
// // 	actualRequest, err := keeper.GetRequest(ctx, 1)
// // 	require.Nil(t, err)
// // 	expectRequest := types.NewRequest(1, calldata,
// // 		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
// // 		2, 1581589790, "clientID", &types.IBCInfo{sourcePort, sourceChannel},
// // 	)
// // 	require.Equal(t, expectRequest, actualRequest)

// // 	require.Equal(t, int64(2), keeper.GetRawRequestCount(ctx, 1))

// // 	// rawRequests := []types.RawDataRequest{
// // 	// 	types.NewRawRequest(1, []byte("band-protocol")), types.NewRawRequest(2, []byte("band-chain")),
// // 	// }
// // 	// require.Equal(t, rawRequests, keeper.GetRawRequests(ctx, 1))
// // 	// check consumed gas must more than 100000
// // 	// TODO: Write a better test than just checking number comparison
// // 	require.GreaterOrEqual(t, afterGas-beforeGas, uint64(50000))
// // }

// func TestRequestInvalidDataSource(t *testing.T) {
// 	// Setup test environment
// 	ctx, keeper := keep.CreateTestInput(t, false)

// 	ctx = ctx.WithBlockHeight(2)
// 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
// 	calldata := []byte("calldata")
// 	sender := sdk.AccAddress([]byte("sender"))

// 	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

// 	_, err := handleMsgRequestData(ctx, keeper, msg)
// 	require.NotNil(t, err)

// 	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
// 	keeper.SetOracleScript(ctx, 1, script)

// 	pubStr := []string{
// 		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
// 		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
// 	}

// 	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
// 	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

// 	_, err = handleMsgRequestData(ctx, keeper, msg)
// 	require.NotNil(t, err)
// }

// // func TestIBCInfoInvalidDataSource(t *testing.T) {
// // 	// Setup test environment
// // 	ctx, keeper := keep.CreateTestInput(t, false)

// // 	ctx = ctx.WithBlockHeight(2)
// // 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
// // 	calldata := []byte("calldata")
// // 	sender := sdk.AccAddress([]byte("sender"))
// // 	sourcePort := "sourcePort"
// // 	sourceChannel := "sourceChannel"

// // 	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

// // 	_, err := handleMsgRequestDataIBC(ctx, keeper, msg, sourcePort, sourceChannel)
// // 	require.NotNil(t, err)

// // 	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
// // 	keeper.SetOracleScript(ctx, 1, script)

// // 	pubStr := []string{
// // 		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
// // 		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
// // 	}

// // 	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
// // 	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

// // 	_, err = handleMsgRequestDataIBC(ctx, keeper, msg, sourcePort, sourceChannel)
// // 	require.NotNil(t, err)
// // }

// func TestRequestWithPrepareGasExceed(t *testing.T) {
// 	// Setup test environment
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

// 	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
// 	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

// 	dataSource := keep.GetTestDataSource()
// 	keeper.SetDataSource(ctx, 1, dataSource)

// 	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

// 	_, err := handleMsgRequestData(ctx, keeper, msg)
// 	require.NotNil(t, err)
// }

// // func TestIBCInfoWithPrepareGasExceed(t *testing.T) {
// // 	// Setup test environment
// // 	ctx, keeper := keep.CreateTestInput(t, false)

// // 	ctx = ctx.WithBlockHeight(2)
// // 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
// // 	calldata := []byte("calldata")
// // 	sender := sdk.AccAddress([]byte("sender"))
// // 	sourcePort := "sourcePort"
// // 	sourceChannel := "sourceChannel"
// // 	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
// // 	keeper.SetOracleScript(ctx, 1, script)

// // 	pubStr := []string{
// // 		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
// // 		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
// // 	}

// // 	keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
// // 	keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

// // 	dataSource := keep.GetTestDataSource()
// // 	keeper.SetDataSource(ctx, 1, dataSource)

// // 	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

// // 	_, err := handleMsgRequestDataIBC(ctx, keeper, msg, sourcePort, sourceChannel)
// // 	require.NotNil(t, err)
// // }

func TestReportSuccess(t *testing.T) {
	// Setup test environment
	_, ctx, k := createTestInput()

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")

	request := types.NewRequest(1, calldata,
		[]sdk.ValAddress{Validator1.ValAddress, Validator2.ValAddress}, 2,
		2, 1581589790, "clientID", nil, []types.ExternalID{1, 42},
	)
	k.SetRequest(ctx, 1, request)

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	msg := types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(1, 0, []byte("data1")),
		types.NewRawReport(42, 0, []byte("data2")),
	}, Validator1.ValAddress, Validator1.Address)

	_, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)

	list := k.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{}, list)

	msg = types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(1, 0, []byte("data3")),
		types.NewRawReport(42, 0, []byte("data4")),
	}, Validator2.ValAddress, Validator2.Address)

	_, err = oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)

	list = k.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{1}, list)
}

func TestReportFailed(t *testing.T) {
	// Setup test environment
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	calldata := []byte("calldata")

	request := types.NewRequest(1, calldata,
		[]sdk.ValAddress{Validator1.ValAddress, Validator2.ValAddress}, 2,
		2, 1581589790, "clientID", nil, []types.ExternalID{42},
	)
	k.SetRequest(ctx, 1, request)

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	// Report by unauthorized reporter
	msg := types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(42, 0, []byte("data1")),
	}, Validator1.ValAddress, Alice.Address)
	_, err := oracle.NewHandler(k)(ctx, msg)
	require.Error(t, err)

	// Send wrong external ids
	msg = types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(41, 0, []byte("data1")),
	}, Validator1.ValAddress, Validator1.Address)

	// Test only 1 failed case, other case tested in keeper/report_test.go
	_, err = oracle.NewHandler(k)(ctx, msg)
	require.Error(t, err)
}

// // func TestEndBlock(t *testing.T) {
// // 	ctx, keeper := keep.CreateTestInput(t, false)

// // 	ctx = ctx.WithBlockHeight(2)
// // 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
// // 	calldata := []byte("calldata")
// // 	sender := sdk.AccAddress([]byte("sender"))

// // 	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
// // 	keeper.SetOracleScript(ctx, 1, script)

// // 	pubStr := []string{
// // 		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
// // 		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
// // 	}

// // 	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)
// // 	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)

// // 	dataSource := keep.GetTestDataSource()
// // 	keeper.SetDataSource(ctx, 1, dataSource)

// // 	msg := types.NewMsgRequestData(1, calldata, 2, 2, "clientID", sender)

// // 	handleMsgRequestData(ctx, keeper, msg)

// // 	keeper.SetReport(ctx, 1, 1, validatorAddress1, types.NewReport(0, []byte("answer1")))
// // 	keeper.SetReport(ctx, 1, 1, validatorAddress2, types.NewReport(0, []byte("answer2")))

// // 	keeper.SetPendingResolveList(ctx, []types.RequestID{1})

// // 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589999), 0))
// // 	handleEndBlock(ctx, keeper)

// // 	require.Equal(t, []types.RequestID{}, keeper.GetPendingResolveList(ctx))

// // 	result, err := keeper.GetResult(ctx, 1, 1, calldata)
// // 	require.Nil(t, err)
// // 	require.Equal(t,
// // 		types.Result{
// // 			RequestTime:              1581589790,
// // 			AggregationTime:          1581589999,
// // 			RequestedValidatorsCount: 2,
// // 			MinCount: 2,
// // 			ReportedValidatorsCount:  0,
// // 			Data:                     []byte("answer2"),
// // 		},
// // 		result,
// // 	)

// // 	actualRequest, err := keeper.GetRequest(ctx, 1)
// // 	require.Nil(t, err)
// // 	require.Equal(t, types.ResolveStatus_Success, actualRequest.ResolveStatus)
// // }

// func TestEndBlockExecuteFailedIfExecuteGasLessThanGasUsed(t *testing.T) {
// 	// TODO: Write this test properly. Pending on having owasm that can easily control gas usage.
// }

// func TestAddAndRemoveOracleAddress(t *testing.T) {
// 	// Setup test environment
// 	ctx, keeper := keep.CreateTestInput(t, false)

// 	ctx = ctx.WithBlockHeight(2)
// 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
// 	calldata := []byte("calldata")

// 	script := keep.GetTestOracleScript("../../pkg/owasm/res/silly.wasm")
// 	keeper.SetOracleScript(ctx, 1, script)

// 	pubStr := []string{
// 		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
// 		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
// 	}

// 	validatorAddress1 := keep.SetupTestValidator(ctx, keeper, pubStr[0], 10)

// 	address1 := keep.GetAddressFromPub(pubStr[0])

// 	validatorAddress2 := keep.SetupTestValidator(ctx, keeper, pubStr[1], 100)
// 	reporterAddress2 := sdk.AccAddress(validatorAddress2)

// 	dataSource := keep.GetTestDataSource()
// 	keeper.SetDataSource(ctx, 1, dataSource)

// 	_, err := keeper.CoinKeeper.AddCoins(ctx, address1, keep.NewUBandCoins(1000000))
// 	require.Nil(t, err)

// 	request := types.NewRequest(1, calldata,
// 		[]sdk.ValAddress{validatorAddress2, validatorAddress1}, 2,
// 		2, 1581589790, "clientID", nil, nil,
// 	)
// 	keeper.SetRequest(ctx, 1, request)
// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(42, 1, []byte("calldata1")))

// 	ctx = ctx.WithBlockHeight(5)
// 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

// 	keeper.AddReporter(ctx, validatorAddress1, reporterAddress2)
// 	err = keeper.AddReporter(ctx, validatorAddress1, reporterAddress2)

// 	require.NotNil(t, err)

// 	msg := types.NewMsgReportData(1, []types.RawReport{
// 		types.NewRawReport(42, 0, []byte("data1")),
// 	}, validatorAddress1, reporterAddress2)

// 	_, err = handleMsgReportData(ctx, keeper, msg)
// 	require.Nil(t, err)
// 	list := keeper.GetPendingResolveList(ctx)
// 	require.Equal(t, []types.RequestID{}, list)

// 	keeper.RemoveReporter(ctx, validatorAddress1, reporterAddress2)
// 	err = keeper.RemoveReporter(ctx, validatorAddress1, reporterAddress2)
// 	require.NotNil(t, err)

// 	msg = types.NewMsgReportData(1, []types.RawReport{
// 		types.NewRawReport(42, 0, []byte("data2")),
// 	}, validatorAddress1, reporterAddress2)

// 	_, err = handleMsgReportData(ctx, keeper, msg)
// 	require.NotNil(t, err)
// }

func TestAddReporterSuccess(t *testing.T) {
	_, ctx, k := createTestInput()

	validatorAddress := Alice.ValAddress
	reporterAddress := Bob.Address
	// Add bob reporter to alice validator
	msg := types.NewMsgAddReporter(validatorAddress, reporterAddress)

	result, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, result)

	events := result.GetEvents()

	// Add reporter events
	require.Equal(t, 1, len(events))
	require.Equal(t, types.EventTypeAddReporter, events[0].Type)
	require.Equal(t, types.AttributeKeyValidator, string(events[0].Attributes[0].Key))
	require.Equal(t, validatorAddress.String(), string(events[0].Attributes[0].Value))
	require.Equal(t, types.AttributeKeyReporter, string(events[0].Attributes[1].Key))
	require.Equal(t, reporterAddress.String(), string(events[0].Attributes[1].Value))

}

func TestAddReporterFail(t *testing.T) {
	_, ctx, k := createTestInput()

	validatorAddress := Alice.ValAddress
	reporterAddress := Alice.Address
	msg := types.NewMsgAddReporter(validatorAddress, reporterAddress)

	// Should fail, validator is always a reporter of himself so we can't add alice reporter to alice validator
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, fmt.Sprintf("val: %s, addr: %s: reporter already exists", validatorAddress.String(), reporterAddress.String()))
	require.Nil(t, result)
}

func TestRemoveReporterSuccess(t *testing.T) {
	_, ctx, k := createTestInput()

	validatorAddress := Alice.ValAddress
	reporterAddress := Alice.Address

	msg := types.NewMsgRemoveReporter(validatorAddress, reporterAddress)
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)

	events := result.GetEvents()
	// Remove reporter events
	require.Equal(t, 1, len(events))
	require.Equal(t, types.EventTypeRemoveReporter, events[0].Type)
	require.Equal(t, types.AttributeKeyValidator, string(events[0].Attributes[0].Key))
	require.Equal(t, validatorAddress.String(), string(events[0].Attributes[0].Value))
	require.Equal(t, types.AttributeKeyReporter, string(events[0].Attributes[1].Key))
	require.Equal(t, reporterAddress.String(), string(events[0].Attributes[1].Value))

}

func TestRemoveReporterFail(t *testing.T) {
	_, ctx, k := createTestInput()

	validatorAddress := Alice.ValAddress
	reporterAddress := Bob.Address

	// Should fail, bob isn't alice validator's reporter
	msg := types.NewMsgRemoveReporter(validatorAddress, reporterAddress)
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, fmt.Sprintf("val: %s, addr: %s: reporter not found", validatorAddress.String(), reporterAddress.String()))
	require.Nil(t, result)
}
