package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
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
	if recipient, err := sdk.AccAddressFromBech32(evMap[bank.EventTypeTransfer+"."+bank.AttributeKeyRecipient][0]); err == nil {
		app.AddAccountsInBlock(recipient)
	}
}
