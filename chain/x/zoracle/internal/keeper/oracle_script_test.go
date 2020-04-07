package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
)

func mockOracleScript(ctx sdk.Context, keeper Keeper) error {
	owner := sdk.AccAddress([]byte("owner"))
	name := "oracle_script"
	description := "description"
	code := []byte("code")
	_, err := keeper.AddOracleScript(ctx, owner, name, description, code)
	return err
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

func TestAddOracleScriptMustReturnCorrectID(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetDataSource(ctx, 1)
	require.NotNil(t, err)

	owner := sdk.AccAddress([]byte("owner"))
	name := "oracle_script"
	description := "description"
	code := []byte("code")

	id, err := keeper.AddOracleScript(ctx, owner, name, description, code)
	require.Nil(t, err)
	require.Equal(t, types.OracleScriptID(1), id)

	id, err = keeper.AddOracleScript(ctx, owner, name, description, code)
	require.Nil(t, err)
	require.Equal(t, types.OracleScriptID(2), id)
}

func TestAddTooLongOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetOracleScript(ctx, 1)
	require.NotNil(t, err)

	// Set MaxOracleScriptCodeSize to 20
	keeper.SetParam(ctx, types.KeyMaxOracleScriptCodeSize, 20)

	owner := sdk.AccAddress([]byte("owner"))
	name := "oracle_script"
	description := "description"
	code := []byte("The number of bytes of this oracle script is 82 which is obviously longer than 20.")

	_, err = keeper.AddOracleScript(ctx, owner, name, description, code)
	require.NotNil(t, err)
}

func TestAddTooLongOracleScriptName(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetOracleScript(ctx, 1)
	require.NotNil(t, err)

	// Set MaxNameLength to 5
	keeper.SetParam(ctx, types.KeyMaxNameLength, 5)

	owner := sdk.AccAddress([]byte("owner"))
	tooLongName := "oracle_script"
	description := "description"
	code := []byte("code")

	_, err = keeper.AddOracleScript(ctx, owner, tooLongName, description, code)
	require.NotNil(t, err)
}

func TestAddTooLongOracleScriptDescription(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	_, err := keeper.GetOracleScript(ctx, 1)
	require.NotNil(t, err)

	// Set MaxNameLength to 5
	keeper.SetParam(ctx, types.KeyMaxDescriptionLength, 5)

	owner := sdk.AccAddress([]byte("owner"))
	name := "oracle_script"
	tooLongDescription := "description"
	code := []byte("code")

	_, err = keeper.AddOracleScript(ctx, owner, name, tooLongDescription, code)
	require.NotNil(t, err)
}
func TestEditOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	err := mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "oracle_script_2"
	newDescription := "description_2"
	newCode := []byte("code_2")

	err = keeper.EditOracleScript(ctx, 1, newOwner, newName, newDescription, newCode)
	require.Nil(t, err)

	expect, err := keeper.GetOracleScript(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, expect.Owner)
	require.Equal(t, newName, expect.Name)
	require.Equal(t, newDescription, expect.Description)
	require.Equal(t, newCode, expect.Code)
}

func TestEditTooLongOracleScript(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Set MaxOracleScriptCodeSize to 20
	keeper.SetParam(ctx, types.KeyMaxOracleScriptCodeSize, 20)
	err := mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "oracle_script_2"
	newDescription := "description_2"
	newTooLongCode := []byte("The number of bytes of this oracle script is 82 which is obviously longer than 20.")

	err = keeper.EditOracleScript(ctx, 1, newOwner, newName, newDescription, newTooLongCode)
	require.NotNil(t, err)
}
func TestEditOracleScriptTooLongName(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Set MaxOracleScriptCodeSize to 20
	keeper.SetParam(ctx, types.KeyMaxNameLength, 20)
	err := mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newTooLongName := "Tooooo Looooooong Nameeeee"
	tooLongDescription := "description"
	newCode := []byte("code")

	err = keeper.EditOracleScript(ctx, 1, newOwner, newTooLongName, tooLongDescription, newCode)
	require.NotNil(t, err)
}
func TestEditOracleScriptTooLongDescription(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// Set MaxDescriptionLength to 11
	keeper.SetParam(ctx, types.KeyMaxDescriptionLength, 11)
	err := mockOracleScript(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "oracle_script_2"
	tooLongDescription := "too long description"
	newCode := []byte("code")

	err = keeper.EditOracleScript(ctx, 1, newOwner, newName, tooLongDescription, newCode)
	require.NotNil(t, err)
}

func TestGetAllOracleScripts(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	oracleScripts := []types.OracleScript{
		types.NewOracleScript(
			sdk.AccAddress([]byte("owner1")),
			"name1",
			"description1",
			[]byte("code1"),
		),
		types.NewOracleScript(
			sdk.AccAddress([]byte("owner2")),
			"name2",
			"description2",
			[]byte("code2"),
		),
	}
	keeper.SetOracleScript(ctx, 1, oracleScripts[0])
	keeper.SetOracleScript(ctx, 2, oracleScripts[1])

	require.Equal(t, oracleScripts, keeper.GetAllOracleScripts(ctx))
}
