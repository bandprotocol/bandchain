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
	encodedResult := types.CalculateEncodedResult(req, res)
	k.SetResult(ctx, types.RequestID(1), encodedResult)

	// Test GetResult func
	reqPacket, resPacket, err := k.GetResult(ctx, types.RequestID(1))
	require.NoError(t, err)
	require.Equal(t, req, reqPacket)
	require.Equal(t, res, resPacket)

	// Test MustGetResult func
	reqPacket, resPacket = k.MustGetResult(ctx, types.RequestID(1))
	require.NoError(t, err)
	require.Equal(t, req, reqPacket)
	require.Equal(t, res, resPacket)

	_, _, err = k.GetResult(ctx, types.RequestID(2))
	require.Error(t, err)

	require.Panics(t, func() { k.GetResult(ctx, types.RequestID(2)) })

	// Test HasResult func
	require.True(t, k.HasResult(ctx, types.RequestID(1)))
	require.False(t, k.HasResult(ctx, types.RequestID(2)))
}

func TestGetAllResults(t *testing.T) {
	_, ctx, k := createTestInput()

	req1 := types.NewOracleRequestPacketData("alice", 1, BasicCalldata, 1, 1)
	res1 := types.NewOracleResponsePacketData("alice", 1, 1, 1589535020, 1589535022, 1, BasicCalldata)
	encodedResult1 := types.CalculateEncodedResult(req1, res1)
	k.SetResult(ctx, types.RequestID(1), encodedResult1)
	req4 := types.NewOracleRequestPacketData("bob", 1, BasicCalldata, 1, 1)
	res4 := types.NewOracleResponsePacketData("bob", 4, 1, 1589535020, 1589535022, 1, BasicCalldata)
	encodedResult4 := types.CalculateEncodedResult(req4, res4)
	k.SetResult(ctx, types.RequestID(4), encodedResult4)

	results := k.GetAllResults(ctx)

	require.Equal(t, 4, len(results))
	require.Equal(t, encodedResult1, results[0])

	// result of reqID 2 and 3 should be nil
	require.Empty(t, results[1])
	require.Empty(t, results[2])

	require.Equal(t, encodedResult4, results[3])
}
