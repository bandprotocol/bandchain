package keeper

import (
	"testing"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetterSetterCode(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	owner := sdk.AccAddress([]byte("owner"))

	code := []byte("This is code")
	codeHash := types.NewStoredCode(code, owner).GetCodeHash()

	_, err := keeper.GetCode(ctx, codeHash)
	require.NotNil(t, err)

	actualCodeHash := keeper.SetCode(ctx, code, owner)
	storedCode, err := keeper.GetCode(ctx, actualCodeHash)
	require.Nil(t, err)
	require.Equal(t, code, storedCode.Code)
	require.Equal(t, owner, storedCode.Owner)
	require.Equal(t, codeHash, actualCodeHash)
}

func TestDeleteCode(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	owner := sdk.AccAddress([]byte("owner"))

	code := []byte("This is code")
	codeHash := types.NewStoredCode(code, owner).GetCodeHash()

	keeper.SetCode(ctx, code, owner)

	keeper.DeleteCode(ctx, codeHash)
	_, err := keeper.GetCode(ctx, codeHash)
	require.NotNil(t, err)
	require.Equal(t, err.Code(), types.CodeInvalidInput)
}

func TestGetCodesIterator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	owner := sdk.AccAddress([]byte("owner"))

	codes := [][]byte{[]byte("This is code"), []byte("This is code2")}

	for _, code := range codes {
		keeper.SetCode(ctx, code, owner)
	}

	iterator := keeper.GetCodesIterator(ctx)
	i := 0
	for ; iterator.Valid(); iterator.Next() {
		require.Equal(t, types.NewStoredCode(codes[i], owner).GetCodeHash(), iterator.Key()[1:])
		i++
	}
	require.Equal(t, len(codes), i)
}
