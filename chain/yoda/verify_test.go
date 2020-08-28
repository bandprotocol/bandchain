package yoda

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestGetSignBytesVerificationMessage(t *testing.T) {
	app.SetBech32AddressPrefixesAndBip44CoinType(sdk.GetConfig())
	validator, _ := sdk.ValAddressFromBech32("bandvaloper13eznuehmqzd3r84fkxu8wklxl22r2qfm8f05zn")
	vmsg := NewVerificationMessage("bandchain", validator, types.RequestID(1), types.ExternalID(2))
	expected, _ := hex.DecodeString("7b22636861696e5f6964223a2262616e64636861696e222c2265787465726e616c5f6964223a2231222c22726571756573745f6964223a2233222c2276616c696461746f72223a2262616e6476616c6f7065723133657a6e7565686d717a6433723834666b787538776b6c786c3232723271666d386630357a6e227d")
	require.Equal(t, expected, vmsg.GetSignBytes())
}
