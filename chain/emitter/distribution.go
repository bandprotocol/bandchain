package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
)

func (app *App) getCurrentRewardAndCurrentRatio(addr sdk.ValAddress) (string, string) {
	currentReward := "0"
	currentRatio := "0"

	reward := app.DistrKeeper.GetValidatorCurrentRewards(app.DeliverContext, addr)
	latestReward := app.DistrKeeper.GetValidatorHistoricalRewards(app.DeliverContext, addr, reward.Period-1)

	if !reward.Rewards.IsZero() {
		currentReward = reward.Rewards[0].Amount.String()
	}
	if !latestReward.CumulativeRewardRatio.IsZero() {
		currentRatio = latestReward.CumulativeRewardRatio[0].Amount.String()
	}

	return currentReward, currentRatio
}

func (app *App) emitUpdateValidatorReward(addr sdk.ValAddress) {
	currentReward, currentRatio := app.getCurrentRewardAndCurrentRatio(addr)
	app.Write("UPDATE_VALIDATOR", JsDict{
		"operator_address": addr.String(),
		"current_reward":   currentReward,
		"current_ratio":    currentRatio,
	})
}

// handleMsgWithdrawDelegatorReward implements emitter handler for MsgWithdrawDelegatorReward.
func (app *App) handleMsgWithdrawDelegatorReward(
	txHash []byte, msg dist.MsgWithdrawDelegatorReward, evMap EvMap, extra JsDict,
) {
	app.AddAccountsInTx(msg.DelegatorAddress.String())
	app.emitUpdateValidatorReward(msg.ValidatorAddress)
}

// handleMsgSetWithdrawAddress implements emitter handler for MsgSetWithdrawAddress.
func (app *App) handleMsgSetWithdrawAddress(
	txHash []byte, msg dist.MsgSetWithdrawAddress, evMap EvMap, extra JsDict,
) {
	app.AddAccountsInTx(msg.WithdrawAddress.String(), msg.DelegatorAddress.String())
}
