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
	"github.com/bandprotocol/bandchain/chain/x/oracle/testapp"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestCreateDataSourceSuccess(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()

	owner := testapp.Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	msg := types.NewMsgCreateDataSource(owner, name, description, executable, testapp.Alice.Address)
	res, err := oracle.NewHandler(keeper)(ctx, msg)

	require.Nil(t, err)
	require.NotNil(t, res)

	dataSource, err := keeper.GetDataSource(ctx, 6)
	require.Nil(t, err)
	require.Equal(t, testapp.Owner.Address, dataSource.Owner)
	require.Equal(t, name, dataSource.Name)
	executableHash := sha256.Sum256(executable)
	expectFilename := hex.EncodeToString(executableHash[:])
	require.Equal(t, expectFilename, dataSource.Filename)

}

func TestCreateGzippedExecutableDataSourceSuccess(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()

	owner := testapp.Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(executable)
	zw.Close()
	gzippedExecutable := buf.Bytes()

	sender := testapp.Alice.Address
	msg := types.NewMsgCreateDataSource(owner, name, description, gzippedExecutable, sender)
	res, err := oracle.NewHandler(keeper)(ctx, msg)

	require.NoError(t, err)
	require.NotNil(t, res)

	dataSource, err := keeper.GetDataSource(ctx, 6)
	require.Nil(t, err)
	require.Equal(t, testapp.Owner.Address, dataSource.Owner)
	require.Equal(t, name, dataSource.Name)
	executableHash := sha256.Sum256(executable)
	expectFilename := hex.EncodeToString(executableHash[:])
	require.Equal(t, expectFilename, dataSource.Filename)
}

func TestCreateGzippedExecutableDataSourceFail(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()

	owner := testapp.Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(executable)
	zw.Close()
	gzippedExecutable := buf.Bytes()[:5]

	sender := testapp.Alice.Address
	msg := types.NewMsgCreateDataSource(owner, name, description, gzippedExecutable, sender)
	res, err := oracle.NewHandler(keeper)(ctx, msg)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestEditDataSourceSuccess(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()

	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	msg := types.NewMsgCreateDataSource(testapp.Owner.Address, name, description, executable, testapp.Alice.Address)
	oracle.NewHandler(keeper)(ctx, msg)

	newName := "beeb"
	newDescription := "new_description"
	newExecutable := []byte("executable2")
	msgEdit := types.NewMsgEditDataSource(
		1, testapp.Owner.Address, newName, newDescription, newExecutable, testapp.Owner.Address,
	)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)

	require.NoError(t, err)
	require.NotNil(t, res)

	dataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, testapp.Owner.Address, dataSource.Owner)
	require.Equal(t, newName, dataSource.Name)
	executableHash := sha256.Sum256(newExecutable)
	expectFilename := hex.EncodeToString(executableHash[:])
	require.Equal(t, expectFilename, dataSource.Filename)
}

func TestEditDataSourceFail(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()

	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	msg := types.NewMsgCreateDataSource(testapp.Owner.Address, name, description, executable, testapp.Alice.Address)
	oracle.NewHandler(keeper)(ctx, msg)

	newName := "beeb"
	newDescription := "new_description"
	newExecutable := []byte("executable2")
	wrongDID := types.DataSourceID(99999)

	msgEdit := types.NewMsgEditDataSource(
		wrongDID, testapp.Owner.Address, newName, newDescription, newExecutable, testapp.Owner.Address,
	)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)
	require.Nil(t, res)

	wrongSender := testapp.Bob.Address
	msgEdit = types.NewMsgEditDataSource(
		1, testapp.Owner.Address, newName, newDescription, newExecutable, wrongSender,
	)
	res, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)
	require.Nil(t, res)

	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(executable)
	zw.Close()
	wrongGzippedExecutable := buf.Bytes()[:5]
	msgEdit = types.NewMsgEditDataSource(
		1, testapp.Owner.Address, newName, newDescription, wrongGzippedExecutable, testapp.Owner.Address,
	)
	res, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestCreateOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()
	nextID := keeper.GetOracleScriptCount(ctx) + 1
	name := "os_1"
	description := "beeb"
	code := testapp.WasmExtra1
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		testapp.Owner.Address, name, description, code, schema, url, testapp.Alice.Address,
	)
	res, err := oracle.NewHandler(keeper)(ctx, msg)
	require.NoError(t, err)
	expectEvents := sdk.Events{
		sdk.NewEvent(types.EventTypeCreateOracleScript, sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", nextID))),
	}
	require.Equal(t, expectEvents, res.Events)
	require.NoError(t, err)

	oracleScript, err := keeper.GetOracleScript(ctx, types.OracleScriptID(nextID))
	require.NoError(t, err)

	// Code is store must be compiled code.
	expectOracleScript := types.NewOracleScript(
		testapp.Owner.Address, name, description, testapp.WasmExtra1FileName, schema, url,
	)
	require.Equal(t, expectOracleScript, oracleScript)
}

func TestCreateOracleScriptFailed(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()
	name := "os_1"
	description := "beeb"
	code := []byte("non wasm code")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		testapp.Owner.Address, name, description, code, schema, url, testapp.Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.Error(t, err)
}
func TestCreateGzippedOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()
	nextID := keeper.GetOracleScriptCount(ctx) + 1
	name := "os_1"
	description := "beeb"
	code := testapp.WasmExtra1
	schema := "schema"
	url := "url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(code)
	zw.Close()
	gzippedCode := buf.Bytes()

	msg := types.NewMsgCreateOracleScript(
		testapp.Owner.Address, name, description, gzippedCode, schema, url, testapp.Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.NoError(t, err)

	oracleScript, err := keeper.GetOracleScript(ctx, types.OracleScriptID(nextID))
	require.NoError(t, err)

	// Code is store must be compiled code.
	expectOracleScript := types.NewOracleScript(
		testapp.Owner.Address, name, description, testapp.WasmExtra1FileName, schema, url,
	)

	require.Equal(t, expectOracleScript, oracleScript)
}

func TestCreateGzippedOracleScriptFail(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()
	name := "os_1"
	description := "beeb"
	code := testapp.WasmExtra1
	schema := "schema"
	url := "url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(code)
	zw.Close()
	gzippedCode := buf.Bytes()[:5]

	msg := types.NewMsgCreateOracleScript(
		testapp.Owner.Address, name, description, gzippedCode, schema, url, testapp.Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.Error(t, err)
}

func TestEditOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()
	name := "os_1"
	description := "beeb"
	code := testapp.WasmExtra1
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		testapp.Owner.Address, name, description, code, schema, url, testapp.Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.NoError(t, err)

	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := testapp.WasmExtra2
	newSchema := "new_schema"
	newURL := "new_url"

	msgEdit := types.NewMsgEditOracleScript(
		oracleScriptID, testapp.Owner.Address, newName, newDescription,
		newCode, newSchema, newURL, testapp.Owner.Address,
	)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)
	require.NoError(t, err)

	expectEvents := sdk.Events{
		sdk.NewEvent(types.EventTypeEditOracleScript, sdk.NewAttribute(types.AttributeKeyID, "1")),
	}
	require.Equal(t, expectEvents, res.Events)

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.NoError(t, err)
	expectOracleScript := types.NewOracleScript(
		testapp.Owner.Address, newName, newDescription,
		testapp.WasmExtra2FileName, newSchema, newURL,
	)

	require.Equal(t, expectOracleScript, oracleScript)
}

func TestEditOracleScriptFail(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()
	name := "os_1"
	description := "beeb"
	code := testapp.WasmExtra1
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		testapp.Owner.Address, name, description, code, schema, url, testapp.Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.NoError(t, err)

	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := testapp.WasmExtra2
	newSchema := "new_schema"
	newURL := "new_url"

	// No oracle script id 99
	msgEdit := types.NewMsgEditOracleScript(
		types.OracleScriptID(99), testapp.Owner.Address, newName, newDescription,
		newCode, newSchema, newURL, testapp.Alice.Address,
	)
	_, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)

	// testapp.Alice can't edit oracle script
	msgEdit = types.NewMsgEditOracleScript(
		oracleScriptID, testapp.Owner.Address, newName, newDescription,
		newCode, newSchema, newURL, testapp.Alice.Address,
	)
	_, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)

	// Cannot send bad owasm code
	msgEdit = types.NewMsgEditOracleScript(
		oracleScriptID, testapp.Owner.Address, newName, newDescription,
		[]byte("code"), newSchema, newURL, testapp.Owner.Address,
	)
	_, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)

}

func TestEditGzippedOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()
	name := "os_1"
	description := "beeb"
	code := testapp.WasmExtra1
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		testapp.Owner.Address, name, description, code, schema, url, testapp.Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.NoError(t, err)

	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := testapp.WasmExtra2
	require.NoError(t, err)
	newSchema := "new_schema"
	newURL := "new_url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(newCode)
	zw.Close()
	gzippedCode := buf.Bytes()

	msgEdit := types.NewMsgEditOracleScript(
		oracleScriptID, testapp.Owner.Address, newName, newDescription,
		gzippedCode, newSchema, newURL, testapp.Owner.Address,
	)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)
	require.NoError(t, err)
	expectEvents := sdk.Events{
		sdk.NewEvent(types.EventTypeEditOracleScript, sdk.NewAttribute(types.AttributeKeyID, "1")),
	}
	require.Equal(t, expectEvents, res.Events)

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.NoError(t, err)

	expectOracleScript := types.NewOracleScript(
		testapp.Owner.Address, newName, newDescription, testapp.WasmExtra2FileName, newSchema, newURL,
	)
	require.Equal(t, expectOracleScript, oracleScript)
}

func TestEditGzippedOracleScriptFail(t *testing.T) {
	_, ctx, keeper := testapp.CreateTestInput()
	name := "os_1"
	description := "beeb"
	code := testapp.WasmExtra1
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		testapp.Owner.Address, name, description, code, schema, url, testapp.Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)

	require.Nil(t, err)
	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := testapp.WasmExtra2
	newSchema := "new_schema"
	newURL := "new_url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(newCode)
	zw.Close()
	gzippedCode := buf.Bytes()[:5]

	msgEdit := types.NewMsgEditOracleScript(
		oracleScriptID, testapp.Owner.Address, newName, newDescription,
		gzippedCode, newSchema, newURL, testapp.Owner.Address,
	)
	_, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)
}

func TestRequestDataSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()

	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))

	oracleScriptID := types.OracleScriptID(1)
	calldata := []byte("beeb")
	msg := types.NewMsgRequestData(oracleScriptID, calldata, 2, 2, "alice", testapp.Alice.Address)

	result, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, result)

	expectEvents := sdk.Events{
		sdk.NewEvent(
			types.EventTypeRequest,
			sdk.NewAttribute(types.AttributeKeyID, "1"),
			sdk.NewAttribute(types.AttributeKeyOracleScriptID, "1"),
			sdk.NewAttribute(types.AttributeKeyCalldata, "62656562"), // "beeb" in hex
			sdk.NewAttribute(types.AttributeKeyAskCount, "2"),
			sdk.NewAttribute(types.AttributeKeyMinCount, "2"),
			sdk.NewAttribute(types.AttributeKeyValidator, testapp.Validator3.ValAddress.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, testapp.Validator1.ValAddress.String()),
		),
		sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, "1"),
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, testapp.DataSources[1].Filename),
			sdk.NewAttribute(types.AttributeKeyExternalID, "1"),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(calldata)),
		),
		sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, "2"),
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, testapp.DataSources[2].Filename),
			sdk.NewAttribute(types.AttributeKeyExternalID, "2"),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(calldata)),
		),
		sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, "3"),
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, testapp.DataSources[3].Filename),
			sdk.NewAttribute(types.AttributeKeyExternalID, "3"),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(calldata)),
		),
	}

	require.Equal(t, expectEvents, result.Events)
}

func TestRequestDataFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()

	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))

	wrongOracleScript := 42

	calldata := []byte("test")
	msg := types.NewMsgRequestData(
		types.OracleScriptID(wrongOracleScript), calldata, 1, 1, "alice", testapp.Alice.Address,
	)

	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, `oracle script not found: id: 42`)
	require.Nil(t, result)

	// Add Oracle Script
	// os, clear := getTestOracleScript()
	// defer clear()

	// oracleScriptID := k.AddOracleScript(ctx, os)
	// msg = types.NewMsgRequestData(
	// 	types.OracleScriptID(oracleScriptID), calldata, 1, 1, "alice", testapp.Alice.Address,
	// )

	// result, err = oracle.NewHandler(k)(ctx, msg)
	// require.EqualError(t, err, `data source not found: id: 1`)
	// require.Nil(t, result)
}

func TestReportSuccess(t *testing.T) {
	// Setup test environment
	_, ctx, k := testapp.CreateTestInput()

	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))
	calldata := []byte("calldata")

	request := types.NewRequest(1, calldata,
		[]sdk.ValAddress{testapp.Validator1.ValAddress, testapp.Validator2.ValAddress}, 2,
		2, testapp.ParseTime(1581589790), "clientID", []types.RawRequest{
			types.NewRawRequest(1, 1, []byte("beeb")),
			types.NewRawRequest(42, 2, []byte("beeb")),
		},
	)
	k.SetRequest(ctx, 1, request)

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(1581589800, 0))

	msg := types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(1, 0, []byte("data1")),
		types.NewRawReport(42, 0, []byte("data2")),
	}, testapp.Validator1.ValAddress, testapp.Validator1.Address)

	_, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)

	list := k.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{}, list)

	msg = types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(1, 0, []byte("data3")),
		types.NewRawReport(42, 0, []byte("data4")),
	}, testapp.Validator2.ValAddress, testapp.Validator2.Address)

	_, err = oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)

	list = k.GetPendingResolveList(ctx)
	require.Equal(t, []types.RequestID{1}, list)
}

func TestReportFailed(t *testing.T) {
	// Setup test environment
	_, ctx, k := testapp.CreateTestInput()
	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))
	calldata := []byte("calldata")

	request := types.NewRequest(1, calldata,
		[]sdk.ValAddress{testapp.Validator1.ValAddress, testapp.Validator2.ValAddress}, 2,
		2, testapp.ParseTime(1581589790), "clientID", []types.RawRequest{types.NewRawRequest(42, 1, []byte("beeb"))},
	)
	k.SetRequest(ctx, 1, request)

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	// Report by unauthorized reporter
	msg := types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(42, 0, []byte("data1")),
	}, testapp.Validator1.ValAddress, testapp.Alice.Address)
	_, err := oracle.NewHandler(k)(ctx, msg)
	require.Error(t, err)

	// Send wrong external ids
	msg = types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(41, 0, []byte("data1")),
	}, testapp.Validator1.ValAddress, testapp.Validator1.Address)

	// Test only 1 failed case, other case tested in keeper/report_test.go
	_, err = oracle.NewHandler(k)(ctx, msg)
	require.Error(t, err)
}

func TestReportOnExpiredRequest(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	ctx = ctx.WithBlockHeight(2)
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))
	calldata := []byte("calldata")

	request := types.NewRequest(1, calldata,
		[]sdk.ValAddress{testapp.Validator1.ValAddress, testapp.Validator2.ValAddress}, 2,
		2, testapp.ParseTime(1581589790), "clientID", []types.RawRequest{types.NewRawRequest(42, 1, []byte("beeb"))},
	)
	k.SetRequest(ctx, 1, request)

	ctx = ctx.WithBlockHeight(5)
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	k.SetRequestLastExpired(ctx, 1)

	msg := types.NewMsgReportData(1, []types.RawReport{
		types.NewRawReport(42, 0, []byte("data1")),
	}, testapp.Validator1.ValAddress, testapp.Validator1.Address)
	_, err := oracle.NewHandler(k)(ctx, msg)
	require.Equal(t, types.ErrRequestAlreadyExpired, err)
}

func TestAddReporterSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()

	validatorAddress := testapp.Alice.ValAddress
	reporterAddress := testapp.Bob.Address
	// Add testapp.Bob reporter to testapp.Alice validator
	msg := types.NewMsgAddReporter(validatorAddress, reporterAddress)

	result, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, result)

	events := result.Events

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
	_, ctx, k := testapp.CreateTestInput()

	validatorAddress := testapp.Alice.ValAddress
	reporterAddress := testapp.Alice.Address
	msg := types.NewMsgAddReporter(validatorAddress, reporterAddress)

	// Should fail, validator is always a reporter of himself so we can't add testapp.Alice reporter to testapp.Alice validator
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, fmt.Sprintf("reporter already exists: val: %s, addr: %s", validatorAddress.String(), reporterAddress.String()))
	require.Nil(t, result)
}

func TestRemoveReporterSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()

	validatorAddress := testapp.Alice.ValAddress
	reporterAddress := testapp.Bob.Address

	// Add testapp.Bob reporter to testapp.Alice validator
	err := k.AddReporter(ctx, validatorAddress, reporterAddress)
	require.NoError(t, err)

	msg := types.NewMsgRemoveReporter(validatorAddress, reporterAddress)
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)

	events := result.Events

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
	_, ctx, k := testapp.CreateTestInput()

	validatorAddress := testapp.Alice.ValAddress
	reporterAddress := testapp.Bob.Address

	// Should fail, testapp.Bob isn't testapp.Alice validator's reporter
	msg := types.NewMsgRemoveReporter(validatorAddress, reporterAddress)
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, fmt.Sprintf("reporter not found: val: %s, addr: %s", validatorAddress.String(), reporterAddress.String()))
	require.Nil(t, result)
}
