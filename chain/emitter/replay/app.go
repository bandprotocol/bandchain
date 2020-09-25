package emitter

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/segmentio/kafka-go"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/emitter/common"
)

// App extends the standard Band Cosmos-SDK application with Kafka emitter
// functionality to act as an event producer for all events in the blockchains.
type App struct {
	*bandapp.BandApp
	// Decoder for unmarshaling []byte into sdk.Tx.
	txDecoder sdk.TxDecoder
	// Main Kafka writer instance.
	writer *kafka.Writer
	// Temporary variables that are reset on every block.
	accsInBlock map[string]bool  // The accounts that need balance update at the end of block.
	accsInTx    map[string]bool  // The accounts related to the current processing transaction.
	msgs        []common.Message // The list of all messages to publish for this block.
}

// NewBandAppWithEmitter creates a new App instance.
func NewBandAppWithEmitter(
	kafkaURI string, logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string,
	disableFeelessReports bool, baseAppOptions ...func(*bam.BaseApp),
) *App {
	app := bandapp.NewBandApp(
		logger, db, traceStore, loadLatest, invCheckPeriod, skipUpgradeHeights,
		home, disableFeelessReports, baseAppOptions...,
	)
	paths := strings.SplitN(kafkaURI, "@", 2)
	return &App{
		BandApp:   app,
		txDecoder: auth.DefaultTxDecoder(app.Codec()),
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      paths[1:],
			Topic:        paths[0],
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 1 * time.Millisecond,
			// Async:    true, // TODO: We may be able to enable async mode on replay
		}),
	}
}

// AddAccountsInBlock adds the given accounts to the list of accounts to update balances end-of-block.
func (app *App) AddAccountsInBlock(accs ...sdk.AccAddress) {
	for _, acc := range accs {
		app.accsInBlock[acc.String()] = true
	}
}

// AddAccountsInTx adds the given accounts to the list of accounts to track related account in transaction.
func (app *App) AddAccountsInTx(accs ...sdk.AccAddress) {
	for _, acc := range accs {
		app.accsInTx[acc.String()] = true
	}
}

// Write adds the given key-value pair to the list of messages to publish during Commit.
func (app *App) Write(key string, val common.JsDict) {
	app.msgs = append(app.msgs, common.Message{Key: key, Value: val})
}

// FlushMessages publishes all pending messages to Kafka. Blocks until completion.
func (app *App) FlushMessages() {
	kafkaMsgs := make([]kafka.Message, len(app.msgs))
	for idx, msg := range app.msgs {
		res, _ := json.Marshal(msg.Value) // Error must always be nil.
		kafkaMsgs[idx] = kafka.Message{Key: []byte(msg.Key), Value: res}
	}
	err := app.writer.WriteMessages(context.Background(), kafkaMsgs...)
	if err != nil {
		panic(err)
	}
}

// BeginBlock calls into the underlying BeginBlock and emits relevant events to Kafka.
func (app *App) BeginBlock(req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	res := app.BandApp.BeginBlock(req)
	app.accsInBlock = make(map[string]bool)
	app.accsInTx = make(map[string]bool)
	app.msgs = []common.Message{}
	common.EmitNewBlock(
		app,
		req.Header.GetHeight(),
		app.DeliverContext.BlockTime().UnixNano(),
		sdk.ConsAddress(req.Header.GetProposerAddress()).String(),
		req.GetHash(),
		app.MintKeeper.GetMinter(app.DeliverContext).Inflation.String(),
		app.SupplyKeeper.GetSupply(app.DeliverContext).GetTotal().String(),
	)
	for _, event := range res.Events {
		app.handleBeginBlockEndBlockEvent(event)
	}

	return res
}

// DeliverTx calls into the underlying DeliverTx and emits relevant events to Kafka.
func (app *App) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	res := app.BandApp.DeliverTx(req)
	app.accsInTx = make(map[string]bool)
	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return res
	}
	stdTx, ok := tx.(auth.StdTx)
	if !ok {
		return res
	}
	txHash := tmhash.Sum(req.Tx)
	var errMsg *string
	if !res.IsOK() {
		errMsg = &res.Log
	}
	txDict := common.JsDict{
		"hash":         txHash,
		"block_height": app.DeliverContext.BlockHeight(),
		"gas_used":     res.GasUsed,
		"gas_limit":    stdTx.Fee.Gas,
		"gas_fee":      stdTx.Fee.Amount.String(),
		"err_msg":      errMsg,
		"sender":       stdTx.GetSigners()[0].String(),
		"success":      res.IsOK(),
		"memo":         stdTx.Memo,
	}
	// NOTE: We add txDict to the list of pending Kafka messages here, but it will still be
	// mutated in the loop below as we know the messages won't get flushed until ABCI Commit.
	app.Write("NEW_TRANSACTION", txDict)
	logs, _ := sdk.ParseABCILogs(res.Log) // Error must always be nil if res.IsOK is true.
	messages := []map[string]interface{}{}
	for idx, msg := range tx.GetMsgs() {
		var extra = make(common.JsDict)
		if res.IsOK() {
			app.handleMsg(txHash, msg, logs[idx], extra)
		}
		messages = append(messages, common.JsDict{
			"msg":   msg,
			"type":  msg.Type(),
			"extra": extra,
		})
	}
	app.AddAccountsInTx(stdTx.GetSigners()...)
	relatedAccounts := make([]sdk.AccAddress, 0, len(app.accsInBlock))
	for accStr, _ := range app.accsInTx {
		acc, _ := sdk.AccAddressFromBech32(accStr)
		relatedAccounts = append(relatedAccounts, acc)
	}
	common.EmitSetRelatedTransaction(app, txHash, relatedAccounts)
	app.AddAccountsInBlock(relatedAccounts...)
	txDict["messages"] = messages
	return res
}

// EndBlock calls into the underlying EndBlock and emits relevant events to Kafka.
func (app *App) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	res := app.BandApp.EndBlock(req)
	for _, event := range res.Events {
		app.handleBeginBlockEndBlockEvent(event)
	}
	common.EmitCommit(app, req.Height)
	return res
}

// Commit makes sure all Kafka messages are broadcasted and then calls into the underlying Commit.
func (app *App) Commit() (res abci.ResponseCommit) {
	app.FlushMessages()
	return app.BandApp.Commit()
}
