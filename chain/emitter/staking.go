package emitter

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	types "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	EventTypeCompleteUnbonding = types.EventTypeCompleteUnbonding
)

func (app *App) emitSetValidator(addr sdk.ValAddress) {
	val, _ := app.StakingKeeper.GetValidator(app.DeliverContext, addr)
	currentReward, currentRatio := app.getCurrentRewardAndCurrentRatio(addr)
	accCommission, _ := app.DistrKeeper.GetValidatorAccumulatedCommission(app.DeliverContext, addr).TruncateDecimal()
	app.Write("SET_VALIDATOR", JsDict{
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
	})
}

func (app *App) emitUpdateValidator(addr sdk.ValAddress) {
	val, _ := app.StakingKeeper.GetValidator(app.DeliverContext, addr)
	currentReward, currentRatio := app.getCurrentRewardAndCurrentRatio(addr)
	app.Write("UPDATE_VALIDATOR", JsDict{
		"operator_address": addr.String(),
		"tokens":           val.Tokens.Uint64(),
		"delegator_shares": val.DelegatorShares.String(),
		"current_reward":   currentReward,
		"current_ratio":    currentRatio,
	})
}

func (app *App) emitUpdateValidatorStatus(addr sdk.ValAddress) {
	status := app.OracleKeeper.GetValidatorStatus(app.DeliverContext, addr)
	app.Write("UPDATE_VALIDATOR", JsDict{
		"operator_address": addr.String(),
		"status":           status.IsActive,
		"status_since":     status.Since.UnixNano(),
	})
}

func (app *App) emitDelegationAfterWithdrawReward(operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	_, ratio := app.getCurrentRewardAndCurrentRatio(operatorAddress)
	app.Write("UPDATE_DELEGATION", JsDict{
		"delegator_address": delegatorAddress,
		"operator_address":  operatorAddress,
		"last_ratio":        ratio,
	})
}

func (app *App) emitDelegation(operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	delegation, found := app.StakingKeeper.GetDelegation(app.DeliverContext, delegatorAddress, operatorAddress)
	if found {
		_, ratio := app.getCurrentRewardAndCurrentRatio(operatorAddress)
		app.Write("SET_DELEGATION", JsDict{
			"delegator_address": delegatorAddress,
			"operator_address":  operatorAddress,
			"shares":            delegation.Shares.String(),
			"last_ratio":        ratio,
		})
	} else {
		app.Write("REMOVE_DELEGATION", JsDict{
			"delegator_address": delegatorAddress,
			"operator_address":  operatorAddress,
		})
	}
}

// handleMsgCreateValidator implements emitter handler for MsgCreateValidator.
func (app *App) handleMsgCreateValidator(
	txHash []byte, msg staking.MsgCreateValidator, evMap EvMap, extra JsDict,
) {
	app.emitSetValidator(msg.ValidatorAddress)
	app.emitDelegation(msg.ValidatorAddress, msg.DelegatorAddress)
}

// handleMsgEditValidator implements emitter handler for MsgEditValidator.
func (app *App) handleMsgEditValidator(
	txHash []byte, msg staking.MsgEditValidator, evMap EvMap, extra JsDict,
) {
	app.emitSetValidator(msg.ValidatorAddress)
}

func (app *App) emitUpdateValidatorAndDelegation(operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	app.emitUpdateValidator(operatorAddress)
	app.emitDelegation(operatorAddress, delegatorAddress)
}

// handleMsgDelegate implements emitter handler for MsgDelegate
func (app *App) handleMsgDelegate(
	txHash []byte, msg staking.MsgDelegate, evMap EvMap, extra JsDict,
) {
	app.emitUpdateValidatorAndDelegation(msg.ValidatorAddress, msg.DelegatorAddress)
}

// handleMsgUndelegate implements emitter handler for MsgUndelegate
func (app *App) handleMsgUndelegate(
	txHash []byte, msg staking.MsgUndelegate, evMap EvMap, extra JsDict,
) {
	app.emitUpdateValidatorAndDelegation(msg.ValidatorAddress, msg.DelegatorAddress)
	app.emitUnbondingDelegation(msg, evMap)
}

func (app *App) emitUnbondingDelegation(msg staking.MsgUndelegate, evMap EvMap) {
	completeTime, _ := time.Parse(time.RFC3339, evMap[types.EventTypeUnbond+"."+types.AttributeKeyCompletionTime][0])
	app.Write("NEW_UNBONDING_DELEGATION", JsDict{
		"delegator_address": msg.DelegatorAddress,
		"operator_address":  msg.ValidatorAddress,
		"creation_height":   app.DeliverContext.BlockHeight(),
		"completion_time":   completeTime.UnixNano(),
		"amount":            evMap[types.EventTypeUnbond+"."+sdk.AttributeKeyAmount][0],
	})
}

// handleMsgBeginRedelegate implements emitter handler for MsgBeginRedelegate
func (app *App) handleMsgBeginRedelegate(
	txHash []byte, msg staking.MsgBeginRedelegate, evMap EvMap, extra JsDict,
) {
	app.emitUpdateValidatorAndDelegation(msg.ValidatorSrcAddress, msg.DelegatorAddress)
	app.emitUpdateValidatorAndDelegation(msg.ValidatorDstAddress, msg.DelegatorAddress)
	app.emitUpdateRedelation(msg.ValidatorSrcAddress, msg.ValidatorDstAddress, msg.DelegatorAddress, evMap)
}

func (app *App) emitUpdateRedelation(operatorSrcAddress sdk.ValAddress, operatorDstAddress sdk.ValAddress, delegatorAddress sdk.AccAddress, evMap EvMap) {
	completeTime, _ := time.Parse(time.RFC3339, evMap[types.EventTypeRedelegate+"."+types.AttributeKeyCompletionTime][0])
	app.Write("NEW_REDELEGATION", JsDict{
		"delegator_address":    delegatorAddress.String(),
		"operator_src_address": operatorSrcAddress.String(),
		"operator_dst_address": operatorDstAddress.String(),
		"completion_time":      completeTime.UnixNano(),
		"amount":               evMap[types.EventTypeRedelegate+"."+sdk.AttributeKeyAmount][0],
	})
}

func (app *App) handleEventTypeCompleteUnbonding(evMap EvMap) {
	acc, _ := sdk.AccAddressFromBech32(evMap[types.EventTypeCompleteUnbonding+"."+types.AttributeKeyDelegator][0])
	app.AddAccountsInBlock(acc)
}
