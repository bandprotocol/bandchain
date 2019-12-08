package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRequestCount(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Initial request count must be 0
	require.Equal(t, uint64(0), keeper.GetRequestCount(ctx))
}

func TestGetNextRequestID(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// First request id must be 1
	require.Equal(t, uint64(1), keeper.GetNextRequestID(ctx))

	// After add new request, request count must be 1
	require.Equal(t, uint64(1), keeper.GetRequestCount(ctx))

	require.Equal(t, uint64(2), keeper.GetNextRequestID(ctx))
	require.Equal(t, uint64(3), keeper.GetNextRequestID(ctx))
	require.Equal(t, uint64(4), keeper.GetNextRequestID(ctx))

	require.Equal(t, uint64(4), keeper.GetRequestCount(ctx))
}
