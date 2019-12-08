package keeper

import (
	"testing"

	"github.com/bandprotocol/bandx/oracle/x/oracle/internal/types"
	"github.com/stretchr/testify/require"
)

func TestGetterSetterRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetRequest(ctx, 1)
	require.NotNil(t, err)

	request := types.NewDataPoint(1, []byte("CodeHash"), 10)

	keeper.SetRequest(ctx, 1, request)
	actualRequest, err := keeper.GetRequest(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, request, actualRequest)
}

// Can get/set pending request correctly and set empty case
func TestGetSetPendingRequests(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	reqIDs := keeper.GetPending(ctx)
	require.Equal(t, []uint64{}, reqIDs)

	keeper.SetPending(ctx, []uint64{1, 2, 3})

	reqIDs = keeper.GetPending(ctx)
	require.Equal(t, []uint64{1, 2, 3}, reqIDs)

	keeper.SetPending(ctx, []uint64{})
	reqIDs = keeper.GetPending(ctx)
	require.Equal(t, []uint64{}, reqIDs)
}

// Can set pending request will set only unique request IDs
func TestGetSetPendingRequestUnique(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	reqIDs := keeper.GetPending(ctx)
	require.Equal(t, []uint64{}, reqIDs)

	keeper.SetPending(ctx, []uint64{3, 2, 3, 1, 2, 1, 3, 2, 1})
	reqIDs = keeper.GetPending(ctx)
	// no guarantee of an order
	require.Equal(t, []uint64{3, 2, 1}, reqIDs)
}
