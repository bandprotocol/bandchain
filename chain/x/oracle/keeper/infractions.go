package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// HandleValidatorReport handles a validator report, must be called once per validator per request.
func (k Keeper) HandleValidatorReport(ctx sdk.Context, val sdk.ValAddress, reported bool) {
	logger := k.Logger(ctx)
	// Fetch the existing report info of this validator.
	info := k.GetValidatorReportInfoWithDefault(ctx, val)
	validator := k.StakingKeeper.Validator(ctx, val)
	if validator == nil || validator.IsJailed() {
		// If validator not found or has been jailed, we reset the consecutive miss counts.
		info := k.GetValidatorReportInfoWithDefault(ctx, val)
		info.ConsecutiveMissed = 0
		k.SetValidatorReportInfo(ctx, val, info)
		logger.Info(fmt.Sprintf("Validator %s missed report, but was either not found in store or already jailed", val))
		return
	}
	// Update the consecutive misses of this validator accordingly.
	if reported {
		info.ConsecutiveMissed = 0
	} else {
		info.ConsecutiveMissed++
	}
	maxMisses := k.GetParam(ctx, types.KeyMaxConsecutiveMisses)
	// if the validator misses reports consecutively more than max misses, then jail him/her!
	if info.ConsecutiveMissed > maxMisses {
		logger.Info(fmt.Sprintf("Validator %s missed report more than %d", val, maxMisses))
		consAddr := validator.GetConsAddr()
		// Emit slashing event to notify that the jail occurs.
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			slashing.EventTypeSlash,
			sdk.NewAttribute(slashing.AttributeKeyJailed, consAddr.String()),
		))
		k.StakingKeeper.Jail(ctx, consAddr)
		info.ConsecutiveMissed = 0
	}
	// Everything is complete. Now let's udpate the validator info accordingly.
	k.SetValidatorReportInfo(ctx, val, info)
}
