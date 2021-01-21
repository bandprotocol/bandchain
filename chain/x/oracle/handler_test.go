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

	"github.com/GeoDB-Limited/odincore/chain/x/oracle"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/testapp"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

func TestCreateDataSourceSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	dsCount := k.GetDataSourceCount(ctx)
	owner := testapp.Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	executableHash := sha256.Sum256(executable)
	filename := hex.EncodeToString(executableHash[:])
	msg := types.NewMsgCreateDataSource(owner, name, description, executable, testapp.Alice.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	ds, err := k.GetDataSource(ctx, types.DataSourceID(dsCount+1))
	require.NoError(t, err)
	require.Equal(t, types.NewDataSource(testapp.Owner.Address, name, description, filename), ds)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeCreateDataSource,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", dsCount+1)),
	)}, res.Events)
}

func TestCreateGzippedExecutableDataSourceSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	dsCount := k.GetDataSourceCount(ctx)
	owner := testapp.Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	executableHash := sha256.Sum256(executable)
	filename := hex.EncodeToString(executableHash[:])
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(executable)
	zw.Close()
	msg := types.NewMsgCreateDataSource(owner, name, description, buf.Bytes(), testapp.Alice.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	ds, err := k.GetDataSource(ctx, types.DataSourceID(dsCount+1))
	require.NoError(t, err)
	require.Equal(t, types.NewDataSource(testapp.Owner.Address, name, description, filename), ds)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeCreateDataSource,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", dsCount+1)),
	)}, res.Events)
}

func TestCreateGzippedExecutableDataSourceFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	owner := testapp.Owner.Address
	name := "data_source_1"
	description := "description"
	executable := []byte("executable")
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(executable)
	zw.Close()
	sender := testapp.Alice.Address
	msg := types.NewMsgCreateDataSource(owner, name, description, buf.Bytes()[:5], sender)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "uncompression failed: unexpected EOF")
	require.Nil(t, res)
}

func TestEditDataSourceSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	newName := "beeb"
	newDescription := "new_description"
	newExecutable := []byte("executable2")
	newExecutableHash := sha256.Sum256(newExecutable)
	newFilename := hex.EncodeToString(newExecutableHash[:])
	msg := types.NewMsgEditDataSource(1, testapp.Owner.Address, newName, newDescription, newExecutable, testapp.Owner.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	ds, err := k.GetDataSource(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, types.NewDataSource(testapp.Owner.Address, newName, newDescription, newFilename), ds)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeEditDataSource,
		sdk.NewAttribute(types.AttributeKeyID, "1"),
	)}, res.Events)
}

func TestEditDataSourceFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	newName := "beeb"
	newDescription := "new_description"
	newExecutable := []byte("executable2")
	// Bad ID
	msg := types.NewMsgEditDataSource(42, testapp.Owner.Address, newName, newDescription, newExecutable, testapp.Owner.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "data source not found: id: 42")
	require.Nil(t, res)
	// Not owner
	msg = types.NewMsgEditDataSource(1, testapp.Owner.Address, newName, newDescription, newExecutable, testapp.Bob.Address)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "editor not authorized")
	require.Nil(t, res)
	// Bad Gzip
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(newExecutable)
	zw.Close()
	msg = types.NewMsgEditDataSource(1, testapp.Owner.Address, newName, newDescription, buf.Bytes()[:5], testapp.Owner.Address)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "uncompression failed: unexpected EOF")
	require.Nil(t, res)
}

func TestCreateOracleScriptSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	osCount := k.GetOracleScriptCount(ctx)
	name := "os_1"
	description := "beeb"
	code := testapp.WasmExtra1
	schema := "schema"
	url := "url"
	msg := types.NewMsgCreateOracleScript(testapp.Owner.Address, name, description, code, schema, url, testapp.Alice.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	os, err := k.GetOracleScript(ctx, types.OracleScriptID(osCount+1))
	require.NoError(t, err)
	require.Equal(t, types.NewOracleScript(testapp.Owner.Address, name, description, testapp.WasmExtra1FileName, schema, url), os)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeCreateOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", osCount+1)),
	)}, res.Events)
}

func TestCreateGzippedOracleScriptSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	osCount := k.GetOracleScriptCount(ctx)
	name := "os_1"
	description := "beeb"
	schema := "schema"
	url := "url"
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(testapp.WasmExtra1)
	zw.Close()
	msg := types.NewMsgCreateOracleScript(testapp.Owner.Address, name, description, buf.Bytes(), schema, url, testapp.Alice.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	os, err := k.GetOracleScript(ctx, types.OracleScriptID(osCount+1))
	require.NoError(t, err)
	require.Equal(t, types.NewOracleScript(testapp.Owner.Address, name, description, testapp.WasmExtra1FileName, schema, url), os)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeCreateOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", osCount+1)),
	)}, res.Events)
}

func TestCreateOracleScriptFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	name := "os_1"
	description := "beeb"
	schema := "schema"
	url := "url"
	// Bad Owasm code
	msg := types.NewMsgCreateOracleScript(testapp.Owner.Address, name, description, []byte("BAD"), schema, url, testapp.Alice.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "owasm compilation failed: with error: wasm code does not pass basic validation")
	require.Nil(t, res)
	// Bad Gzip
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(testapp.WasmExtra1)
	zw.Close()
	msg = types.NewMsgCreateOracleScript(testapp.Owner.Address, name, description, buf.Bytes()[:5], schema, url, testapp.Alice.Address)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "uncompression failed: unexpected EOF")
	require.Nil(t, res)
}

func TestEditOracleScriptSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := testapp.WasmExtra2
	newSchema := "new_schema"
	newURL := "new_url"
	msg := types.NewMsgEditOracleScript(1, testapp.Owner.Address, newName, newDescription, newCode, newSchema, newURL, testapp.Owner.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	os, err := k.GetOracleScript(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, types.NewOracleScript(testapp.Owner.Address, newName, newDescription, testapp.WasmExtra2FileName, newSchema, newURL), os)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeEditOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, "1"),
	)}, res.Events)
}

func TestEditOracleScriptFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	newName := "os_2"
	newDescription := "beebbeeb"
	newCode := testapp.WasmExtra2
	newSchema := "new_schema"
	newURL := "new_url"
	// Bad ID
	msg := types.NewMsgEditOracleScript(999, testapp.Owner.Address, newName, newDescription, newCode, newSchema, newURL, testapp.Owner.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "oracle script not found: id: 999")
	require.Nil(t, res)
	// Not owner
	msg = types.NewMsgEditOracleScript(1, testapp.Owner.Address, newName, newDescription, newCode, newSchema, newURL, testapp.Bob.Address)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "editor not authorized")
	require.Nil(t, res)
	// Bad Owasm code
	msg = types.NewMsgEditOracleScript(1, testapp.Owner.Address, newName, newDescription, []byte("BAD_CODE"), newSchema, newURL, testapp.Owner.Address)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "owasm compilation failed: with error: wasm code does not pass basic validation")
	require.Nil(t, res)
	// Bad Gzip
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(testapp.WasmExtra2)
	zw.Close()
	msg = types.NewMsgEditOracleScript(1, testapp.Owner.Address, newName, newDescription, buf.Bytes()[:5], newSchema, newURL, testapp.Owner.Address)
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "uncompression failed: unexpected EOF")
	require.Nil(t, res)
}

func TestRequestDataSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	ctx = ctx.WithBlockHeight(124).WithBlockTime(testapp.ParseTime(1581589790))
	msg := types.NewMsgRequestData(1, []byte("beeb"), 2, 2, "CID", testapp.Alice.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, types.NewRequest(
		1,
		[]byte("beeb"),
		[]sdk.ValAddress{testapp.Validator3.ValAddress, testapp.Validator1.ValAddress},
		2,
		124,
		testapp.ParseTime(1581589790),
		"CID",
		[]types.RawRequest{
			types.NewRawRequest(1, 1, []byte("beeb")),
			types.NewRawRequest(2, 2, []byte("beeb")),
			types.NewRawRequest(3, 3, []byte("beeb")),
		},
	), k.MustGetRequest(ctx, 1))
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeRequest,
		sdk.NewAttribute(types.AttributeKeyID, "1"),
		sdk.NewAttribute(types.AttributeKeyClientID, "CID"),
		sdk.NewAttribute(types.AttributeKeyOracleScriptID, "1"),
		sdk.NewAttribute(types.AttributeKeyCalldata, "62656562"), // "beeb" in hex
		sdk.NewAttribute(types.AttributeKeyAskCount, "2"),
		sdk.NewAttribute(types.AttributeKeyMinCount, "2"),
		sdk.NewAttribute(types.AttributeKeyGasUsed, "785"),
		sdk.NewAttribute(types.AttributeKeyValidator, testapp.Validator3.ValAddress.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, testapp.Validator1.ValAddress.String()),
	), sdk.NewEvent(
		types.EventTypeRawRequest,
		sdk.NewAttribute(types.AttributeKeyDataSourceID, "1"),
		sdk.NewAttribute(types.AttributeKeyDataSourceHash, testapp.DataSources[1].Filename),
		sdk.NewAttribute(types.AttributeKeyExternalID, "1"),
		sdk.NewAttribute(types.AttributeKeyCalldata, "beeb"),
	), sdk.NewEvent(
		types.EventTypeRawRequest,
		sdk.NewAttribute(types.AttributeKeyDataSourceID, "2"),
		sdk.NewAttribute(types.AttributeKeyDataSourceHash, testapp.DataSources[2].Filename),
		sdk.NewAttribute(types.AttributeKeyExternalID, "2"),
		sdk.NewAttribute(types.AttributeKeyCalldata, "beeb"),
	), sdk.NewEvent(
		types.EventTypeRawRequest,
		sdk.NewAttribute(types.AttributeKeyDataSourceID, "3"),
		sdk.NewAttribute(types.AttributeKeyDataSourceHash, testapp.DataSources[3].Filename),
		sdk.NewAttribute(types.AttributeKeyExternalID, "3"),
		sdk.NewAttribute(types.AttributeKeyCalldata, "beeb"),
	)}, res.Events)
}

func TestRequestDataFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	// No active oracle validators
	res, err := oracle.NewHandler(k)(ctx, types.NewMsgRequestData(1, []byte("beeb"), 2, 2, "CID", testapp.Alice.Address))
	require.EqualError(t, err, "insufficent available validators: 0 < 2")
	require.Nil(t, res)
	k.Activate(ctx, testapp.Validator1.ValAddress)
	k.Activate(ctx, testapp.Validator2.ValAddress)
	// Too high ask count
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgRequestData(1, []byte("beeb"), 3, 2, "CID", testapp.Alice.Address))
	require.EqualError(t, err, "insufficent available validators: 2 < 3")
	require.Nil(t, res)
	// Bad oracle script ID
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgRequestData(999, []byte("beeb"), 2, 2, "CID", testapp.Alice.Address))
	require.EqualError(t, err, "oracle script not found: id: 999")
	require.Nil(t, res)
}

func TestReportSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Set up a mock request asking 3 validators with min count 2.
	k.SetRequest(ctx, 42, types.NewRequest(
		1,
		[]byte("beeb"),
		[]sdk.ValAddress{testapp.Validator3.ValAddress, testapp.Validator2.ValAddress, testapp.Validator1.ValAddress},
		2,
		124,
		testapp.ParseTime(1581589790),
		"CID",
		[]types.RawRequest{
			types.NewRawRequest(1, 1, []byte("beeb")),
			types.NewRawRequest(2, 2, []byte("beeb")),
		},
	))
	// Common raw reports for everyone.
	reports := []types.RawReport{types.NewRawReport(1, 0, []byte("data1")), types.NewRawReport(2, 0, []byte("data2"))}
	// Validator1 reports data.
	res, err := oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Validator1.ValAddress, testapp.Validator1.Address))
	require.NoError(t, err)
	require.Equal(t, []types.RequestID{}, k.GetPendingResolveList(ctx))
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeReport,
		sdk.NewAttribute(types.AttributeKeyID, "42"),
		sdk.NewAttribute(types.AttributeKeyValidator, testapp.Validator1.ValAddress.String()),
	)}, res.Events)
	// Validator2 reports data. Now the request should move to pending resolve.
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Validator2.ValAddress, testapp.Validator2.Address))
	require.NoError(t, err)
	require.Equal(t, []types.RequestID{42}, k.GetPendingResolveList(ctx))
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeReport,
		sdk.NewAttribute(types.AttributeKeyID, "42"),
		sdk.NewAttribute(types.AttributeKeyValidator, testapp.Validator2.ValAddress.String()),
	)}, res.Events)
	// Even if we resolve the request, validator3 should still be able to report.
	k.SetPendingResolveList(ctx, []types.RequestID{})
	k.ResolveSuccess(ctx, 42, []byte("RESOLVE_RESULT!"), 1234)
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Validator3.ValAddress, testapp.Validator3.Address))
	require.NoError(t, err)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeReport,
		sdk.NewAttribute(types.AttributeKeyID, "42"),
		sdk.NewAttribute(types.AttributeKeyValidator, testapp.Validator3.ValAddress.String()),
	)}, res.Events)
	// Check the reports of this request. We should see 3 reports, with report from validator3 comes after resolve.
	require.Contains(t, k.GetReports(ctx, 42), types.NewReport(testapp.Validator1.ValAddress, true, reports))
	require.Contains(t, k.GetReports(ctx, 42), types.NewReport(testapp.Validator2.ValAddress, true, reports))
	require.Contains(t, k.GetReports(ctx, 42), types.NewReport(testapp.Validator3.ValAddress, false, reports))
}

func TestReportFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Set up a mock request asking 3 validators with min count 2.
	k.SetRequest(ctx, 42, types.NewRequest(
		1,
		[]byte("beeb"),
		[]sdk.ValAddress{testapp.Validator3.ValAddress, testapp.Validator2.ValAddress, testapp.Validator1.ValAddress},
		2,
		124,
		testapp.ParseTime(1581589790),
		"CID",
		[]types.RawRequest{
			types.NewRawRequest(1, 1, []byte("beeb")),
			types.NewRawRequest(2, 2, []byte("beeb")),
		},
	))
	// Common raw reports for everyone.
	reports := []types.RawReport{types.NewRawReport(1, 0, []byte("data1")), types.NewRawReport(2, 0, []byte("data2"))}
	// Bad ID
	res, err := oracle.NewHandler(k)(ctx, types.NewMsgReportData(999, reports, testapp.Validator1.ValAddress, testapp.Validator1.Address))
	require.EqualError(t, err, "request not found: id: 999")
	require.Nil(t, res)
	// Not-asked validator
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Alice.ValAddress, testapp.Alice.Address))
	require.EqualError(t, err, fmt.Sprintf("validator not requested: reqID: 42, val: %s", testapp.Alice.ValAddress.String()))
	require.Nil(t, res)
	// Not an authorized reporter
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Validator1.ValAddress, testapp.Alice.Address))
	require.EqualError(t, err, "reporter not authorized")
	require.Nil(t, res)
	// Not having all raw reports
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, []types.RawReport{types.NewRawReport(1, 0, []byte("data1"))}, testapp.Validator1.ValAddress, testapp.Validator1.Address))
	require.EqualError(t, err, "invalid report size")
	require.Nil(t, res)
	// Incorrect external IDs
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, []types.RawReport{types.NewRawReport(1, 0, []byte("data1")), types.NewRawReport(42, 0, []byte("data2"))}, testapp.Validator1.ValAddress, testapp.Validator1.Address))
	require.EqualError(t, err, "raw request not found: reqID: 42, extID: 42")
	require.Nil(t, res)
	// Request already expired
	k.SetRequestLastExpired(ctx, 42)
	res, err = oracle.NewHandler(k)(ctx, types.NewMsgReportData(42, reports, testapp.Validator1.ValAddress, testapp.Validator1.Address))
	require.EqualError(t, err, "request already expired")
	require.Nil(t, res)
}

func TestActivateSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	ctx = ctx.WithBlockTime(testapp.ParseTime(1000000))
	require.Equal(t,
		types.NewValidatorStatus(false, time.Time{}),
		k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress),
	)
	msg := types.NewMsgActivate(testapp.Validator1.ValAddress)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	require.Equal(t,
		types.NewValidatorStatus(true, testapp.ParseTime(1000000)),
		k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress),
	)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeActivate,
		sdk.NewAttribute(types.AttributeKeyValidator, testapp.Validator1.ValAddress.String()),
	)}, res.Events)
}

func TestActivateFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	msg := types.NewMsgActivate(testapp.Validator1.ValAddress)
	// Already active.
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "validator already active")
	require.Nil(t, res)
	// Too soon to activate.
	ctx = ctx.WithBlockTime(testapp.ParseTime(100000))
	k.MissReport(ctx, testapp.Validator1.ValAddress, testapp.ParseTime(99999))
	ctx = ctx.WithBlockTime(testapp.ParseTime(100001))
	res, err = oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, "too soon to activate")
	require.Nil(t, res)
	// OK
	ctx = ctx.WithBlockTime(testapp.ParseTime(200000))
	_, err = oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
}

func TestAddReporterSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	require.False(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address))
	// Add testapp.Bob to a reporter of testapp.Alice validator.
	msg := types.NewMsgAddReporter(testapp.Alice.ValAddress, testapp.Bob.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	require.True(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address))
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeAddReporter,
		sdk.NewAttribute(types.AttributeKeyValidator, testapp.Alice.ValAddress.String()),
		sdk.NewAttribute(types.AttributeKeyReporter, testapp.Bob.Address.String()),
	)}, res.Events)
}

func TestAddReporterFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	// Should fail when you try to add yourself as your reporter.
	msg := types.NewMsgAddReporter(testapp.Alice.ValAddress, testapp.Alice.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, fmt.Sprintf("reporter already exists: val: %s, addr: %s", testapp.Alice.ValAddress.String(), testapp.Alice.Address.String()))
	require.Nil(t, res)
}

func TestRemoveReporterSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	// Add testapp.Bob to a reporter of testapp.Alice validator.
	err := k.AddReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.True(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address))
	require.NoError(t, err)
	// Now remove testapp.Bob from the set of testapp.Alice's reporters.
	msg := types.NewMsgRemoveReporter(testapp.Alice.ValAddress, testapp.Bob.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.NoError(t, err)
	require.False(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address))
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeRemoveReporter,
		sdk.NewAttribute(types.AttributeKeyValidator, testapp.Alice.ValAddress.String()),
		sdk.NewAttribute(types.AttributeKeyReporter, testapp.Bob.Address.String()),
	)}, res.Events)
}

func TestRemoveReporterFail(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	// Should fail because testapp.Bob isn't testapp.Alice validator's reporter.
	msg := types.NewMsgRemoveReporter(testapp.Alice.ValAddress, testapp.Bob.Address)
	res, err := oracle.NewHandler(k)(ctx, msg)
	require.EqualError(t, err, fmt.Sprintf("reporter not found: val: %s, addr: %s", testapp.Alice.ValAddress.String(), testapp.Bob.Address.String()))
	require.Nil(t, res)
}
