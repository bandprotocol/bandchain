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

func TestSetResult(t *testing.T) {
	_, ctx, k := createTestInput()

	req := types.NewOracleRequestPacketData("alice", 1, BasicCalldata, 1, 1)
	res := types.NewOracleResponsePacketData("alice", 1, 1, 1589535020, 1589535022, 1, BasicCalldata)
	resultHash := types.CalculateResultHash(req, res)
	k.SetResult(ctx, types.RequestID(1), resultHash)

	resultHashReqID1, err := k.GetResult(ctx, types.RequestID(1))
	require.NoError(t, err)
	require.Equal(t, resultHash, resultHashReqID1)
}
