package keeper_test

import (
	"encoding/hex"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestPrepareRequestSuccess(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))

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
	askCount := uint64(1)
	minCount := uint64(1)
	clientID := "beeb"
	requestHeight := int64(0)
	rawRequests := []types.RawRequest{
		types.NewRawRequest(1, 1, []byte("beeb")),
		types.NewRawRequest(2, 2, []byte("beeb")),
		types.NewRawRequest(3, 3, []byte("beeb")),
	}

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m)
	require.NoError(t, err)

	req, err := k.GetRequest(ctx, 1)
	require.NoError(t, err)
	expectReq := types.NewRequest(oracleScriptID, calldata, []sdk.ValAddress{Validator1.ValAddress}, minCount,
		requestHeight, int64(1581589790), clientID, rawRequests)
	require.Equal(t, expectReq, req)
	expectEvents := sdk.Events{
		sdk.NewEvent(
			types.EventTypeRequest,
			sdk.NewAttribute(types.AttributeKeyID, "1"),
			sdk.NewAttribute(types.AttributeKeyOracleScriptID, "1"),
			sdk.NewAttribute(types.AttributeKeyCalldata, "62656562"), // "beeb" in hex
			sdk.NewAttribute(types.AttributeKeyAskCount, "1"),
			sdk.NewAttribute(types.AttributeKeyMinCount, "1"),
			sdk.NewAttribute(types.AttributeKeyValidator, Validator1.ValAddress.String()),
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
	require.Equal(t, expectEvents, ctx.EventManager().Events())
}

func TestPrepareRequestInvalidAskCountFail(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))
	k.SetParam(ctx, types.KeyMaxAskCount, 1000) // Set MaxAskCount 1000

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
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := uint64(100000) // Set ask count 100000
	minCount := uint64(2)
	clientID := "beeb"

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m)
	require.Error(t, err)
}

func TestPrepareRequestBaseRequestFeePanic(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(90000)) // Set Gas Meter 90000

	baseRequestGas := uint64(100000)
	k.SetParam(ctx, types.KeyBaseRequestGas, baseRequestGas) // Set BaseRequestGas 100000

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
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := uint64(1)
	minCount := uint64(2)
	clientID := "beeb"

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	require.Panics(t, func() { k.PrepareRequest(ctx, &m) })

	ctx = ctx.WithGasMeter(sdk.NewGasMeter(200000))
	err := k.PrepareRequest(ctx, &m)
	require.NoError(t, err)

}

func TestPrepareRequestPerValidatorRequestFeePanic(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(150000)) //Set Gas Meter 150000

	baseRequestGas := uint64(100000)
	k.SetParam(ctx, types.KeyBaseRequestGas, baseRequestGas) // Set BaseRequestGas 100000
	perValidatorRequestGas := uint64(100000)
	k.SetParam(ctx, types.KeyPerValidatorRequestGas, perValidatorRequestGas) // Set PerValidatorRequestGas 100000

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
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := uint64(1)
	minCount := uint64(2)
	clientID := "beeb"

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)

	// PrepareRequest panics because set gas meter at 150000
	// but PrepareRequest consume gas more than 200000
	// (baseRequestGas + askCount*perValidatorRequestGas = 200000)
	require.Panics(t, func() { k.PrepareRequest(ctx, &m) })
}

func TestPrepareRequestGetRandomValidatorsFail(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))

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
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := uint64(15)
	minCount := uint64(2)
	clientID := "beeb"

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)

	err := k.PrepareRequest(ctx, &m)
	require.Error(t, err)
}

func TestPrepareRequestGetOracleScriptFail(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))

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

	k.AddOracleScript(ctx, os)
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := uint64(1)
	minCount := uint64(2)
	clientID := "beeb"

	m := types.NewMsgRequestData(9999, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m)
	require.Error(t, err)
}

func TestPrepareRequestBadWasmExecutionFail(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))

	os, clear4 := getBadOracleScript()
	defer clear4()

	oracleScriptID := k.AddOracleScript(ctx, os)
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := uint64(1)
	minCount := uint64(2)
	clientID := "beeb"

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m, nil)
	require.EqualError(t, err, `bad wasm execution: failed to prepare request with error: get external data on non-execution period`)
}

func TestPrepareRequestGetDataSourceFail(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))

	os, clear4 := getTestOracleScript()
	defer clear4()

	oracleScriptID := k.AddOracleScript(ctx, os)
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := uint64(1)
	minCount := uint64(2)
	clientID := "beeb"

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m)
	require.Error(t, err)
}

func TestResolveRequestSuccess(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))

	os, clear := getTestOracleScript()
	defer clear()
	oracleScriptID := k.AddOracleScript(ctx, os)

	calldata := []byte("calldata")
	minCount := uint64(1)
	clientID := "owasm test"
	rawRequests := []types.RawRequest{
		types.NewRawRequest(1, 1, []byte("beeb")),
		types.NewRawRequest(2, 2, []byte("beeb")),
		types.NewRawRequest(3, 3, []byte("beeb")),
	}
	vals := []sdk.ValAddress{Validator1.ValAddress, Validator2.ValAddress, Validator3.ValAddress}
	requestHeight := int64(4000)
	requestTime := int64(1581589700)

	req := types.NewRequest(
		oracleScriptID, calldata, vals, minCount, requestHeight,
		requestTime, clientID, rawRequests)
	reqID := k.AddRequest(ctx, req)

	// Set report validator
	k.AddReport(ctx, types.RequestID(1), types.NewReport(
		Validator1.ValAddress,
		true,
		[]types.RawReport{
			types.NewRawReport(1, 0, []byte("answer1")),
			types.NewRawReport(2, 0, []byte("answer2")),
			types.NewRawReport(3, 0, []byte("answer3")),
		}),
	)

	k.ResolveRequest(ctx, reqID)
	res, err := k.GetResult(ctx, reqID)
	require.NoError(t, err)

	reqPacket := types.NewOracleRequestPacketData(
		clientID, types.OracleScriptID(1), calldata, uint64(len(vals)), minCount,
	)

	resPacket := types.NewOracleResponsePacketData(
		clientID, reqID, 1, requestTime,
		ctx.BlockTime().Unix(), types.ResolveStatus_Success, []byte("beeb"),
	)
	require.Equal(t, types.Result{RequestPacketData: reqPacket, ResponsePacketData: resPacket}, res)

}

// TODO: Patch to "Bad" wasm code that is a valid wasm code.
// func TestResolveRequestFail(t *testing.T) {
// 	_, ctx, k := createTestInput()
// 	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))

// 	ds1, clear1 := getTestDataSource("code1")
// 	defer clear1()
// 	k.AddDataSource(ctx, ds1)

// 	ds2, clear2 := getTestDataSource("code2")
// 	defer clear2()
// 	k.AddDataSource(ctx, ds2)

// 	ds3, clear3 := getTestDataSource("code3")
// 	defer clear3()
// 	k.AddDataSource(ctx, ds3)

// 	os, clear4 := getBadOracleScript()
// 	defer clear4()

// 	oracleScriptID := k.AddOracleScript(ctx, os)
// 	calldata, _ := hex.DecodeString("00")
// 	minCount := uint64(2)
// 	clientID := "beeb"
// 	rawRequestID := []types.ExternalID{1, 2}
// 	vals := []sdk.ValAddress{Validator1.ValAddress}
// 	requestHeight := int64(4000)
// 	requestTime := int64(1581589700)

// 	req := types.NewRequest(
// 		oracleScriptID, calldata, vals, minCount, requestHeight,
// 		requestTime, clientID, nil, rawRequestID)
// 	reqID := k.AddRequest(ctx, req)
// 	k.ResolveRequest(ctx, reqID)

// 	res, err := k.GetResult(ctx, reqID)
// 	require.NoError(t, err)

// 	r := k.MustGetRequest(ctx, reqID)
// 	resPacket := types.NewOracleResponsePacketData(
// 		r.ClientID, reqID, k.GetReportCount(ctx, reqID), r.RequestTime,
// 		ctx.BlockTime().Unix(), types.ResolveStatus_Failure, nil,
// 	)
// 	reqPacket := types.NewOracleRequestPacketData(
// 		r.ClientID, r.OracleScriptID, r.Calldata, r.MinCount, uint64(len(r.RequestedValidators)),
// 	)
// 	expecetRes := obi.MustEncode(reqPacket, resPacket)

// 	require.Equal(t, expecetRes, res)
// }
