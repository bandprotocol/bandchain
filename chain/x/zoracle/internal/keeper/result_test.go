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
