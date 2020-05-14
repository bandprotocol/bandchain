package keeper_test

import (
	"bytes"
	"sort"
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
}

func TestGetAllValidatorReportInfos(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetValidatorReportInfo(ctx, Validator1.ValAddress, types.NewValidatorReportInfo(Validator1.ValAddress, 3))
	k.SetValidatorReportInfo(ctx, Validator2.ValAddress, types.NewValidatorReportInfo(Validator2.ValAddress, 6))

	expectedInfos := []types.ValidatorReportInfo{
		types.NewValidatorReportInfo(Validator1.ValAddress, 3),
		types.NewValidatorReportInfo(Validator2.ValAddress, 6),
	}

	sort.Slice(expectedInfos, func(i, j int) bool { return bytes.Compare(expectedInfos[i].Validator, expectedInfos[j].Validator) < 0 })
	infos := k.GetAllValidatorReportInfos(ctx)
	require.Equal(t, expectedInfos, infos)
}
