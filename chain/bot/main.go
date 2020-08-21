package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	sdkCtx "github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	httpclient "github.com/tendermint/tendermint/rpc/client/http"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

var (
	cdc = app.MakeCodec()
)

func main() {
	appConfig := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(appConfig)
	appConfig.Seal()

	keybase, _ := keys.NewKeyring("bandchain", "test", os.ExpandEnv("${HOME}/.bandcli"), nil)
	keys, _ := keybase.List()
	calldata, _ := hex.DecodeString("0000000462616e64")
	msg := types.NewMsgRequestData(types.OracleScriptID(3), calldata, 16, 12, "Test", keys[0].GetAddress())
	client, _ := httpclient.New("http://rpc-guanyu-testnet2.bandchain.org:26657", "/websocket")
	cliCtx := sdkCtx.CLIContext{Client: client, TrustNode: true, Codec: cdc}
	acc, _ := auth.NewAccountRetriever(cliCtx).GetAccount(keys[0].GetAddress())
	accNum := acc.GetAccountNumber()
	seq := acc.GetSequence()
	for i := 0; i < 10; i++ {
		txBldr := auth.NewTxBuilder(
			auth.DefaultTxEncoder(cdc), accNum, seq,
			5000000, 1, false, "band-guanyu-testnet2", "", sdk.NewCoins(), sdk.NewDecCoins(),
		)
		out, _ := txBldr.WithKeybase(keybase).BuildAndSign(keys[0].GetName(), ckeys.DefaultKeyPass, []sdk.Msg{msg, msg, msg, msg, msg})
		res, _ := cliCtx.BroadcastTxSync(out)
		for {
			time.Sleep(1)
			txRes, err := utils.QueryTx(cliCtx, res.TxHash)
			if err != nil {
				continue
			}
			if txRes.Code != 0 {
				log.Fatalf(":exploding_head: Tx returned nonzero code %d with log %s, tx hash: %s", txRes.Code, txRes.RawLog, txRes.TxHash)
			}
			fmt.Printf(":smiling_face_with_sunglasses: Successfully broadcast tx with hash: %s\n", res.TxHash)
			seq++
			break
		}
	}
}
