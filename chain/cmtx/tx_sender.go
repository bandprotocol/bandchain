package cmtx

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/tendermint/tendermint/crypto"
)

// TxSender contain codec and create context to create tx
type TxSender struct {
	cdc     *codec.Codec
	addr    sdk.AccAddress
	privKey crypto.PrivKey
}

const Bech32MainPrefix = "band"

func SetBech32AddressPrefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32MainPrefix, Bech32MainPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator, Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
}

func privKeyToBandAccAddress(privKey crypto.PrivKey) sdk.AccAddress {
	config := sdk.GetConfig()
	SetBech32AddressPrefixes(config)
	return sdk.AccAddress(privKey.PubKey().Address().Bytes())
}

func NewTxSender(privKey crypto.PrivKey) TxSender {
	return TxSender{
		cdc:     NewCodec(),
		addr:    privKeyToBandAccAddress(privKey),
		privKey: privKey,
	}
}

func (tx TxSender) Sender() sdk.AccAddress {
	return tx.addr
}

func (tx TxSender) SendTransaction(msg sdk.Msg) (string, error) {
	cliCtx := NewCLIContext(tx.addr).WithCodec(tx.cdc)
	txBldr := NewTxBuilder(utils.GetTxEncoder(tx.cdc))

	return completeAndBroadcastTxCLI(cliCtx, txBldr, []sdk.Msg{msg}, tx.privKey)
}

func (tx TxSender) QueryTx(hashHexStr string) (sdk.TxResponse, error) {
	cliCtx := NewCLIContext(tx.addr).WithCodec(tx.cdc)

	return utils.QueryTx(cliCtx, hashHexStr)
}
