package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// HandleValidatorReport handles a validator report, must be called once per validator per request.
func (k Keeper) HandleValidatorReport(ctx sdk.Context, requestID types.RequestID, address sdk.ValAddress, reported bool) {
	logger := k.Logger(ctx)

	// If validator not found or has been jailed, we will reset consecutive miss count if report info exists.
	validator := k.StakingKeeper.Validator(ctx, address)
	if validator == nil || validator.IsJailed() {
		// Validator was (a) not found or (b) already jailed
		reportInfo := k.GetValidatorReportInfoWithDefault(ctx, address)
		reportInfo.ConsecutiveMissed = 0
		k.SetValidatorReportInfo(ctx, address, reportInfo)

		logger.Info(
			fmt.Sprintf("Validator %s missed report, but was either not found in store or already jailed", address),
		)
		return
	}
	// fetch report info
	reportInfo := k.GetValidatorReportInfoWithDefault(ctx, address)

	maxMisses := k.GetParam(ctx, types.KeyMaxConsecutiveMisses)
	if reported {
		reportInfo.ConsecutiveMissed = 0
	} else {
		reportInfo.ConsecutiveMissed++
	}

	// if validator missed report consecutively more than max misses, then jail him.
	if reportInfo.ConsecutiveMissed > maxMisses {
		// Downtime confirmed: jail the validator
		logger.Info(fmt.Sprintf("Validator %s missed report more than %d",
			address, maxMisses))

		consAddr := validator.GetConsAddr()

		// Emit slash to notify jailed event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				slashing.EventTypeSlash,
				sdk.NewAttribute(slashing.AttributeKeyJailed, consAddr.String()),
			),
		)

		k.StakingKeeper.Jail(ctx, consAddr)

		// Reset consecutive miss count
		reportInfo.ConsecutiveMissed = 0
	}
	// Set the updated report info
	k.SetValidatorReportInfo(ctx, address, reportInfo)
}
