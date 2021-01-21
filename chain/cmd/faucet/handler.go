package main

import (
	"fmt"
	"net/http"

	"github.com/GeoDB-Limited/odincore/chain/app"
	sdkctx "github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/gin-gonic/gin"
)

type Request struct {
	Address string `json:"address" binding:"required"`
}

type Response struct {
	TxHash string `json:"txHash"`
}

var (
	cdc = app.MakeCodec()
)

func handleRequest(gc *gin.Context, c *Context) {
	key := <-c.keys
	defer func() {
		c.keys <- key
	}()

	var req Request
	if err := gc.ShouldBindJSON(&req); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	to, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg := bank.NewMsgSend(key.GetAddress(), to, c.amount)
	if err := msg.ValidateBasic(); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cliCtx := sdkctx.CLIContext{Client: c.client}
	acc, err := auth.NewAccountRetriever(cliCtx).GetAccount(key.GetAddress())
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	txBldr := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc), acc.GetAccountNumber(), acc.GetSequence(),
		200000, 1, false, cfg.ChainID, "", sdk.NewCoins(), c.gasPrices,
	)
	out, err := txBldr.WithKeybase(keybase).BuildAndSign(key.GetName(), ckeys.DefaultKeyPass, []sdk.Msg{msg})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, err := cliCtx.BroadcastTxCommit(out)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.Code != 0 {
		gc.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf(":exploding_head: Tx returned nonzero code %d with log %s, tx hash: %s",
				res.Code, res.RawLog, res.TxHash,
			)})
		return
	}
	gc.JSON(200, Response{
		TxHash: res.TxHash,
	})
}
