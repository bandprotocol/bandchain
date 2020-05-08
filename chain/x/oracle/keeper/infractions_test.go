package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestHandleValidatorReport(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetParam(ctx, types.KeyReportedWindow, 5)
	k.SetParam(ctx, types.KeyMinReportedPerWindow, 2)

	k.HandleValidatorReport(ctx, types.RequestID(1), Validator1.ValAddress, false)
	k.HandleValidatorReport(ctx, types.RequestID(3), Validator1.ValAddress, true)
	k.HandleValidatorReport(ctx, types.RequestID(4), Validator1.ValAddress, false)
	k.HandleValidatorReport(ctx, types.RequestID(6), Validator1.ValAddress, false)

	// This validator has been jailed because it doesn't complete report window period
	info, err := k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, false, 4, 3), info)

	k.HandleValidatorReport(ctx, types.RequestID(8), Validator1.ValAddress, true)
	info, err = k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, true, 5, 3), info)

	// Shift window still have report enough minimum report (2/5)
	k.HandleValidatorReport(ctx, types.RequestID(9), Validator1.ValAddress, false)
	info, err = k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, true, 6, 3), info)
	v := k.StakingKeeper.Validator(ctx, Validator1.ValAddress)
	require.NotNil(t, v)
	require.False(t, v.IsJailed())

	// Miss one more report, so he will be jailed (1/5)
	k.HandleValidatorReport(ctx, types.RequestID(12), Validator1.ValAddress, false)
	info, err = k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, false, 0, 0), info)
	v = k.StakingKeeper.Validator(ctx, Validator1.ValAddress)
	require.NotNil(t, v)
	require.True(t, v.IsJailed())
}

func TestHandleValidatorReportOnNonValidator(t *testing.T) {
	_, ctx, k := createTestInput()

	// Must not found report info on non-validator
	k.HandleValidatorReport(ctx, types.RequestID(1), Alice.ValAddress, false)
	_, err := k.GetValidatorReportInfo(ctx, Alice.ValAddress)
	require.Error(t, err)

	emptyReportInfo := types.NewValidatorReportInfo(Alice.ValAddress, false, 0, 0)
	k.SetValidatorReportInfo(ctx, Alice.ValAddress, emptyReportInfo)

	// If report info has existed, but this validator doesn't be validator anymore.
	// This function should ignore it.
	k.HandleValidatorReport(ctx, types.RequestID(1), Alice.ValAddress, false)
	info, err := k.GetValidatorReportInfo(ctx, Alice.ValAddress)
	require.Nil(t, err)
	require.Equal(t, emptyReportInfo, info)
}

func TestHandleValidatorReportOnJailedValidator(t *testing.T) {
	app, ctx, k := createTestInput()

	// Set time not be zero for unjail
	ctx = ctx.WithBlockTime(time.Unix(int64(1581589800), 0))

	k.HandleValidatorReport(ctx, types.RequestID(1), Validator1.ValAddress, true)
	k.HandleValidatorReport(ctx, types.RequestID(2), Validator1.ValAddress, false)
	info, err := k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, false, 2, 1), info)

	validator := app.StakingKeeper.Validator(ctx, Validator1.ValAddress)
	app.StakingKeeper.Jail(ctx, validator.GetConsAddr())

	// Jailed validator report info should not be updated.
	k.HandleValidatorReport(ctx, types.RequestID(2), Validator1.ValAddress, false)
	info, err = k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, false, 2, 1), info)

	// Try to unjail via slashing module
	err = app.SlashingKeeper.Unjail(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	validator = app.StakingKeeper.Validator(ctx, Validator1.ValAddress)
	require.False(t, validator.IsJailed())

	// Report info must be freeze before jailed
	info, err = k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, false, 2, 1), info)

	k.HandleValidatorReport(ctx, types.RequestID(10), Validator1.ValAddress, false)
	info, err = k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, false, 3, 2), info)
}

func TestHandleValidatorReportUpdateBitArray(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetParam(ctx, types.KeyReportedWindow, 5)
	k.SetParam(ctx, types.KeyMinReportedPerWindow, 2)

	k.HandleValidatorReport(ctx, types.RequestID(1), Validator1.ValAddress, false)
	k.HandleValidatorReport(ctx, types.RequestID(2), Validator1.ValAddress, true)
	k.HandleValidatorReport(ctx, types.RequestID(3), Validator1.ValAddress, true)
	k.HandleValidatorReport(ctx, types.RequestID(4), Validator1.ValAddress, true)
	k.HandleValidatorReport(ctx, types.RequestID(5), Validator1.ValAddress, true)

	info, err := k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, true, 5, 1), info)

	k.HandleValidatorReport(ctx, types.RequestID(6), Validator1.ValAddress, true)
	info, err = k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, true, 6, 0), info)
}
