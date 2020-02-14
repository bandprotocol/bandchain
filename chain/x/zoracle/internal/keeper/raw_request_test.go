package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

func TestGettterSetterRawDataRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	err := keeper.SetRawDataRequest(ctx, 1, 1, 0, []byte("calldata1"))
	require.NotNil(t, err)

	_, err = keeper.GetRawDataRequest(ctx, 1, 1)
	require.NotNil(t, err)

	request := newDefaultRequest()
	keeper.SetRequest(ctx, 1, request)

	err = keeper.SetRawDataRequest(ctx, 1, 1, 0, []byte("calldata1"))
	require.NotNil(t, err)

	dataSource := types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 10)),
		[]byte("executable"),
	)
	keeper.SetDataSource(ctx, 1, dataSource)

	err = keeper.SetRawDataRequest(ctx, 1, 42, 1, []byte("calldata1"))
	require.Nil(t, err)

	rawRequest, err := keeper.GetRawDataRequest(ctx, 1, 42)
	require.Nil(t, err)

	expect := types.NewRawDataRequest(1, []byte("calldata1"))
	require.Equal(t, expect, rawRequest)

	_, err = keeper.GetRawDataRequest(ctx, 1, 3)
	require.Equal(t, types.CodeRequestNotFound, err.Code())

	// Add new datasource
	dataSource2 := types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source2",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 20)),
		[]byte("executable2"),
	)
	keeper.SetDataSource(ctx, 2, dataSource2)

	// Cannot set to existed external id
	err = keeper.SetRawDataRequest(ctx, 1, 42, 2, []byte("calldata3"))
	require.NotNil(t, err)
}

func TestGetRawDataRequestCount(t *testing.T) {
	// TODO: Write test
}

func TestGetRawDataRequests(t *testing.T) {
	// TODO: Write test
}
