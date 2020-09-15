package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/emitter/common"
)

// handleMsg handles the given message by publishing relevant events and populates accounts
// that need balance update in 'app.accs'. Also fills in extra info for this message.
func (app *App) handleMsg(txHash []byte, msg sdk.Msg, log sdk.ABCIMessageLog, extra common.JsDict) {
	evMap := common.ParseEvents(log.Events)
	switch msg := msg.(type) {
	case staking.MsgCreateValidator:
		app.handleMsgCreateValidator(txHash, msg, evMap, extra)
	case staking.MsgEditValidator:
		app.handleMsgEditValidator(txHash, msg, evMap, extra)
	case staking.MsgDelegate:
		app.handleMsgDelegate(txHash, msg, evMap, extra)
	case staking.MsgUndelegate:
		app.handleMsgUndelegate(txHash, msg, evMap, extra)
	case staking.MsgBeginRedelegate:
		app.handleMsgBeginRedelegate(txHash, msg, evMap, extra)
	case bank.MsgSend:
		app.handleMsgSend(txHash, msg, evMap, extra)
	case bank.MsgMultiSend:
		app.handleMsgMultiSend(txHash, msg, evMap, extra)
	case dist.MsgWithdrawDelegatorReward:
		app.handleMsgWithdrawDelegatorReward(txHash, msg, evMap, extra)
	case dist.MsgSetWithdrawAddress:
		app.handleMsgSetWithdrawAddress(txHash, msg, evMap, extra)
	case dist.MsgWithdrawValidatorCommission:
		app.handleMsgWithdrawValidatorCommission(txHash, msg, evMap, extra)
	case slashing.MsgUnjail:
		app.handleMsgUnjail(txHash, msg, evMap, extra)
	case gov.MsgSubmitProposal:
		app.handleMsgSubmitProposal(txHash, msg, evMap, extra)
	case gov.MsgVote:
		app.handleMsgVote(txHash, msg, evMap, extra)
	case gov.MsgDeposit:
		app.handleMsgDeposit(txHash, msg, evMap, extra)
	}
}

func (app *App) handleBeginBlockEndBlockEvent(event abci.Event) {
	events := sdk.StringifyEvents([]abci.Event{event})
	evMap := common.ParseEvents(events)
	switch event.Type {
	case slashing.EventTypeSlash:
		app.handleEventSlash(evMap)
	case EventTypeCompleteUnbonding:
		app.handleEventTypeCompleteUnbonding(evMap)
	case EventTypeInactiveProposal:
		app.handleEventInactiveProposal(evMap)
	case EventTypeActiveProposal:
		app.handleEventTypeActiveProposal(evMap)
	case bank.EventTypeTransfer:
		app.handleEventTypeTransfer(evMap)
	default:
		break
	}
}
