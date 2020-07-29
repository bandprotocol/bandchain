package yoda

import (
	"fmt"
	"time"

	sdkCtx "github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/bandprotocol/bandchain/chain/app"
	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

var (
	cdc              = app.MakeCodec()
	MaxTry           = 3
	SleepTime        = 1 * time.Second
	BroadCastTimeout = 7 * time.Second
)

func SubmitReport(c *Context, l *Logger, id otypes.RequestID, reps []otypes.RawReport) {
	key := <-c.keys
	defer func() {
		c.keys <- key
	}()

	msg := otypes.NewMsgReportData(otypes.RequestID(id), reps, c.validator, key.GetAddress())
	if err := msg.ValidateBasic(); err != nil {
		l.Error(":exploding_head: Failed to validate basic with error: %s", err.Error())
		return
	}
	cliCtx := sdkCtx.CLIContext{Client: c.client, TrustNode: true, Codec: cdc}
	acc, err := auth.NewAccountRetriever(cliCtx).GetAccount(key.GetAddress())
	if err != nil {
		l.Error(":exploding_head: Failed to retreive account with error: %s", err.Error())
		return
	}

	txBldr := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc), acc.GetAccountNumber(), acc.GetSequence(),
		200000, 1, false, cfg.ChainID, "", sdk.NewCoins(), c.gasPrices,
	)
	// txBldr, err = authclient.EnrichWithGas(txBldr, cliCtx, []sdk.Msg{msg})
	// if err != nil {
	// 	l.Error(":exploding_head: Failed to enrich with gas with error: %s", err.Error())
	// 	return
	// }
	out, err := txBldr.WithKeybase(keybase).BuildAndSign(key.GetName(), ckeys.DefaultKeyPass, []sdk.Msg{msg})
	if err != nil {
		l.Error(":exploding_head: Failed to build tx with error: %s", err.Error())
		return
	}

	for try := 1; try <= MaxTry; try++ {
		res, err := cliCtx.BroadcastTxSync(out)
		l.Info("Try to broadcast: %d/%d", try, MaxTry)
		if err != nil {
			l.Error(":exploding_head: Failed to broadcast tx with error: %s", err.Error())
			time.Sleep(SleepTime)
			continue
		}
		for start := time.Now(); time.Since(start) < BroadCastTimeout; {
			time.Sleep(SleepTime)
			txRes, err := utils.QueryTx(cliCtx, res.TxHash)
			if err != nil {
				l.Error(":exploding_head: Failed to query tx with error: %s", err.Error())
				continue
			}
			if txRes.Code != 0 {
				l.Error(":exploding_head: Tx returned nonzero code %d with log %s, tx hash: %s", txRes.Code, txRes.RawLog, txRes.TxHash)
				break
			}
			l.Info(":smiling_face_with_sunglasses: Successfully broadcast tx with hash: %s", res.TxHash)
			return
		}
		l.Error(":exploding_head: Failed to broadcast on %d try", try)
	}
	l.Error(":exploding_head: Failed to broadcast tx with max try = %d", MaxTry)
}

// GetExecutable fetches data source executable using the provided client.
func GetExecutable(c *Context, l *Logger, hash string) ([]byte, error) {
	resValue, err := c.fileCache.GetFile(hash)
	if err != nil {
		l.Debug(":magnifying_glass_tilted_left: Fetching data source hash: %s from bandchain querier", hash)
		res, err := c.client.ABCIQueryWithOptions(fmt.Sprintf("custom/%s/%s/%s", otypes.StoreKey, otypes.QueryData, hash), nil, rpcclient.ABCIQueryOptions{})
		if err != nil {
			l.Error(":exploding_head: Failed to get data source with error: %s", err.Error())
			return nil, err
		}
		resValue = res.Response.GetValue()
		c.fileCache.AddFile(resValue)
	} else {
		l.Debug(":card_file_box: Found data source hash: %s in cache file", hash)
	}

	l.Debug(":balloon: Received data source hash: %s content: %q", hash, resValue[:32])
	return resValue, nil
}
