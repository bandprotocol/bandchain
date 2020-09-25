package emitter

import (
	"github.com/bandprotocol/bandchain/chain/emitter/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
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

func (app *App) emitUpdateValidatorRewardAndAccumulatedCommission(addr sdk.ValAddress) {
	currentReward, currentRatio := app.getCurrentRewardAndCurrentRatio(addr)
	accCommission, _ := app.DistrKeeper.GetValidatorAccumulatedCommission(app.DeliverContext, addr).TruncateDecimal()
	app.Write("UPDATE_VALIDATOR", common.JsDict{
		"operator_address":       addr.String(),
		"current_reward":         currentReward,
		"current_ratio":          currentRatio,
		"accumulated_commission": accCommission.String(),
	})
}

func (app *App) emitUpdateValidatorReward(addr sdk.ValAddress) {
	currentReward, currentRatio := app.getCurrentRewardAndCurrentRatio(addr)
	app.Write("UPDATE_VALIDATOR", common.JsDict{
		"operator_address": addr.String(),
		"current_reward":   currentReward,
		"current_ratio":    currentRatio,
	})
}

// handleMsgWithdrawDelegatorReward implements emitter handler for MsgWithdrawDelegatorReward.
func (app *App) handleMsgWithdrawDelegatorReward(
	txHash []byte, msg dist.MsgWithdrawDelegatorReward, evMap common.EvMap, extra common.JsDict,
) {
	withdrawAddr := app.DistrKeeper.GetDelegatorWithdrawAddr(app.DeliverContext, msg.DelegatorAddress)
	app.AddAccountsInTx(withdrawAddr)
	app.emitUpdateValidatorReward(msg.ValidatorAddress)
	app.emitDelegationAfterWithdrawReward(msg.ValidatorAddress, withdrawAddr)
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
	app.emitUpdateValidatorRewardAndAccumulatedCommission(msg.ValidatorAddress)
	extra["commission_amount"] = evMap[types.EventTypeWithdrawCommission+"."+sdk.AttributeKeyAmount][0]
}
