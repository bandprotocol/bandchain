package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	"github.com/bandprotocol/bandchain/chain/x/oracle/testapp"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestResultBasicFunctions(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()

	req := types.NewOracleRequestPacketData("alice", 1, BasicCalldata, 1, 1)
	res := types.NewOracleResponsePacketData("alice", 1, 1, 1589535020, 1589535022, 1, BasicCalldata)
	encodedResult := obi.MustEncode(req, res)
	k.SetResult(ctx, types.RequestID(1), encodedResult)

	// Test GetResult func
	result, err := k.GetResult(ctx, types.RequestID(1))
	require.NoError(t, err)
	require.Equal(t, result, types.Result{RequestPacketData: req, ResponsePacketData: res})

	// Test MustGetResult func
	result = k.MustGetResult(ctx, types.RequestID(1))
	require.NoError(t, err)
	require.Equal(t, result, types.Result{RequestPacketData: req, ResponsePacketData: res})

	_, err = k.GetResult(ctx, types.RequestID(2))
	require.Error(t, err)

	require.Panics(t, func() { k.MustGetResult(ctx, types.RequestID(2)) })

	// Test HasResult func
	require.True(t, k.HasResult(ctx, types.RequestID(1)))
	require.False(t, k.HasResult(ctx, types.RequestID(2)))
}
