package main

import (
	"fmt"

	sdkCtx "github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

var (
	cdc      = codecstd.MakeCodec(app.ModuleBasics)
	appCodec = codecstd.NewAppCodec(cdc)
)

func init() {
	authclient.Codec = appCodec
}

func BroadCastMsgs(client *rpchttp.HTTP, key keyring.Info, msgs []sdk.Msg) {
	// TODO: Make this a queue. Make it better.
	cliCtx := sdkCtx.CLIContext{Client: client}
	acc, err := auth.NewAccountRetriever(appCodec, cliCtx).GetAccount(key.GetAddress())
	if err != nil {
		fmt.Println("ERR1", err)
		return
	}

	out, err := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc), acc.GetAccountNumber(), acc.GetSequence(),
		1000000, 1, false, chainID, "", sdk.NewCoins(), sdk.NewDecCoins(),
	).WithKeybase(keybase).BuildAndSign(key.GetName(), ckeys.DefaultKeyPass, msgs)
	if err != nil {
		fmt.Println("ERR2", err)
		return
	}

	res, err := cliCtx.BroadcastTxCommit(out)
	if err != nil {
		fmt.Println("ERR3", err)
		return
	}

	fmt.Println("EZ", res)
}

// GetExecutable fetches data source executable using the provided client.
func GetExecutable(client *rpchttp.HTTP, id int) ([]byte, error) {
	logger.Debug("⛏ Fetching data source #%d from the remote node", id)
	res, _, err := sdkCtx.CLIContext{Client: client}.Query(
		fmt.Sprintf("custom/oracle/%s/%d", oracle.QueryDataSourceByID, id),
	)
	if err != nil {
		return nil, err
	}

	var dataSource oracle.DataSourceQuerierInfo
	err = cdc.UnmarshalJSON(res, &dataSource)
	if err != nil {
		return nil, err
	}
	logger.Debug("👀 Received data source #%d content: 0x%X...", id, dataSource.Executable[:32])
	return dataSource.Executable, nil
}
