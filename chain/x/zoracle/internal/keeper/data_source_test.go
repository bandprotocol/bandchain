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

	// Size of "executable" is  10 bytes
	executable := []byte("executable")
	return keeper.AddDataSource(ctx, owner, name, fee, executable)
}

func TestGetterSetterDataSource(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetDataSource(ctx, 1)
	require.NotNil(t, err)

	// Set MaxDataSourceExecutableSize to 20
	keeper.SetMaxDataSourceExecutableSize(ctx, 20)
	err = mockDataSource(ctx, keeper)
	require.Nil(t, err)

	actualDataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, sdk.AccAddress([]byte("owner")), actualDataSource.Owner)
	require.Equal(t, "data_source", actualDataSource.Name)
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 10)), actualDataSource.Fee)
	require.Equal(t, []byte("executable"), actualDataSource.Executable)
}

func TestAddTooLongDataSource(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetDataSource(ctx, 1)
	require.NotNil(t, err)

	// Set MaxDataSourceExecutableSize to 20
	keeper.SetMaxDataSourceExecutableSize(ctx, 20)

	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	tooLongExecutable := []byte("The number of bytes of this data source is 80 which is obviously longer than 20.")

	err = keeper.AddDataSource(ctx, owner, name, fee, tooLongExecutable)
	require.NotNil(t, err)
}

func TestEditDataSource(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Set MaxDataSourceExecutableSize to 20
	keeper.SetMaxDataSourceExecutableSize(ctx, 20)
	err := mockDataSource(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "data_source_2"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 1))
	newExecutable := []byte("executable_2")

	err = keeper.EditDataSource(ctx, 1, newOwner, newName, newFee, newExecutable)
	require.Nil(t, err)

	actualDataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, actualDataSource.Owner)
	require.Equal(t, newName, actualDataSource.Name)
	require.Equal(t, newFee, actualDataSource.Fee)
	require.Equal(t, newExecutable, actualDataSource.Executable)
}

func TestEditTooLongDataSource(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Set MaxDataSourceExecutableSize to 20
	keeper.SetMaxDataSourceExecutableSize(ctx, 20)
	err := mockDataSource(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "data_source_2"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 1))
	newTooLongExecutable := []byte("The number of bytes of this data source is 80 which is obviously longer than 20.")

	err = keeper.EditDataSource(ctx, 1, newOwner, newName, newFee, newTooLongExecutable)
	require.NotNil(t, err)
}
