package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

func TestGettterSetterRawDataRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	keeper.SetRawDataRequest(ctx, 1, 1, types.NewRawDataRequest(0, []byte("calldata1")))
	keeper.SetRawDataRequest(ctx, 1, 2, types.NewRawDataRequest(1, []byte("calldata2")))

	rawRequest, err := keeper.GetRawDataRequest(ctx, 1, 1)
	require.Nil(t, err)
	expect := types.NewRawDataRequest(0, []byte("calldata1"))
	require.Equal(t, expect, rawRequest)

	_, err = keeper.GetRawDataRequest(ctx, 1, 3)
	require.Equal(t, types.CodeRequestNotFound, err.Code())
}

func TestAddNewRawDataRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.SetMaxCalldataSize(ctx, 20)
	keeper.SetMaxDataSourceCountPerRequest(ctx, 1)

	err := keeper.AddNewRawDataRequest(ctx, 1, 1, 0, []byte("calldata1"))
	require.NotNil(t, err)

	_, err = keeper.GetRawDataRequest(ctx, 1, 1)
	require.NotNil(t, err)

	request := newDefaultRequest()
	keeper.SetRequest(ctx, 1, request)

	err = keeper.AddNewRawDataRequest(ctx, 1, 1, 0, []byte("calldata1"))
	require.NotNil(t, err)

	dataSource := types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 10)),
		[]byte("executable"),
	)
	keeper.SetDataSource(ctx, 1, dataSource)

	err = keeper.AddNewRawDataRequest(ctx, 1, 42, 1, []byte("calldata1"))
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
	err = keeper.AddNewRawDataRequest(ctx, 1, 42, 2, []byte("calldata3"))
	require.NotNil(t, err)
}

func TestGetRawDataRequestCount(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	keeper.SetRawDataRequest(ctx, 1, 1, types.NewRawDataRequest(0, []byte("calldata1")))
	keeper.SetRawDataRequest(ctx, 1, 2, types.NewRawDataRequest(1, []byte("calldata2")))

	keeper.SetRawDataRequest(ctx, 2, 1, types.NewRawDataRequest(0, []byte("calldata1")))
	keeper.SetRawDataRequest(ctx, 2, 2, types.NewRawDataRequest(1, []byte("calldata2")))
	keeper.SetRawDataRequest(ctx, 2, 3, types.NewRawDataRequest(2, []byte("calldata3")))

	require.Equal(t, int64(2), keeper.GetRawDataRequestCount(ctx, 1))
	require.Equal(t, int64(3), keeper.GetRawDataRequestCount(ctx, 2))
	require.Equal(t, int64(0), keeper.GetRawDataRequestCount(ctx, 3))
	require.Equal(t, int64(0), keeper.GetRawDataRequestCount(ctx, -1))
}

func TestGetRawDataRequests(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	keeper.SetRawDataRequest(ctx, 1, 1, types.NewRawDataRequest(0, []byte("calldata1")))
	keeper.SetRawDataRequest(ctx, 1, 2, types.NewRawDataRequest(1, []byte("calldata2")))

	keeper.SetRawDataRequest(ctx, 2, 1, types.NewRawDataRequest(0, []byte("calldata1")))
	keeper.SetRawDataRequest(ctx, 2, 3, types.NewRawDataRequest(2, []byte("calldata3")))
	keeper.SetRawDataRequest(ctx, 2, 2, types.NewRawDataRequest(1, []byte("calldata2")))

	ans1 := []types.RawDataRequest{
		types.NewRawDataRequest(0, []byte("calldata1")),
		types.NewRawDataRequest(1, []byte("calldata2")),
	}

	ans2 := []types.RawDataRequest{
		types.NewRawDataRequest(0, []byte("calldata1")),
		types.NewRawDataRequest(1, []byte("calldata2")),
		types.NewRawDataRequest(2, []byte("calldata3")),
	}
	require.Equal(t, ans1, keeper.GetRawDataRequests(ctx, 1))
	require.Equal(t, ans2, keeper.GetRawDataRequests(ctx, 2))
	require.Equal(t, []types.RawDataRequest{}, keeper.GetRawDataRequests(ctx, 3))
	require.Equal(t, []types.RawDataRequest{}, keeper.GetRawDataRequests(ctx, -1))
}
