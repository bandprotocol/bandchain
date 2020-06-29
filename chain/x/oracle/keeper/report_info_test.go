package keeper_test

import (
	"bytes"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestGetSetReportInfo(t *testing.T) {
	_, ctx, k := createTestInput()

	found := k.HasReportInfo(ctx, Alice.ValAddress)
	require.False(t, found)
	newInfo := types.NewReportInfo(Alice.ValAddress, 5)
	k.SetReportInfo(ctx, Alice.ValAddress, newInfo)
	info := k.GetReportInfoWithDefault(ctx, Alice.ValAddress)
	require.Equal(t, uint64(5), info.ConsecutiveMissed)

	info = k.GetReportInfoWithDefault(ctx, Bob.ValAddress)
	require.Equal(t, types.NewReportInfo(Bob.ValAddress, 0), info)
}

func TestGetAllReportInfos(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetReportInfo(ctx, Validator1.ValAddress, types.NewReportInfo(Validator1.ValAddress, 3))
	k.SetReportInfo(ctx, Validator2.ValAddress, types.NewReportInfo(Validator2.ValAddress, 6))

	expectedInfos := []types.ReportInfo{
		types.NewReportInfo(Validator1.ValAddress, 3),
		types.NewReportInfo(Validator2.ValAddress, 6),
	}

	sort.Slice(expectedInfos, func(i, j int) bool { return bytes.Compare(expectedInfos[i].Validator, expectedInfos[j].Validator) < 0 })
	infos := k.GetAllReportInfos(ctx)
	require.Equal(t, expectedInfos, infos)
}
