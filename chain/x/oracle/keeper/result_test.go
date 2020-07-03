package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/testapp"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
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
