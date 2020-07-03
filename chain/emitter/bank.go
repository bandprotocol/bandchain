package emitter

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// handleMsgSend implements emitter handler for MsgSend.
func (app *App) handleMsgSend(msg bank.MsgSend) {
	app.AddAccounts(msg.FromAddress, msg.ToAddress)
}

// handleMsgMultiSend implements emitter handler for MsgMultiSend.
func (app *App) handleMsgMultiSend(msg bank.MsgMultiSend) {
	for _, input := range msg.Inputs {
		app.AddAccounts(input.Address)
	}
	for _, output := range msg.Outputs {
		app.AddAccounts(output.Address)
	}
}
