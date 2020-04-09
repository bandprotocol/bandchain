package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/zoracle/types"
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
	require.NotNil(t, err)
}

func TestAddNewRawDataRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

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
		"description",
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
	require.NotNil(t, err)

	// Add new datasource
	dataSource2 := types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source2",
		"description2",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 20)),
		[]byte("executable2"),
	)
	keeper.SetDataSource(ctx, 2, dataSource2)

	// Cannot set to existed external id
	err = keeper.AddNewRawDataRequest(ctx, 1, 42, 2, []byte("calldata3"))
	require.NotNil(t, err)
}

func TestGasConsumeByAddNewRawDataRequest(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	request := newDefaultRequest()

	// Set GasPerRawDataRequestPerValidator to 10000
	keeper.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, 10000)
	keeper.SetRequest(ctx, 1, request)

	dataSource := types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source",
		"description",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 10)),
		[]byte("executable"),
	)
	keeper.SetDataSource(ctx, 1, dataSource)

	beforeGas := ctx.GasMeter().GasConsumed()
	err := keeper.AddNewRawDataRequest(ctx, 1, 42, 1, []byte("calldata1"))
	require.Nil(t, err)

	gasUsed := ctx.GasMeter().GasConsumed() - beforeGas

	// Set GasPerRawDataRequestPerValidator to 25000
	keeper.SetParam(ctx, types.KeyGasPerRawDataRequestPerValidator, 25000)
	keeper.SetRequest(ctx, 2, request)
	beforeGas = ctx.GasMeter().GasConsumed()
	err = keeper.AddNewRawDataRequest(ctx, 2, 42, 1, []byte("calldata1"))
	require.Nil(t, err)

	gasUsed2 := ctx.GasMeter().GasConsumed() - beforeGas
	require.Equal(t, uint64(30000), gasUsed2-gasUsed)
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

func TestGetRawDataRequestWithExternalIDs(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	keeper.SetRawDataRequest(ctx, 1, 1, types.NewRawDataRequest(0, []byte("calldata1")))
	keeper.SetRawDataRequest(ctx, 1, 2, types.NewRawDataRequest(1, []byte("calldata2")))

	keeper.SetRawDataRequest(ctx, 2, 1, types.NewRawDataRequest(0, []byte("calldata1")))
	keeper.SetRawDataRequest(ctx, 2, 3, types.NewRawDataRequest(1, []byte("calldata2")))
	keeper.SetRawDataRequest(ctx, 2, 2, types.NewRawDataRequest(2, []byte("calldata3")))

	ans1 := []types.RawDataRequestWithExternalID{
		types.NewRawDataRequestWithExternalID(
			1,
			types.NewRawDataRequest(0, []byte("calldata1")),
		),
		types.NewRawDataRequestWithExternalID(
			2,
			types.NewRawDataRequest(1, []byte("calldata2")),
		),
	}

	ans2 := []types.RawDataRequestWithExternalID{
		types.NewRawDataRequestWithExternalID(
			1,
			types.NewRawDataRequest(0, []byte("calldata1")),
		),
		types.NewRawDataRequestWithExternalID(
			2,
			types.NewRawDataRequest(2, []byte("calldata3")),
		),
		types.NewRawDataRequestWithExternalID(
			3,
			types.NewRawDataRequest(1, []byte("calldata2")),
		),
	}
	require.Equal(t, ans1, keeper.GetRawDataRequestWithExternalIDs(ctx, 1))
	require.Equal(t, ans2, keeper.GetRawDataRequestWithExternalIDs(ctx, 2))
	require.Equal(t, []types.RawDataRequestWithExternalID{}, keeper.GetRawDataRequestWithExternalIDs(ctx, 3))
	require.Equal(t, []types.RawDataRequestWithExternalID{}, keeper.GetRawDataRequestWithExternalIDs(ctx, -1))
}
