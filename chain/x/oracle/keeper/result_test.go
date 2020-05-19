package keeper_test

import (
	"testing"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/stretchr/testify/require"
)

// import (
// 	"testing"
// 	"time"

// 	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/require"
// )

// func TestGetterSetterResult(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	_, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
// 	require.NotNil(t, err)

// 	keeper.SetResult(ctx, 1, 1, []byte("calldata"), types.NewResult(1, 2, 3, 2, 2, []byte("result")))
// 	actualResult, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
// 	require.Nil(t, err)
// 	require.Equal(t, types.Result{
// 		RequestTime:              1,
// 		AggregationTime:          2,
// 		RequestedValidatorsCount: 3,
// 		MinCount: 2,
// 		ReportedValidatorsCount:  2,
// 		Data:                     []byte("result"),
// 	}, actualResult)
// }

// func TestAddResultSuccess(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	_, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
// 	require.NotNil(t, err)

// 	keeper.SetRequest(ctx, 1, types.NewRequest(
// 		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))}, 1, 0, 0, "clientID",
// 	))

// 	ctx = ctx.WithBlockTime(time.Unix(int64(1581589999), 0))

// 	err = keeper.AddResult(ctx, 1, 1, []byte("calldata"), []byte("result"))
// 	require.Nil(t, err)

// 	actualResult, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
// 	require.Nil(t, err)
// 	require.Equal(
// 		t,
// 		types.Result{
// 			RequestTime:              0,
// 			AggregationTime:          1581589999,
// 			RequestedValidatorsCount: 1,
// 			MinCount: 1,
// 			ReportedValidatorsCount:  0,
// 			Data:                     []byte("result"),
// 		},
// 		actualResult,
// 	)
// }

// func TestAddResultFailWithExceedResultSize(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	_, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
// 	require.NotNil(t, err)

// 	keeper.SetParam(ctx, types.KeyMaxResultSize, 1)
// 	err = keeper.AddResult(ctx, 1, 1, []byte("calldata"), []byte("result"))
// 	require.NotNil(t, err)
// }

func TestGetAllResults(t *testing.T) {
	_, ctx, k := createTestInput()

	reqPacket1 := types.NewOracleRequestPacketData("alice", 1, BasicCalldata, 1, 1)
	resPacket1 := types.NewOracleResponsePacketData("alice", 1, 1, 1589535020, 1589535022, 1, BasicCalldata)

	resultHashReqID1, err := k.AddResult(ctx, types.RequestID(1), reqPacket1, resPacket1)
	require.NoError(t, err)

	reqPacket4 := types.NewOracleRequestPacketData("bob", 1, BasicCalldata, 1, 1)
	resPacket4 := types.NewOracleResponsePacketData("bob", 4, 1, 1589535020, 1589535022, 1, BasicCalldata)

	resultHashReqID4, err := k.AddResult(ctx, types.RequestID(4), reqPacket4, resPacket4)
	require.NoError(t, err)

	results := k.GetAllResults(ctx)

	require.Equal(t, 4, len(results))
	require.Equal(t, resultHashReqID1, results[0])

	// result of reqID 2 and 3 should be empty byte array
	require.Empty(t, results[1])
	require.Empty(t, results[2])

	require.Equal(t, resultHashReqID4, results[3])
}

func TestSetResult(t *testing.T) {
	_, ctx, k := createTestInput()

	reqPacket := types.NewOracleRequestPacketData("alice", 1, BasicCalldata, 1, 1)
	resPacket := types.NewOracleResponsePacketData("alice", 1, 1, 1589535020, 1589535022, 1, BasicCalldata)

	resultHash, err := k.AddResult(ctx, types.RequestID(1), reqPacket, resPacket)
	require.NoError(t, err)

	// Set result for request ID 2
	k.SetResult(ctx, types.RequestID(2), resultHash)
	resultHashReqID2, err := k.GetResult(ctx, types.RequestID(2))
	require.NoError(t, err)

	require.Equal(t, resultHash, resultHashReqID2)
}
