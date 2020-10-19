package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
)

func (h *EmitterHook) getCurrentRewardAndCurrentRatio(ctx sdk.Context, addr sdk.ValAddress) (string, string) {
	currentReward := "0"
	currentRatio := "0"

	reward := h.distrKeeper.GetValidatorCurrentRewards(ctx, addr)
	latestReward := h.distrKeeper.GetValidatorHistoricalRewards(ctx, addr, reward.Period-1)

	if !reward.Rewards.IsZero() {
		currentReward = reward.Rewards[0].Amount.String()
	}
	if !latestReward.CumulativeRewardRatio.IsZero() {
		currentRatio = latestReward.CumulativeRewardRatio[0].Amount.String()
	}

	return currentReward, currentRatio
}

func (h *EmitterHook) emitUpdateValidatorRewardAndAccumulatedCommission(ctx sdk.Context, addr sdk.ValAddress) {
	currentReward, currentRatio := h.getCurrentRewardAndCurrentRatio(ctx, addr)
	accCommission, _ := h.distrKeeper.GetValidatorAccumulatedCommission(ctx, addr).TruncateDecimal()
	h.Write("UPDATE_VALIDATOR", common.JsDict{
		"operator_address":       addr.String(),
		"current_reward":         currentReward,
		"current_ratio":          currentRatio,
		"accumulated_commission": accCommission.String(),
	})
}

func (h *EmitterHook) emitUpdateValidatorReward(ctx sdk.Context, addr sdk.ValAddress) {
	currentReward, currentRatio := h.getCurrentRewardAndCurrentRatio(ctx, addr)
	h.Write("UPDATE_VALIDATOR", common.JsDict{
		"operator_address": addr.String(),
		"current_reward":   currentReward,
		"current_ratio":    currentRatio,
	})
}

// handleMsgWithdrawDelegatorReward implements emitter handler for MsgWithdrawDelegatorReward.
func (h *EmitterHook) handleMsgWithdrawDelegatorReward(
	ctx sdk.Context, msg dist.MsgWithdrawDelegatorReward, evMap common.EvMap, extra common.JsDict,
) {
	withdrawAddr := h.distrKeeper.GetDelegatorWithdrawAddr(ctx, msg.DelegatorAddress)
	h.AddAccountsInTx(withdrawAddr)
	h.emitUpdateValidatorReward(ctx, msg.ValidatorAddress)
	h.emitDelegationAfterWithdrawReward(ctx, msg.ValidatorAddress, withdrawAddr)
	extra["reward_amount"] = evMap[dist.EventTypeWithdrawRewards+"."+sdk.AttributeKeyAmount][0]
}

// handleMsgSetWithdrawAddress implements emitter handler for MsgSetWithdrawAddress.
func (h *EmitterHook) handleMsgSetWithdrawAddress(msg dist.MsgSetWithdrawAddress) {
	h.AddAccountsInTx(msg.WithdrawAddress)
}

// handleMsgWithdrawValidatorCommission implements emitter handler for MsgWithdrawValidatorCommission.
func (h *EmitterHook) handleMsgWithdrawValidatorCommission(
	ctx sdk.Context, msg dist.MsgWithdrawValidatorCommission, evMap common.EvMap, extra common.JsDict,
) {
	withdrawAddr := h.distrKeeper.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(msg.ValidatorAddress))
	h.AddAccountsInTx(withdrawAddr)
	h.emitUpdateValidatorRewardAndAccumulatedCommission(ctx, msg.ValidatorAddress)
	extra["commission_amount"] = evMap[types.EventTypeWithdrawCommission+"."+sdk.AttributeKeyAmount][0]
}
