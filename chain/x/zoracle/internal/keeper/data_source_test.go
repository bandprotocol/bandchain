package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func mockDataSource(ctx sdk.Context, keeper Keeper) sdk.Error {
	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	executable := []byte("executable")
	return keeper.AddDataSource(ctx, owner, name, fee, executable)
}

func TestGetterSetterDataSource(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetDataSource(ctx, 1)
	require.NotNil(t, err)

	err = mockDataSource(ctx, keeper)
	require.Nil(t, err)
	actualDataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, sdk.AccAddress([]byte("owner")), actualDataSource.Owner)
	require.Equal(t, "data_source", actualDataSource.Name)
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 10)), actualDataSource.Fee)
	require.Equal(t, []byte("executable"), actualDataSource.Executable)
}

func TestEditDataSource(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	err := mockDataSource(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner"))
	newName := "data_source"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	newExecutable := []byte("executable")

	err = keeper.EditDataSource(ctx, 1, newOwner, newName, newFee, newExecutable)
	require.Nil(t, err)

	actualDataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, actualDataSource.Owner)
	require.Equal(t, newName, actualDataSource.Name)
	require.Equal(t, newFee, actualDataSource.Fee)
	require.Equal(t, newExecutable, actualDataSource.Executable)
}
