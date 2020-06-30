package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestHasReport(t *testing.T) {
	_, ctx, k := createTestInput()
	// We should not have a report to request ID 42 from Alice without setting it.
	require.False(t, k.HasReport(ctx, 42, Alice.ValAddress))
	// After we set it, we should be able to find it.
	k.SetReport(ctx, 42, types.NewReport(Alice.ValAddress, true, nil))
	require.True(t, k.HasReport(ctx, 42, Alice.ValAddress))
}

func TestAddReportSuccess(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.NoError(t, err)
	require.Equal(t, []types.Report{types.NewReport(Validator1.ValAddress, true, []types.RawReport{
		types.NewRawReport(2, 1, []byte("data1/1")),
		types.NewRawReport(10, 0, []byte("data2/1")),
	})}, k.GetReports(ctx, 1))
}

func TestReportOnInvalidRequest(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 2, types.NewReport(
		Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.Error(t, err)
}

func TestReportByNotRequestedValidator(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		Alice.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.Error(t, err)
}

func TestDuplicateReport(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.NoError(t, err)
	err = k.AddReport(ctx, 1, types.NewReport(
		Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 0, []byte("new data 1")),
			types.NewRawReport(10, 0, []byte("new data 2")),
		},
	))
	require.Error(t, err)
}

func TestReportInvalidDataSourceCount(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
		},
	))
	require.Error(t, err)
}

func TestReportInvalidExternalIDs(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequests = []types.RawRequest{
		types.NewRawRequest(2, 1, []byte("calldata")),
		types.NewRawRequest(10, 1, []byte("calldata")),
	}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		Validator1.ValAddress, true, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(11, 1, []byte("data1/1")),
		},
	))
	require.Error(t, err)
}

func TestGetReportCount(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetReport(ctx, types.RequestID(1), types.NewReport(Validator1.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(1), types.NewReport(Validator2.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Alice.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Bob.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Carol.ValAddress, true, []types.RawReport{}))

	require.Equal(t, uint64(2), k.GetReportCount(ctx, types.RequestID(1)))
	require.Equal(t, uint64(3), k.GetReportCount(ctx, types.RequestID(2)))
}

func TestDeleteReports(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetReport(ctx, types.RequestID(1), types.NewReport(Validator1.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(1), types.NewReport(Validator2.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Alice.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Bob.ValAddress, true, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Carol.ValAddress, true, []types.RawReport{}))

	require.True(t, k.HasReport(ctx, types.RequestID(1), Validator1.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), Alice.ValAddress))

	k.DeleteReports(ctx, types.RequestID(1))
	require.False(t, k.HasReport(ctx, types.RequestID(1), Validator1.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), Alice.ValAddress))
}
