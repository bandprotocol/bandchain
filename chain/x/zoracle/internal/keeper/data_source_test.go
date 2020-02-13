package keeper

import (
	"testing"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	"github.com/stretchr/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGetterSetterDataSource(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetDataSource(ctx, 1)
	require.NotNil(t, err)

	dataSource := types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source", 
		sdk.NewCoins(sdk.NewInt64Coin("uband", 10)), 
		[]byte("executable"),
	)

	keeper.SetDataSource(ctx, 1, dataSource)
	actualDataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, dataSource, actualDataSource)
}

