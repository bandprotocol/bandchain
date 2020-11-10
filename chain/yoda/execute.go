package yoda

import (
	"fmt"
	"strings"
	"time"

	sdkCtx "github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

var (
	cdc = app.MakeCodec()
)

func signAndBroadcast(
	c *Context, key keys.Info, msgs []sdk.Msg, gasLimit uint64, memo string,
) (string, error) {
	cliCtx := sdkCtx.CLIContext{Client: c.client, TrustNode: true, Codec: cdc}
	acc, err := auth.NewAccountRetriever(cliCtx).GetAccount(key.GetAddress())
	if err != nil {
		return "", fmt.Errorf("Failed to retreive account with error: %s", err.Error())
	}

	txBldr := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc), acc.GetAccountNumber(), acc.GetSequence(),
		gasLimit, 1, false, cfg.ChainID, memo, sdk.NewCoins(), c.gasPrices,
	)
	// txBldr, err = authclient.EnrichWithGas(txBldr, cliCtx, []sdk.Msg{msg})
	// if err != nil {
	// 	l.Error(":exploding_head: Failed to enrich with gas with error: %s", err.Error())
	// 	return
	// }

	out, err := txBldr.WithKeybase(keybase).BuildAndSign(key.GetName(), ckeys.DefaultKeyPass, msgs)
	if err != nil {
		return "", fmt.Errorf("Failed to build tx with error: %s", err.Error())
	}

	res, err := cliCtx.BroadcastTxSync(out)
	if err != nil {
		return "", fmt.Errorf("Failed to broadcast tx with error: %s", err.Error())
	}
	return res.TxHash, nil
}

func SubmitReport(c *Context, l *Logger, keyIndex int64, reports []ReportMsgWithKey) {
	// Return key when done with SubmitReport whether successfully or not.
	defer func() {
		c.freeKeys <- keyIndex
	}()

	// Summarize execute version
	versionMap := make(map[string]bool)
	msgs := make([]sdk.Msg, len(reports))
	ids := make([]types.RequestID, len(reports))
	feeEstimations := make([]FeeEstimationData, len(reports))

	for i, report := range reports {
		if err := report.msg.ValidateBasic(); err != nil {
			l.Error(":exploding_head: Failed to validate basic with error: %s", err.Error())
			return
		}
		msgs[i] = report.msg
		ids[i] = report.msg.RequestID
		feeEstimations[i] = report.feeEstimationData
		for _, exec := range report.execVersion {
			versionMap[exec] = true
		}
	}
	l = l.With("rids", ids)

	versions := make([]string, 0, len(versionMap))
	for exec := range versionMap {
		versions = append(versions, exec)
	}
	memo := fmt.Sprintf("yoda:%s/exec:%s", version.Version, strings.Join(versions, ","))
	key := c.keys[keyIndex]
	cliCtx := sdkCtx.CLIContext{Client: c.client, TrustNode: true, Codec: cdc}
	gasLimit := estimateGas(*c, msgs, feeEstimations)
	// We want to resend transaction only if tx returns Out of gas error.
	for sendAttempt := uint64(1); sendAttempt <= c.maxTry; sendAttempt++ {
		var txHash string
		l.Info(":e-mail: Sending report transaction attempt: (%d/%d)", sendAttempt, c.maxTry)
		for broadcastTry := uint64(1); broadcastTry <= c.maxTry; broadcastTry++ {
			l.Info(":writing_hand: Try to sign and broadcast report transaction(%d/%d)", broadcastTry, c.maxTry)
			hash, err := signAndBroadcast(c, key, msgs, gasLimit, memo)
			if err != nil {
				// Use info level because this error can happen and retry process can solve this error.
				l.Info(":warning: %s", err.Error())
				time.Sleep(c.rpcPollInterval)
				continue
			}
			// Transaction passed CheckTx process and wait to include in block.
			txHash = hash
			break
		}
		if txHash == "" {
			l.Error(":exploding_head: Cannot try to broadcast more than %d try", c.maxTry)
			return
		}
		txFound := false
	FindTx:
		for start := time.Now(); time.Since(start) < c.broadcastTimeout; {
			time.Sleep(c.rpcPollInterval)
			txRes, err := utils.QueryTx(cliCtx, txHash)
			if err != nil {
				l.Debug(":warning: Failed to query tx with error: %s", err.Error())
				continue
			}
			switch txRes.Code {
			case 0:
				l.Info(":smiling_face_with_sunglasses: Successfully broadcast tx with hash: %s", txHash)
				return
			case sdkerrors.ErrOutOfGas.ABCICode():
				// Increase gas limit and try to broadcast again
				gasLimit = gasLimit * 110 / 100
				l.Info(":fuel_pump: Tx(%s) is out of gas and will be rebroadcasted with %d gas", txHash, gasLimit)
				txFound = true
				break FindTx
			default:
				l.Error(":exploding_head: Tx returned nonzero code %d with log %s, tx hash: %s", txRes.Code, txRes.RawLog, txRes.TxHash)
				return
			}
		}
		if !txFound {
			l.Error(":question_mark: Cannot get transaction response from hash: %s transaction might be included in the next few blocks or check your node's health.", txHash)
			return
		}
	}
	l.Error(":anxious_face_with_sweat: Cannot send reports with adjusted gas: %d", gasLimit)
	return
}

// GetExecutable fetches data source executable using the provided client.
func GetExecutable(c *Context, l *Logger, hash string) ([]byte, error) {
	resValue, err := c.fileCache.GetFile(hash)
	if err != nil {
		l.Debug(":magnifying_glass_tilted_left: Fetching data source hash: %s from bandchain querier", hash)
		res, err := c.client.ABCIQueryWithOptions(fmt.Sprintf("custom/%s/%s/%s", types.StoreKey, types.QueryData, hash), nil, rpcclient.ABCIQueryOptions{})
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
