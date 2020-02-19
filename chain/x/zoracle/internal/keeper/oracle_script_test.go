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

	// Set MaxOracleScriptCodeSize to 20
	keeper.SetMaxOracleScriptCodeSize(ctx, 20)
	err = mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	actualOracleScript, err := keeper.GetOracleScript(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, sdk.AccAddress([]byte("owner")), actualOracleScript.Owner)
	require.Equal(t, "oracle_script", actualOracleScript.Name)
	require.Equal(t, []byte("code"), actualOracleScript.Code)
}

func TestAddTooLongOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetOracleScript(ctx, 1)
	require.NotNil(t, err)

	// Set MaxOracleScriptCodeSize to 20
	keeper.SetMaxOracleScriptCodeSize(ctx, 20)

	owner := sdk.AccAddress([]byte("owner"))
	name := "oracle_script"
	code := []byte("The number of bytes of this oracle script is 82 which is obviously longer than 20.")

	err = keeper.AddOracleScript(ctx, owner, name, code)
	require.NotNil(t, err)
}

func TestEditOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Set MaxOracleScriptCodeSize to 20
	keeper.SetMaxOracleScriptCodeSize(ctx, 20)
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

func TestEditTooLongOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Set MaxOracleScriptCodeSize to 20
	keeper.SetMaxOracleScriptCodeSize(ctx, 20)
	err := mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "oracle_script_2"
	newTooLongCode := []byte("The number of bytes of this oracle script is 82 which is obviously longer than 20.")

	err = keeper.EditOracleScript(ctx, 1, newOwner, newName, newTooLongCode)
	require.NotNil(t, err)
}
