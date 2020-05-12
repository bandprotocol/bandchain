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
	newInfo := types.NewValidatorReportInfo(Alice.ValAddress, 5)
	k.SetValidatorReportInfo(ctx, Alice.ValAddress, newInfo)
	info, err := k.GetValidatorReportInfo(ctx, Alice.ValAddress)
	require.Nil(t, err)
	require.Equal(t, info.ConsecutiveMissed, uint64(5))

	_, err = k.GetValidatorReportInfo(ctx, Bob.ValAddress)
	require.Error(t, err)

	require.Panics(t, func() { k.MustGetValidatorReportInfo(ctx, Bob.ValAddress) })
}
