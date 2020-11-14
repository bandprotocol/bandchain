package replay

import (
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
)

// handleMsgSend implements emitter handler for MsgSend.
func (h *Hook) handleMsgSend(
	txHash []byte, msg bank.MsgSend, evMap common.EvMap, extra common.JsDict,
) {
	h.AddAccountsInTx(msg.ToAddress)
}

// handleMsgMultiSend implements emitter handler for MsgMultiSend.
func (h *Hook) handleMsgMultiSend(
	txHash []byte, msg bank.MsgMultiSend, evMap common.EvMap, extra common.JsDict,
) {
	for _, output := range msg.Outputs {
		h.AddAccountsInTx(output.Address)
	}
}
