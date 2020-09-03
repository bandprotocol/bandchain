package emitter

import (
	"github.com/bandprotocol/bandchain/chain/emitter/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
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

func (app *App) handleEventTypeTransfer(evMap common.EvMap) {
	if recipient, err := sdk.AccAddressFromBech32(evMap[bank.EventTypeTransfer+"."+bank.AttributeKeyRecipient][0]); err == nil {
		app.AddAccountsInBlock(recipient)
	}
}
