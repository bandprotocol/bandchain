package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetterSetterResult(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
	require.NotNil(t, err)

	keeper.SetResult(ctx, 1, 1, []byte("calldata"), []byte("result"))
	actualResult, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
	require.Nil(t, err)
	require.Equal(t, []byte("result"), actualResult)
}

func TestAddResultSuccess(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
	require.NotNil(t, err)

	keeper.SetMaxResultSize(ctx, int64(1024))
	err = keeper.AddResult(ctx, 1, 1, []byte("calldata"), []byte("result"))
	require.Nil(t, err)

	actualResult, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
	require.Nil(t, err)
	require.Equal(t, []byte("result"), actualResult)
}

func TestAddResultFailWithExceedResultSize(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetResult(ctx, 1, 1, []byte("calldata"))
	require.NotNil(t, err)

	keeper.SetMaxResultSize(ctx, int64(1))
	err = keeper.AddResult(ctx, 1, 1, []byte("calldata"), []byte("result"))
	require.NotNil(t, err)
}
