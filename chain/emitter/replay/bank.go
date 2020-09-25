package emitter

import (
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/bandprotocol/bandchain/chain/emitter/common"
)

// handleMsgSend implements emitter handler for MsgSend.
func (app *App) handleMsgSend(
	txHash []byte, msg bank.MsgSend, evMap common.EvMap, extra common.JsDict,
) {
	app.AddAccountsInTx(msg.ToAddress)
}

// handleMsgMultiSend implements emitter handler for MsgMultiSend.
func (app *App) handleMsgMultiSend(
	txHash []byte, msg bank.MsgMultiSend, evMap common.EvMap, extra common.JsDict,
) {
	for _, output := range msg.Outputs {
		app.AddAccountsInTx(output.Address)
	}
}
