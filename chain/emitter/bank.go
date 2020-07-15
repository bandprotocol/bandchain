package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	types "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// handleMsgSend implements emitter handler for MsgSend.
func (app *App) handleMsgSend(
	txHash []byte, msg bank.MsgSend, evMap EvMap, extra JsDict,
) {
	app.AddAccountsInTx(msg.ToAddress)
}

// handleMsgMultiSend implements emitter handler for MsgMultiSend.
func (app *App) handleMsgMultiSend(
	txHash []byte, msg bank.MsgMultiSend, evMap EvMap, extra JsDict,
) {
	for _, output := range msg.Outputs {
		app.AddAccountsInTx(output.Address)
	}
}

func (app *App) handleEventTypeTransfer(evMap EvMap) {
	acc, _ := sdk.AccAddressFromBech32(evMap[types.EventTypeTransfer+"."+types.AttributeKeyRecipient][0])
	app.AddAccountsInBlock(acc)
}
