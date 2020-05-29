package keeper_test

import (
	"testing"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestResultBasicFunctions(t *testing.T) {
	_, ctx, k := createTestInput()

	req := types.NewOracleRequestPacketData("alice", 1, BasicCalldata, 1, 1)
	res := types.NewOracleResponsePacketData("alice", 1, 1, 1589535020, 1589535022, 1, BasicCalldata)
	resultHash := types.CalculateResultHash(req, res)
	k.SetResult(ctx, types.RequestID(1), resultHash)

	// Test GetResult func
	resultHashReqID1, err := k.GetResult(ctx, types.RequestID(1))
	require.NoError(t, err)
	require.Equal(t, resultHash, resultHashReqID1)

	_, err = k.GetResult(ctx, types.RequestID(2))
	require.Error(t, err)

	// Test HasResult func
	require.True(t, k.HasResult(ctx, types.RequestID(1)))
	require.False(t, k.HasResult(ctx, types.RequestID(2)))
}

func TestGetAllResults(t *testing.T) {
	_, ctx, k := createTestInput()

	req1 := types.NewOracleRequestPacketData("alice", 1, BasicCalldata, 1, 1)
	res1 := types.NewOracleResponsePacketData("alice", 1, 1, 1589535020, 1589535022, 1, BasicCalldata)
	resultHash1 := types.CalculateResultHash(req1, res1)
	k.SetResult(ctx, types.RequestID(1), resultHash1)
	req4 := types.NewOracleRequestPacketData("bob", 1, BasicCalldata, 1, 1)
	res4 := types.NewOracleResponsePacketData("bob", 4, 1, 1589535020, 1589535022, 1, BasicCalldata)
	resultHash4 := types.CalculateResultHash(req4, res4)
	k.SetResult(ctx, types.RequestID(4), resultHash4)

	results := k.GetAllResults(ctx)

	require.Equal(t, 4, len(results))
	require.Equal(t, resultHash1, results[0])

	// result of reqID 2 and 3 should be nil
	require.Empty(t, results[1])
	require.Empty(t, results[2])

	require.Equal(t, resultHash4, results[3])
}
