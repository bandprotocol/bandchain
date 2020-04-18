package keeper

// import (
// 	"testing"

// 	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/require"
// )

// func TestGetterSetterRawDataReport(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	validatorAddress1, _ := sdk.ValAddressFromHex("4aea6cfc5bd14f2308954d544e1dc905268357db")

// 	keeper.SetReport(ctx, 1, 3, validatorAddress1, types.NewRawDataReport(0, []byte("Data1/3")))

// 	report, err := keeper.GetRawDataReport(ctx, 1, 3, validatorAddress1)
// 	require.Nil(t, err)
// 	require.Equal(t, types.NewRawDataReport(0, []byte("Data1/3")), report)

// 	_, err = keeper.GetRawDataReport(ctx, 2, 3, validatorAddress1)
// 	require.NotNil(t, err)

// 	_, err = keeper.GetRawDataReport(ctx, 1, 1, validatorAddress1)
// 	require.NotNil(t, err)

// 	_, err = keeper.GetRawDataReport(ctx, 1, 3, sdk.ValAddress([]byte("val1")))
// 	require.NotNil(t, err)
// }

// func TestAddReportSuccess(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	request := newDefaultRequest()
// 	keeper.SetRequest(ctx, 1, request)

// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(2, 1, []byte("calldata1")))
// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(10, 2, []byte("calldata2")))

// 	err := keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 1, []byte("data1/1")),
// 		types.NewRawDataReportWithID(10, 0, []byte("data2/1")),
// 	}, sdk.ValAddress([]byte("validator1")), sdk.AccAddress([]byte("validator1")))

// 	require.Nil(t, err)

// 	req, err := keeper.GetRequest(ctx, 1)
// 	require.Nil(t, err)
// 	require.Equal(t, []sdk.ValAddress{sdk.ValAddress([]byte("validator1"))}, req.ReceivedValidators)

// 	report, err := keeper.GetRawDataReport(ctx, 1, 2, sdk.ValAddress([]byte("validator1")))
// 	require.Nil(t, err)
// 	require.Equal(t, types.NewRawDataReport(1, []byte("data1/1")), report)

// 	report, err = keeper.GetRawDataReport(ctx, 1, 10, sdk.ValAddress([]byte("validator1")))
// 	require.Nil(t, err)
// 	require.Equal(t, types.NewRawDataReport(0, []byte("data2/1")), report)

// 	list := keeper.GetPendingResolveList(ctx)
// 	require.Equal(t, []types.RequestID{}, list)

// 	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 0, []byte("data1/2")),
// 		types.NewRawDataReportWithID(10, 2, []byte("data2/2")),
// 	}, sdk.ValAddress([]byte("validator2")), sdk.AccAddress([]byte("validator2")))
// 	require.Nil(t, err)

// 	report, err = keeper.GetRawDataReport(ctx, 1, 2, sdk.ValAddress([]byte("validator2")))
// 	require.Nil(t, err)
// 	require.Equal(t, types.NewRawDataReport(0, []byte("data1/2")), report)

// 	report, err = keeper.GetRawDataReport(ctx, 1, 10, sdk.ValAddress([]byte("validator2")))
// 	require.Nil(t, err)
// 	require.Equal(t, types.NewRawDataReport(2, []byte("data2/2")), report)

// 	list = keeper.GetPendingResolveList(ctx)
// 	require.Equal(t, []types.RequestID{1}, list)
// }

// func TestAddReportFailed(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	// Send report on invalid request.
// 	err := keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 0, []byte("data1/1")),
// 		types.NewRawDataReportWithID(10, 0, []byte("data2/1")),
// 	}, sdk.ValAddress([]byte("validator1")), sdk.AccAddress([]byte("validator1")))
// 	require.NotNil(t, err)

// 	// Send report on resolved request.
// 	request := newDefaultRequest()
// 	// request.ResolveStatus = types.Success
// 	// keeper.SetRequest(ctx, 1, request)

// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(2, 1, []byte("calldata1")))
// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(10, 2, []byte("calldata2")))

// 	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 0, []byte("data1/1")),
// 		types.NewRawDataReportWithID(10, 0, []byte("data2/1")),
// 	}, sdk.ValAddress([]byte("validator1")), sdk.AccAddress([]byte("validator1")))
// 	require.NotNil(t, err)

// 	// Send report by invalid validator.
// 	request = newDefaultRequest()
// 	keeper.SetRequest(ctx, 1, request)

// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(2, 1, []byte("calldata1")))
// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(10, 2, []byte("calldata2")))

// 	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 0, []byte("data1/1")),
// 		types.NewRawDataReportWithID(10, 0, []byte("data2/1")),
// 	}, sdk.ValAddress([]byte("nonvalidator")), sdk.AccAddress([]byte("nonvalidator")))
// 	require.NotNil(t, err)

// 	// Send report on expired request.
// 	request = newDefaultRequest()
// 	request.ExpirationHeight = 5
// 	keeper.SetRequest(ctx, 1, request)

// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(2, 1, []byte("calldata1")))
// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(10, 2, []byte("calldata2")))

// 	ctx = ctx.WithBlockHeight(6)
// 	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 0, []byte("data1/1")),
// 		types.NewRawDataReportWithID(10, 0, []byte("data2/1")),
// 	}, sdk.ValAddress([]byte("validator1")), sdk.AccAddress([]byte("validator1")))
// 	require.NotNil(t, err)

// 	// Send incomplete report on request.
// 	request = newDefaultRequest()
// 	keeper.SetRequest(ctx, 1, request)

// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(2, 1, []byte("calldata1")))
// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(10, 2, []byte("calldata2")))

// 	ctx = ctx.WithBlockHeight(2)
// 	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 0, []byte("data1/1")),
// 	}, sdk.ValAddress([]byte("validator1")), sdk.AccAddress([]byte("validator1")))
// 	require.NotNil(t, err)

// 	// Send invalid order report.
// 	request = newDefaultRequest()
// 	keeper.SetRequest(ctx, 1, request)

// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(2, 1, []byte("calldata1")))
// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(10, 2, []byte("calldata2")))

// 	// Send invalid external id.
// 	request = newDefaultRequest()
// 	keeper.SetRequest(ctx, 1, request)

// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(2, 1, []byte("calldata1")))
// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(10, 2, []byte("calldata2")))

// 	ctx = ctx.WithBlockHeight(2)
// 	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(3, 0, []byte("data2/1")),
// 		types.NewRawDataReportWithID(10, 0, []byte("data1/1")),
// 	}, sdk.ValAddress([]byte("validator1")), sdk.AccAddress([]byte("validator1")))
// 	require.NotNil(t, err)

// 	// Cannot report in same request id.
// 	request = newDefaultRequest()
// 	keeper.SetRequest(ctx, 1, request)

// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(2, 1, []byte("calldata1")))
// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(10, 2, []byte("calldata2")))

// 	ctx = ctx.WithBlockHeight(2)
// 	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 0, []byte("OldValue1")),
// 		types.NewRawDataReportWithID(10, 0, []byte("OldValue2")),
// 	}, sdk.ValAddress([]byte("validator2")), sdk.AccAddress([]byte("validator2")))
// 	require.Nil(t, err)

// 	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 0, []byte("NewValue1")),
// 		types.NewRawDataReportWithID(10, 0, []byte("NewValue2")),
// 	}, sdk.ValAddress([]byte("validator2")), sdk.AccAddress([]byte("validator2")))
// 	require.NotNil(t, err)

// }

// func TestAddRawRequestCallDataSizeTooBig(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	request := newDefaultRequest()
// 	keeper.SetRequest(ctx, 1, request)

// 	owner := sdk.AccAddress([]byte("owner"))
// 	name := "data_source"
// 	description := "description"
// 	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
// 	executable := []byte("executable")
// 	keeper.AddDataSource(ctx, owner, name, description, fee, executable)

// 	// Set MaxCalldataSize to 0
// 	// AddRawRequest should fail because size of "calldata" is > 0
// 	keeper.SetParam(ctx, types.KeyMaxCalldataSize, 0)
// 	err := keeper.AddRawRequest(ctx, 1, types.NewRawRequest(1, 1, []byte("calldata")))
// 	require.NotNil(t, err)

// 	// Set MaxCalldataSize to 20
// 	// AddRawRequest should pass because size of "calldata" is < 20
// 	keeper.SetParam(ctx, types.KeyMaxCalldataSize, 20)
// 	err = keeper.AddRawRequest(ctx, 1, types.NewRawRequest(1, 1, []byte("calldata")))
// 	require.Nil(t, err)
// }

// func TestAddReportReportSizeExceedMaxRawDataReportSize(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)
// 	keeper.SetParam(ctx, types.KeyMaxRawDataReportSize, 20)

// 	request := newDefaultRequest()
// 	keeper.SetRequest(ctx, 1, request)

// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(2, 1, []byte("calldata")))

// 	// Size of "short report" is 12 bytes which is shorter than 20 bytes.
// 	err := keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 0, []byte("short report")),
// 	}, sdk.ValAddress([]byte("validator1")), sdk.AccAddress([]byte("validator1")))
// 	require.Nil(t, err)

// 	// Size of "a report that obviously longer than 20 bytes" is 44 bytes.
// 	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 0, []byte("a report that obviously longer than 20 bytes")),
// 	}, sdk.ValAddress([]byte("validator2")), sdk.AccAddress([]byte("validator2")))
// 	require.NotNil(t, err)

// }

// func TestInvalidReport(t *testing.T) {
// 	ctx, keeper := CreateTestInput(t, false)

// 	request := newDefaultRequest()
// 	keeper.SetRequest(ctx, 1, request)

// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(2, 1, []byte("calldata1")))
// 	keeper.SetRawRequest(ctx, 1, types.NewRawRequest(10, 2, []byte("calldata2")))

// 	err := keeper.AddReport(ctx, 1, []types.RawDataReportWithID{
// 		types.NewRawDataReportWithID(2, 1, []byte("data1/1")),
// 		types.NewRawDataReportWithID(10, 0, []byte("data2/1")),
// 	}, sdk.ValAddress([]byte("validator1")), sdk.AccAddress([]byte("nonValidator")))

// 	require.NotNil(t, err)
// }

// // func TestGetReportsIterator(t *testing.T) {
// // 	ctx, keeper := CreateTestInput(t, false)

// // 	validatorAddress1, _ := sdk.ValAddressFromHex("4aea6cfc5bd14f2308954d544e1dc905268357db")
// // 	validatorAddress2, _ := sdk.ValAddressFromHex("4bca6cfc5bd14f2308954d544e1dc905268357db")

// // 	data1 := []types.RawDataReport{
// // 		types.NewRawDataReport(1, []byte("data1:1")),
// // 		types.NewRawDataReport(2, []byte("data2:1")),
// // 	}
// // 	data2 := []types.RawDataReport{
// // 		types.NewRawDataReport(1, []byte("data1:2")),
// // 		types.NewRawDataReport(2, []byte("data2:2")),
// // 	}

// // 	keeper.SetReport(ctx, 1, validatorAddress1, data1)
// // 	keeper.SetReport(ctx, 1, validatorAddress2, data2)

// // 	iterator := keeper.GetReportsIterator(ctx, 1)
// // 	var i int
// // 	for i = 0; iterator.Valid(); iterator.Next() {
// // 		i++
// // 	}
// // 	require.Equal(t, 2, i)
// // }

// // func TestGetDataReports(t *testing.T) {
// // 	ctx, keeper := CreateTestInput(t, false)

// // 	validatorAddress1, _ := sdk.ValAddressFromHex("4aea6cfc5bd14f2308954d544e1dc905268357db")
// // 	validatorAddress2, _ := sdk.ValAddressFromHex("4bca6cfc5bd14f2308954d544e1dc905268357db")

// // 	data1 := []types.RawDataReport{
// // 		types.NewRawDataReport(1, []byte("data1:1")),
// // 		types.NewRawDataReport(2, []byte("data2:1")),
// // 	}
// // 	data2 := []types.RawDataReport{
// // 		types.NewRawDataReport(1, []byte("data1:2")),
// // 		types.NewRawDataReport(2, []byte("data2:2")),
// // 	}

// // 	datas := [][]types.RawDataReport{data1, data2}

// // 	keeper.SetReport(ctx, 1, validatorAddress1, data1)
// // 	keeper.SetReport(ctx, 1, validatorAddress2, data2)

// // 	packedData := keeper.GetDataReports(ctx, 1)
// // 	var i int
// // 	for _, report := range packedData {
// // 		require.Equal(t, report.Data, datas[i])
// // 		i++
// // 	}
// // 	require.Equal(t, 2, i)
// // }
