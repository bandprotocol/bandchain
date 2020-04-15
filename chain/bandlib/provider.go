package bandlib

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/bandprotocol/bandchain/chain/app"
)

// BandProvider contains context, txBuilder, private key, and address
type BandProvider struct {
	cliCtx  context.CLIContext
	txBldr  authtypes.TxBuilder
	addr    sdk.AccAddress
	privKey crypto.PrivKey
}

var (
	cdc      = codecstd.MakeCodec(app.ModuleBasics)
	appCodec = codecstd.NewAppCodec(cdc)
)

// NewBandProvider creates new BandProvider create new cliCtx and txBldr
func NewBandProvider(nodeURI string, privKey crypto.PrivKey, chainID string) (BandProvider, error) {
	// TODO: This is a mess from Cosmos
	authclient.Codec = appCodec

	addr := sdk.AccAddress(privKey.PubKey().Address())
	cliCtx := NewCLIContext(nodeURI, addr).WithCodec(cdc)
	num, _, err := authtypes.NewAccountRetriever(appCodec, cliCtx).GetAccountNumberSequence(addr)
	if err != nil {
		return BandProvider{}, err
	}

	return BandProvider{
		cliCtx:  cliCtx,
		txBldr:  NewTxBuilder(authclient.GetTxEncoder(cdc), chainID).WithAccountNumber(num),
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
	return authclient.QueryTx(provider.cliCtx, hashHexStr)
}

func (provider *BandProvider) GetSequenceFromChain() (uint64, error) {
	_, seq, err := authtypes.NewAccountRetriever(appCodec, provider.cliCtx).GetAccountNumberSequence(provider.addr)
	return seq, err
}
