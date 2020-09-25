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
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/segmentio/kafka-go"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/emitter/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
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
	accsInBlock    map[string]bool  // The accounts that need balance update at the end of block.
	accsInTx       map[string]bool  // The accounts related to the current processing transaction.
	msgs           []common.Message // The list of all messages to publish for this block.
	emitStartState bool             // If emitStartState is true will emit all non historical state to Kafka
}

// NewBandAppWithEmitter creates a new App instance.
func NewBandAppWithEmitter(
	kafkaURI string, logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string,
	disableFeelessReports bool, enableFastSync bool, baseAppOptions ...func(*bam.BaseApp),
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
		emitStartState: enableFastSync,
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

// InitChain calls into the underlying InitChain and emits relevant events to Kafka.
func (app *App) InitChain(req abci.RequestInitChain) abci.ResponseInitChain {
	res := app.BandApp.InitChain(req)
	var genesisState bandapp.GenesisState
	app.Codec().MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	// Auth module
	var genaccountsState auth.GenesisState
	auth.ModuleCdc.MustUnmarshalJSON(genesisState[auth.ModuleName], &genaccountsState)
	for _, account := range genaccountsState.Accounts {
		app.Write("SET_ACCOUNT", common.JsDict{
			"address": account.GetAddress(),
			"balance": app.BankKeeper.GetCoins(app.DeliverContext, account.GetAddress()).String(),
		})
	}
	// GenUtil module for create validator genesis transactions.
	var genutilState genutil.GenesisState
	app.Codec().MustUnmarshalJSON(genesisState[genutil.ModuleName], &genutilState)
	for _, genTx := range genutilState.GenTxs {
		var tx auth.StdTx
		app.Codec().MustUnmarshalJSON(genTx, &tx)
		for _, msg := range tx.Msgs {
			if msg, ok := msg.(staking.MsgCreateValidator); ok {
				app.handleMsgCreateValidator(nil, msg, nil, nil)
			}
		}
	}

	// Staking module
	var stakingState staking.GenesisState
	app.Codec().MustUnmarshalJSON(genesisState[staking.ModuleName], &stakingState)
	for _, val := range stakingState.Validators {
		app.emitSetValidator(val.OperatorAddress)
	}

	for _, del := range stakingState.Delegations {
		app.emitDelegation(del.ValidatorAddress, del.DelegatorAddress)
	}
	for _, unbonding := range stakingState.UnbondingDelegations {
		for _, entry := range unbonding.Entries {
			common.EmitNewUnbondingDelegation(
				app,
				unbonding.DelegatorAddress,
				unbonding.ValidatorAddress,
				entry.CompletionTime.UnixNano(),
				entry.Balance,
			)

		}
	}
	for _, redelegate := range stakingState.Redelegations {
		for _, entry := range redelegate.Entries {
			common.EmitNewRedelegation(app,
				redelegate.DelegatorAddress,
				redelegate.ValidatorSrcAddress,
				redelegate.ValidatorDstAddress,
				entry.CompletionTime.UnixNano(),
				entry.InitialBalance,
			)
		}
	}

	// Oracle module
	var oracleState oracle.GenesisState
	app.Codec().MustUnmarshalJSON(genesisState[oracle.ModuleName], &oracleState)
	for idx, ds := range oracleState.DataSources {
		id := types.DataSourceID(idx + 1)
		app.emitSetDataSource(id, ds, nil)
		common.EmitNewDataSourceRequest(app, id)
	}
	for idx, os := range oracleState.OracleScripts {
		id := types.OracleScriptID(idx + 1)
		app.emitSetOracleScript(id, os, nil)
		common.EmitNewOracleScriptRequest(app, id)
	}
	app.FlushMessages()
	return res
}

func (app *App) emitNonHistoricalState() {
	app.emitAccountModule()
	app.emitStakingModule()
	app.emitOracleModule()
	common.EmitCommit(app, -1)
	app.FlushMessages()
	app.msgs = []common.Message{}
}

// BeginBlock calls into the underlying BeginBlock and emits relevant events to Kafka.
func (app *App) BeginBlock(req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	res := app.BandApp.BeginBlock(req)
	app.accsInBlock = make(map[string]bool)
	app.accsInTx = make(map[string]bool)
	app.msgs = []common.Message{}
	if app.emitStartState {
		app.emitStartState = false
		app.emitNonHistoricalState()
	} else {
		for _, val := range req.GetLastCommitInfo().Votes {
			validator := app.StakingKeeper.ValidatorByConsAddr(app.DeliverContext, val.GetValidator().Address)
			common.EmitNewValidatorVote(
				app,
				validator.GetConsAddr().String(),
				req.Header.GetHeight()-1,
				val.GetSignedLastBlock(),
			)
			app.emitUpdateValidatorRewardAndAccumulatedCommission(validator.GetOperator())
		}
	}
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
	// Update balances of all affected accounts on this block.
	// Index 0 is message NEW_BLOCK, we insert SET_ACCOUNT messages right after it.
	modifiedMsgs := []common.Message{app.msgs[0]}
	for accStr, _ := range app.accsInBlock {
		acc, _ := sdk.AccAddressFromBech32(accStr)
		modifiedMsgs = append(modifiedMsgs, common.Message{
			Key: "SET_ACCOUNT",
			Value: common.JsDict{
				"address": acc,
				"balance": app.BankKeeper.GetCoins(app.DeliverContext, acc).String(),
			}})
	}
	app.msgs = append(modifiedMsgs, app.msgs[1:]...)
	common.EmitCommit(app, req.Height)
	return res
}

// Commit makes sure all Kafka messages are broadcasted and then calls into the underlying Commit.
func (app *App) Commit() (res abci.ResponseCommit) {
	app.FlushMessages()
	return app.BandApp.Commit()
}
