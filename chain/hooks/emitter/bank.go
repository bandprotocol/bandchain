package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
)

// handleMsgSend implements emitter handler for MsgSend.
func (h *Hook) handleMsgSend(msg bank.MsgSend) {
	h.AddAccountsInTx(msg.ToAddress)
}

// handleMsgMultiSend implements emitter handler for MsgMultiSend.
func (h *Hook) handleMsgMultiSend(msg bank.MsgMultiSend) {
	for _, output := range msg.Outputs {
		h.AddAccountsInTx(output.Address)
	}
}

func (h *Hook) handleEventTypeTransfer(evMap common.EvMap) {
	if recipient, err := sdk.AccAddressFromBech32(evMap[bank.EventTypeTransfer+"."+bank.AttributeKeyRecipient][0]); err == nil {
		h.AddAccountsInBlock(recipient)
	}
}
