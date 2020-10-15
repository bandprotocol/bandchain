package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func parseEvents(events sdk.StringEvents) common.EvMap {
	evMap := make(common.EvMap)
	for _, event := range events {
		for _, kv := range event.Attributes {
			key := event.Type + "." + kv.Key
			evMap[key] = append(evMap[key], kv.Value)
		}
	}
	return evMap
}

// handleMsg handles the given message by publishing relevant events and populates accounts
// that need balance update in 'h.accs'. Also fills in extra info for this message.
func (h *EmitterHook) handleMsg(ctx sdk.Context, txHash []byte, msg sdk.Msg, log sdk.ABCIMessageLog, extra common.JsDict) {
	evMap := parseEvents(log.Events)
	switch msg := msg.(type) {
	case oracle.MsgRequestData:
		h.handleMsgRequestData(ctx, txHash, msg, evMap, extra)
	case oracle.MsgReportData:
		h.handleMsgReportData(ctx, txHash, msg, evMap, extra)
	case oracle.MsgCreateDataSource:
		h.handleMsgCreateDataSource(ctx, txHash, msg, evMap, extra)
	case oracle.MsgCreateOracleScript:
		h.handleMsgCreateOracleScript(ctx, txHash, msg, evMap, extra)
	case oracle.MsgEditDataSource:
		h.handleMsgEditDataSource(ctx, txHash, msg, evMap, extra)
	case oracle.MsgEditOracleScript:
		h.handleMsgEditOracleScript(ctx, txHash, msg, evMap, extra)
	case oracle.MsgAddReporter:
		h.handleMsgAddReporter(ctx, txHash, msg, evMap, extra)
	case oracle.MsgRemoveReporter:
		h.handleMsgRemoveReporter(ctx, txHash, msg, evMap, extra)
	case oracle.MsgActivate:
		h.handleMsgActivate(ctx, txHash, msg, evMap, extra)
	case staking.MsgCreateValidator:
		h.handleMsgCreateValidator(ctx, txHash, msg, evMap, extra)
	case staking.MsgEditValidator:
		h.handleMsgEditValidator(ctx, txHash, msg, evMap, extra)
	case staking.MsgDelegate:
		h.handleMsgDelegate(ctx, txHash, msg, evMap, extra)
	case staking.MsgUndelegate:
		h.handleMsgUndelegate(ctx, txHash, msg, evMap, extra)
	case staking.MsgBeginRedelegate:
		h.handleMsgBeginRedelegate(ctx, txHash, msg, evMap, extra)
	case bank.MsgSend:
		h.handleMsgSend(ctx, txHash, msg, evMap, extra)
	case bank.MsgMultiSend:
		h.handleMsgMultiSend(ctx, txHash, msg, evMap, extra)
	case dist.MsgWithdrawDelegatorReward:
		h.handleMsgWithdrawDelegatorReward(ctx, txHash, msg, evMap, extra)
	case dist.MsgSetWithdrawAddress:
		h.handleMsgSetWithdrawAddress(ctx, txHash, msg, evMap, extra)
	case dist.MsgWithdrawValidatorCommission:
		h.handleMsgWithdrawValidatorCommission(ctx, txHash, msg, evMap, extra)
	case slashing.MsgUnjail:
		h.handleMsgUnjail(ctx, txHash, msg, evMap, extra)
	case gov.MsgSubmitProposal:
		h.handleMsgSubmitProposal(ctx, txHash, msg, evMap, extra)
	case gov.MsgVote:
		h.handleMsgVote(ctx, txHash, msg, evMap, extra)
	case gov.MsgDeposit:
		h.handleMsgDeposit(ctx, txHash, msg, evMap, extra)
	}
}

func (h *EmitterHook) handleBeginBlockEndBlockEvent(ctx sdk.Context, event abci.Event) {
	events := sdk.StringifyEvents([]abci.Event{event})
	evMap := parseEvents(events)
	switch event.Type {
	case types.EventTypeResolve:
		h.handleEventRequestExecute(ctx, evMap)
	case slashing.EventTypeSlash:
		h.handleEventSlash(ctx, evMap)
	case types.EventTypeDeactivate:
		h.handleEventDeactivate(ctx, evMap)
	case EventTypeCompleteUnbonding:
		h.handleEventTypeCompleteUnbonding(ctx, evMap)
	case EventTypeCompleteRedelegation:
		h.handEventTypeCompleteRedelegation(ctx, evMap)
	case EventTypeInactiveProposal:
		h.handleEventInactiveProposal(ctx, evMap)
	case EventTypeActiveProposal:
		h.handleEventTypeActiveProposal(ctx, evMap)
	case bank.EventTypeTransfer:
		h.handleEventTypeTransfer(ctx, evMap)
	default:
		break
	}
}
