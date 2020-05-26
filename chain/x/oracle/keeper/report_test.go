package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/simapp"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestHasReport(t *testing.T) {
	_, ctx, k := createTestInput()
	// We should not have a report to request ID 42 from Alice without setting it.
	require.False(t, k.HasReport(ctx, 42, Alice.ValAddress))
	// After we set it, we should be able to find it.
	k.SetReport(ctx, 42, types.NewReport(Alice.ValAddress, nil))
	require.True(t, k.HasReport(ctx, 42, Alice.ValAddress))
}

func TestAddReportSuccess(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequestIDs = []types.ExternalID{2, 10}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		Validator1.ValAddress, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.NoError(t, err)
	require.Equal(t, []types.Report{types.NewReport(Validator1.ValAddress, []types.RawReport{
		types.NewRawReport(2, 1, []byte("data1/1")),
		types.NewRawReport(10, 0, []byte("data2/1")),
	})}, k.GetReports(ctx, 1))
}

func TestReportOnInvalidRequest(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequestIDs = []types.ExternalID{2, 10}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 2, types.NewReport(
		Validator1.ValAddress, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.Error(t, err)
}

func TestReportByNotRequestedValidator(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequestIDs = []types.ExternalID{2, 10}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		Alice.ValAddress, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.Error(t, err)
}

func TestDuplicateReport(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequestIDs = []types.ExternalID{2, 10}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		Validator1.ValAddress, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(10, 0, []byte("data2/1")),
		},
	))

	require.NoError(t, err)
	err = k.AddReport(ctx, 1, types.NewReport(
		Validator1.ValAddress, []types.RawReport{
			types.NewRawReport(2, 0, []byte("new data 1")),
			types.NewRawReport(10, 0, []byte("new data 2")),
		},
	))
	require.Error(t, err)
}

func TestReportInvalidDataSourceCount(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequestIDs = []types.ExternalID{2, 10}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		Validator1.ValAddress, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
		},
	))
	require.Error(t, err)
}

func TestReportInvalidExternalIDs(t *testing.T) {
	_, ctx, k := createTestInput()
	request := newDefaultRequest()
	request.RawRequestIDs = []types.ExternalID{2, 10}
	k.SetRequest(ctx, 1, request)

	err := k.AddReport(ctx, 1, types.NewReport(
		Validator1.ValAddress, []types.RawReport{
			types.NewRawReport(2, 1, []byte("data1/1")),
			types.NewRawReport(11, 1, []byte("data1/1")),
		},
	))
	require.Error(t, err)
}

func TestGetReportCount(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetReport(ctx, types.RequestID(1), types.NewReport(Validator1.ValAddress, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(1), types.NewReport(Validator2.ValAddress, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Alice.ValAddress, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Bob.ValAddress, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Carol.ValAddress, []types.RawReport{}))

	require.Equal(t, int64(2), k.GetReportCount(ctx, types.RequestID(1)))
	require.Equal(t, int64(3), k.GetReportCount(ctx, types.RequestID(2)))
}

func TestDeleteReports(t *testing.T) {
	_, ctx, k := createTestInput()

	k.SetReport(ctx, types.RequestID(1), types.NewReport(Validator1.ValAddress, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(1), types.NewReport(Validator2.ValAddress, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Alice.ValAddress, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Bob.ValAddress, []types.RawReport{}))
	k.SetReport(ctx, types.RequestID(2), types.NewReport(Carol.ValAddress, []types.RawReport{}))

	require.True(t, k.HasReport(ctx, types.RequestID(1), Validator1.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), Alice.ValAddress))

	k.DeleteReports(ctx, types.RequestID(1))
	require.False(t, k.HasReport(ctx, types.RequestID(1), Validator1.ValAddress))
	require.True(t, k.HasReport(ctx, types.RequestID(2), Alice.ValAddress))
}
func TestUpdateReportInfos(t *testing.T) {
	_, ctx, k := createTestInput()

	ctx = ctx.WithBlockHeight(100)
	request := types.NewRequest(
		types.OracleScriptID(1), []byte("calldata"),
		[]sdk.ValAddress{Validator1.ValAddress, Validator2.ValAddress},
		2, 100, 100, "test", nil, nil,
	)
	k.SetRequest(ctx, types.RequestID(1), request)

	// 2 Validators report
	k.SetReport(ctx, types.RequestID(1), types.NewReport(
		Validator1.ValAddress, []types.RawReport{},
	))

	k.SetReport(ctx, types.RequestID(1), types.NewReport(
		Validator2.ValAddress, []types.RawReport{},
	))

	// Update report info
	k.UpdateReportInfos(ctx, types.RequestID(1))
	cases := []struct {
		validator simapp.Account
		info      types.ValidatorReportInfo
	}{
		{Validator1, types.NewValidatorReportInfo(Validator1.ValAddress, 0)},
		{Validator2, types.NewValidatorReportInfo(Validator2.ValAddress, 0)},
	}
	for _, tc := range cases {
		info := k.GetValidatorReportInfoWithDefault(ctx, tc.validator.ValAddress)
		require.Equal(t, tc.info, info)
	}

	k.SetRequest(ctx, types.RequestID(2), request)

	// Only Validator1 report
	k.SetReport(ctx, types.RequestID(2), types.NewReport(
		Validator1.ValAddress, []types.RawReport{},
	))

	// Update report info
	k.UpdateReportInfos(ctx, types.RequestID(2))
	cases = []struct {
		validator simapp.Account
		info      types.ValidatorReportInfo
	}{
		{Validator1, types.NewValidatorReportInfo(Validator1.ValAddress, 0)},
		{Validator2, types.NewValidatorReportInfo(Validator2.ValAddress, 1)},
	}
	for _, tc := range cases {
		info := k.GetValidatorReportInfoWithDefault(ctx, tc.validator.ValAddress)
		require.Equal(t, tc.info, info)
	}

	k.SetRequest(ctx, types.RequestID(3), request)
	// Only Validator2 report
	k.SetReport(ctx, types.RequestID(3), types.NewReport(
		Validator2.ValAddress, []types.RawReport{},
	))

	// Update report info
	k.UpdateReportInfos(ctx, types.RequestID(3))
	cases = []struct {
		validator simapp.Account
		info      types.ValidatorReportInfo
	}{
		{Validator1, types.NewValidatorReportInfo(Validator1.ValAddress, 1)},
		{Validator2, types.NewValidatorReportInfo(Validator2.ValAddress, 0)},
	}
	for _, tc := range cases {
		info := k.GetValidatorReportInfoWithDefault(ctx, tc.validator.ValAddress)
		require.Equal(t, tc.info, info)
	}

	k.SetRequest(ctx, types.RequestID(4), request)
	// Noone reports
	k.UpdateReportInfos(ctx, types.RequestID(4))
	cases = []struct {
		validator simapp.Account
		info      types.ValidatorReportInfo
	}{
		{Validator1, types.NewValidatorReportInfo(Validator1.ValAddress, 2)},
		{Validator2, types.NewValidatorReportInfo(Validator2.ValAddress, 1)},
	}
	for _, tc := range cases {
		info := k.GetValidatorReportInfoWithDefault(ctx, tc.validator.ValAddress)
		require.Equal(t, tc.info, info)
	}
}

func TestGetJailedUpdateReportInfos(t *testing.T) {
	app, ctx, k := createTestInput()
	k.SetParam(ctx, types.KeyMaxConsecutiveMisses, 3)

	ctx = ctx.WithBlockHeight(100)
	request := types.NewRequest(
		types.OracleScriptID(1), []byte("calldata"),
		[]sdk.ValAddress{Validator1.ValAddress, Validator2.ValAddress},
		2, 100, 100, "test", nil, nil,
	)
	k.SetRequest(ctx, types.RequestID(1), request)

	// 2 Validators report
	k.SetReport(ctx, types.RequestID(1), types.NewReport(
		Validator1.ValAddress, []types.RawReport{},
	))

	k.SetReport(ctx, types.RequestID(1), types.NewReport(
		Validator2.ValAddress, []types.RawReport{},
	))
	k.UpdateReportInfos(ctx, types.RequestID(1))

	// Validator2 reported downed
	for i := 2; i <= 5; i++ {
		k.SetRequest(ctx, types.RequestID(i), request)
		k.SetReport(ctx, types.RequestID(i), types.NewReport(
			Validator1.ValAddress, []types.RawReport{},
		))
		k.UpdateReportInfos(ctx, types.RequestID(i))
	}

	cases := []struct {
		validator simapp.Account
		info      types.ValidatorReportInfo
		jailed    bool
	}{
		{Validator1, types.NewValidatorReportInfo(Validator1.ValAddress, 0), false},
		{Validator2, types.NewValidatorReportInfo(Validator2.ValAddress, 0), true},
	}
	for _, tc := range cases {
		info := k.GetValidatorReportInfoWithDefault(ctx, tc.validator.ValAddress)
		require.Equal(t, tc.info, info)
		validator := app.StakingKeeper.Validator(ctx, tc.validator.ValAddress)
		require.Equal(t, tc.jailed, validator.IsJailed())
	}
}
