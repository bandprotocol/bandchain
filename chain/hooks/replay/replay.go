package replay

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/segmentio/kafka-go"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
)

// Hook uses Kafka functionality to act as an event producer for all events in the blockchains.
type Hook struct {
	cdc       *codec.Codec
	txDecoder sdk.TxDecoder
	// Main Kafka writer instance.
	writer *kafka.Writer
	// Temporary variables that are reset on every block.
	accsInBlock map[string]bool  // The accounts that need balance update at the end of block.
	accsInTx    map[string]bool  // The accounts related to the current processing transaction.
	msgs        []common.Message // The list of all messages to publish for this block.

	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	supplyKeeper  supply.Keeper
	stakingKeeper staking.Keeper
	mintKeeper    mint.Keeper
	distrKeeper   distr.Keeper
	govKeeper     gov.Keeper
	oracleKeeper  oracle.Keeper
}

// NewHook creates an emitter hook instance that will be added in Band App.
func NewHook(
	cdc *codec.Codec, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, supplyKeeper supply.Keeper,
	stakingKeeper staking.Keeper, mintKeeper mint.Keeper, distrKeeper distr.Keeper, govKeeper gov.Keeper,
	oracleKeeper keeper.Keeper, kafkaURI string,
) *Hook {
	paths := strings.SplitN(kafkaURI, "@", 2)
	return &Hook{
		cdc:       cdc,
		txDecoder: auth.DefaultTxDecoder(cdc),
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      paths[1:],
			Topic:        paths[0],
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 1 * time.Millisecond,
			// Async:    true, // TODO: We may be able to enable async mode on replay
		}),
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		supplyKeeper:  supplyKeeper,
		stakingKeeper: stakingKeeper,
		mintKeeper:    mintKeeper,
		distrKeeper:   distrKeeper,
		govKeeper:     govKeeper,
		oracleKeeper:  oracleKeeper,
	}
}

// AddAccountsInBlock adds the given accounts to the list of accounts to update balances end-of-block.
func (h *Hook) AddAccountsInBlock(accs ...sdk.AccAddress) {
	for _, acc := range accs {
		h.accsInBlock[acc.String()] = true
	}
}

// AddAccountsInTx adds the given accounts to the list of accounts to track related account in transaction.
func (h *Hook) AddAccountsInTx(accs ...sdk.AccAddress) {
	for _, acc := range accs {
		h.accsInTx[acc.String()] = true
	}
}

// Write adds the given key-value pair to the list of messages to publish during Commit.
func (h *Hook) Write(key string, val common.JsDict) {
	h.msgs = append(h.msgs, common.Message{Key: key, Value: val})
}

// FlushMessages publishes all pending messages to Kafka. Blocks until completion.
func (h *Hook) FlushMessages() {
	kafkaMsgs := make([]kafka.Message, len(h.msgs))
	for idx, msg := range h.msgs {
		res, _ := json.Marshal(msg.Value) // Error must always be nil.
		kafkaMsgs[idx] = kafka.Message{Key: []byte(msg.Key), Value: res}
	}
	err := h.writer.WriteMessages(context.Background(), kafkaMsgs...)
	if err != nil {
		panic(err)
	}
}

// AfterInitChain specify actions need to do after chain initialization (app.Hook interface).
func (h *Hook) AfterInitChain(ctx sdk.Context, req abci.RequestInitChain, res abci.ResponseInitChain) {
}

// AfterBeginBlock specify actions need to do after begin block period (app.Hook interface).
func (h *Hook) AfterBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock) {
	h.accsInBlock = make(map[string]bool)
	h.accsInTx = make(map[string]bool)
	h.msgs = []common.Message{}
	h.Write("NEW_BLOCK", common.JsDict{
		"height":    req.Header.GetHeight(),
		"timestamp": ctx.BlockTime().UnixNano(),
		"proposer":  sdk.ConsAddress(req.Header.GetProposerAddress()).String(),
		"hash":      req.GetHash(),
		"inflation": h.mintKeeper.GetMinter(ctx).Inflation.String(),
		"supply":    h.supplyKeeper.GetSupply(ctx).GetTotal().String(),
	})
	for _, event := range res.Events {
		h.handleBeginBlockEndBlockEvent(ctx, event)
	}
}

// AfterDeliverTx specify actions need to do after transaction has been processed (app.Hook interface).
func (h *Hook) AfterDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx, res abci.ResponseDeliverTx) {
	if ctx.BlockHeight() == 0 {
		return
	}
	h.accsInTx = make(map[string]bool)
	tx, err := h.txDecoder(req.Tx)
	if err != nil {
		return
	}
	stdTx, ok := tx.(auth.StdTx)
	if !ok {
		return
	}
	txHash := tmhash.Sum(req.Tx)
	var errMsg *string
	if !res.IsOK() {
		errMsg = &res.Log
	}
	txDict := common.JsDict{
		"hash":         txHash,
		"block_height": ctx.BlockHeight(),
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
	h.Write("NEW_TRANSACTION", txDict)
	logs, _ := sdk.ParseABCILogs(res.Log) // Error must always be nil if res.IsOK is true.
	messages := []map[string]interface{}{}
	for idx, msg := range tx.GetMsgs() {
		var extra = make(common.JsDict)
		if res.IsOK() {
			h.handleMsg(ctx, txHash, msg, logs[idx], extra)
		}
		messages = append(messages, common.JsDict{
			"msg":   msg,
			"type":  msg.Type(),
			"extra": extra,
		})
	}
	h.AddAccountsInTx(stdTx.GetSigners()...)
	relatedAccounts := make([]sdk.AccAddress, 0, len(h.accsInBlock))
	for accStr := range h.accsInTx {
		acc, _ := sdk.AccAddressFromBech32(accStr)
		relatedAccounts = append(relatedAccounts, acc)
	}
	common.EmitSetRelatedTransaction(h, txHash, relatedAccounts)
	h.AddAccountsInBlock(relatedAccounts...)
	txDict["messages"] = messages
}

// AfterEndBlock specify actions need to do after end block period (app.Hook interface).
func (h *Hook) AfterEndBlock(ctx sdk.Context, req abci.RequestEndBlock, res abci.ResponseEndBlock) {
	for _, event := range res.Events {
		h.handleBeginBlockEndBlockEvent(ctx, event)
	}
	h.Write("COMMIT", common.JsDict{"height": req.Height})
}

// ApplyQuery catch the custom query that matches specific paths (app.Hook interface).
func (h *Hook) ApplyQuery(req abci.RequestQuery) (res abci.ResponseQuery, stop bool) {
	return abci.ResponseQuery{}, false
}

// BeforeCommit specify actions need to do before commit block (app.Hook interface).
func (h *Hook) BeforeCommit() {
	h.FlushMessages()
}
