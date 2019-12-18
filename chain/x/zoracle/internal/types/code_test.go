package types

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetCodeHash(t *testing.T) {

	code := []byte("This is code")
	expectedCodeHash, _ := hex.DecodeString("15e0082d52b1a8fd85ce99ae62d1b0ac236709d4ae3ab7554bba52931231bd8a")
	sc := NewStoredCode(code, sdk.AccAddress("owner"))

	require.Equal(t, expectedCodeHash, sc.GetCodeHash())
}
