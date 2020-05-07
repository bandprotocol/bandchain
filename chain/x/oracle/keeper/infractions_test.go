package keeper_test

import (
	"testing"

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
