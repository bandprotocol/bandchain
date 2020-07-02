package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/testapp"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func newDefaultRequest() types.Request {
	return types.NewRequest(
		1,
		[]byte("calldata"),
		[]sdk.ValAddress{testapp.Validator1.ValAddress, testapp.Validator2.ValAddress},
		2,
		0,
		1581503227,
		"clientID",
		[]types.RawRequest{types.NewRawRequest(42, 1, []byte("calldata"))},
	)
}

func TestHasReport(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	// We should not have a report to request ID 42 from Alice without setting it.
	require.False(t, k.HasReport(ctx, 42, testapp.Alice.ValAddress))
	// After we set it, we should be able to find it.
	k.SetReport(ctx, 42, types.NewReport(testapp.Alice.ValAddress, true, nil))
	require.True(t, k.HasReport(ctx, 42, testapp.Alice.ValAddress))
}

func TestAddReportSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.NoError(t, err)
	require.Equal(t, []types.Report{types.NewReport(testapp.Validator1.ValAddress, true, []types.RawReport{
		types.NewRawReport(2, 1, []byte("data1/1")),
		types.NewRawReport(10, 0, []byte("data2/1")),
	})}, k.GetReports(ctx, 1))
}

func TestReportOnInvalidRequest(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 2, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.Error(t, err)
}

func TestReportByNotRequestedValidator(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Alice.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.Error(t, err)
}

func TestDuplicateReport(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.NoError(t, err)
	err = k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 0, []byte("new data 1")),
			types.NewRawReport(10, 0, []byte("new data 2")),
		},
	))
	require.Error(t, err)
}

func TestReportInvalidDataSourceCount(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
		},
	))
	require.Error(t, err)
}

func TestReportInvalidExternalIDs(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		testapp.Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(11, 1, []byte("data1/1")),
		},
	))
	require.Error(t, err)
}

func TestGetReportCount(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()

	k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Validator1.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Validator2.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Alice.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Bob.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Carol.ValAddress, true, []types.RawReport{}))

	require.Equal(t, uint64(2), k.GetReportCount(ctx, types.RequestID(1)))
	require.Equal(t, uint64(3), k.GetReportCount(ctx, types.RequestID(2)))
}

func TestDeleteReports(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()

	k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Validator1.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(1), types.NewReport(testapp.Validator2.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Alice.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Bob.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(testapp.Carol.ValAddress, true, []types.RawReport{}))

	require.True(t, k.HasReport(ctx, types.RequestID(1), testapp.Validator1.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), testapp.Alice.ValAddress))

	k.DeleteReports(ctx, types.RequestID(1))
	require.False(t, k.HasReport(ctx, types.RequestID(1), testapp.Validator1.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), testapp.Alice.ValAddress))
}
