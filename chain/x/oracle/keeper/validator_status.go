package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

// valWithPower is an internal type to track validator with voting power inside of AllocateTokens.
type valWithPower struct {
	val   stakingexported.ValidatorI
	power int64
}

// AllocateTokens allocates a portion of fee collected in the previous blocks to validators that
// that are actively performing oracle tasks. Note that this reward is also subjected to comm tax.
func (k Keeper) AllocateTokens(ctx sdk.Context, previousVotes []abci.VoteInfo) {
	toReward := []valWithPower{}
	totalPower := int64(0)
	for _, vote := range previousVotes {
		val := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
		if k.GetValidatorStatus(ctx, val.GetOperator()).IsActive {
			toReward = append(toReward, valWithPower{val: val, power: vote.Validator.Power})
			totalPower += vote.Validator.Power
		}
	}
	if totalPower == 0 {
		// No active validators performing oracle tasks, nothing needs to be done here.
		return
	}
	feeCollector := k.supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	totalFee := sdk.NewDecCoinsFromCoins(feeCollector.GetCoins()...)
	// Compute the fee allocated for oracle module to distribute to active validators.
	oracleRewardRatio := sdk.NewDecWithPrec(int64(k.GetParam(ctx, types.KeyOracleRewardPercentage)), 2)
	oracleRewardInt, _ := totalFee.MulDecTruncate(oracleRewardRatio).TruncateDecimal()
	// Transfer the oracle reward portion from fee collector to distr module.
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, distr.ModuleName, oracleRewardInt)
	if err != nil {
		panic(err)
	}
	// Convert the transfered tokens back to DecCoins for internal distr allocations.
	oracleReward := sdk.NewDecCoinsFromCoins(oracleRewardInt...)
	remaining := oracleReward
	rewardMultiplier := sdk.OneDec().Sub(k.distrKeeper.GetCommunityTax(ctx))
	// Allocate non-community pool tokens to active validators weighted by voting power.
	for _, each := range toReward {
		powerFraction := sdk.NewDec(each.power).QuoTruncate(sdk.NewDec(totalPower))
		reward := oracleReward.MulDecTruncate(rewardMultiplier).MulDecTruncate(powerFraction)
		k.distrKeeper.AllocateTokensToValidator(ctx, each.val, reward)
		remaining = remaining.Sub(reward)
	}
	// Allocate the remaining coins to the community pool.
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.distrKeeper.SetFeePool(ctx, feePool)
}

// GetValidatorStatus returns the validator status for the given validator. Note that validator
// status is default to [inactive, 0], so new validators start with inactive state.
func (k Keeper) GetValidatorStatus(ctx sdk.Context, val sdk.ValAddress) types.ValidatorStatus {
	bz := ctx.KVStore(k.storeKey).Get(types.ValidatorStatusStoreKey(val))
	if bz == nil {
		return types.NewValidatorStatus(false, time.Time{})
	}
	var status types.ValidatorStatus
	k.cdc.MustUnmarshalBinaryBare(bz, &status)
	return status
}

// SetValidatorStatus sets the validator status for the given validator.
func (k Keeper) SetValidatorStatus(ctx sdk.Context, val sdk.ValAddress, status types.ValidatorStatus) {
	ctx.KVStore(k.storeKey).Set(types.ValidatorStatusStoreKey(val), k.cdc.MustMarshalBinaryBare(status))
}

// Activate changes the given validator's status to active. Returns error if the validator is
// already active or was deactivated recently, as specified by InactivePenaltyDuration parameter.
func (k Keeper) Activate(ctx sdk.Context, val sdk.ValAddress) error {
	status := k.GetValidatorStatus(ctx, val)
	if status.IsActive {
		return types.ErrValidatorAlreadyActive
	}
	penaltyDuration := time.Duration(k.GetParam(ctx, types.KeyInactivePenaltyDuration))
	if !status.Since.IsZero() && status.Since.Add(penaltyDuration).After(ctx.BlockHeader().Time) {
		return types.ErrTooSoonToActivate
	}
	k.SetValidatorStatus(ctx, val, types.NewValidatorStatus(true, ctx.BlockHeader().Time))
	return nil
}

// MissReport changes the given validator's status to inactive. No-op if already inactive or
// if the validator was active after the time the request happened.
func (k Keeper) MissReport(ctx sdk.Context, val sdk.ValAddress, requestTime time.Time) {
	status := k.GetValidatorStatus(ctx, val)
	if status.IsActive && status.Since.Before(requestTime) {
		k.SetValidatorStatus(ctx, val, types.NewValidatorStatus(false, ctx.BlockHeader().Time))
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeDeactivate,
			sdk.NewAttribute(types.AttributeKeyValidator, val.String()),
		))
	}
}
