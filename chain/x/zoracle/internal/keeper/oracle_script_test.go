package keeper

import (
	"testing"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetterSetterOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetOracleScript(ctx, 1)
	require.NotNil(t, err)

	oracleScript := types.NewOracleScript(
		sdk.AccAddress([]byte("owner")),
		"oracle_script",
		[]byte("code"),
	)

	keeper.SetOracleScript(ctx, 1, oracleScript)
	actualOracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, oracleScript, actualOracleScript)
}
