package bandlib

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/bandprotocol/d3n/chain/x/zoracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	zoracle.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
