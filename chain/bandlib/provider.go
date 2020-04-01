package bandlib

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto"
)

// BandProvider contains context, txBuilder, private key, and address
type BandProvider struct {
	cliCtx  context.CLIContext
	txBldr  authtypes.TxBuilder
	addr    sdk.AccAddress
	privKey crypto.PrivKey
}

func privKeyToBandAccAddress(privKey crypto.PrivKey) sdk.AccAddress {
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(config)
	return sdk.AccAddress(privKey.PubKey().Address())
}

// NewBandProvider creates new BandProvider create new cliCtx and txBldr
func NewBandProvider(nodeURI string, privKey crypto.PrivKey) (BandProvider, error) {
	cdc := NewCodec()
	addr := privKeyToBandAccAddress(privKey)
	cliCtx := NewCLIContext(nodeURI, addr).WithCodec(cdc)
	num, _, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(addr)
	if err != nil {
		return BandProvider{}, err
	}

	return BandProvider{
		cliCtx:  cliCtx,
		txBldr:  NewTxBuilder(utils.GetTxEncoder(cdc)).WithAccountNumber(num),
		addr:    addr,
		privKey: privKey,
	}, nil
}

func (provider *BandProvider) Sender() sdk.AccAddress {
	return provider.addr
}

func (provider *BandProvider) SendTransaction(
	msgs []sdk.Msg, seq, gas uint64,
	memo, fees, broadcastMode string,
) (sdk.TxResponse, error) {
	cliCtx := provider.cliCtx.WithBroadcastMode(broadcastMode)
	txBldr := provider.txBldr.WithSequence(seq).WithGas(gas).WithMemo(memo).WithFees(fees)
	return completeAndBroadcastTxCLI(cliCtx, txBldr, msgs, provider.privKey)
}

func (provider *BandProvider) QueryTx(hashHexStr string) (sdk.TxResponse, error) {
	return utils.QueryTx(provider.cliCtx, hashHexStr)
}

func (provider *BandProvider) GetSequenceFromChain() (uint64, error) {
	_, seq, err := authtypes.NewAccountRetriever(provider.cliCtx).GetAccountNumberSequence(provider.addr)
	return seq, err
}

func (provider *BandProvider) SignStdTx(seq uint64, msgs []sdk.Msg) (authtypes.StdTx, error) {
	txBldr := provider.txBldr.WithSequence(seq).WithGas(200000).WithMemo("").WithFees("").WithAccountNumber(1)
	fmt.Println(txBldr)
	// build and sign the transaction
	signMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return authtypes.StdTx{}, err
	}

	sigBytes, err := provider.privKey.Sign(signMsg.Bytes())
	if err != nil {
		return authtypes.StdTx{}, err
	}
	sig := authtypes.StdSignature{
		PubKey:    provider.privKey.PubKey(),
		Signature: sigBytes,
	}
	return authtypes.NewStdTx(signMsg.Msgs, signMsg.Fee, []authtypes.StdSignature{sig}, signMsg.Memo), nil
}
