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
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := uint64(1)
	minCount := uint64(2)
	clientID := "beeb"
	requestHeight := int64(0)
	rawRequestID := []types.ExternalID{1, 2, 3}

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m, nil)
	require.NoError(t, err)

	events := ctx.EventManager().Events()

	req, err := k.GetRequest(ctx, 1)
	require.NoError(t, err)

	expectReq := types.NewRequest(oracleScriptID, calldata, []sdk.ValAddress{Validator1.ValAddress}, minCount,
		requestHeight, int64(1581589790), clientID, nil, rawRequestID)
	require.Equal(t, expectReq, req)

	require.Equal(t, 4, len(events))
	require.Equal(t, types.EventTypeRequest, events[0].Type)
	require.Equal(t, []byte(types.AttributeKeyID), events[0].Attributes[0].Key)
	require.Equal(t, []byte("1"), events[0].Attributes[0].Value)

	require.Equal(t, []byte(types.AttributeKeyValidator), events[0].Attributes[1].Key)
	require.Equal(t, []byte(Validator1.ValAddress.String()), events[0].Attributes[1].Value)

	expectEvents := []struct {
		eventType string
		dsID      []byte
		filname   []byte
		exID      []byte
		calldata  []byte
	}{
		{types.EventTypeRawRequest, []byte("1"), []byte(ds1.Filename), []byte("1"), []byte("beeb")},
		{types.EventTypeRawRequest, []byte("2"), []byte(ds2.Filename), []byte("2"), []byte("beeb")},
		{types.EventTypeRawRequest, []byte("3"), []byte(ds3.Filename), []byte("3"), []byte("beeb")},
	}

	for idx, expectEvent := range expectEvents {
		require.Equal(t, expectEvent.eventType, events[idx+1].Type)
		require.Equal(t, []byte(types.AttributeKeyDataSourceID), events[idx+1].Attributes[0].Key)
		require.Equal(t, expectEvent.dsID, events[idx+1].Attributes[0].Value)
		require.Equal(t, []byte(types.AttributeKeyDataSourceHash), events[idx+1].Attributes[1].Key)
		require.Equal(t, expectEvent.filname, events[idx+1].Attributes[1].Value)
		require.Equal(t, []byte(types.AttributeKeyExternalID), events[idx+1].Attributes[2].Key)
		require.Equal(t, expectEvent.exID, events[idx+1].Attributes[2].Value)
		require.Equal(t, []byte(types.AttributeKeyCalldata), events[idx+1].Attributes[3].Key)
		require.Equal(t, expectEvent.calldata, events[idx+1].Attributes[3].Value)
	}
}

func TestPrepareRequestInvalidAskCountFail(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))
	k.SetParam(ctx, types.KeyMaxAskCount, 1000)

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
	askCount := uint64(100000)
	minCount := uint64(2)
	clientID := "beeb"

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m, nil)
	require.Error(t, err)
}

func TestPrepareRequestBaseRequestFeePanic(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(100000))

	baseRequestGas := uint64(100000)
	k.SetParam(ctx, types.KeyBaseRequestGas, baseRequestGas)

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

	require.Panics(t, func() { k.PrepareRequest(ctx, &m, nil) })
}

func TestPrepareRequestPerValidatorRequestFeePanic(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(150000))

	baseRequestGas := uint64(100000)
	k.SetParam(ctx, types.KeyBaseRequestGas, baseRequestGas)
	perValidatorRequestGas := uint64(100000)
	k.SetParam(ctx, types.KeyPerValidatorRequestGas, perValidatorRequestGas)

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

	require.Panics(t, func() { k.PrepareRequest(ctx, &m, nil) })
}

func TestPrepareRequestGetRandomValidatorsFail(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(200000000))

	k.SetParam(ctx, types.KeyMaxAskCount, 16)
	baseRequestGas := uint64(100000)
	k.SetParam(ctx, types.KeyBaseRequestGas, baseRequestGas)
	perValidatorRequestGas := uint64(100000)
	k.SetParam(ctx, types.KeyPerValidatorRequestGas, perValidatorRequestGas)

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

	// test consume gas for data requests
	expectGasNotLessThan := baseRequestGas + askCount*perValidatorRequestGas
	beforeGas := ctx.GasMeter().GasConsumed()
	err := k.PrepareRequest(ctx, &m, nil)
	afterGas := ctx.GasMeter().GasConsumed()
	require.Greater(t, afterGas-beforeGas, expectGasNotLessThan)

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
	err := k.PrepareRequest(ctx, &m, nil)
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
	require.Error(t, err)
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
	err := k.PrepareRequest(ctx, &m, nil)
	require.Error(t, err)
}

func TestResolveRequestSuccess(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(1581589790, 0))

	ds1, clear1 := getTestDataSource("code1")
	defer clear1()
	k.AddDataSource(ctx, ds1)

	ds2, clear2 := getTestDataSource("code2")
	defer clear2()
	k.AddDataSource(ctx, ds2)

	os, clear3 := getTestOracleScript()
	defer clear3()

	oracleScriptID := k.AddOracleScript(ctx, os)
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	minCount := uint64(2)
	clientID := "beeb"
	ibcInfo := types.NewIBCInfo("source_port", "source_channel")
	rawRequestID := []types.ExternalID{1, 2}
	vals := []sdk.ValAddress{Validator1.ValAddress}
	requestHeight := int64(4000)
	requestTime := int64(1581589700)

	req := types.NewRequest(
		oracleScriptID, calldata, vals, minCount, requestHeight,
		requestTime, clientID, &ibcInfo, rawRequestID)
	reqID := k.AddRequest(ctx, req)
	k.ResolveRequest(ctx, reqID)

	res, err := k.GetResult(ctx, reqID)
	require.NoError(t, err)

	r := k.MustGetRequest(ctx, reqID)
	resPacket := types.NewOracleResponsePacketData(
		r.ClientID, reqID, k.GetReportCount(ctx, reqID), r.RequestTime,
		ctx.BlockTime().Unix(), types.ResolveStatus_Success, []byte("beeb"),
	)
	reqPacket := types.NewOracleRequestPacketData(
		r.ClientID, r.OracleScriptID, r.Calldata, r.MinCount, uint64(len(r.RequestedValidators)),
	)
	expecetRes := k.AddResult(ctx, reqID, reqPacket, resPacket)

	require.Equal(t, expecetRes, res)
}

func TestResolveRequestFail(t *testing.T) {
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

	os, clear4 := getBadOracleScript()
	defer clear4()

	oracleScriptID := k.AddOracleScript(ctx, os)
	calldata, _ := hex.DecodeString("00")
	minCount := uint64(2)
	clientID := "beeb"
	rawRequestID := []types.ExternalID{1, 2}
	vals := []sdk.ValAddress{Validator1.ValAddress}
	requestHeight := int64(4000)
	requestTime := int64(1581589700)

	req := types.NewRequest(
		oracleScriptID, calldata, vals, minCount, requestHeight,
		requestTime, clientID, nil, rawRequestID)
	reqID := k.AddRequest(ctx, req)
	k.ResolveRequest(ctx, reqID)

	res, err := k.GetResult(ctx, reqID)
	require.NoError(t, err)

	r := k.MustGetRequest(ctx, reqID)
	resPacket := types.NewOracleResponsePacketData(
		r.ClientID, reqID, k.GetReportCount(ctx, reqID), r.RequestTime,
		ctx.BlockTime().Unix(), types.ResolveStatus_Failure, nil,
	)
	reqPacket := types.NewOracleRequestPacketData(
		r.ClientID, r.OracleScriptID, r.Calldata, r.MinCount, uint64(len(r.RequestedValidators)),
	)
	expecetRes := k.AddResult(ctx, reqID, reqPacket, resPacket)

	require.Equal(t, expecetRes, res)
}
