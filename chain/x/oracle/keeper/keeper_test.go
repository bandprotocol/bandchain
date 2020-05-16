package keeper

import (
	"testing"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestGetRequestCount(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Initial request count must be 0
	require.Equal(t, int64(0), keeper.GetRequestCount(ctx))
}

func TestGetNextRequestID(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// First request id must be 1
	require.Equal(t, types.RequestID(1), keeper.GetNextRequestID(ctx))

	// After add new request, request count must be 1
	require.Equal(t, int64(1), keeper.GetRequestCount(ctx))

	require.Equal(t, types.RequestID(2), keeper.GetNextRequestID(ctx))
	require.Equal(t, types.RequestID(3), keeper.GetNextRequestID(ctx))
	require.Equal(t, types.RequestID(4), keeper.GetNextRequestID(ctx))

	require.Equal(t, int64(4), keeper.GetRequestCount(ctx))
}

func TestGetSetMaxRawRequestCount(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetParam(ctx, types.KeyMaxRawRequestCount, 1)
	require.Equal(t, uint64(1), keeper.GetParam(ctx, types.KeyMaxRawRequestCount))
	keeper.SetParam(ctx, types.KeyMaxRawRequestCount, 2)
	require.Equal(t, uint64(2), keeper.GetParam(ctx, types.KeyMaxRawRequestCount))
}

func TestGetSetMaxRawDataReportSize(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetParam(ctx, types.KeyMaxRawDataReportSize, 1)
	require.Equal(t, uint64(1), keeper.GetParam(ctx, types.KeyMaxRawDataReportSize))
	keeper.SetParam(ctx, types.KeyMaxRawDataReportSize, 2)
	require.Equal(t, uint64(2), keeper.GetParam(ctx, types.KeyMaxRawDataReportSize))
}

func TestGetSetMaxResultSize(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetParam(ctx, types.KeyMaxResultSize, 1)
	require.Equal(t, uint64(1), keeper.GetParam(ctx, types.KeyMaxResultSize))
	keeper.SetParam(ctx, types.KeyMaxResultSize, 2)
	require.Equal(t, uint64(2), keeper.GetParam(ctx, types.KeyMaxResultSize))
}

func TestGetSetGasPerRawDataRequestPerValidator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, uint64(3000))
	require.Equal(t, uint64(3000), keeper.GetParam(ctx, types.KeyGasPerRawDataRequestPerValidator))
	keeper.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, uint64(5000))
	require.Equal(t, uint64(5000), keeper.GetParam(ctx, types.KeyGasPerRawDataRequestPerValidator))
}

func TestGetSetParams(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	keeper.SetParam(ctx, types.KeyMaxRawRequestCount, 1)
	keeper.SetParam(ctx, types.KeyMaxRawDataReportSize, 1)
	keeper.SetParam(ctx, types.KeyMaxResultSize, 1)
	keeper.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, 1000)
	keeper.SetParam(ctx, types.KeyExpirationBlockCount, 30)
	require.Equal(t, types.NewParams(1, 1, 1, 1000, 30), keeper.GetParams(ctx))

	keeper.SetParam(ctx, types.KeyMaxRawRequestCount, 2)
	keeper.SetParam(ctx, types.KeyMaxRawDataReportSize, 2)
	keeper.SetParam(ctx, types.KeyMaxResultSize, 2)
	keeper.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, 2000)
	keeper.SetParam(ctx, types.KeyExpirationBlockCount, 40)
	require.Equal(t, types.NewParams(2, 2, 2, 2000, 40), keeper.GetParams(ctx))
}
