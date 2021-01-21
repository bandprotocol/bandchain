package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/testapp"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

func TestHasRequest(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// We should not have a request ID 42 without setting it.
	require.False(t, k.HasRequest(ctx, 42))
	// After we set it, we should be able to find it.
	k.SetRequest(ctx, 42, types.NewRequest(1, BasicCalldata, nil, 1, 1, testapp.ParseTime(0), "", nil))
	require.True(t, k.HasRequest(ctx, 42))
}

func TestDeleteRequest(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// After we set it, we should be able to find it.
	k.SetRequest(ctx, 42, types.NewRequest(1, BasicCalldata, nil, 1, 1, testapp.ParseTime(0), "", nil))
	require.True(t, k.HasRequest(ctx, 42))
	// After we delete it, we should not find it anymore.
	k.DeleteRequest(ctx, 42)
	require.False(t, k.HasRequest(ctx, 42))
	_, err := k.GetRequest(ctx, 42)
	require.Error(t, err)
	require.Panics(t, func() { _ = k.MustGetRequest(ctx, 42) })
}

func TestSetterGetterRequest(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Getting a non-existent request should return error.
	_, err := k.GetRequest(ctx, 42)
	require.Error(t, err)
	require.Panics(t, func() { _ = k.MustGetRequest(ctx, 42) })
	// Creates some basic requests.
	req1 := types.NewRequest(1, BasicCalldata, nil, 1, 1, testapp.ParseTime(0), "", nil)
	req2 := types.NewRequest(2, BasicCalldata, nil, 1, 1, testapp.ParseTime(0), "", nil)
	// Sets id 42 with request 1 and id 42 with request 2.
	k.SetRequest(ctx, 42, req1)
	k.SetRequest(ctx, 43, req2)
	// Checks that Get and MustGet perform correctly.
	req1Res, err := k.GetRequest(ctx, 42)
	require.Nil(t, err)
	require.Equal(t, req1, req1Res)
	require.Equal(t, req1, k.MustGetRequest(ctx, 42))
	req2Res, err := k.GetRequest(ctx, 43)
	require.Nil(t, err)
	require.Equal(t, req2, req2Res)
	require.Equal(t, req2, k.MustGetRequest(ctx, 43))
	// Replaces id 42 with another request.
	k.SetRequest(ctx, 42, req2)
	require.NotEqual(t, req1, k.MustGetRequest(ctx, 42))
	require.Equal(t, req2, k.MustGetRequest(ctx, 42))
}

func TestSetterGettterPendingResolveList(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Initially, we should get an empty list of pending resolves.
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{})
	// After we set something, we should get that thing back.
	k.SetPendingResolveList(ctx, []types.RequestID{5, 6, 7, 8})
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{5, 6, 7, 8})
	// Let's also try setting it back to empty list.
	k.SetPendingResolveList(ctx, []types.RequestID{})
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{})
	// Nil should also works.
	k.SetPendingResolveList(ctx, nil)
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{})
}

func TestAddDataSourceBasic(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// We start by setting an oracle request available at ID 42.
	k.SetOracleScript(ctx, 42, types.NewOracleScript(
		testapp.Owner.Address, BasicName, BasicDesc, BasicFilename, BasicSchema, BasicSourceCodeURL,
	))
	// Adding the first request should return ID 1.
	id := k.AddRequest(ctx, types.NewRequest(42, BasicCalldata, nil, 1, 1, testapp.ParseTime(0), "", nil))
	require.Equal(t, id, types.RequestID(1))
	// Adding another request should return ID 2.
	id = k.AddRequest(ctx, types.NewRequest(42, BasicCalldata, nil, 1, 1, testapp.ParseTime(0), "", nil))
	require.Equal(t, id, types.RequestID(2))
}

func TestAddPendingResolveList(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// Initially, we should get an empty list of pending resolves.
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{})
	// Everytime we append a new request ID, it should show up.
	k.AddPendingRequest(ctx, 42)
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{42})
	k.AddPendingRequest(ctx, 43)
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RequestID{42, 43})
}

func TestProcessExpiredRequests(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	k.SetParam(ctx, types.KeyExpirationBlockCount, 3)
	// Set some initial requests. All requests are asked to validators 1 & 2.
	req1 := defaultRequest()
	req1.RequestHeight = 5
	req2 := defaultRequest()
	req2.RequestHeight = 6
	req3 := defaultRequest()
	req3.RequestHeight = 6
	req4 := defaultRequest()
	req4.RequestHeight = 10
	k.AddRequest(ctx, req1)
	k.AddRequest(ctx, req2)
	k.AddRequest(ctx, req3)
	k.AddRequest(ctx, req4)
	// Initially all validators are active.
	require.True(t, k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress).IsActive)
	require.True(t, k.GetValidatorStatus(ctx, testapp.Validator2.ValAddress).IsActive)
	// Validator 1 reports all requests. Validator 2 misses request#3.
	rawReports := []types.RawReport{types.NewRawReport(42, 0, BasicReport), types.NewRawReport(43, 0, BasicReport)}
	k.AddReport(ctx, 1, types.NewReport(testapp.Validator1.ValAddress, false, rawReports))
	k.AddReport(ctx, 2, types.NewReport(testapp.Validator1.ValAddress, true, rawReports))
	k.AddReport(ctx, 3, types.NewReport(testapp.Validator1.ValAddress, false, rawReports))
	k.AddReport(ctx, 4, types.NewReport(testapp.Validator1.ValAddress, true, rawReports))
	k.AddReport(ctx, 1, types.NewReport(testapp.Validator2.ValAddress, true, rawReports))
	k.AddReport(ctx, 2, types.NewReport(testapp.Validator2.ValAddress, true, rawReports))
	k.AddReport(ctx, 4, types.NewReport(testapp.Validator2.ValAddress, true, rawReports))
	// Request 1, 2 and 4 gets resolved. Request 3 does not.
	k.ResolveSuccess(ctx, 1, BasicResult, 1234)
	k.ResolveFailure(ctx, 2, "ARBITRARY_REASON")
	k.ResolveSuccess(ctx, 4, BasicResult, 1234)
	// Initially, last expired request ID should be 0.
	require.Equal(t, types.RequestID(0), k.GetRequestLastExpired(ctx))
	// At block 7, nothing should happen.
	ctx = ctx.WithBlockHeight(7).WithBlockTime(testapp.ParseTime(7000)).WithEventManager(sdk.NewEventManager())
	k.ProcessExpiredRequests(ctx)
	require.Equal(t, sdk.Events{}, ctx.EventManager().Events())
	require.Equal(t, types.RequestID(0), k.GetRequestLastExpired(ctx))
	require.True(t, k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress).IsActive)
	require.True(t, k.GetValidatorStatus(ctx, testapp.Validator2.ValAddress).IsActive)
	// At block 8, now last request ID should move to 1. No events should be emitted.
	ctx = ctx.WithBlockHeight(8).WithBlockTime(testapp.ParseTime(8000)).WithEventManager(sdk.NewEventManager())
	k.ProcessExpiredRequests(ctx)
	require.Equal(t, sdk.Events{}, ctx.EventManager().Events())
	require.Equal(t, types.RequestID(1), k.GetRequestLastExpired(ctx))
	require.True(t, k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress).IsActive)
	require.True(t, k.GetValidatorStatus(ctx, testapp.Validator2.ValAddress).IsActive)
	// At block 9, request#3 is expired and validator 2 becomes inactive.
	ctx = ctx.WithBlockHeight(9).WithBlockTime(testapp.ParseTime(9000)).WithEventManager(sdk.NewEventManager())
	k.ProcessExpiredRequests(ctx)
	require.Equal(t, sdk.Events{sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, "3"),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, "3"),
	), sdk.NewEvent(
		types.EventTypeDeactivate,
		sdk.NewAttribute(types.AttributeKeyValidator, testapp.Validator2.ValAddress.String()),
	)}, ctx.EventManager().Events())
	require.Equal(t, types.RequestID(3), k.GetRequestLastExpired(ctx))
	require.True(t, k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress).IsActive)
	require.False(t, k.GetValidatorStatus(ctx, testapp.Validator2.ValAddress).IsActive)
	require.Equal(t, types.NewOracleResponsePacketData(
		BasicClientID, 3, 1, req3.RequestTime.Unix(), testapp.ParseTime(9000).Unix(),
		types.ResolveStatus_Expired, []byte{},
	), k.MustGetResult(ctx, 3).ResponsePacketData)
	// At block 10, nothing should happen
	ctx = ctx.WithBlockHeight(10).WithBlockTime(testapp.ParseTime(10000)).WithEventManager(sdk.NewEventManager())
	k.ProcessExpiredRequests(ctx)
	require.Equal(t, sdk.Events{}, ctx.EventManager().Events())
	require.Equal(t, types.RequestID(3), k.GetRequestLastExpired(ctx))
	// At block 13, last expired request becomes 4.
	ctx = ctx.WithBlockHeight(13).WithBlockTime(testapp.ParseTime(13000)).WithEventManager(sdk.NewEventManager())
	k.ProcessExpiredRequests(ctx)
	require.Equal(t, sdk.Events{}, ctx.EventManager().Events())
	require.Equal(t, types.RequestID(4), k.GetRequestLastExpired(ctx))
}
