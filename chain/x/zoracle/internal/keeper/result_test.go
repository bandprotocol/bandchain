package keeper

import (
	"testing"
	"time"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetterSetterResult(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
	require.NotNil(t, err)

	keeper.SetResult(ctx, 1, 1, []byte("calldata"), types.NewResult(1, 2, 3, 2, 2, []byte("result")))
	actualResult, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
	require.Nil(t, err)
	require.Equal(t, types.Result{
		RequestTime:              1,
		AggregationTime:          2,
		RequestedValidatorsCount: 3,
		SufficientValidatorCount: 2,
		ReportedValidatorsCount:  2,
		Data:                     []byte("result"),
	}, actualResult)
}

func TestAddResultSuccess(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
	require.NotNil(t, err)

	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))}, 1, 0, 0, 100, 10000,
	))

	ctx = ctx.WithBlockTime(time.Unix(int64(1581589999), 0))

	err = keeper.AddResult(ctx, 1, 1, []byte("calldata"), []byte("result"))
	require.Nil(t, err)

	actualResult, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
	require.Nil(t, err)
	require.Equal(
		t,
		types.Result{
			RequestTime:              0,
			AggregationTime:          1581589999,
			RequestedValidatorsCount: 1,
			SufficientValidatorCount: 1,
			ReportedValidatorsCount:  0,
			Data:                     []byte("result"),
		},
		actualResult,
	)
}

func TestAddResultFailWithExceedResultSize(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
	require.NotNil(t, err)

	keeper.SetMaxResultSize(ctx, int64(1))
	err = keeper.AddResult(ctx, 1, 1, []byte("calldata"), []byte("result"))
	require.NotNil(t, err)
}
