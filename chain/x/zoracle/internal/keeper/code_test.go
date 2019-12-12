package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestGetterSetterCode(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	code := []byte("This is code")
	codeHash := crypto.Keccak256(code)
	_, err := keeper.GetCode(ctx, codeHash)
	require.NotNil(t, err)

	actualCodeHash := keeper.SetCode(ctx, code)
	actualCode, err := keeper.GetCode(ctx, actualCodeHash)
	require.Nil(t, err)
	require.Equal(t, code, actualCode)
	require.Equal(t, codeHash, actualCodeHash)
}
