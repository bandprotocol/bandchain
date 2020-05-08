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

	// If validator not found or has been jailed, we will freeze his stat before jail.
	// This assumption help validator to be jailed again when node up.
	validator := k.StakingKeeper.Validator(ctx, address)
	if validator == nil || validator.IsJailed() {
		// Validator was (a) not found or (b) already jailed
		logger.Info(
			fmt.Sprintf("Validator %s missed report, but was either not found in store or already jailed", address),
		)
		return
	}
	// fetch report info
	reportInfo := k.MustGetValidatorReportInfo(ctx, address)

	// this is a relative index, so it counts reports the validator *should* have reported
	// will use the 0-value default report info if not present
	reportedBlockWindow := k.GetParam(ctx, types.KeyReportedWindow)
	minReport := k.GetParam(ctx, types.KeyMinReportedPerWindow)
	index := reportInfo.IndexOffset % reportedBlockWindow
	reportInfo.IndexOffset++

	// Update report bit array & counter
	// This counter just tracks the sum of the bit array
	// That way we avoid needing to read/write the whole array each time
	previous := k.GetValidatorMissedReportBitArray(ctx, address, index)
	missed := !reported
	switch {
	case !previous && missed:
		// Array value has changed from not missed to missed, increment counter
		k.SetValidatorMissedReportBitArray(ctx, address, index, true)
		reportInfo.MissedReportsCounter++
	case previous && !missed:
		// Array value has changed from missed to not missed, decrement counter
		k.SetValidatorMissedReportBitArray(ctx, address, index, false)
		reportInfo.MissedReportsCounter--
	default:
		// Array value at this index has not changed, no need to update counter
	}

	// Emit event if validator missed report?
	// if missed {
	// 	ctx.EventManager().EmitEvent(
	// 		sdk.NewEvent(
	// 			types.EventTypeLiveness,
	// 			sdk.NewAttribute(types.AttributeKeyAddress, consAddr.String()),
	// 			sdk.NewAttribute(types.AttributeKeyMissedBlocks, fmt.Sprintf("%d", signInfo.MissedBlocksCounter)),
	// 			sdk.NewAttribute(types.AttributeKeyHeight, fmt.Sprintf("%d", height)),
	// 		),
	// 	)

	// 	logger.Info(
	// 		fmt.Sprintf("Absent validator %s at height %d, %d missed, threshold %d", consAddr, height, signInfo.MissedBlocksCounter, k.MinSignedPerWindow(ctx)))
	// }

	// Update IsFullTime of report info
	if reportInfo.IndexOffset == reportedBlockWindow {
		reportInfo.IsFullTime = true
	}

	maxMissed := reportedBlockWindow - minReport

	// if we are past the window request and the validator has missed too many blocks, punish them
	if reportInfo.IsFullTime && reportInfo.MissedReportsCounter > maxMissed {
		// Downtime confirmed: jail the validator
		logger.Info(fmt.Sprintf("Validator %s missed report more than %d",
			address, minReport))

		consAddr := validator.GetConsAddr()

		// Emit slash to notify jailed event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				slashing.EventTypeSlash,
				sdk.NewAttribute(slashing.AttributeKeyJailed, consAddr.String()),
			),
		)

		k.StakingKeeper.Jail(ctx, consAddr)

		// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
		reportInfo.MissedReportsCounter = 0
		reportInfo.IndexOffset = 0
		reportInfo.IsFullTime = false
		k.clearValidatorMissedReportBitArray(ctx, address)
	}
	// Set the updated report info
	k.SetValidatorReportInfo(ctx, address, reportInfo)
}
