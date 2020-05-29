package oracle_test

import (
	"bytes"
	gz "compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestCreateDataSourceSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")

	owner := Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	msg := types.NewMsgCreateDataSource(owner, name, description, executable, Alice.Address)
	res, err := oracle.NewHandler(keeper)(ctx, msg)

	filename := keeper.MustGetDataSource(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

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
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")

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
	filename := keeper.MustGetDataSource(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

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
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")

	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	msg := types.NewMsgCreateDataSource(Owner.Address, name, description, executable, Alice.Address)
	oracle.NewHandler(keeper)(ctx, msg)

	filename := keeper.MustGetDataSource(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

	newName := "beeb"
	newDescription := "new_description"
	newExecutable := []byte("executable2")
	msgEdit := types.NewMsgEditDataSource(1, Owner.Address, newName, newDescription, newExecutable, Owner.Address)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)

	filename = keeper.MustGetDataSource(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

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
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")

	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	msg := types.NewMsgCreateDataSource(Owner.Address, name, description, executable, Alice.Address)
	oracle.NewHandler(keeper)(ctx, msg)
	filename := keeper.MustGetDataSource(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

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

func TestCreateOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := []byte("code")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(Owner.Address, name, description, code, schema, url, Alice.Address)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	defer deleteFile(code)

	require.NoError(t, err)

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.NoError(t, err)
	codeHash := sha256.Sum256(code)
	expectFilename := hex.EncodeToString(codeHash[:])
	expectOracleScript := types.NewOracleScript(Owner.Address, name, description, expectFilename, schema, url)

	require.Equal(t, expectOracleScript, oracleScript)
}
func TestCreateGzippedOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := []byte("OWASMCODE")
	schema := "schema"
	url := "url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(code)
	zw.Close()
	gzippedCode := buf.Bytes()

	msg := types.NewMsgCreateOracleScript(Owner.Address, name, description, gzippedCode, schema, url, Alice.Address)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	defer deleteFile(code)

	require.NoError(t, err)

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.NoError(t, err)
	codeHash := sha256.Sum256(code)
	expectFilename := hex.EncodeToString(codeHash[:])
	expectOracleScript := types.NewOracleScript(Owner.Address, name, description, expectFilename, schema, url)

	require.Equal(t, expectOracleScript, oracleScript)
}

func TestCreateGzippedOracleScriptFail(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := []byte("OWASMCODE")
	schema := "schema"
	url := "url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(code)
	zw.Close()
	gzippedCode := buf.Bytes()[:5]

	msg := types.NewMsgCreateOracleScript(Owner.Address, name, description, gzippedCode, schema, url, Alice.Address)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.Error(t, err)
}

func TestEditOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := []byte("code")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(Owner.Address, name, description, code, schema, url, Alice.Address)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	defer deleteFile(code)

	require.Nil(t, err)
	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := []byte("codecode")
	newSchema := "new_schema"
	newUrl := "new_url"

	msgEdit := types.NewMsgEditOracleScript(oracleScriptID, Owner.Address, newName, newDescription, newCode, newSchema, newUrl, Owner.Address)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)
	defer deleteFile(newCode)
	require.NoError(t, err)
	expectEvents := sdk.Events{
		sdk.NewEvent(types.EventTypeEditOracleScript, sdk.NewAttribute(types.AttributeKeyID, "1")),
	}
	require.Equal(t, expectEvents, res.GetEvents())

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.NoError(t, err)
	codeHash := sha256.Sum256(newCode)
	expectFilename := hex.EncodeToString(codeHash[:])
	expectOracleScript := types.NewOracleScript(Owner.Address, newName, newDescription, expectFilename, newSchema, newUrl)

	require.Equal(t, expectOracleScript, oracleScript)
}

func TestEditOracleScriptFail(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := []byte("code")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(Owner.Address, name, description, code, schema, url, Alice.Address)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	defer deleteFile(code)

	require.Nil(t, err)
	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := []byte("codecode")
	newSchema := "new_schema"
	newUrl := "new_url"

	// No oracle script id 99
	msgEdit := types.NewMsgEditOracleScript(types.OracleScriptID(99), Owner.Address, newName, newDescription, newCode, newSchema, newUrl, Alice.Address)
	_, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)

	// Alice can't edit oracle script
	msgEdit = types.NewMsgEditOracleScript(oracleScriptID, Owner.Address, newName, newDescription, newCode, newSchema, newUrl, Alice.Address)
	_, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)

}

func TestEditGzippedOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := []byte("code")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(Owner.Address, name, description, code, schema, url, Alice.Address)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	defer deleteFile(code)

	require.Nil(t, err)
	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := []byte("codecode")
	newSchema := "new_schema"
	newUrl := "new_url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(newCode)
	zw.Close()
	gzippedCode := buf.Bytes()

	msgEdit := types.NewMsgEditOracleScript(oracleScriptID, Owner.Address, newName, newDescription, gzippedCode, newSchema, newUrl, Owner.Address)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)
	require.NoError(t, err)
	expectEvents := sdk.Events{
		sdk.NewEvent(types.EventTypeEditOracleScript, sdk.NewAttribute(types.AttributeKeyID, "1")),
	}
	require.Equal(t, expectEvents, res.GetEvents())

	defer deleteFile(newCode)

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.NoError(t, err)
	codeHash := sha256.Sum256(newCode)
	expectFilename := hex.EncodeToString(codeHash[:])
	expectOracleScript := types.NewOracleScript(Owner.Address, newName, newDescription, expectFilename, newSchema, newUrl)

	require.Equal(t, expectOracleScript, oracleScript)
}

func TestEditGzippedOracleScriptFail(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := []byte("code")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(Owner.Address, name, description, code, schema, url, Alice.Address)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	defer deleteFile(code)

	require.Nil(t, err)
	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := []byte("codecode")
	newSchema := "new_schema"
	newUrl := "new_url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(newCode)
	zw.Close()
	gzippedCode := buf.Bytes()[:5]

	msgEdit := types.NewMsgEditOracleScript(oracleScriptID, Owner.Address, newName, newDescription, gzippedCode, newSchema, newUrl, Owner.Address)
	_, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)
}

func TestRequestDataSuccess(t *testing.T) {
	_, ctx, k := createTestInput()

	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))

	ds1, clear1 := getTestDataSource("code1")
	defer clear1()
	k.AddDataSource(ctx, ds1)

	ds2, clear2 := getTestDataSource("code2")
	defer clear2()
	k.AddDataSource(ctx, ds2)

	ds3, clear3 := getTestDataSource("code3")
	defer clear3()
	k.AddDataSource(ctx, ds3)

	os, clear4 := getTestOracleScript()
	defer clear4()

	oracleScriptID := k.AddOracleScript(ctx, os)

	calldata := []byte("beeb")
	msg := types.NewMsgRequestData(oracleScriptID, calldata, 2, 2, "alice", Alice.Address)

	result, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, result)

	expectEvents := sdk.Events{
		sdk.NewEvent(
			types.EventTypeRequest,
			sdk.NewAttribute(types.AttributeKeyID, "1"),
			sdk.NewAttribute(types.AttributeKeyValidator, Validator1.ValAddress.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, Validator3.ValAddress.String()),
		),
		sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, "1"),
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, ds1.Filename),
			sdk.NewAttribute(types.AttributeKeyExternalID, "1"),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(calldata)),
		),
		sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, "2"),
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, ds2.Filename),
			sdk.NewAttribute(types.AttributeKeyExternalID, "2"),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(calldata)),
		),
		sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, "3"),
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, ds3.Filename),
			sdk.NewAttribute(types.AttributeKeyExternalID, "3"),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(calldata)),
		),
	}

	require.Equal(t, expectEvents, result.GetEvents())
}

func TestRequestDataFail(t *testing.T) {
	_, ctx, k := createTestInput()

	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))

	wrongOracleScript := 1

	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	msg := types.NewMsgRequestData(types.OracleScriptID(wrongOracleScript), calldata, 1, 1, "alice", Alice.Address)

	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, `id: 1: oracle script not found`)
	require.Nil(t, result)

	// Add Oracle Script
	os, clear := getTestOracleScript()
	defer clear()

	oracleScriptID := k.AddOracleScript(ctx, os)
	msg = types.NewMsgRequestData(types.OracleScriptID(oracleScriptID), calldata, 1, 1, "alice", Alice.Address)

	result, err = oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, `id: 1: data source not found`)
	require.Nil(t, result)
}

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
	// Add Bob reporter to Alice validator
	msg := types.NewMsgAddReporter(validatorAddress, reporterAddress)

	result, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, result)

	events := result.GetEvents()

	expectedEvent := sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddReporter,
			sdk.NewAttribute(types.AttributeKeyValidator, validatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyReporter, reporterAddress.String()),
		),
	}

	require.Equal(t, expectedEvent, events)
}

func TestAddReporterFail(t *testing.T) {
	_, ctx, k := createTestInput()

	validatorAddress := Alice.ValAddress
	reporterAddress := Alice.Address
	msg := types.NewMsgAddReporter(validatorAddress, reporterAddress)

	// Should fail, validator is always a reporter of himself so we can't add Alice reporter to Alice validator
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, fmt.Sprintf("val: %s, addr: %s: reporter already exists", validatorAddress.String(), reporterAddress.String()))
	require.Nil(t, result)
}

func TestRemoveReporterSuccess(t *testing.T) {
	_, ctx, k := createTestInput()

	validatorAddress := Alice.ValAddress
	reporterAddress := Bob.Address

	// Add Bob reporter to Alice validator
	err := k.AddReporter(ctx, validatorAddress, reporterAddress)
	require.NoError(t, err)

	msg := types.NewMsgRemoveReporter(validatorAddress, reporterAddress)
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)

	events := result.GetEvents()

	expectedEvent := sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveReporter,
			sdk.NewAttribute(types.AttributeKeyValidator, validatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyReporter, reporterAddress.String()),
		),
	}

	require.Equal(t, expectedEvent, events)
}

func TestRemoveReporterFail(t *testing.T) {
	_, ctx, k := createTestInput()

	validatorAddress := Alice.ValAddress
	reporterAddress := Bob.Address

	// Should fail, Bob isn't Alice validator's reporter
	msg := types.NewMsgRemoveReporter(validatorAddress, reporterAddress)
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, fmt.Sprintf("val: %s, addr: %s: reporter not found", validatorAddress.String(), reporterAddress.String()))
	require.Nil(t, result)
}
