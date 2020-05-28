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
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := int64(1)
	minCount := int64(2)
	clientID := "beeb"
	ibcInfo := types.NewIBCInfo("source_port", "source_channel")

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m, &ibcInfo)
	require.NoError(t, err)
}

func TestPrepareRequestGetRandomValidatorFail(t *testing.T) {
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
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := int64(10000)
	minCount := int64(2)
	clientID := "beeb"
	ibcInfo := types.NewIBCInfo("source_port", "source_channel")

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m, &ibcInfo)
	require.Error(t, err)
}

func TestPrepareRequestGetOracleScriptFail(t *testing.T) {
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

	k.AddOracleScript(ctx, os)
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := int64(1)
	minCount := int64(2)
	clientID := "beeb"
	ibcInfo := types.NewIBCInfo("source_port", "source_channel")

	m := types.NewMsgRequestData(9999, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m, &ibcInfo)
	require.Error(t, err)
}

func TestPrepareRequestBadWasmExecutionFail(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))

	os, clear4 := getTestOracleScript()
	defer clear4()

	oracleScriptID := k.AddOracleScript(ctx, os)
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := int64(1)
	minCount := int64(2)
	clientID := "beeb"
	ibcInfo := types.NewIBCInfo("source_port", "source_channel")

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m, &ibcInfo)
	require.Error(t, err)
}

func TestPrepareRequestGetDataSourceFail(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))

	os, clear4 := getBadOracleScript()
	defer clear4()

	oracleScriptID := k.AddOracleScript(ctx, os)
	calldata, _ := hex.DecodeString("030000004254436400000000000000")
	askCount := int64(1)
	minCount := int64(2)
	clientID := "beeb"
	ibcInfo := types.NewIBCInfo("source_port", "source_channel")

	m := types.NewMsgRequestData(oracleScriptID, calldata, askCount, minCount, clientID, Alice.Address)
	err := k.PrepareRequest(ctx, &m, &ibcInfo)
	require.Error(t, err)
}

func TestResolveRequestSuccess(t *testing.T) {
	_, ctx, k := createTestInput()
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))

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
	minCount := int64(2)
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
		r.ClientID, reqID, int64(k.GetReportCount(ctx, reqID)), r.RequestTime,
		ctx.BlockTime().Unix(), types.ResolveStatus_Success, []byte("beeb"),
	)
	reqPacket := types.NewOracleRequestPacketData(
		r.ClientID, r.OracleScriptID, r.Calldata, r.MinCount, int64(len(r.RequestedValidators)),
	)
	expecetRes := k.AddResult(ctx, reqID, reqPacket, resPacket)

	require.Equal(t, expecetRes, res)
}

func TestResolveRequestFail(t *testing.T) {
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

	os, clear4 := getBadOracleScript()
	defer clear4()

	oracleScriptID := k.AddOracleScript(ctx, os)
	calldata, _ := hex.DecodeString("00")
	minCount := int64(2)
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
		r.ClientID, reqID, int64(k.GetReportCount(ctx, reqID)), r.RequestTime,
		ctx.BlockTime().Unix(), types.ResolveStatus_Failure, nil,
	)
	reqPacket := types.NewOracleRequestPacketData(
		r.ClientID, r.OracleScriptID, r.Calldata, r.MinCount, int64(len(r.RequestedValidators)),
	)
	expecetRes := k.AddResult(ctx, reqID, reqPacket, resPacket)

	require.Equal(t, expecetRes, res)
}
