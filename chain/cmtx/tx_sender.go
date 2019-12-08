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

func NewTxSender(privKey crypto.PrivKey) TxSender {
	return TxSender{
		cdc:     NewCodec(),
		addr:    sdk.AccAddress(privKey.PubKey().Address().Bytes()),
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
