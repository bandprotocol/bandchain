package bandlib

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

func NewCLIContext(nodeURI string, fromAddress sdk.AccAddress) context.CLIContext {
	rpc, err := rpchttp.New(nodeURI, "/websocket")
	if err != nil {
		panic(err)
	}

	return context.CLIContext{
		Client:        rpc,
		Output:        os.Stdout,
		NodeURI:       nodeURI,
		OutputFormat:  "json",
		Height:        0,
		TrustNode:     true,
		UseLedger:     false,
		BroadcastMode: "sync",
		Simulate:      false,
		GenerateOnly:  false,
		FromAddress:   fromAddress,
		// FromName:      from,
		Indent:      true,
		SkipConfirm: true,
	}
}

func NewTxBuilder(txEncoder sdk.TxEncoder, chainID string) authtypes.TxBuilder {
	fee, _ := sdk.ParseCoins("")
	gasPrices, _ := sdk.ParseDecCoins("")
	// TODO: Remove hard code gas limit and gas adjustment
	return authtypes.NewTxBuilder(txEncoder, 0, 0, 20000000, 1, false, chainID, "", fee, gasPrices)
}

func completeAndBroadcastTxCLI(
	cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder,
	msgs []sdk.Msg,
	privKey crypto.PrivKey,
) (sdk.TxResponse, error) {
	txBldr, err := authclient.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	if txBldr.Gas() == 0 {
		txBldr, err = authclient.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return sdk.TxResponse{}, err
		}
	}

	// build and sign the transaction
	signMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	sigBytes, err := privKey.Sign(signMsg.Bytes())
	if err != nil {
		return sdk.TxResponse{}, err
	}
	sig := authtypes.StdSignature{
		PubKey:    privKey.PubKey().Bytes(),
		Signature: sigBytes,
	}

	txBytes, err := txBldr.TxEncoder()(
		authtypes.NewStdTx(signMsg.Msgs, signMsg.Fee, []authtypes.StdSignature{sig}, signMsg.Memo),
	)

	if err != nil {
		return sdk.TxResponse{}, err
	}
	// broadcast to a Tendermint node
	return cliCtx.BroadcastTx(txBytes)
}
