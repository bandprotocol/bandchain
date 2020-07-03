package emitter

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// handleMsgSend implements emitter handler for MsgSend.
func (app *App) handleMsgSend(msg bank.MsgSend) {
	app.accs = append(app.accs, msg.FromAddress, msg.ToAddress)
}

// handleMsgMultiSend implements emitter handler for MsgMultiSend.
func (app *App) handleMsgMultiSend(msg bank.MsgMultiSend) {
	for _, input := range msg.Inputs {
		app.accs = append(app.accs, input.Address)
	}
	for _, output := range msg.Outputs {
		app.accs = append(app.accs, output.Address)
	}
}
