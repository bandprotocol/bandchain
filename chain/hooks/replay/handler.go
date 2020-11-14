package replay

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// handleMsg handles the given message by publishing relevant events and populates accounts
// that need balance update in 'app.accs'. Also fills in extra info for this message.
func (h *Hook) handleMsg(ctx sdk.Context, txHash []byte, msg sdk.Msg, log sdk.ABCIMessageLog, extra common.JsDict) {
	evMap := common.ParseEvents(log.Events)
	switch msg := msg.(type) {
	case oracle.MsgRequestData:
		h.handleMsgRequestData(ctx, txHash, msg, evMap, extra)
	case oracle.MsgReportData:
		h.handleMsgReportData(txHash, msg, evMap, extra)
	case oracle.MsgCreateDataSource:
		h.handleMsgCreateDataSource(txHash, msg, evMap, extra)
	case oracle.MsgCreateOracleScript:
		h.handleMsgCreateOracleScript(txHash, msg, evMap, extra)
	case oracle.MsgEditDataSource:
		h.handleMsgEditDataSource(txHash, msg, evMap, extra)
	case oracle.MsgActivate:
		h.handleMsgActivate(ctx, txHash, msg, evMap, extra)
	case oracle.MsgAddReporter:
		h.handleMsgAddReporter(ctx, txHash, msg, evMap, extra)
	case oracle.MsgRemoveReporter:
		h.handleMsgRemoveReporter(ctx, txHash, msg, evMap, extra)
	case staking.MsgCreateValidator:
		h.handleMsgCreateValidator(ctx, txHash, msg, evMap, extra)
	case staking.MsgDelegate:
		h.handleMsgDelegate(ctx, txHash, msg, evMap, extra)
	case staking.MsgUndelegate:
		h.handleMsgUndelegate(ctx, txHash, msg, evMap, extra)
	case staking.MsgBeginRedelegate:
		h.handleMsgBeginRedelegate(ctx, txHash, msg, evMap, extra)
	case bank.MsgSend:
		h.handleMsgSend(txHash, msg, evMap, extra)
	case bank.MsgMultiSend:
		h.handleMsgMultiSend(txHash, msg, evMap, extra)
	case dist.MsgWithdrawDelegatorReward:
		h.handleMsgWithdrawDelegatorReward(ctx, txHash, msg, evMap, extra)
	case dist.MsgSetWithdrawAddress:
		h.handleMsgSetWithdrawAddress(txHash, msg, evMap, extra)
	case dist.MsgWithdrawValidatorCommission:
		h.handleMsgWithdrawValidatorCommission(ctx, txHash, msg, evMap, extra)
	default:
		break
	}
}

func (h *Hook) handleBeginBlockEndBlockEvent(ctx sdk.Context, event abci.Event) {
	events := sdk.StringifyEvents([]abci.Event{event})
	evMap := common.ParseEvents(events)
	switch event.Type {
	case types.EventTypeDeactivate:
		h.handleEventDeactivate(ctx, evMap)
	default:
		break
	}
}
