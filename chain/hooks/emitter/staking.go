package emitter

import (
	"time"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	types "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	EventTypeCompleteUnbonding    = types.EventTypeCompleteUnbonding
	EventTypeCompleteRedelegation = types.EventTypeCompleteRedelegation
)

func (h *Hook) emitStakingModule(ctx sdk.Context) {
	h.stakingKeeper.IterateValidators(ctx, func(_ int64, val exported.ValidatorI) (stop bool) {
		h.emitSetValidator(ctx, val.GetOperator())
		return false
	})

	h.stakingKeeper.IterateAllDelegations(ctx, func(delegation types.Delegation) (stop bool) {
		h.emitDelegation(ctx, delegation.ValidatorAddress, delegation.DelegatorAddress)
		return false
	})
	h.stakingKeeper.IterateRedelegations(ctx, func(_ int64, red types.Redelegation) (stop bool) {
		for _, entry := range red.Entries {
			h.Write("NEW_REDELEGATION", common.JsDict{
				"delegator_address":    red.DelegatorAddress,
				"operator_src_address": red.ValidatorSrcAddress,
				"operator_dst_address": red.ValidatorDstAddress,
				"completion_time":      entry.CompletionTime.UnixNano(),
				"amount":               entry.SharesDst.String(),
			})
		}
		return false
	})
	h.stakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd types.UnbondingDelegation) (stop bool) {
		for _, entry := range ubd.Entries {
			h.Write("NEW_UNBONDING_DELEGATION", common.JsDict{
				"delegator_address": ubd.DelegatorAddress,
				"operator_address":  ubd.ValidatorAddress,
				"completion_time":   entry.CompletionTime.UnixNano(),
				"amount":            entry.Balance.String(),
			})
		}
		return false
	})
}

func (h *Hook) emitSetValidator(ctx sdk.Context, addr sdk.ValAddress) {
	val, _ := h.stakingKeeper.GetValidator(ctx, addr)
	currentReward, currentRatio := h.getCurrentRewardAndCurrentRatio(ctx, addr)
	accCommission, _ := h.distrKeeper.GetValidatorAccumulatedCommission(ctx, addr).TruncateDecimal()
	h.Write("SET_VALIDATOR", common.JsDict{
		"operator_address":       addr.String(),
		"delegator_address":      sdk.AccAddress(addr).String(),
		"consensus_address":      sdk.ConsAddress(val.ConsPubKey.Address()).String(),
		"consensus_pubkey":       sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, val.ConsPubKey),
		"moniker":                val.Description.Moniker,
		"identity":               val.Description.Identity,
		"website":                val.Description.Website,
		"details":                val.Description.Details,
		"commission_rate":        val.Commission.Rate.String(),
		"commission_max_rate":    val.Commission.MaxRate.String(),
		"commission_max_change":  val.Commission.MaxChangeRate.String(),
		"min_self_delegation":    val.MinSelfDelegation.String(),
		"tokens":                 val.Tokens.Uint64(),
		"jailed":                 val.Jailed,
		"delegator_shares":       val.DelegatorShares.String(),
		"current_reward":         currentReward,
		"current_ratio":          currentRatio,
		"accumulated_commission": accCommission.String(),
		"last_update":            ctx.BlockTime().UnixNano(),
	})
}

func (h *Hook) emitUpdateValidator(ctx sdk.Context, addr sdk.ValAddress) {
	val, _ := h.stakingKeeper.GetValidator(ctx, addr)
	currentReward, currentRatio := h.getCurrentRewardAndCurrentRatio(ctx, addr)
	h.Write("UPDATE_VALIDATOR", common.JsDict{
		"operator_address": addr.String(),
		"tokens":           val.Tokens.Uint64(),
		"delegator_shares": val.DelegatorShares.String(),
		"current_reward":   currentReward,
		"current_ratio":    currentRatio,
		"last_update":      ctx.BlockTime().UnixNano(),
	})
}

func (h *Hook) emitUpdateValidatorStatus(ctx sdk.Context, addr sdk.ValAddress) {
	status := h.oracleKeeper.GetValidatorStatus(ctx, addr)
	h.Write("UPDATE_VALIDATOR", common.JsDict{
		"operator_address": addr.String(),
		"status":           status.IsActive,
		"status_since":     status.Since.UnixNano(),
	})
}

func (h *Hook) emitDelegationAfterWithdrawReward(ctx sdk.Context, operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	_, ratio := h.getCurrentRewardAndCurrentRatio(ctx, operatorAddress)
	h.Write("UPDATE_DELEGATION", common.JsDict{
		"delegator_address": delegatorAddress,
		"operator_address":  operatorAddress,
		"last_ratio":        ratio,
	})
}

func (h *Hook) emitDelegation(ctx sdk.Context, operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	delegation, found := h.stakingKeeper.GetDelegation(ctx, delegatorAddress, operatorAddress)
	if found {
		_, ratio := h.getCurrentRewardAndCurrentRatio(ctx, operatorAddress)
		h.Write("SET_DELEGATION", common.JsDict{
			"delegator_address": delegatorAddress,
			"operator_address":  operatorAddress,
			"shares":            delegation.Shares.String(),
			"last_ratio":        ratio,
		})
	} else {
		h.Write("REMOVE_DELEGATION", common.JsDict{
			"delegator_address": delegatorAddress,
			"operator_address":  operatorAddress,
		})
	}
}

// handleMsgCreateValidator implements emitter handler for MsgCreateValidator.
func (h *Hook) handleMsgCreateValidator(
	ctx sdk.Context, msg staking.MsgCreateValidator,
) {
	h.emitSetValidator(ctx, msg.ValidatorAddress)
	h.emitDelegation(ctx, msg.ValidatorAddress, msg.DelegatorAddress)
}

// handleMsgEditValidator implements emitter handler for MsgEditValidator.
func (h *Hook) handleMsgEditValidator(
	ctx sdk.Context, msg staking.MsgEditValidator,
) {
	h.emitSetValidator(ctx, msg.ValidatorAddress)
}

func (h *Hook) emitUpdateValidatorAndDelegation(ctx sdk.Context, operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	h.emitUpdateValidator(ctx, operatorAddress)
	h.emitDelegation(ctx, operatorAddress, delegatorAddress)
}

// handleMsgDelegate implements emitter handler for MsgDelegate
func (h *Hook) handleMsgDelegate(
	ctx sdk.Context, msg staking.MsgDelegate,
) {
	h.emitUpdateValidatorAndDelegation(ctx, msg.ValidatorAddress, msg.DelegatorAddress)
}

// handleMsgUndelegate implements emitter handler for MsgUndelegate
func (h *Hook) handleMsgUndelegate(
	ctx sdk.Context, msg staking.MsgUndelegate, evMap common.EvMap,
) {
	h.emitUpdateValidatorAndDelegation(ctx, msg.ValidatorAddress, msg.DelegatorAddress)
	h.emitUnbondingDelegation(ctx, msg, evMap)
}

func (h *Hook) emitUnbondingDelegation(ctx sdk.Context, msg staking.MsgUndelegate, evMap common.EvMap) {
	completeTime, _ := time.Parse(time.RFC3339, evMap[types.EventTypeUnbond+"."+types.AttributeKeyCompletionTime][0])
	h.Write("NEW_UNBONDING_DELEGATION", common.JsDict{
		"delegator_address": msg.DelegatorAddress,
		"operator_address":  msg.ValidatorAddress,
		"creation_height":   ctx.BlockHeight(),
		"completion_time":   completeTime.UnixNano(),
		"amount":            evMap[types.EventTypeUnbond+"."+sdk.AttributeKeyAmount][0],
	})
}

// handleMsgBeginRedelegate implements emitter handler for MsgBeginRedelegate
func (h *Hook) handleMsgBeginRedelegate(
	ctx sdk.Context, msg staking.MsgBeginRedelegate, evMap common.EvMap,
) {
	h.emitUpdateValidatorAndDelegation(ctx, msg.ValidatorSrcAddress, msg.DelegatorAddress)
	h.emitUpdateValidatorAndDelegation(ctx, msg.ValidatorDstAddress, msg.DelegatorAddress)
	h.emitUpdateRedelation(msg.ValidatorSrcAddress, msg.ValidatorDstAddress, msg.DelegatorAddress, evMap)
}

func (h *Hook) emitUpdateRedelation(operatorSrcAddress sdk.ValAddress, operatorDstAddress sdk.ValAddress, delegatorAddress sdk.AccAddress, evMap common.EvMap) {
	completeTime, _ := time.Parse(time.RFC3339, evMap[types.EventTypeRedelegate+"."+types.AttributeKeyCompletionTime][0])
	h.Write("NEW_REDELEGATION", common.JsDict{
		"delegator_address":    delegatorAddress.String(),
		"operator_src_address": operatorSrcAddress.String(),
		"operator_dst_address": operatorDstAddress.String(),
		"completion_time":      completeTime.UnixNano(),
		"amount":               evMap[types.EventTypeRedelegate+"."+sdk.AttributeKeyAmount][0],
	})
}

func (h *Hook) handleEventTypeCompleteUnbonding(ctx sdk.Context, evMap common.EvMap) {
	acc, _ := sdk.AccAddressFromBech32(evMap[types.EventTypeCompleteUnbonding+"."+types.AttributeKeyDelegator][0])
	h.Write("REMOVE_UNBONDING", common.JsDict{"timestamp": ctx.BlockTime().UnixNano()})
	h.AddAccountsInBlock(acc)
}

func (h *Hook) handEventTypeCompleteRedelegation(ctx sdk.Context) {
	h.Write("REMOVE_REDELEGATION", common.JsDict{"timestamp": ctx.BlockTime().UnixNano()})
}
