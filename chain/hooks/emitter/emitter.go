package emitter

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
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/segmentio/kafka-go"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// Emitter hook uses Kafka functionality to act as an event producer for all events in the blockchains.
type EmitterHook struct {
	cdc       *codec.Codec
	txDecoder sdk.TxDecoder
	// Main Kafka writer instance.
	writer *kafka.Writer
	// Temporary variables that are reset on every block.
	accsInBlock    map[string]bool  // The accounts that need balance update at the end of block.
	accsInTx       map[string]bool  // The accounts related to the current processing transaction.
	msgs           []common.Message // The list of all messages to publish for this block.
	emitStartState bool             // If emitStartState is true will emit all non historical state to Kafka

	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	supplyKeeper  supply.Keeper
	stakingKeeper staking.Keeper
	mintKeeper    mint.Keeper
	distrKeeper   distr.Keeper
	govKeeper     gov.Keeper
	oracleKeeper  oracle.Keeper
}

func NewEmitterHook(
	cdc *codec.Codec, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, supplyKeeper supply.Keeper,
	stakingKeeper staking.Keeper, mintKeeper mint.Keeper, distrKeeper distr.Keeper, govKeeper gov.Keeper,
	oracleKeeper keeper.Keeper, kafkaURI string, emitStartState bool,
) *EmitterHook {
	paths := strings.SplitN(kafkaURI, "@", 2)
	return &EmitterHook{
		cdc:       cdc,
		txDecoder: auth.DefaultTxDecoder(cdc),
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      paths[1:],
			Topic:        paths[0],
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 1 * time.Millisecond,
			// Async:    true, // TODO: We may be able to enable async mode on replay
		}),
		emitStartState: emitStartState,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		supplyKeeper:   supplyKeeper,
		stakingKeeper:  stakingKeeper,
		mintKeeper:     mintKeeper,
		distrKeeper:    distrKeeper,
		govKeeper:      govKeeper,
		oracleKeeper:   oracleKeeper,
	}
}

// AddAccountsInBlock adds the given accounts to the list of accounts to update balances end-of-block.
func (h *EmitterHook) AddAccountsInBlock(accs ...sdk.AccAddress) {
	for _, acc := range accs {
		h.accsInBlock[acc.String()] = true
	}
}

// AddAccountsInTx adds the given accounts to the list of accounts to track related account in transaction.
func (h *EmitterHook) AddAccountsInTx(accs ...sdk.AccAddress) {
	for _, acc := range accs {
		h.accsInTx[acc.String()] = true
	}
}

// Write adds the given key-value pair to the list of messages to publish during Commit.
func (h *EmitterHook) Write(key string, val common.JsDict) {
	h.msgs = append(h.msgs, common.Message{Key: key, Value: val})
}

// FlushMessages publishes all pending messages to Kafka. Blocks until completion.
func (h *EmitterHook) FlushMessages() {
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

func (h *EmitterHook) AfterInitChain(ctx sdk.Context, req abci.RequestInitChain, res abci.ResponseInitChain) {
	var genesisState bandapp.GenesisState
	h.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	// Auth module
	var genaccountsState auth.GenesisState
	auth.ModuleCdc.MustUnmarshalJSON(genesisState[auth.ModuleName], &genaccountsState)
	for _, account := range genaccountsState.Accounts {
		h.Write("SET_ACCOUNT", common.JsDict{
			"address": account.GetAddress(),
			"balance": h.bankKeeper.GetCoins(ctx, account.GetAddress()).String(),
		})
	}
	// GenUtil module for create validator genesis transactions.
	var genutilState genutil.GenesisState
	h.cdc.MustUnmarshalJSON(genesisState[genutil.ModuleName], &genutilState)
	for _, genTx := range genutilState.GenTxs {
		var tx auth.StdTx
		h.cdc.MustUnmarshalJSON(genTx, &tx)
		for _, msg := range tx.Msgs {
			if msg, ok := msg.(staking.MsgCreateValidator); ok {
				h.handleMsgCreateValidator(ctx, nil, msg, nil, nil)
			}
		}
	}

	// Staking module
	var stakingState staking.GenesisState
	h.cdc.MustUnmarshalJSON(genesisState[staking.ModuleName], &stakingState)
	for _, val := range stakingState.Validators {
		h.emitSetValidator(ctx, val.OperatorAddress)
	}

	for _, del := range stakingState.Delegations {
		h.emitDelegation(ctx, del.ValidatorAddress, del.DelegatorAddress)
	}

	for _, unbonding := range stakingState.UnbondingDelegations {
		for _, entry := range unbonding.Entries {
			h.Write("NEW_UNBONDING_DELEGATION", common.JsDict{
				"delegator_address": unbonding.DelegatorAddress,
				"operator_address":  unbonding.ValidatorAddress,
				"completion_time":   entry.CompletionTime.UnixNano(),
				"amount":            entry.Balance,
			})
		}
	}

	for _, redelegate := range stakingState.Redelegations {
		for _, entry := range redelegate.Entries {
			h.Write("NEW_REDELEGATION", common.JsDict{
				"delegator_address":    redelegate.DelegatorAddress,
				"operator_src_address": redelegate.ValidatorSrcAddress,
				"operator_dst_address": redelegate.ValidatorDstAddress,
				"completion_time":      entry.CompletionTime.UnixNano(),
				"amount":               entry.InitialBalance,
			})
		}
	}

	// Gov module
	var govState gov.GenesisState
	h.cdc.MustUnmarshalJSON(genesisState[gov.ModuleName], &govState)
	for _, proposal := range govState.Proposals {
		h.Write("NEW_PROPOSAL", common.JsDict{
			"id":               proposal.ProposalID,
			"proposer":         nil,
			"type":             proposal.ProposalType(),
			"title":            proposal.Content.GetTitle(),
			"description":      proposal.Content.GetDescription(),
			"proposal_route":   proposal.Content.ProposalRoute(),
			"status":           int(proposal.Status),
			"submit_time":      proposal.SubmitTime.UnixNano(),
			"deposit_end_time": proposal.DepositEndTime.UnixNano(),
			"total_deposit":    proposal.TotalDeposit.String(),
			"voting_time":      proposal.VotingStartTime.UnixNano(),
			"voting_end_time":  proposal.VotingEndTime.UnixNano(),
		})
	}
	for _, deposit := range govState.Deposits {
		h.Write("SET_DEPOSIT", common.JsDict{
			"proposal_id": deposit.ProposalID,
			"depositor":   deposit.Depositor,
			"amount":      deposit.Amount.String(),
			"tx_hash":     nil,
		})
	}
	for _, vote := range govState.Votes {
		h.Write("SET_VOTE", common.JsDict{
			"proposal_id": vote.ProposalID,
			"voter":       vote.Voter,
			"answer":      int(vote.Option),
			"tx_hash":     nil,
		})
	}

	// Oracle module
	var oracleState oracle.GenesisState
	h.cdc.MustUnmarshalJSON(genesisState[oracle.ModuleName], &oracleState)
	for idx, ds := range oracleState.DataSources {
		h.emitSetDataSource(types.DataSourceID(idx+1), ds, nil)
	}
	for idx, os := range oracleState.OracleScripts {
		h.emitSetOracleScript(types.OracleScriptID(idx+1), os, nil)
	}
	h.Write("COMMIT", common.JsDict{"height": 0})
	h.FlushMessages()
}

func (h *EmitterHook) emitNonHistoricalState(ctx sdk.Context) {
	h.emitAuthModule(ctx)
	h.emitStakingModule(ctx)
	h.emitOracleModule(ctx)
	h.Write("COMMIT", common.JsDict{"height": -1})
	h.FlushMessages()
	h.msgs = []common.Message{}
}

func (h *EmitterHook) AfterBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock) {
	h.accsInBlock = make(map[string]bool)
	h.accsInTx = make(map[string]bool)
	h.msgs = []common.Message{}
	if h.emitStartState {
		h.emitStartState = false
		h.emitNonHistoricalState(ctx)
	} else {
		{
			for _, val := range req.GetLastCommitInfo().Votes {
				validator := h.stakingKeeper.ValidatorByConsAddr(ctx, val.GetValidator().Address)
				h.Write("NEW_VALIDATOR_VOTE", common.JsDict{
					"consensus_address": validator.GetConsAddr().String(),
					"block_height":      req.Header.GetHeight() - 1,
					"voted":             val.GetSignedLastBlock(),
				})
				h.emitUpdateValidatorRewardAndAccumulatedCommission(ctx, validator.GetOperator())
			}
		}
	}
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

func (h *EmitterHook) AfterDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx, res abci.ResponseDeliverTx) {
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
	for accStr, _ := range h.accsInTx {
		acc, _ := sdk.AccAddressFromBech32(accStr)
		relatedAccounts = append(relatedAccounts, acc)
	}

	txDict["related_accounts"] = relatedAccounts
	h.AddAccountsInBlock(relatedAccounts...)
	txDict["messages"] = messages
}

func (h *EmitterHook) AfterEndBlock(ctx sdk.Context, req abci.RequestEndBlock, res abci.ResponseEndBlock) {
	for _, event := range res.Events {
		h.handleBeginBlockEndBlockEvent(ctx, event)
	}
	// Update balances of all affected accounts on this block.
	// Index 0 is message NEW_BLOCK, we insert SET_ACCOUNT messages right after it.
	modifiedMsgs := []common.Message{h.msgs[0]}
	for accStr, _ := range h.accsInBlock {
		acc, _ := sdk.AccAddressFromBech32(accStr)
		modifiedMsgs = append(modifiedMsgs, common.Message{
			Key: "SET_ACCOUNT",
			Value: common.JsDict{
				"address": acc,
				"balance": h.bankKeeper.GetCoins(ctx, acc).String(),
			}})
	}
	h.msgs = append(modifiedMsgs, h.msgs[1:]...)
	h.Write("COMMIT", common.JsDict{"height": req.Height})
}

func (h *EmitterHook) ApplyQuery(req abci.RequestQuery) (res abci.ResponseQuery, stop bool) {
	return abci.ResponseQuery{}, false
}

func (h *EmitterHook) BeforeCommit() {
	h.FlushMessages()
}
