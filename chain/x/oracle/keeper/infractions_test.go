package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestHandleValidatorReport(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetParam(ctx, types.KeyMaxConsecutiveMisses, 2)

	k.HandleValidatorReport(ctx, types.RequestID(1), Validator1.ValAddress, false)
	k.HandleValidatorReport(ctx, types.RequestID(3), Validator1.ValAddress, true)
	k.HandleValidatorReport(ctx, types.RequestID(4), Validator1.ValAddress, false)
	k.HandleValidatorReport(ctx, types.RequestID(6), Validator1.ValAddress, false)

	// This validator hasn't been jailed because consecutive misses doesn't reach maximum misses.
	info := k.GetValidatorReportInfoWithDefault(ctx, Validator1.ValAddress)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, 2), info)

	k.HandleValidatorReport(ctx, types.RequestID(8), Validator1.ValAddress, true)
	info = k.GetValidatorReportInfoWithDefault(ctx, Validator1.ValAddress)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, 0), info)

	k.HandleValidatorReport(ctx, types.RequestID(9), Validator1.ValAddress, false)
	k.HandleValidatorReport(ctx, types.RequestID(10), Validator1.ValAddress, false)
	info = k.GetValidatorReportInfoWithDefault(ctx, Validator1.ValAddress)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, 2), info)
	v := k.StakingKeeper.Validator(ctx, Validator1.ValAddress)
	require.NotNil(t, v)
	require.False(t, v.IsJailed())

	// Miss one more report, so he will be jailed (3 consecutive misses)
	k.HandleValidatorReport(ctx, types.RequestID(12), Validator1.ValAddress, false)
	info = k.GetValidatorReportInfoWithDefault(ctx, Validator1.ValAddress)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, 0), info)
	v = k.StakingKeeper.Validator(ctx, Validator1.ValAddress)
	require.NotNil(t, v)
	require.True(t, v.IsJailed())
}

func TestHandleValidatorReportOnNonValidator(t *testing.T) {
	_, ctx, k := createTestInput()

	// Must not found report info on non-validator
	k.HandleValidatorReport(ctx, types.RequestID(1), Alice.ValAddress, false)
	info := k.GetValidatorReportInfoWithDefault(ctx, Alice.ValAddress)
	require.Equal(t, types.NewValidatorReportInfo(Alice.ValAddress, 0), info)

	emptyReportInfo := types.NewValidatorReportInfo(Alice.ValAddress, 0)
	k.SetValidatorReportInfo(ctx, Alice.ValAddress, emptyReportInfo)

	// If report info has existed, but this validator doesn't be validator anymore.
	// This function should ignore it.
	k.HandleValidatorReport(ctx, types.RequestID(1), Alice.ValAddress, false)
	info = k.GetValidatorReportInfoWithDefault(ctx, Alice.ValAddress)
	require.Equal(t, emptyReportInfo, info)
}

func TestHandleValidatorReportOnJailedValidator(t *testing.T) {
	app, ctx, k := createTestInput()

	// Set time not be zero for unjail
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	k.HandleValidatorReport(ctx, types.RequestID(1), Validator1.ValAddress, true)
	k.HandleValidatorReport(ctx, types.RequestID(2), Validator1.ValAddress, false)
	info := k.GetValidatorReportInfoWithDefault(ctx, Validator1.ValAddress)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, 1), info)

	validator := app.StakingKeeper.Validator(ctx, Validator1.ValAddress)
	app.StakingKeeper.Jail(ctx, validator.GetConsAddr())

	// Jailed validator report info should be reset.
	k.HandleValidatorReport(ctx, types.RequestID(2), Validator1.ValAddress, false)
	info = k.GetValidatorReportInfoWithDefault(ctx, Validator1.ValAddress)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, 0), info)

	// Try to unjail via slashing module
	err := app.SlashingKeeper.Unjail(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	validator = app.StakingKeeper.Validator(ctx, Validator1.ValAddress)
	require.False(t, validator.IsJailed())

	k.HandleValidatorReport(ctx, types.RequestID(10), Validator1.ValAddress, false)
	info = k.GetValidatorReportInfoWithDefault(ctx, Validator1.ValAddress)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, 1), info)
}
