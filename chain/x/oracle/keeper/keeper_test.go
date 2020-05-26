package keeper_test

import (
	"testing"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestGetRequestCount(t *testing.T) {
	_, ctx, k := createTestInput()

	// Initial request count must be 0
	require.Equal(t, int64(0), k.GetRequestCount(ctx))
}

func TestGetNextRequestID(t *testing.T) {
	_, ctx, k := createTestInput()

	// First request id must be 1
	require.Equal(t, types.RequestID(1), k.GetNextRequestID(ctx))

	// After add new request, request count must be 1
	require.Equal(t, int64(1), k.GetRequestCount(ctx))

	require.Equal(t, types.RequestID(2), k.GetNextRequestID(ctx))
	require.Equal(t, types.RequestID(3), k.GetNextRequestID(ctx))
	require.Equal(t, types.RequestID(4), k.GetNextRequestID(ctx))

	require.Equal(t, int64(4), k.GetRequestCount(ctx))
}

func TestGetSetMaxRawRequestCount(t *testing.T) {
	_, ctx, k := createTestInput()
	k.SetParam(ctx, types.KeyMaxRawRequestCount, 1)
	require.Equal(t, uint64(1), k.GetParam(ctx, types.KeyMaxRawRequestCount))
	k.SetParam(ctx, types.KeyMaxRawRequestCount, 2)
	require.Equal(t, uint64(2), k.GetParam(ctx, types.KeyMaxRawRequestCount))
}

func TestGetSetGasPerRawDataRequestPerValidator(t *testing.T) {
	_, ctx, k := createTestInput()
	k.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, uint64(3000))
	require.Equal(t, uint64(3000), k.GetParam(ctx, types.KeyGasPerRawDataRequestPerValidator))
	k.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, uint64(5000))
	require.Equal(t, uint64(5000), k.GetParam(ctx, types.KeyGasPerRawDataRequestPerValidator))
}

func TestGetSetParams(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetParam(ctx, types.KeyMaxRawRequestCount, 1)
	k.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, 1000)
	k.SetParam(ctx, types.KeyExpirationBlockCount, 30)
	k.SetParam(ctx, types.KeyMaxConsecutiveMisses, 10)
	require.Equal(t, types.NewParams(1, 1000, 30, 10), k.GetParams(ctx))

	k.SetParam(ctx, types.KeyMaxRawRequestCount, 2)
	k.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, 2000)
	k.SetParam(ctx, types.KeyExpirationBlockCount, 40)
	k.SetParam(ctx, types.KeyMaxConsecutiveMisses, 20)
	require.Equal(t, types.NewParams(2, 2000, 40, 20), k.GetParams(ctx))
}
