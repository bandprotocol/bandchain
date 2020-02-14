package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetterSetterResult(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetResult(ctx, 1)
	require.NotNil(t, err)

	keeper.SetResult(ctx, 1, []byte("result"))
	actualResult, err := keeper.GetResult(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, []byte("result"), actualResult)
}
