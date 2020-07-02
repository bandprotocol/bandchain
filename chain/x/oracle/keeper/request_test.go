package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/testapp"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
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
	_, ctx, k := testapp.CreateTestInput()

	k.SetParam(ctx, types.KeyExpirationBlockCount, 3)
	currentTime := time.Unix(1000, 0)
	currentBlock := int64(1)
	ctx = ctx.WithBlockHeight(currentBlock).WithBlockTime(currentTime)
	r := defaultRequest()
	r.RequestHeight = currentBlock
	r.RequestTime = currentTime
	k.SetRequest(ctx, 1, r)
	k.SetRequestCount(ctx, 1)

	currentTime = currentTime.Add(3 * time.Second)
	currentBlock += 1
	k.AddReport(ctx, 1, types.NewReport(testapp.Validator1.ValAddress, true, []types.RawReport{
		types.NewRawReport(42, 0, BasicReport),
		types.NewRawReport(43, 0, BasicReport),
	}))

	// Nothing happen
	k.ProcessExpiredRequests(ctx)
	require.Equal(t, types.RequestID(0), k.GetRequestLastExpired(ctx))
	require.False(t, k.HasResult(ctx, 1))

	currentTime = currentTime.Add(7 * time.Second)
	currentBlock += 2
	ctx = ctx.WithBlockHeight(currentBlock).WithBlockTime(currentTime)
	k.ProcessExpiredRequests(ctx)

	result, err := k.GetResult(ctx, 1)
	require.NoError(t, err)
	req := types.NewOracleRequestPacketData(
		r.ClientID, r.OracleScriptID, r.Calldata, uint64(len(r.RequestedValidators)), r.MinCount,
	)
	res := types.NewOracleResponsePacketData(
		r.ClientID, 1, 1, r.RequestTime.Unix(),
		currentTime.Unix(), types.ResolveStatus_Expired, []byte{},
	)
	require.Equal(t, types.NewResult(req, res), result)

	// Check validator status
	require.True(t, k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress).IsActive)
	require.False(t, k.GetValidatorStatus(ctx, testapp.Validator2.ValAddress).IsActive)
}

func TestProcessSuccessRequests(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()

	k.SetParam(ctx, types.KeyExpirationBlockCount, 3)
	currentTime := time.Unix(1000, 0)
	currentBlock := int64(1)
	ctx = ctx.WithBlockHeight(currentBlock).WithBlockTime(currentTime)
	r := defaultRequest()
	r.RequestHeight = currentBlock
	r.RequestTime = currentTime
	k.SetRequest(ctx, 1, r)
	k.SetRequestCount(ctx, 1)

	currentTime = currentTime.Add(3 * time.Second)
	currentBlock += 1
	k.AddReport(ctx, 1, types.NewReport(testapp.Validator1.ValAddress, true, []types.RawReport{
		types.NewRawReport(42, 0, BasicReport),
		types.NewRawReport(43, 0, BasicReport),
	}))

	k.AddReport(ctx, 1, types.NewReport(testapp.Validator2.ValAddress, true, []types.RawReport{
		types.NewRawReport(42, 0, BasicReport),
		types.NewRawReport(43, 0, BasicReport),
	}))

	// Nothing happen
	k.ProcessExpiredRequests(ctx)
	require.Equal(t, types.RequestID(0), k.GetRequestLastExpired(ctx))
	require.False(t, k.HasResult(ctx, 1))

	currentTime = currentTime.Add(7 * time.Second)
	currentBlock += 2
	ctx = ctx.WithBlockHeight(currentBlock).WithBlockTime(currentTime)
	k.ProcessExpiredRequests(ctx)

	// Expired status must not be saved in store.
	require.False(t, k.HasResult(ctx, 1))

	// Check validator status
	require.True(t, k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress).IsActive)
	require.True(t, k.GetValidatorStatus(ctx, testapp.Validator2.ValAddress).IsActive)
	require.Equal(t, types.RequestID(1), k.GetRequestLastExpired(ctx))
}
