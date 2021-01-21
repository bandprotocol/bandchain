package yoda

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/GeoDB-Limited/odincore/chain/app"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

func TestGetSignBytesVerificationMessage(t *testing.T) {
	app.SetBech32AddressPrefixesAndBip44CoinType(sdk.GetConfig())
	validator, _ := sdk.ValAddressFromBech32("bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec")
	vmsg := NewVerificationMessage("bandchain", validator, types.RequestID(1), types.ExternalID(1))
	expected, _ := hex.DecodeString("7b22636861696e5f6964223a2262616e64636861696e222c2265787465726e616c5f6964223a2231222c22726571756573745f6964223a2231222c2276616c696461746f72223a2262616e6476616c6f706572317034307968337a6b6d6863763065637170336d63617a7938337361353772676a646536776563227d")
	require.Equal(t, expected, vmsg.GetSignBytes())
}
