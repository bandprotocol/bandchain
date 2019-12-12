package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetReportsIterator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	validatorAddress1, _ := sdk.ValAddressFromHex("4aea6cfc5bd14f2308954d544e1dc905268357db")
	validatorAddress2, _ := sdk.ValAddressFromHex("4bca6cfc5bd14f2308954d544e1dc905268357db")

	data1 := []byte("This is data1")
	data2 := []byte("This is data2")

	keeper.SetReport(ctx, 1, validatorAddress1, data1)
	keeper.SetReport(ctx, 1, validatorAddress2, data2)

	iterator := keeper.GetReportsIterator(ctx, 1)
	var i int
	for i = 0; iterator.Valid(); iterator.Next() {
		i++
	}
	require.Equal(t, 2, i)
}

func TestGetDataReports(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	validatorAddress1, _ := sdk.ValAddressFromHex("4aea6cfc5bd14f2308954d544e1dc905268357db")
	validatorAddress2, _ := sdk.ValAddressFromHex("4bca6cfc5bd14f2308954d544e1dc905268357db")

	data1 := []byte("This is data1")
	data2 := []byte("This is data2")

	datas := [][]byte{data1, data2}

	keeper.SetReport(ctx, 1, validatorAddress1, data1)
	keeper.SetReport(ctx, 1, validatorAddress2, data2)

	packedData := keeper.GetDataReports(ctx, 1)
	var i int
	for _, report := range packedData {
		require.Equal(t, report.Data, datas[i])
		i++
	}
	require.Equal(t, 2, i)
}
