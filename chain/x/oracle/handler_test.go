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
	msgEdit := types.NewMsgEditDataSource(
		1, Owner.Address, newName, newDescription, newExecutable, Owner.Address,
	)
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

	msgEdit := types.NewMsgEditDataSource(
		wrongDID, Owner.Address, newName, newDescription, newExecutable, Owner.Address,
	)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)
	require.Nil(t, res)

	wrongSender := Bob.Address
	msgEdit = types.NewMsgEditDataSource(
		1, Owner.Address, newName, newDescription, newExecutable, wrongSender,
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
		1, Owner.Address, newName, newDescription, wrongGzippedExecutable, Owner.Address,
	)
	res, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestCreateOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := mustGetOwasmCode("beeb.wat")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		Owner.Address, name, description, code, schema, url, Alice.Address,
	)
	res, err := oracle.NewHandler(keeper)(ctx, msg)
	require.NoError(t, err)
	expectEvents := sdk.Events{
		sdk.NewEvent(types.EventTypeCreateOracleScript, sdk.NewAttribute(types.AttributeKeyID, "1")),
	}
	require.Equal(t, expectEvents, res.Events)

	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	filename := keeper.MustGetOracleScript(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

	require.NoError(t, err)

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.NoError(t, err)

	// Code is store must be compiled code.
	codeHash := sha256.Sum256(mustCompileOwasm(code))
	expectFilename := hex.EncodeToString(codeHash[:])
	expectOracleScript := types.NewOracleScript(
		Owner.Address, name, description, expectFilename, schema, url,
	)
	require.Equal(t, expectOracleScript, oracleScript)
}

func TestCreateOracleScriptFailed(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := []byte("non wasm code")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		Owner.Address, name, description, code, schema, url, Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.Error(t, err)
}
func TestCreateGzippedOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := mustGetOwasmCode("beeb.wat")
	schema := "schema"
	url := "url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(code)
	zw.Close()
	gzippedCode := buf.Bytes()

	msg := types.NewMsgCreateOracleScript(
		Owner.Address, name, description, gzippedCode, schema, url, Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.NoError(t, err)

	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	filename := keeper.MustGetOracleScript(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.NoError(t, err)

	// Code is store must be compiled code.
	codeHash := sha256.Sum256(mustCompileOwasm(code))
	expectFilename := hex.EncodeToString(codeHash[:])
	expectOracleScript := types.NewOracleScript(
		Owner.Address, name, description, expectFilename, schema, url,
	)

	require.Equal(t, expectOracleScript, oracleScript)
}

func TestCreateGzippedOracleScriptFail(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := mustGetOwasmCode("beeb.wat")
	schema := "schema"
	url := "url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(code)
	zw.Close()
	gzippedCode := buf.Bytes()[:5]

	msg := types.NewMsgCreateOracleScript(
		Owner.Address, name, description, gzippedCode, schema, url, Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.Error(t, err)
}

func TestEditOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := mustGetOwasmCode("beeb.wat")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		Owner.Address, name, description, code, schema, url, Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.NoError(t, err)

	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	filename := keeper.MustGetOracleScript(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := mustGetOwasmCode("edited_beeb.wat")
	newSchema := "new_schema"
	newURL := "new_url"

	msgEdit := types.NewMsgEditOracleScript(
		oracleScriptID, Owner.Address, newName, newDescription,
		newCode, newSchema, newURL, Owner.Address,
	)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)
	require.NoError(t, err)

	filename = keeper.MustGetOracleScript(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

	expectEvents := sdk.Events{
		sdk.NewEvent(types.EventTypeEditOracleScript, sdk.NewAttribute(types.AttributeKeyID, "1")),
	}
	require.Equal(t, expectEvents, res.Events)

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.NoError(t, err)
	// Code is store must be compiled code.
	codeHash := sha256.Sum256(mustCompileOwasm(newCode))

	expectFilename := hex.EncodeToString(codeHash[:])
	expectOracleScript := types.NewOracleScript(
		Owner.Address, newName, newDescription,
		expectFilename, newSchema, newURL,
	)

	require.Equal(t, expectOracleScript, oracleScript)
}

func TestEditOracleScriptFail(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := mustGetOwasmCode("beeb.wat")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		Owner.Address, name, description, code, schema, url, Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.NoError(t, err)

	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	filename := keeper.MustGetOracleScript(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := mustGetOwasmCode("edited_beeb.wat")
	newSchema := "new_schema"
	newURL := "new_url"

	// No oracle script id 99
	msgEdit := types.NewMsgEditOracleScript(
		types.OracleScriptID(99), Owner.Address, newName, newDescription,
		newCode, newSchema, newURL, Alice.Address,
	)
	_, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)

	// Alice can't edit oracle script
	msgEdit = types.NewMsgEditOracleScript(
		oracleScriptID, Owner.Address, newName, newDescription,
		newCode, newSchema, newURL, Alice.Address,
	)
	_, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)

	// Cannot send bad owasm code
	msgEdit = types.NewMsgEditOracleScript(
		oracleScriptID, Owner.Address, newName, newDescription,
		[]byte("code"), newSchema, newURL, Owner.Address,
	)
	_, err = oracle.NewHandler(keeper)(ctx, msgEdit)
	require.Error(t, err)

}

func TestEditGzippedOracleScriptSuccess(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := mustGetOwasmCode("beeb.wat")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		Owner.Address, name, description, code, schema, url, Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	require.NoError(t, err)

	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	filename := keeper.MustGetOracleScript(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := mustGetOwasmCode("edited_beeb.wat")
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
		oracleScriptID, Owner.Address, newName, newDescription,
		gzippedCode, newSchema, newURL, Owner.Address,
	)
	res, err := oracle.NewHandler(keeper)(ctx, msgEdit)
	require.NoError(t, err)
	expectEvents := sdk.Events{
		sdk.NewEvent(types.EventTypeEditOracleScript, sdk.NewAttribute(types.AttributeKeyID, "1")),
	}
	require.Equal(t, expectEvents, res.Events)
	filename = keeper.MustGetOracleScript(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

	oracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.NoError(t, err)

	codeHash := sha256.Sum256(mustCompileOwasm(newCode))
	expectFilename := hex.EncodeToString(codeHash[:])
	expectOracleScript := types.NewOracleScript(
		Owner.Address, newName, newDescription, expectFilename, newSchema, newURL,
	)
	require.Equal(t, expectOracleScript, oracleScript)
}

func TestEditGzippedOracleScriptFail(t *testing.T) {
	_, ctx, keeper := createTestInput()
	name := "os_1"
	description := "beeb"
	code := mustGetOwasmCode("beeb.wat")
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(
		Owner.Address, name, description, code, schema, url, Alice.Address,
	)
	_, err := oracle.NewHandler(keeper)(ctx, msg)
	dir := filepath.Join(viper.GetString(cli.HomeFlag), "files")
	filename := keeper.MustGetOracleScript(ctx, 1).Filename
	defer deleteFile(filepath.Join(dir, filename))

	require.Nil(t, err)
	oracleScriptID := types.OracleScriptID(1)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := mustGetOwasmCode("edited_beeb.wat")
	newSchema := "new_schema"
	newURL := "new_url"

	// Gzipped executable file
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(newCode)
	zw.Close()
	gzippedCode := buf.Bytes()[:5]

	msgEdit := types.NewMsgEditOracleScript(
		oracleScriptID, Owner.Address, newName, newDescription,
		gzippedCode, newSchema, newURL, Owner.Address,
	)
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
			sdk.NewAttribute(types.AttributeKeyOracleScriptID, "1"),
			sdk.NewAttribute(types.AttributeKeyCalldata, "62656562"), // "beeb" in hex
			sdk.NewAttribute(types.AttributeKeyAskCount, "2"),
			sdk.NewAttribute(types.AttributeKeyMinCount, "2"),
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

	require.Equal(t, expectEvents, result.Events)
}

func TestRequestDataFail(t *testing.T) {
	_, ctx, k := createTestInput()

	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))

	wrongOracleScript := 1

	calldata := []byte("test")
	msg := types.NewMsgRequestData(
		types.OracleScriptID(wrongOracleScript), calldata, 1, 1, "alice", Alice.Address,
	)

	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, `oracle script not found: id: 1`)
	require.Nil(t, result)

	// Add Oracle Script
	os, clear := getTestOracleScript()
	defer clear()

	oracleScriptID := k.AddOracleScript(ctx, os)
	msg = types.NewMsgRequestData(
		types.OracleScriptID(oracleScriptID), calldata, 1, 1, "alice", Alice.Address,
	)

	result, err = oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, `data source not found: id: 1`)
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
		2, 1581589790, "clientID", nil, []types.RawRequest{
			types.NewRawRequest(1, 1, []byte("beeb")),
			types.NewRawRequest(42, 2, []byte("beeb")),
		},
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
		2, 1581589790, "clientID", nil, []types.RawRequest{types.NewRawRequest(42, 1, []byte("beeb"))},
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

func TestAddReporterSuccess(t *testing.T) {
	_, ctx, k := createTestInput()

	validatorAddress := Alice.ValAddress
	reporterAddress := Bob.Address
	// Add Bob reporter to Alice validator
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
	_, ctx, k := createTestInput()

	validatorAddress := Alice.ValAddress
	reporterAddress := Alice.Address
	msg := types.NewMsgAddReporter(validatorAddress, reporterAddress)

	// Should fail, validator is always a reporter of himself so we can't add Alice reporter to Alice validator
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, fmt.Sprintf("reporter already exists: val: %s, addr: %s", validatorAddress.String(), reporterAddress.String()))
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
	_, ctx, k := createTestInput()

	validatorAddress := Alice.ValAddress
	reporterAddress := Bob.Address

	// Should fail, Bob isn't Alice validator's reporter
	msg := types.NewMsgRemoveReporter(validatorAddress, reporterAddress)
	result, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, fmt.Sprintf("reporter not found: val: %s, addr: %s", validatorAddress.String(), reporterAddress.String()))
	require.Nil(t, result)
}
