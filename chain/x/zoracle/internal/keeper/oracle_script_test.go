package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func mockOracleScript(ctx sdk.Context, keeper Keeper) sdk.Error {
	owner := sdk.AccAddress([]byte("owner"))
	name := "oracle_script"
	code := []byte("code")
	return keeper.AddOracleScript(ctx, owner, name, code)
}

func TestGetterSetterOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetOracleScript(ctx, 1)
	require.NotNil(t, err)

	err = mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	actualOracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, sdk.AccAddress([]byte("owner")), actualOracleScript.Owner)
	require.Equal(t, "oracle_script", actualOracleScript.Name)
	require.Equal(t, []byte("code"), actualOracleScript.Code)
}

func TestEditOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	err := mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "oracle_script_2"
	newCode := []byte("code_2")

	err = keeper.EditOracleScript(ctx, 1, newOwner, newName, newCode)
	require.Nil(t, err)

	expect, err := keeper.GetOracleScript(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, expect.Owner)
	require.Equal(t, newName, expect.Name)
	require.Equal(t, newCode, expect.Code)
}
