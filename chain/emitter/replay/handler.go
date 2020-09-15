package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/emitter/common"
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
// that need balance update in 'app.accs'. Also fills in extra info for this message.
func (app *App) handleMsg(txHash []byte, msg sdk.Msg, log sdk.ABCIMessageLog, extra common.JsDict) {
	evMap := common.ParseEvents(log.Events)
	switch msg := msg.(type) {
	case staking.MsgCreateValidator:
		app.handleMsgCreateValidator(txHash, msg, evMap, extra)
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
	default:
		break
	}
}

func (app *App) handleBeginBlockEndBlockEvent(event abci.Event) {

}
