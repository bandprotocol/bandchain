package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/testapp"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

func TestResultBasicFunctions(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// We start by setting result of request#1.
	req := types.NewOracleRequestPacketData("alice", 1, BasicCalldata, 1, 1)
	res := types.NewOracleResponsePacketData("alice", 1, 1, 1589535020, 1589535022, 1, BasicCalldata)
	k.SetResult(ctx, 1, types.NewResult(req, res))
	// GetResult and MustGetResult should return what we set.
	result, err := k.GetResult(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, result, types.NewResult(req, res))
	result = k.MustGetResult(ctx, 1)
	require.Equal(t, result, types.NewResult(req, res))
	// GetResult of another request should return error.
	_, err = k.GetResult(ctx, 2)
	require.Error(t, err)
	require.Panics(t, func() { k.MustGetResult(ctx, 2) })
	// HasResult should also perform correctly.
	require.True(t, k.HasResult(ctx, 1))
	require.False(t, k.HasResult(ctx, 2))
}

func TestSaveResultOK(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	ctx = ctx.WithBlockTime(testapp.ParseTime(200))
	k.SetRequest(ctx, 42, defaultRequest()) // See report_test.go
	k.SetReport(ctx, 42, types.NewReport(testapp.Validator1.ValAddress, true, nil))
	k.SaveResult(ctx, 42, types.ResolveStatus_Success, BasicResult)
	expect := types.NewResult(types.NewOracleRequestPacketData(
		BasicClientID, 1, BasicCalldata, 2, 2,
	), types.NewOracleResponsePacketData(
		BasicClientID, 42, 1, testapp.ParseTime(0).Unix(), testapp.ParseTime(200).Unix(),
		types.ResolveStatus_Success, BasicResult,
	))
	result, err := k.GetResult(ctx, 42)
	require.NoError(t, err)
	require.Equal(t, expect, result)
}

func TestResolveSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	k.SetRequest(ctx, 42, defaultRequest()) // See report_test.go
	k.SetReport(ctx, 42, types.NewReport(testapp.Validator1.ValAddress, true, nil))
	k.ResolveSuccess(ctx, 42, BasicResult, 1234)
	require.Equal(t, types.ResolveStatus_Success, k.MustGetResult(ctx, 42).ResponsePacketData.ResolveStatus)
	require.Equal(t, BasicResult, k.MustGetResult(ctx, 42).ResponsePacketData.Result)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, "42"),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, "1"),
		sdk.NewAttribute(types.AttributeKeyResult, "42415349435f524553554c54"), // BASIC_RESULT
		sdk.NewAttribute(types.AttributeKeyGasUsed, "1234"),
	)}, ctx.EventManager().Events())
}

func TestResolveFailure(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	k.SetRequest(ctx, 42, defaultRequest()) // See report_test.go
	k.SetReport(ctx, 42, types.NewReport(testapp.Validator1.ValAddress, true, nil))
	k.ResolveFailure(ctx, 42, "REASON")
	require.Equal(t, types.ResolveStatus_Failure, k.MustGetResult(ctx, 42).ResponsePacketData.ResolveStatus)
	require.Equal(t, []byte{}, k.MustGetResult(ctx, 42).ResponsePacketData.Result)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, "42"),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, "2"),
		sdk.NewAttribute(types.AttributeKeyReason, "REASON"),
	)}, ctx.EventManager().Events())
}

func TestResolveExpired(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	k.SetRequest(ctx, 42, defaultRequest()) // See report_test.go
	k.SetReport(ctx, 42, types.NewReport(testapp.Validator1.ValAddress, true, nil))
	k.ResolveExpired(ctx, 42)
	require.Equal(t, types.ResolveStatus_Expired, k.MustGetResult(ctx, 42).ResponsePacketData.ResolveStatus)
	require.Equal(t, []byte{}, k.MustGetResult(ctx, 42).ResponsePacketData.Result)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, "42"),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, "3"),
	)}, ctx.EventManager().Events())
}
