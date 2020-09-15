package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/bandprotocol/bandchain/chain/emitter/common"
)

// handleMsgWithdrawDelegatorReward implements emitter handler for MsgWithdrawDelegatorReward.
func (app *App) handleMsgWithdrawDelegatorReward(
	txHash []byte, msg dist.MsgWithdrawDelegatorReward, evMap common.EvMap, extra common.JsDict,
) {
	withdrawAddr := app.DistrKeeper.GetDelegatorWithdrawAddr(app.DeliverContext, msg.DelegatorAddress)
	app.AddAccountsInTx(withdrawAddr)
	extra["reward_amount"] = evMap[dist.EventTypeWithdrawRewards+"."+sdk.AttributeKeyAmount][0]
}

// handleMsgSetWithdrawAddress implements emitter handler for MsgSetWithdrawAddress.
func (app *App) handleMsgSetWithdrawAddress(
	txHash []byte, msg dist.MsgSetWithdrawAddress, evMap common.EvMap, extra common.JsDict,
) {
	app.AddAccountsInTx(msg.WithdrawAddress)
}

// handleMsgWithdrawValidatorCommission implements emitter handler for MsgWithdrawValidatorCommission.
func (app *App) handleMsgWithdrawValidatorCommission(
	txHash []byte, msg dist.MsgWithdrawValidatorCommission, evMap common.EvMap, extra common.JsDict,
) {
	withdrawAddr := app.DistrKeeper.GetDelegatorWithdrawAddr(app.DeliverContext, sdk.AccAddress(msg.ValidatorAddress))
	app.AddAccountsInTx(withdrawAddr)
	extra["commission_amount"] = evMap[types.EventTypeWithdrawCommission+"."+sdk.AttributeKeyAmount][0]
}
