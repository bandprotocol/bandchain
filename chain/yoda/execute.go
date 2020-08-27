package yoda

import (
	"fmt"
	"strings"
	"time"

	sdkCtx "github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// TODO: Improve precision of equation.
// const (
// 	rawReportMultiplier   = uint64(10)
// 	dataSizeMultiplier    = uint64(10)
// 	msgReportDataConstant = uint64(10)
// 	txSizeConstant        = uint64(10)
// 	baseTransaction       = uint64(30000)
// )

// Constant used to estimate gas price of reports transaction.
const (
	dataSizeMultiplier    = uint64(75)
	msgReportDataConstant = uint64(16000)
	txSizeConstant        = uint64(10) // Using DefaultTxSizeCostPerByte
	baseTransaction       = uint64(40000)
	pendingRequests       = uint64(4000)
)

func estimatedReportsGas(msgs []sdk.Msg) uint64 {
	est := baseTransaction
	txSize := uint64(0)
	for _, msg := range msgs {
		msg, ok := msg.(types.MsgReportData)
		if !ok {
			panic("Don't support non-report data message")
		}
		calldataSize := uint64(0)
		for _, c := range msg.RawReports {
			calldataSize += uint64(len(c.Data))
		}
		est += dataSizeMultiplier*calldataSize + msgReportDataConstant
		txSize += uint64(len(msg.GetSignBytes()))
	}
	return est + txSize*txSizeConstant + pendingRequests
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

	for i, report := range reports {
		if err := report.msg.ValidateBasic(); err != nil {
			l.Error(":exploding_head: Failed to validate basic with error: %s", err.Error())
			return
		}
		msgs[i] = report.msg
		ids[i] = report.msg.RequestID
		for _, exec := range report.execVersion {
			versionMap[exec] = true
		}
	}
	l = l.With("rids", ids)

	versions := make([]string, 0, len(versionMap))
	for exec := range versionMap {
		versions = append(versions, exec)
	}
	executorVersion := strings.Join(versions, ",")
	key := c.keys[keyIndex]
	cliCtx := sdkCtx.CLIContext{Client: c.client, TrustNode: true, Codec: cdc}
	txHash := ""
	for try := uint64(1); try <= c.maxTry; try++ {
		l.Info(":e-mail: Try to broadcast report transaction(%d/%d)", try, c.maxTry)
		acc, err := auth.NewAccountRetriever(cliCtx).GetAccount(key.GetAddress())
		if err != nil {
			l.Info(":warning: Failed to retreive account with error: %s", err.Error())
			time.Sleep(c.rpcPollInterval)
			continue
		}

		txBldr := auth.NewTxBuilder(
			auth.DefaultTxEncoder(cdc), acc.GetAccountNumber(), acc.GetSequence(),
			estimatedReportsGas(msgs), 1, false, cfg.ChainID, fmt.Sprintf("yoda:%s/exec:%s", version.Version, executorVersion), sdk.NewCoins(), c.gasPrices,
		)
		// txBldr, err = authclient.EnrichWithGas(txBldr, cliCtx, []sdk.Msg{msg})
		// if err != nil {
		// 	l.Error(":exploding_head: Failed to enrich with gas with error: %s", err.Error())
		// 	return
		// }

		out, err := txBldr.WithKeybase(keybase).BuildAndSign(key.GetName(), ckeys.DefaultKeyPass, msgs)
		if err != nil {
			l.Info(":warning: Failed to build tx with error: %s", err.Error())
			time.Sleep(c.rpcPollInterval)
			continue
		}
		res, err := cliCtx.BroadcastTxSync(out)
		if err == nil {
			txHash = res.TxHash
			break
		}
		l.Info(":warning: Failed to broadcast tx with error: %s", err.Error())
		time.Sleep(c.rpcPollInterval)
	}
	if txHash == "" {
		l.Error(":exploding_head: Cannot try to broadcast more than %d try", c.maxTry)
		return
	}
	for start := time.Now(); time.Since(start) < c.broadcastTimeout; {
		time.Sleep(c.rpcPollInterval)
		txRes, err := utils.QueryTx(cliCtx, txHash)
		if err != nil {
			l.Debug(":warning: Failed to query tx with error: %s", err.Error())
			continue
		}
		if txRes.Code != 0 {
			l.Error(":exploding_head: Tx returned nonzero code %d with log %s, tx hash: %s", txRes.Code, txRes.RawLog, txRes.TxHash)
			return
		}
		l.Info(":smiling_face_with_sunglasses: Successfully broadcast tx with hash: %s", txHash)
		return
	}
	l.Info(":question_mark: Cannot get transaction response from hash: %s transaction might be included in the next few blocks or check your node's health.", txHash)
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
