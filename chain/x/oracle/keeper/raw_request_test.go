package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	me "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func createRawRequestTestInput() (*bandapp.BandApp, sdk.Context, me.Keeper) {
	app, ctx, k := createTestInput()
	// Before running any test in this file:
	// - We set a mock request available at ID 20.
	// - We set a mock data source available at ID 42 and 43.
	k.SetDataSource(ctx, 42, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband, BasicExec,
	))
	k.SetDataSource(ctx, 43, types.NewDataSource(
		Owner.Address, BasicName, BasicDesc, Coins10uband, BasicExec,
	))
	return app, ctx, k
}

func TestHasRawRequest(t *testing.T) {
	_, ctx, k := createTestInput()
	// We should not have a raw request ID (42, 42) without setting it.
	require.False(t, k.HasRawRequest(ctx, 42, 42))
	// After we set it, we should be able to find it.
	k.SetRawRequest(ctx, 42, 42, types.NewRawRequest(1, BasicCalldata))
	require.True(t, k.HasRawRequest(ctx, 42, 42))
}

func TestSetterGetterRawRequest(t *testing.T) {
	_, ctx, k := createRawRequestTestInput()
	// Getting a non-existent raw request should return error.
	_, err := k.GetRawRequest(ctx, 10, 20)
	require.Error(t, err)
	// Sets some basic raw requests.
	rawRequest1 := types.NewRawRequest(1, []byte("RAW_REQUEST_CALLDATA_1"))
	rawRequest2 := types.NewRawRequest(2, []byte("RAW_REQUEST_CALLDATA_2"))
	k.SetRawRequest(ctx, 10, 20, rawRequest1)
	k.SetRawRequest(ctx, 11, 21, rawRequest2)
	// Checks that Get performs correctly.
	rawRequest1Res, err := k.GetRawRequest(ctx, 10, 20)
	require.Nil(t, err)
	require.Equal(t, rawRequest1, rawRequest1Res)
	rawRequest2Res, err := k.GetRawRequest(ctx, 11, 21)
	require.Nil(t, err)
	require.Equal(t, rawRequest2, rawRequest2Res)
	// Replaces id (10, 20) with another raw request.
	k.SetRawRequest(ctx, 10, 20, rawRequest2)
	rawRequestRes, err := k.GetRawRequest(ctx, 10, 20)
	require.Nil(t, err)
	require.NotEqual(t, rawRequest1, rawRequestRes)
	require.Equal(t, rawRequest2, rawRequestRes)
}

// TODO: Add more tests related to adding raw requests once
