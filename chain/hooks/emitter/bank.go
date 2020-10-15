package emitter

import (
	"github.com/bandprotocol/bandchain/chain/hooks/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// handleMsgSend implements emitter handler for MsgSend.
func (h *EmitterHook) handleMsgSend(
	ctx sdk.Context, txHash []byte, msg bank.MsgSend, evMap common.EvMap, extra common.JsDict,
) {
	h.AddAccountsInTx(msg.ToAddress)
}

// handleMsgMultiSend implements emitter handler for MsgMultiSend.
func (h *EmitterHook) handleMsgMultiSend(
	ctx sdk.Context, txHash []byte, msg bank.MsgMultiSend, evMap common.EvMap, extra common.JsDict,
) {
	for _, output := range msg.Outputs {
		h.AddAccountsInTx(output.Address)
	}
}

func (h *EmitterHook) handleEventTypeTransfer(ctx sdk.Context, evMap common.EvMap) {
	if recipient, err := sdk.AccAddressFromBech32(evMap[bank.EventTypeTransfer+"."+bank.AttributeKeyRecipient][0]); err == nil {
		h.AddAccountsInBlock(recipient)
	}
}
