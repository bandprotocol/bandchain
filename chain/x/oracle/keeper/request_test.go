package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestHasRequest(t *testing.T) {
	_, ctx, k := createTestInput()
	// We should not have a request ID 42 without setting it.
	require.False(t, k.HasRequest(ctx, 42))
	// After we set it, we should be able to find it.
	k.SetRequest(ctx, 42, types.NewRequest(1, BasicCalldata, nil, 1, 1, 1, ""))
	require.True(t, k.HasRequest(ctx, 42))
}

func TestDeleteRequest(t *testing.T) {
	_, ctx, k := createTestInput()
	// After we set it, we should be able to find it.
	k.SetRequest(ctx, 42, types.NewRequest(1, BasicCalldata, nil, 1, 1, 1, ""))
	require.True(t, k.HasRequest(ctx, 42))
	// After we delete it, we should not find it anymore.
	k.DeleteRequest(ctx, 42)
	require.False(t, k.HasRequest(ctx, 42))
	_, err := k.GetRequest(ctx, 42)
	require.Error(t, err)
	require.Panics(t, func() { _ = k.MustGetRequest(ctx, 42) })
}

func TestSetterGetterRequest(t *testing.T) {
	_, ctx, k := createTestInput()
	// Getting a non-existent request should return error.
	_, err := k.GetRequest(ctx, 42)
	require.Error(t, err)
	require.Panics(t, func() { _ = k.MustGetRequest(ctx, 42) })
	// Creates some basic requests.
	req1 := types.NewRequest(1, BasicCalldata, nil, 1, 1, 1, "")
	req2 := types.NewRequest(2, BasicCalldata, nil, 1, 1, 1, "")
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
	_, ctx, k := createTestInput()
	// Initially, we should get an empty list of pending resolves.
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RID{})
	// After we set something, we should get that thing back.
	k.SetPendingResolveList(ctx, []types.RID{5, 6, 7, 8})
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RID{5, 6, 7, 8})
	// Let's also try setting it back to empty list.
	k.SetPendingResolveList(ctx, []types.RID{})
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RID{})
	// Nil should also works.
	k.SetPendingResolveList(ctx, nil)
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RID{})
}

func TestAddDataSourceBasic(t *testing.T) {
	_, ctx, k := createTestInput()
	// We start by setting an oracle request available at ID 42.
	k.SetOracleScript(ctx, 42, types.NewOracleScript(
		Owner.Address, BasicName, BasicDesc, BasicCode, BasicSchema, BasicSourceCodeURL,
	))
	// Adding the first request should return ID 1.
	id, err := k.AddRequest(ctx, types.NewRequest(42, BasicCalldata, nil, 1, 1, 1, ""))
	require.Nil(t, err)
	require.Equal(t, id, types.RID(1))
	// Adding another request should return ID 2.
	id, err = k.AddRequest(ctx, types.NewRequest(42, BasicCalldata, nil, 1, 1, 1, ""))
	require.Nil(t, err)
	require.Equal(t, id, types.RID(2))
}

func TestAddRequestNoOracleScript(t *testing.T) {
	_, ctx, k := createTestInput()
	// Adding a new request that refers a non-existent oracle script should fail.
	require.False(t, k.HasOracleScript(ctx, 42))
	_, err := k.AddRequest(ctx, types.NewRequest(42, BasicCalldata, nil, 1, 1, 1, ""))
	require.Error(t, err)
}

func TestAddRequestTooBigCalldata(t *testing.T) {
	_, ctx, k := createTestInput()
	// We start by setting an oracle request available at ID 42.
	k.SetOracleScript(ctx, 42, types.NewOracleScript(
		Owner.Address, BasicName, BasicDesc, BasicCode, BasicSchema, BasicSourceCodeURL,
	))
	require.True(t, k.HasOracleScript(ctx, 42))
	// We set max calldata sie to 41, so adding a request with calldata size 42 should fail.
	k.SetParam(ctx, types.KeyMaxCalldataSize, 41)
	_, err := k.AddRequest(ctx, types.NewRequest(42,
		[]byte("________THIS_STRING_HAS_SIZE_OF_42________"), nil, 1, 1, 1, ""),
	)
	require.Error(t, err)
	// We now change max calldata size to 42, so it should be fine.
	k.SetParam(ctx, types.KeyMaxCalldataSize, 42)
	_, err = k.AddRequest(ctx, types.NewRequest(42,
		[]byte("________THIS_STRING_HAS_SIZE_OF_42________"), nil, 1, 1, 1, ""),
	)
	require.Nil(t, err)
}

func TestAddPendingResolveList(t *testing.T) {
	_, ctx, k := createTestInput()
	// Initially, we should get an empty list of pending resolves.
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RID{})
	// Everytime we append a new request ID, it should show up.
	k.AddPendingRequest(ctx, 42)
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RID{42})
	k.AddPendingRequest(ctx, 43)
	require.Equal(t, k.GetPendingResolveList(ctx), []types.RID{42, 43})
}

func TestGetRandomValidatorsSuccess(t *testing.T) {
	// TODO: Update this test once GetRandomValidators is actually random
}

func TestGetRandomValidatorsTooBigSize(t *testing.T) {
	// TODO: Update this test once GetRandomValidators is actually random
}
