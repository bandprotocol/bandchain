package emitter

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// handleMsgSend implements emitter handler for MsgSend.
func (app *App) handleMsgSend(
	txHash []byte, msg bank.MsgSend, evMap EvMap, extra JsDict,
) {
	app.AddAccountsInTx(msg.FromAddress, msg.ToAddress)
}

// handleMsgMultiSend implements emitter handler for MsgMultiSend.
func (app *App) handleMsgMultiSend(
	txHash []byte, msg bank.MsgMultiSend, evMap EvMap, extra JsDict,
) {
	for _, input := range msg.Inputs {
		app.AddAccountsInTx(input.Address)
	}
	for _, output := range msg.Outputs {
		app.AddAccountsInTx(output.Address)
	}
}
