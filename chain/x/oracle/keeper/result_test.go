package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
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

	reqPacket1 := oracle.OracleRequestPacketData{
		ClientID:       "alice",
		OracleScriptID: 1,
		Calldata:       BasicCalldata,
		AskCount:       1,
		MinCount:       1,
	}

	resPacket1 := oracle.OracleResponsePacketData{
		ClientID:      "alice",
		RequestID:     1,
		AnsCount:      1,
		RequestTime:   1589535020,
		ResolveTime:   1589535022,
		ResolveStatus: 1,
		Result:        BasicCalldata,
	}

	_, err := k.AddResult(ctx, types.RequestID(1), reqPacket1, resPacket1)
	require.NoError(t, err)

	reqPacket4 := oracle.OracleRequestPacketData{
		ClientID:       "bob",
		OracleScriptID: 1,
		Calldata:       BasicCalldata,
		AskCount:       1,
		MinCount:       1,
	}

	resPacket4 := oracle.OracleResponsePacketData{
		ClientID:      "bob",
		RequestID:     4,
		AnsCount:      1,
		RequestTime:   1589535020,
		ResolveTime:   1589535022,
		ResolveStatus: 1,
		Result:        BasicCalldata,
	}

	_, err = k.AddResult(ctx, types.RequestID(4), reqPacket4, resPacket4)
	require.NoError(t, err)

	results := k.GetAllResults(ctx)

	require.Equal(t, 4, len(results))
	require.NotEmpty(t, results[0])

	// result of reqID 2 and 3 should be empty byte array
	require.Empty(t, results[1])
	require.Empty(t, results[2])

	require.NotEmpty(t, results[3])
}
