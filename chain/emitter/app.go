package emitter

import (
	"context"
	"encoding/json"
	"io"
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
	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

type Message struct {
	Key   string
	Value JsDict
}

// App extends the standard Band Cosmos-SDK application with Kafka emitter
// functionality to act as an event producer for all events in the blockchains.
type App struct {
	*bandapp.BandApp
	// Main Kafka writer instance.
	writer *kafka.Writer
	// Temporary variables that are reset on every block.
	txIdx int              // The current transaction's index on the current block starting from 1.
	accs  []sdk.AccAddress // The accounts that need balance update at the end of block.
	msgs  []Message        // The list of all messages to publish for this block.
}

// NewBandAppWithEmitter creates a new App instance.
func NewBandAppWithEmitter(
	topic string, logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string,
	baseAppOptions ...func(*bam.BaseApp),
) *App {
	return &App{
		BandApp: bandapp.NewBandApp(
			logger, db, traceStore, loadLatest, invCheckPeriod, skipUpgradeHeights,
			home, baseAppOptions...,
		),
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{"localhost:9092"}, // TODO: Remove hardcode
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 1 * time.Millisecond,
			// Async:    true, // TODO: We may be able to enable async mode on replay
		}),
	}
}

// AddAccount adds the given account to the list of accounts to update balances end-of-block.
func (app *App) AddAccount(acc sdk.AccAddress) {
	app.accs = append(app.accs, acc)
}

// Write adds the given key-value pair to the list of messages to publish during Commit.
func (app *App) Write(key string, val JsDict) {
	app.msgs = append(app.msgs, Message{Key: key, Value: val})
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

// InitChain calls into the underlying InitChain and emits relevant events to Kafka.
func (app *App) InitChain(req abci.RequestInitChain) abci.ResponseInitChain {
	res := app.BandApp.InitChain(req)
	var genesisState bandapp.GenesisState
	app.Codec().MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	// Auth module
	var genaccountsState auth.GenesisState
	auth.ModuleCdc.MustUnmarshalJSON(genesisState[auth.ModuleName], &genaccountsState)
	for _, account := range genaccountsState.Accounts {
		app.Write("SET_ACCOUNT", JsDict{
			"address": account.GetAddress(),
			"balance": app.BankKeeper.GetCoins(app.DeliverContext, account.GetAddress()).String(),
		})
	}
	// Staking module
	// TODO
	// Oracle module
	var oracleState oracle.GenesisState
	app.Codec().MustUnmarshalJSON(genesisState[oracle.ModuleName], &oracleState)
	for idx, ds := range oracleState.DataSources {
		app.Write("NEW_DATA_SOURCE", JsDict{
			"id":          idx + 1,
			"name":        ds.Name,
			"description": ds.Description,
			"owner":       ds.Owner.String(),
			"executable":  app.OracleKeeper.GetFile(ds.Filename),
		})
	}
	for idx, os := range oracleState.OracleScripts {
		app.Write("NEW_ORACLE_SCRIPT", JsDict{
			"id":              idx + 1,
			"name":            os.Name,
			"description":     os.Description,
			"owner":           os.Owner.String(),
			"schema":          os.Schema,
			"codehash":        os.Filename,
			"source_code_url": os.SourceCodeURL,
		})
	}
	app.FlushMessages()
	return res
}

// BeginBlock calls into the underlying BeginBlock and emits relevant events to Kafka.
func (app *App) BeginBlock(req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	res := app.BandApp.BeginBlock(req)
	app.txIdx = 0
	app.accs = []sdk.AccAddress{}
	app.msgs = []Message{}
	app.Write("NEW_BLOCK", JsDict{
		"height":    req.Header.GetHeight(),
		"timestamp": app.DeliverContext.BlockTime().UnixNano(),
		"proposer":  req.Header.GetProposerAddress(),
		"hash":      req.GetHash(),
		"inflation": app.MintKeeper.GetMinter(app.DeliverContext).Inflation.String(),
		"supply":    app.SupplyKeeper.GetSupply(app.DeliverContext).GetTotal().String(),
	})
	// TODO: Handle begin block event
	return res
}

// DeliverTx calls into the underlying DeliverTx and emits relevant events to Kafka.
func (app *App) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	res := app.BandApp.DeliverTx(req)
	tx, err := app.TxDecoder(req.Tx)
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
	app.txIdx++
	txDict := JsDict{
		"hash":         txHash,
		"index":        app.txIdx,
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
		var extra = make(JsDict)
		if res.IsOK() {
			app.handleMsg(txHash, msg, logs[idx], extra)
		}
		messages = append(messages, JsDict{
			"msg":   msg,
			"type":  msg.Type(),
			"extra": extra,
		})
	}
	txDict["messages"] = messages
	app.AddAccount(stdTx.GetSigners()[0])
	return res
}

// EndBlock calls into the underlying EndBlock and emits relevant events to Kafka.
func (app *App) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	res := app.BandApp.EndBlock(req)
	// Update balances of all affected accounts on this block.
	accMap := make(map[string]bool)
	for _, acc := range app.accs {
		accStr := string(acc)
		if accMap[accStr] {
			continue
		}
		accMap[accStr] = true
		app.Write("SET_ACCOUNT", JsDict{
			"address": acc,
			"balance": app.BankKeeper.GetCoins(app.DeliverContext, acc).String(),
		})
	}

	for _, event := range res.Events {
		app.handleEndBlock(event)
	}

	app.Write("COMMIT", JsDict{"height": req.Height})
	return res
}

// Commit makes sure all Kafka messages are broadcasted and then calls into the underlying Commit.
func (app *App) Commit() (res abci.ResponseCommit) {
	app.FlushMessages()
	return app.BandApp.Commit()
}
