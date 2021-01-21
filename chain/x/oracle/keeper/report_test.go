package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/testapp"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

func defaultRequest() types.Request {
	return types.NewRequest(
		1, BasicCalldata,
		[]sdk.ValAddress{testapp.Validator1.ValAddress, testapp.Validator2.ValAddress},
		2, 0, testapp.ParseTime(0),
		BasicClientID, []types.RawRequest{
			types.NewRawRequest(42, 1, BasicCalldata),
			types.NewRawRequest(43, 2, BasicCalldata),
		},
	)
}

func TestHasReport(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// We should not have a report to request ID 42 from Alice without setting it.
	require.False(t, k.HasReport(ctx, 42, testapp.Alice.ValAddress))
	// After we set it, we should be able to find it.
	k.SetReport(ctx, 42, types.NewReport(testapp.Alice.ValAddress, true, nil))
	require.True(t, k.HasReport(ctx, 42, testapp.Alice.ValAddress))
}

func TestAddReportSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	k.SetRequest(ctx, 1, defaultRequest())
	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		},
	))
	require.NoError(t, err)
	require.Equal(t, []types.Report{
		types.NewReport(testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		}),
	}, k.GetReports(ctx, 1))
}

func TestReportOnNonExistingRequest(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		},
	))
	require.Error(t, err)
}

func TestReportByNotRequestedValidator(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	k.SetRequest(ctx, 1, defaultRequest())
	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Alice.ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		},
	))
	require.Error(t, err)
}

func TestDuplicateReport(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	k.SetRequest(ctx, 1, defaultRequest())
	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		},
	))
	require.NoError(t, err)
	err = k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(43, 1, []byte("data2/1")),
		},
	))
	require.Error(t, err)
}

func TestReportInvalidDataSourceCount(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	k.SetRequest(ctx, 1, defaultRequest())
	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
		},
	))
	require.Error(t, err)
}

func TestReportInvalidExternalIDs(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	k.SetRequest(ctx, 1, defaultRequest())
	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(42, 0, []byte("data1/1")),
			types.NewRawReport(44, 1, []byte("data2/1")), // BAD EXTERNAL ID!
		},
	))
	require.Error(t, err)
}

func TestGetReportCount(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// We start by setting some aribrary reports.
	k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Alice.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Bob.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Alice.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Bob.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Carol.ValAddress, true, []types.RawReport{}))
	// GetReportCount should return the correct values.
	require.Equal(t, uint64(2), k.GetReportCount(ctx, types.RequestID(1)))
	require.Equal(t, uint64(3), k.GetReportCount(ctx, types.RequestID(2)))
}

func TestDeleteReports(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)
	// We start by setting some arbitrary reports.
	k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Alice.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Bob.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Alice.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Bob.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Carol.ValAddress, true, []types.RawReport{}))
	// All reports should exist on the state.
	require.True(t, k.HasReport(ctx, types.RequestID(1), testapp.Alice.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(1), testapp.Bob.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), testapp.Alice.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), testapp.Bob.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), testapp.Carol.ValAddress))
	// After we delete reports related to request#1, they must disappear.
	k.DeleteReports(ctx, types.RequestID(1))
	require.False(t, k.HasReport(ctx, types.RequestID(1), testapp.Alice.ValAddress))
	require.False(t, k.HasReport(ctx, types.RequestID(1), testapp.Bob.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), testapp.Alice.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), testapp.Bob.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), testapp.Carol.ValAddress))
}
