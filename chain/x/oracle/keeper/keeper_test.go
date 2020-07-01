package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestGetRequestCount(t *testing.T) {
	_, ctx, k := createTestInput()

	// Initial request count must be 0
	require.Equal(t, int64(0), k.GetRequestCount(ctx))
}

func TestGetSetRequestLastExpiredID(t *testing.T) {
	_, ctx, k := createTestInput()
	// Initial last expired request must be 0
	require.Equal(t, int64(0), k.GetRequestLastExpired(ctx))
	k.SetRequestLastExpired(ctx, 20)
	require.Equal(t, int64(20), k.GetRequestLastExpired(ctx))
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

func TestGetSetParams(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetParam(ctx, types.KeyMaxRawRequestCount, 1)
	k.SetParam(ctx, types.KeyMaxAskCount, 10)
	k.SetParam(ctx, types.KeyExpirationBlockCount, 30)
	k.SetParam(ctx, types.KeyBaseRequestGas, 50000)
	k.SetParam(ctx, types.KeyPerValidatorRequestGas, 3000)
	k.SetParam(ctx, types.KeySamplingTryCount, 3)
	k.SetParam(ctx, types.KeyOracleRewardPercentage, 50)
	k.SetParam(ctx, types.KeyInactivePenaltyDuration, 1000)
	require.Equal(t, types.NewParams(1, 10, 30, 50000, 3000, 3, 50, 1000), k.GetParams(ctx))

	k.SetParam(ctx, types.KeyMaxRawRequestCount, 2)
	k.SetParam(ctx, types.KeyMaxAskCount, 20)
	k.SetParam(ctx, types.KeyExpirationBlockCount, 40)
	k.SetParam(ctx, types.KeyBaseRequestGas, 150000)
	k.SetParam(ctx, types.KeyPerValidatorRequestGas, 30000)
	k.SetParam(ctx, types.KeySamplingTryCount, 5)
	k.SetParam(ctx, types.KeyOracleRewardPercentage, 80)
	k.SetParam(ctx, types.KeyInactivePenaltyDuration, 10000)
	require.Equal(t, types.NewParams(2, 20, 40, 150000, 30000, 5, 80, 10000), k.GetParams(ctx))
}
