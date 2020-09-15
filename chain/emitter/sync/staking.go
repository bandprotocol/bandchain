package emitter

import (
	"time"

	"github.com/bandprotocol/bandchain/chain/emitter/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	types "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	EventTypeCompleteUnbonding = types.EventTypeCompleteUnbonding
)

func (app *App) emitStakingModule() {
	app.StakingKeeper.IterateValidators(app.DeliverContext, func(_ int64, val exported.ValidatorI) (stop bool) {
		app.emitSetValidator(val.GetOperator())
		return false
	})

	app.StakingKeeper.IterateAllDelegations(app.DeliverContext, func(delegation types.Delegation) (stop bool) {
		app.emitDelegation(delegation.ValidatorAddress, delegation.DelegatorAddress)
		return false
	})
	app.StakingKeeper.IterateRedelegations(app.DeliverContext, func(_ int64, red types.Redelegation) (stop bool) {
		for _, entry := range red.Entries {
			app.Write("NEW_REDELEGATION", common.JsDict{
				"delegator_address":    red.DelegatorAddress,
				"operator_src_address": red.ValidatorSrcAddress,
				"operator_dst_address": red.ValidatorDstAddress,
				"completion_time":      entry.CompletionTime.UnixNano(),
				"amount":               entry.SharesDst.String(),
			})
		}
		return false
	})
	app.StakingKeeper.IterateUnbondingDelegations(app.DeliverContext, func(_ int64, ubd types.UnbondingDelegation) (stop bool) {
		for _, entry := range ubd.Entries {
			app.Write("NEW_UNBONDING_DELEGATION", common.JsDict{
				"delegator_address": ubd.DelegatorAddress,
				"operator_address":  ubd.ValidatorAddress,
				"completion_time":   entry.CompletionTime.UnixNano(),
				"amount":            entry.Balance.String(),
			})
		}
		return false
	})
}

func (app *App) emitSetValidator(addr sdk.ValAddress) {
	val, _ := app.StakingKeeper.GetValidator(app.DeliverContext, addr)
	currentReward, currentRatio := app.getCurrentRewardAndCurrentRatio(addr)
	accCommission, _ := app.DistrKeeper.GetValidatorAccumulatedCommission(app.DeliverContext, addr).TruncateDecimal()
	app.Write("SET_VALIDATOR", common.JsDict{
		"operator_address":       addr.String(),
		"delegator_address":      sdk.AccAddress(addr).String(),
		"consensus_address":      sdk.ConsAddress(val.ConsPubKey.Address()).String(),
		"consensus_pubkey":       sdk.MustBech32ifyConsPub(val.ConsPubKey),
		"moniker":                val.Description.Moniker,
		"identity":               val.Description.Identity,
		"website":                val.Description.Website,
		"details":                val.Description.Details,
		"commission_rate":        val.Commission.Rate.String(),
		"commission_max_rate":    val.Commission.MaxRate.String(),
		"commission_max_change":  val.Commission.MaxChangeRate.String(),
		"min_self_delegation":    val.MinSelfDelegation.String(),
		"tokens":                 val.Tokens.BigInt().Uint64(),
		"jailed":                 val.Jailed,
		"delegator_shares":       val.DelegatorShares.String(),
		"current_reward":         currentReward,
		"current_ratio":          currentRatio,
		"accumulated_commission": accCommission.String(),
	})
	common.EmitSetHistoricalBondedTokenOnValidator(app, addr, val.Tokens.BigInt().Uint64(), app.DeliverContext.BlockTime().UnixNano())
}

func (app *App) emitUpdateValidator(addr sdk.ValAddress) {
	val, _ := app.StakingKeeper.GetValidator(app.DeliverContext, addr)
	currentReward, currentRatio := app.getCurrentRewardAndCurrentRatio(addr)
	app.Write("UPDATE_VALIDATOR", common.JsDict{
		"operator_address": addr.String(),
		"tokens":           val.Tokens.BigInt().Uint64(),
		"delegator_shares": val.DelegatorShares.String(),
		"current_reward":   currentReward,
		"current_ratio":    currentRatio,
	})
	common.EmitSetHistoricalBondedTokenOnValidator(app, addr, val.Tokens.BigInt().Uint64(), app.DeliverContext.BlockTime().UnixNano())

}

func (app *App) emitDelegationAfterWithdrawReward(operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	_, ratio := app.getCurrentRewardAndCurrentRatio(operatorAddress)
	app.Write("UPDATE_DELEGATION", common.JsDict{
		"delegator_address": delegatorAddress,
		"operator_address":  operatorAddress,
		"last_ratio":        ratio,
	})
}

func (app *App) emitDelegation(operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	delegation, found := app.StakingKeeper.GetDelegation(app.DeliverContext, delegatorAddress, operatorAddress)
	if found {
		_, ratio := app.getCurrentRewardAndCurrentRatio(operatorAddress)
		app.Write("SET_DELEGATION", common.JsDict{
			"delegator_address": delegatorAddress,
			"operator_address":  operatorAddress,
			"shares":            delegation.Shares.String(),
			"last_ratio":        ratio,
		})
	} else {
		app.Write("REMOVE_DELEGATION", common.JsDict{
			"delegator_address": delegatorAddress,
			"operator_address":  operatorAddress,
		})
	}
}

// handleMsgCreateValidator implements emitter handler for MsgCreateValidator.
func (app *App) handleMsgCreateValidator(
	txHash []byte, msg staking.MsgCreateValidator, evMap common.EvMap, extra common.JsDict,
) {
	app.emitSetValidator(msg.ValidatorAddress)
	app.emitDelegation(msg.ValidatorAddress, msg.DelegatorAddress)
}

// handleMsgEditValidator implements emitter handler for MsgEditValidator.
func (app *App) handleMsgEditValidator(
	txHash []byte, msg staking.MsgEditValidator, evMap common.EvMap, extra common.JsDict,
) {
	app.emitSetValidator(msg.ValidatorAddress)
}

func (app *App) emitUpdateValidatorAndDelegation(operatorAddress sdk.ValAddress, delegatorAddress sdk.AccAddress) {
	app.emitUpdateValidator(operatorAddress)
	app.emitDelegation(operatorAddress, delegatorAddress)
}

// handleMsgDelegate implements emitter handler for MsgDelegate
func (app *App) handleMsgDelegate(
	txHash []byte, msg staking.MsgDelegate, evMap common.EvMap, extra common.JsDict,
) {
	app.emitUpdateValidatorAndDelegation(msg.ValidatorAddress, msg.DelegatorAddress)
}

// handleMsgUndelegate implements emitter handler for MsgUndelegate
func (app *App) handleMsgUndelegate(
	txHash []byte, msg staking.MsgUndelegate, evMap common.EvMap, extra common.JsDict,
) {
	app.emitUpdateValidatorAndDelegation(msg.ValidatorAddress, msg.DelegatorAddress)
	app.emitUnbondingDelegation(msg, evMap)
}

func (app *App) emitUnbondingDelegation(msg staking.MsgUndelegate, evMap common.EvMap) {
	completeTime, _ := time.Parse(time.RFC3339, evMap[types.EventTypeUnbond+"."+types.AttributeKeyCompletionTime][0])
	common.EmitNewUnbondingDelegation(
		app,
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		completeTime.UnixNano(),
		sdk.NewInt(common.Atoi(evMap[types.EventTypeUnbond+"."+sdk.AttributeKeyAmount][0])),
	)
}

// handleMsgBeginRedelegate implements emitter handler for MsgBeginRedelegate
func (app *App) handleMsgBeginRedelegate(
	txHash []byte, msg staking.MsgBeginRedelegate, evMap common.EvMap, extra common.JsDict,
) {
	app.emitUpdateValidatorAndDelegation(msg.ValidatorSrcAddress, msg.DelegatorAddress)
	app.emitUpdateValidatorAndDelegation(msg.ValidatorDstAddress, msg.DelegatorAddress)
	app.emitUpdateRedelation(msg.ValidatorSrcAddress, msg.ValidatorDstAddress, msg.DelegatorAddress, evMap)
}

func (app *App) emitUpdateRedelation(operatorSrcAddress sdk.ValAddress, operatorDstAddress sdk.ValAddress, delegatorAddress sdk.AccAddress, evMap common.EvMap) {
	completeTime, _ := time.Parse(time.RFC3339, evMap[types.EventTypeRedelegate+"."+types.AttributeKeyCompletionTime][0])
	common.EmitNewRedelegation(
		app,
		delegatorAddress,
		operatorSrcAddress,
		operatorDstAddress,
		completeTime.UnixNano(),
		sdk.NewInt(common.Atoi(evMap[types.EventTypeRedelegate+"."+sdk.AttributeKeyAmount][0])),
	)
}

func (app *App) handleEventTypeCompleteUnbonding(evMap common.EvMap) {
	acc, _ := sdk.AccAddressFromBech32(evMap[types.EventTypeCompleteUnbonding+"."+types.AttributeKeyDelegator][0])
	app.Write("REMOVE_UNBONDING", common.JsDict{"timestamp": app.DeliverContext.BlockTime().UnixNano()})
	app.AddAccountsInBlock(acc)
}

func (app *App) handEventTypeCompleteRedelegation(evMap common.EvMap) {
	app.Write("REMOVE_REDELEGATION", common.JsDict{"timestamp": app.DeliverContext.BlockTime().UnixNano()})
}
