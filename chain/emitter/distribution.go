package emitter

import (
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
)

// handleMsgWithdrawDelegatorReward implements emitter handler for MsgWithdrawDelegatorReward.
func (app *App) handleMsgWithdrawDelegatorReward(
	txHash []byte, msg dist.MsgWithdrawDelegatorReward, evMap EvMap, extra JsDict,
) {
	app.emitUpdateValidatorReward(msg.ValidatorAddress)
}
