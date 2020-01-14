package types

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetCodeHash(t *testing.T) {
	code := []byte("This is code")
	expectedCodeHash1, _ := hex.DecodeString("9a794cce93b8fa9b2477f1498bc63f1b602b2dd8edd06747183298822183e935")
	expectedCodeHash2, _ := hex.DecodeString("3619be5b7c53a74dec4d0ca50825681accd4c9ed7471538b2d1438794bf1cd4c")
	sc1 := NewStoredCode(code, "script1", sdk.AccAddress("owner"))
	sc2 := NewStoredCode(code, "script2", sdk.AccAddress("owner"))

	require.Equal(t, expectedCodeHash1, sc1.GetCodeHash())
	require.Equal(t, expectedCodeHash2, sc2.GetCodeHash())
}
