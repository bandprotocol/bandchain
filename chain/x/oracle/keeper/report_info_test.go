package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestGetSetValidatorReportInfo(t *testing.T) {
	_, ctx, k := createTestInput()

	found := k.HasValidatorReportInfo(ctx, Alice.ValAddress)
	require.False(t, found)
	newInfo := types.NewValidatorReportInfo(
		Alice.ValAddress,
		false,
		3,
		10,
	)
	k.SetValidatorReportInfo(ctx, Alice.ValAddress, newInfo)
	info, err := k.GetValidatorReportInfo(ctx, Alice.ValAddress)
	require.Nil(t, err)
	require.False(t, info.IsFullTime)
	require.Equal(t, info.IndexOffset, uint64(3))
	require.Equal(t, info.MissedReportsCounter, uint64(10))
}

func TestGetSetValidatorMissedBlockBitArray(t *testing.T) {
	_, ctx, k := createTestInput()

	missed := k.GetValidatorMissedReportBitArray(ctx, Alice.ValAddress, 0)
	require.False(t, missed) // treat empty key as not missed
	k.SetValidatorMissedReportBitArray(ctx, Alice.ValAddress, 0, true)
	missed = k.GetValidatorMissedReportBitArray(ctx, Alice.ValAddress, 0)
	require.True(t, missed) // now should be missed
}
