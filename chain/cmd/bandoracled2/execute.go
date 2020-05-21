package main

import (
	"fmt"

	sdkCtx "github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/bandprotocol/bandchain/chain/app"
	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

var (
	cdc      = codecstd.MakeCodec(app.ModuleBasics)
	appCodec = codecstd.NewAppCodec(cdc)
)

func init() {
	authclient.Codec = appCodec
}

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

	cliCtx := sdkCtx.CLIContext{Client: c.client}
	acc, err := auth.NewAccountRetriever(appCodec, cliCtx).GetAccount(key.GetAddress())
	if err != nil {
		l.Error(":exploding_head: Failed to retreive account with error: %s", err.Error())
		return
	}

	out, err := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc), acc.GetAccountNumber(), acc.GetSequence(),
		1000000, 1, false, cfg.ChainID, "", sdk.NewCoins(), c.gasPrices,
	).WithKeybase(keybase).BuildAndSign(key.GetName(), ckeys.DefaultKeyPass, []sdk.Msg{msg})
	if err != nil {
		l.Error(":exploding_head: Failed to build tx with error: %s", err.Error())
		return
	}

	res, err := cliCtx.BroadcastTxCommit(out)
	if err != nil {
		l.Error(":exploding_head: Failed to broadcast tx with error: %s", err.Error())
		return
	}
	if res.Code != 0 {
		l.Error(":exploding_head: Tx returned nonzero code %d with log %s", res.Code, res.RawLog)
		return
	}
	l.Info(":smiling_face_with_sunglasses: Successfully broadcast tx with hash: %s", res.TxHash)
}

// GetExecutable fetches data source executable using the provided client.
func GetExecutable(c *Context, l *Logger, id int) ([]byte, error) {
	l.Debug(":magnifying_glass_tilted_left: Fetching data source #%d from the remote node", id)
	res, _, err := sdkCtx.CLIContext{Client: c.client}.Query(
		fmt.Sprintf("custom/oracle/%s/%d", otypes.QueryDataSourceByID, id),
	)
	if err != nil {
		return nil, err
	}

	var dataSource otypes.DataSourceQuerierInfo
	err = cdc.UnmarshalJSON(res, &dataSource)
	if err != nil {
		return nil, err
	}
	l.Debug(":balloon: Received data source #%d content: %q", id, dataSource.Executable[:32])
	return dataSource.Executable, nil
}
