package keeper

import (
	"testing"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
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

func TestGetSetMaxDataSourceExecutableSize(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxDataSourceExecutableSize(ctx, int64(1))
	require.Equal(t, int64(1), keeper.MaxDataSourceExecutableSize(ctx))
	keeper.SetMaxDataSourceExecutableSize(ctx, int64(2))
	require.Equal(t, int64(2), keeper.MaxDataSourceExecutableSize(ctx))
}

func TestGetSetMaxOracleScriptCodeSize(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxOracleScriptCodeSize(ctx, int64(1))
	require.Equal(t, int64(1), keeper.MaxOracleScriptCodeSize(ctx))
	keeper.SetMaxOracleScriptCodeSize(ctx, int64(2))
	require.Equal(t, int64(2), keeper.MaxOracleScriptCodeSize(ctx))
}

func TestGetSetMaxCalldataSize(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxCalldataSize(ctx, int64(1))
	require.Equal(t, int64(1), keeper.MaxCalldataSize(ctx))
	keeper.SetMaxCalldataSize(ctx, int64(2))
	require.Equal(t, int64(2), keeper.MaxCalldataSize(ctx))
}

func TestGetSetMaxDataSourceCountPerRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxDataSourceCountPerRequest(ctx, int64(1))
	require.Equal(t, int64(1), keeper.MaxDataSourceCountPerRequest(ctx))
	keeper.SetMaxDataSourceCountPerRequest(ctx, int64(2))
	require.Equal(t, int64(2), keeper.MaxDataSourceCountPerRequest(ctx))
}

func TestGetSetMaxRawDataReportSize(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxRawDataReportSize(ctx, int64(1))
	require.Equal(t, int64(1), keeper.MaxRawDataReportSize(ctx))
	keeper.SetMaxRawDataReportSize(ctx, int64(2))
	require.Equal(t, int64(2), keeper.MaxRawDataReportSize(ctx))
}

func TestGetSetMaxResultSize(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxResultSize(ctx, int64(1))
	require.Equal(t, int64(1), keeper.MaxResultSize(ctx))
	keeper.SetMaxResultSize(ctx, int64(2))
	require.Equal(t, int64(2), keeper.MaxResultSize(ctx))
}

func TestGetSetEndBlockExecuteGasLimit(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetEndBlockExecuteGasLimit(ctx, uint64(3000))
	require.Equal(t, uint64(3000), keeper.EndBlockExecuteGasLimit(ctx))
	keeper.SetEndBlockExecuteGasLimit(ctx, uint64(5000))
	require.Equal(t, uint64(5000), keeper.EndBlockExecuteGasLimit(ctx))
}

func TestGetSetParams(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	keeper.SetMaxDataSourceExecutableSize(ctx, int64(1))
	keeper.SetMaxOracleScriptCodeSize(ctx, int64(1))
	keeper.SetMaxCalldataSize(ctx, int64(1))
	keeper.SetMaxDataSourceCountPerRequest(ctx, int64(1))
	keeper.SetMaxRawDataReportSize(ctx, int64(1))
	keeper.SetMaxResultSize(ctx, int64(1))
	keeper.SetEndBlockExecuteGasLimit(ctx, uint64(200000))
	keeper.SetMaxNameLength(ctx, int64(1))
	require.Equal(t, types.NewParams(1, 1, 1, 1, 1, 1, 200000, 1), keeper.GetParams(ctx))

	keeper.SetMaxDataSourceExecutableSize(ctx, int64(2))
	keeper.SetMaxOracleScriptCodeSize(ctx, int64(2))
	keeper.SetMaxCalldataSize(ctx, int64(2))
	keeper.SetMaxDataSourceCountPerRequest(ctx, int64(2))
	keeper.SetMaxRawDataReportSize(ctx, int64(2))
	keeper.SetMaxResultSize(ctx, int64(2))
	keeper.SetEndBlockExecuteGasLimit(ctx, uint64(300000))
	keeper.SetMaxNameLength(ctx, int64(2))
	require.Equal(t, types.NewParams(2, 2, 2, 2, 2, 2, 300000, 2), keeper.GetParams(ctx))
}
