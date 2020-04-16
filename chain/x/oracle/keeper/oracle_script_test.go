package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	me "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

var BASIC_OS_NAME = "OS_NAME"
var BASIC_OS_DESC = "OS_DESCRIPTION"
var BASIC_OS_CODE = []byte("OS_WASM_CODE")

func addBasicOracleScript(ctx sdk.Context, k me.Keeper) types.OID {
	id, err := k.AddOracleScript(ctx, Owner.Address, BASIC_OS_NAME, BASIC_OS_DESC, BASIC_OS_CODE)
	if err != nil {
		panic(err)
	}
	return id
}

func TestSetterGetterOracleScript(t *testing.T) {
	_, ctx, k := createTestInput()
	// Getting a non-existent oracle script should return error.
	_, err := k.GetOracleScript(ctx, 42)
	require.Error(t, err)
	require.Panics(t, func() { _ = k.MustGetOracleScript(ctx, 42) })
	// Creates some basic oracle scripts
	oracleScript1 := types.NewOracleScript(Alice.Address, "NAME1", "DESCRIPTION1", []byte("code1"))
	oracleScript2 := types.NewOracleScript(Bob.Address, "NAME2", "DESCRIPTION2", []byte("code2"))
	// Sets id 42 with oracle script 1 and id 42 with oracle script 2.
	k.SetOracleScript(ctx, 42, oracleScript1)
	k.SetOracleScript(ctx, 43, oracleScript2)
	// Checks that Get and MustGet perform correctly.
	oracleScript1Res, err := k.GetOracleScript(ctx, 42)
	require.Nil(t, err)
	require.Equal(t, oracleScript1, oracleScript1Res)
	require.Equal(t, oracleScript1, k.MustGetOracleScript(ctx, 42))
	oracleScript2Res, err := k.GetOracleScript(ctx, 43)
	require.Nil(t, err)
	require.Equal(t, oracleScript2, oracleScript2Res)
	require.Equal(t, oracleScript2, k.MustGetOracleScript(ctx, 43))
	// Replaces id 42 with another oracle script.
	k.SetOracleScript(ctx, 42, oracleScript2)
	require.NotEqual(t, oracleScript1, k.MustGetOracleScript(ctx, 42))
	require.Equal(t, oracleScript2, k.MustGetOracleScript(ctx, 42))
}

func TestAddEditOracleScriptBasic(t *testing.T) {
	_, ctx, k := createTestInput()
	// Creates some basic oracle scripts
	oracleScript1 := types.NewOracleScript(Alice.Address, "NAME1", "DESCRIPTION1", []byte("code1"))
	oracleScript2 := types.NewOracleScript(Bob.Address, "NAME2", "DESCRIPTION2", []byte("code2"))
	// Adds a new oracle script to the store. We should be able to retreive it back.
	id, err := k.AddOracleScript(ctx,
		oracleScript1.Owner, oracleScript1.Name, oracleScript1.Description, oracleScript1.Code,
	)
	require.Nil(t, err)
	require.Equal(t, oracleScript1, k.MustGetOracleScript(ctx, id))
	require.NotEqual(t, oracleScript2, k.MustGetOracleScript(ctx, id))
	// Edits the oracle script. We should get the updated oracle script.
	err = k.EditOracleScript(ctx, id,
		oracleScript2.Owner, oracleScript2.Name, oracleScript2.Description, oracleScript2.Code,
	)
	require.Nil(t, err)
	require.NotEqual(t, oracleScript1, k.MustGetOracleScript(ctx, id))
	require.Equal(t, oracleScript2, k.MustGetOracleScript(ctx, id))
}

func TestAddOracleScriptMustReturnCorrectID(t *testing.T) {
	_, ctx, k := createTestInput()
	// Initially we expect the oracle script count to be zero.
	count := k.GetOracleScriptCount(ctx)
	require.Equal(t, count, int64(0))
	// Every new oracle script we add should return a new ID.
	id1, err := k.AddOracleScript(ctx, Owner.Address, BASIC_OS_NAME, BASIC_OS_DESC, BASIC_OS_CODE)
	require.Nil(t, err)
	require.Equal(t, id1, types.OID(1))
	// Adds another oracle script so now ID should be 2.
	id2, err := k.AddOracleScript(ctx, Owner.Address, BASIC_OS_NAME, BASIC_OS_DESC, BASIC_OS_CODE)
	require.Nil(t, err)
	require.Equal(t, id2, types.OID(2))
	// Finally we expect the oracle script to increase to 2 since we added two oracle scripts.
	count = k.GetOracleScriptCount(ctx)
	require.Equal(t, count, int64(2))
}

func TestEditNonExistentOracleScript(t *testing.T) {
	_, ctx, k := createTestInput()
	// Editing a non-existent data should return error.
	err := k.EditOracleScript(ctx, types.OID(42),
		Owner.Address, BASIC_OS_NAME, BASIC_OS_DESC, BASIC_OS_CODE,
	)
	require.Error(t, err)
}

func TestAddOracleScriptTooLongName(t *testing.T) {
	_, ctx, k := createTestInput()
	// Sets max name length to 9. We should fail to add oracle script with name length 10.
	k.SetParam(ctx, types.KeyMaxNameLength, 9)
	_, err := k.AddOracleScript(ctx, Owner.Address, "0123456789", BASIC_OS_DESC, BASIC_OS_CODE)
	require.Error(t, err)
	// Sets max name length to 10. We should now be able to add the oracle script.
	k.SetParam(ctx, types.KeyMaxNameLength, 10)
	_, err = k.AddOracleScript(ctx, Owner.Address, "0123456789", BASIC_OS_DESC, BASIC_OS_CODE)
	require.Nil(t, err)
}

func TestEditOracleScriptTooLongName(t *testing.T) {
	_, ctx, k := createTestInput()
	id := addBasicOracleScript(ctx, k)
	oracleScript := k.MustGetOracleScript(ctx, id)
	// Sets max name length to 9. We should fail to edit oracle script with name length 10.
	k.SetParam(ctx, types.KeyMaxNameLength, 9)
	err := k.EditOracleScript(ctx, id,
		oracleScript.Owner, "0123456789", oracleScript.Description, oracleScript.Code,
	)
	require.Error(t, err)
	// Sets max name length to 10. We should now be able to edit the oracle script.
	k.SetParam(ctx, types.KeyMaxNameLength, 10)
	err = k.EditOracleScript(ctx, id,
		oracleScript.Owner, "0123456789", oracleScript.Description, oracleScript.Code,
	)
	require.Nil(t, err)
}

func TestAddOracleScriptTooLongDescription(t *testing.T) {
	_, ctx, k := createTestInput()
	// Sets max desc length to 41. We should fail to add oracle script with desc length 42.
	k.SetParam(ctx, types.KeyMaxDescriptionLength, 41)
	_, err := k.AddOracleScript(ctx,
		Owner.Address, BASIC_OS_NAME, "________THIS_STRING_HAS_SIZE_OF_42________", BASIC_OS_CODE,
	)
	require.Error(t, err)
	// Sets max desc length to 42. We should now be able to add the oracle script.
	k.SetParam(ctx, types.KeyMaxDescriptionLength, 42)
	_, err = k.AddOracleScript(ctx,
		Owner.Address, BASIC_OS_NAME, "________THIS_STRING_HAS_SIZE_OF_42________", BASIC_OS_CODE,
	)
	require.Nil(t, err)
}

func TestEditOracleScriptTooLongDescription(t *testing.T) {
	_, ctx, k := createTestInput()
	id := addBasicOracleScript(ctx, k)
	oracleScript := k.MustGetOracleScript(ctx, id)
	// Sets max desc length to 41. We should fail to edit oracle script with name length 42.
	k.SetParam(ctx, types.KeyMaxDescriptionLength, 41)
	err := k.EditOracleScript(ctx, id,
		oracleScript.Owner, oracleScript.Name, "________THIS_STRING_HAS_SIZE_OF_42________",
		oracleScript.Code,
	)
	require.Error(t, err)
	// Sets max desc length to 42. We should now be able to edit the oracle script.
	k.SetParam(ctx, types.KeyMaxDescriptionLength, 42)
	err = k.EditOracleScript(ctx, id,
		oracleScript.Owner, oracleScript.Name, "________THIS_STRING_HAS_SIZE_OF_42________",
		oracleScript.Code,
	)
	require.Nil(t, err)
}

func TestAddOracleScriptTooBigCode(t *testing.T) {
	_, ctx, k := createTestInput()
	// Sets max code size to 40. We should fail to add oracle script with exec size 42.
	k.SetParam(ctx, types.KeyMaxOracleScriptCodeSize, 40)
	_, err := k.AddOracleScript(ctx,
		Owner.Address, BASIC_OS_NAME, BASIC_OS_DESC,
		[]byte("________THIS_STRING_HAS_SIZE_OF_42________"),
	)
	require.Error(t, err)
	// Sets max code size to 50. We should now be able to add the oracle script.
	k.SetParam(ctx, types.KeyMaxOracleScriptCodeSize, 50)
	_, err = k.AddOracleScript(ctx,
		Owner.Address, BASIC_OS_NAME, BASIC_OS_DESC,
		[]byte("________THIS_STRING_HAS_SIZE_OF_42________"),
	)
	require.Nil(t, err)
}

func TestEditOracleScriptTooBigCode(t *testing.T) {
	_, ctx, k := createTestInput()
	id := addBasicOracleScript(ctx, k)
	oracleScript := k.MustGetOracleScript(ctx, id)
	// Sets max code size to 40. We should fail to edit oracle script with exec size 42.
	k.SetParam(ctx, types.KeyMaxOracleScriptCodeSize, 40)
	err := k.EditOracleScript(ctx, id,
		oracleScript.Owner, oracleScript.Name, oracleScript.Description,
		[]byte("________THIS_STRING_HAS_SIZE_OF_42________"),
	)
	require.Error(t, err)
	// Sets max code size to 50. We should now be able to edit the oracle script.
	k.SetParam(ctx, types.KeyMaxOracleScriptCodeSize, 50)
	err = k.EditOracleScript(ctx, id,
		oracleScript.Owner, oracleScript.Name, oracleScript.Description,
		[]byte("________THIS_STRING_HAS_SIZE_OF_42________"),
	)
	require.Nil(t, err)
}

func TestGetAllOracleScripts(t *testing.T) {
	_, ctx, k := createTestInput()
	// Sets the oracle scripts to the storage.
	oracleScripts := []types.OracleScript{
		types.NewOracleScript(Alice.Address, "NAME1", "DESCRIPTION1", []byte("code1")),
		types.NewOracleScript(Bob.Address, "NAME2", "DESCRIPTION2", []byte("code2")),
	}
	k.SetOracleScript(ctx, 1, oracleScripts[0])
	k.SetOracleScript(ctx, 2, oracleScripts[1])
	// We should now be able to get all the existing oracle scripts.
	require.Equal(t, oracleScripts, k.GetAllOracleScripts(ctx))
}
