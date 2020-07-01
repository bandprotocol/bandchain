package emitter

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func (app *App) emitDelegation(operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	fmt.Println("yo")
	fmt.Println(operatorAddress)
	fmt.Println(delegatorAddress)

	delegation, found := app.StakingKeeper.GetDelegation(app.DeliverContext, delegatorAddress, operatorAddress)
	if found {
		info := app.DistrKeeper.GetDelegatorStartingInfo(app.DeliverContext, operatorAddress, delegatorAddress)
		latestReward := app.DistrKeeper.GetValidatorHistoricalRewards(app.DeliverContext, operatorAddress, info.PreviousPeriod)
		app.Write("SET_DELEGATION", JsDict{
			"delegator_address": delegatorAddress,
			"operator_address":  operatorAddress,
			"shares":            delegation.Shares.String(),
			"last_ratio":        latestReward.CumulativeRewardRatio[0].Amount.String(),
			"active":            true,
		})
	} else {
		app.Write("REMOVE_DELEGATION", JsDict{
			"delegator_address": delegatorAddress,
			"operator_address":  operatorAddress,
			"active":            false,
		})
	}
}

func (app *App) emitUpdateValidatorAndDelegation(operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	app.emitUpdateValidator(operatorAddress)
	app.emitDelegation(operatorAddress, delegatorAddress)
}

// handleMsgDelegate implements emitter handler for MsgDelegate
func (app *App) handleMsgDelegate(msg staking.MsgDelegate) {
	app.emitUpdateValidatorAndDelegation(msg.ValidatorAddress, msg.DelegatorAddress)
}

// handleMsgUndelegate implements emitter handler for MsgUndelegate
func (app *App) handleMsgUndelegate(msg staking.MsgUndelegate) {
	app.emitUpdateValidatorAndDelegation(msg.ValidatorAddress, msg.DelegatorAddress)

}

// handleMsgBeginRedelegate implements emitter handler for MsgBeginRedelegate
func (app *App) handleMsgBeginRedelegate(msg staking.MsgBeginRedelegate) {
	app.emitUpdateValidatorAndDelegation(msg.ValidatorSrcAddress, msg.DelegatorAddress)
	app.emitUpdateValidatorAndDelegation(msg.ValidatorDstAddress, msg.DelegatorAddress)
}
