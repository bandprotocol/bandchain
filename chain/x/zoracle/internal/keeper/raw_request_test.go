package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

func TestGettterSetterRawDataRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	keeper.SetRawDataRequest(ctx, 1, 1, 0, []byte("calldata1"))
	keeper.SetRawDataRequest(ctx, 1, 2, 1, []byte("calldata2"))

	rawRequest, err := keeper.GetRawDataRequest(ctx, 1, 1)
	require.Nil(t, err)
	expect := types.NewRawDataRequest(0, []byte("calldata1"))
	require.Equal(t, expect, rawRequest)

	_, err = keeper.GetRawDataRequest(ctx, 1, 3)
	require.Equal(t, types.CodeRequestNotFound, err.Code())
}
