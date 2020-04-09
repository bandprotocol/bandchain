package oracle

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keep "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestNewExecutionEnvironment(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	_, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.NotNil(t, err)

	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))}, 1, 0, 0, 100, 10000, "clientID",
	))

	_, err = NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
}

func TestGetCurrentRequestID(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))}, 1, 0, 0, 100, 10000, "clientID",
	))

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(1), env.GetCurrentRequestID())
}

func TestGetRequestedValidatorCount(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"),
		[]sdk.ValAddress{sdk.ValAddress([]byte("val1")), sdk.ValAddress([]byte("val2"))},
		1, 0, 0, 100, 10000, "clientID",
	))

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(2), env.GetRequestedValidatorCount())
}

func TestGetSufficientValidatorCount(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"),
		[]sdk.ValAddress{
			sdk.ValAddress([]byte("val1")),
			sdk.ValAddress([]byte("val2")),
			sdk.ValAddress([]byte("val3")),
			sdk.ValAddress([]byte("val4")),
		},
		3, 0, 0, 100, 10000, "clientID",
	))

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(3), env.GetSufficientValidatorCount())
}

func TestGetReceivedValidatorCount(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"),
		[]sdk.ValAddress{sdk.ValAddress([]byte("val1")), sdk.ValAddress([]byte("val2"))},
		1, 0, 0, 100, 10000, "clientID",
	))

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(0), env.GetReceivedValidatorCount())

	keeper.AddReport(ctx, 1, []types.RawDataReportWithID{}, sdk.ValAddress([]byte("val1")), sdk.AccAddress([]byte("val1")))

	env, err = NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(1), env.GetReceivedValidatorCount())

}

func TestGetPrepareBlockTime(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))},
		1, 20, 1581589790, 100, 10000, "clientID",
	))

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(1581589790), env.GetPrepareBlockTime())
}

func TestGetAggregateBlockTime(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))},
		1, 0, 0, 100, 10000, "clientID",
	))

	ctx = ctx.WithBlockTime(time.Unix(int64(1581589790), 0))
	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(0), env.GetAggregateBlockTime())

	// Add received validator
	err = keeper.AddReport(ctx, 1, []types.RawDataReportWithID{}, sdk.ValAddress([]byte("val1")), sdk.AccAddress([]byte("val1")))
	require.Nil(t, err)

	// After report is greater or equal SufficientValidatorCount, it will resolve in current block time.
	env, err = NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	require.Equal(t, int64(1581589790), env.GetAggregateBlockTime())
}

func TestGetValidatorPubKey(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	pubStr := []string{
		"03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f",
		"03f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909",
	}
	validatorAddress1 := keep.SetupTestValidator(
		ctx,
		keeper,
		pubStr[0],
		10,
	)
	validatorAddress2 := keep.SetupTestValidator(
		ctx,
		keeper,
		pubStr[1],
		10,
	)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{validatorAddress1, validatorAddress2},
		1, 0, 0, 100, 10000, "clientID",
	))

	env, errSDK := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, errSDK)

	addr1, err := env.GetValidatorAddress(0)
	require.Nil(t, err)
	require.Equal(t, validatorAddress1, sdk.ValAddress(addr1))

	addr2, err := env.GetValidatorAddress(1)
	require.Nil(t, err)
	require.Equal(t, validatorAddress2, sdk.ValAddress(addr2))

	_, err = env.GetValidatorAddress(2)
	require.NotNil(t, err)

	_, err = env.GetValidatorAddress(-1)
	require.NotNil(t, err)
}

func TestRequestExternalData(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	// Set Request
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))},
		1, 0, 0, 100, 10000, "clientID",
	))

	// Set Datasource
	dataSource := types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source",
		"description",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 10)),
		[]byte("executable"),
	)
	keeper.SetDataSource(ctx, 1, dataSource)

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)
	envErr := env.RequestExternalData(1, 42, []byte("prepare32"))
	require.Nil(t, envErr)
	err = env.SaveRawDataRequests(ctx, keeper)
	require.Nil(t, err)

	rawRequest, err := keeper.GetRawDataRequest(ctx, 1, 42)
	require.Nil(t, err)
	require.Equal(t, types.NewRawDataRequest(1, []byte("prepare32")), rawRequest)
}

func TestRequestExternalDataExceedMaxDataSourceCountPerRequest(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)

	// Set Request
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"), []sdk.ValAddress{sdk.ValAddress([]byte("val1"))},
		1, 0, 0, 100, 10000, "clientID",
	))

	// Set Datasource
	dataSource := types.NewDataSource(
		sdk.AccAddress([]byte("owner")),
		"data_source",
		"description",
		sdk.NewCoins(sdk.NewInt64Coin("uband", 10)),
		[]byte("executable"),
	)
	keeper.SetDataSource(ctx, 1, dataSource)

	// Set MaxDataSourceCountPerRequest to 5
	keeper.SetParam(ctx, KeyMaxDataSourceCountPerRequest, 5)
	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)

	reqErr := env.RequestExternalData(1, 41, []byte("prepare32"))
	require.Nil(t, reqErr)
	reqErr = env.RequestExternalData(1, 42, []byte("prepare32"))
	require.Nil(t, reqErr)
	reqErr = env.RequestExternalData(1, 43, []byte("prepare32"))
	require.Nil(t, reqErr)
	reqErr = env.RequestExternalData(1, 44, []byte("prepare32"))
	require.Nil(t, reqErr)
	reqErr = env.RequestExternalData(1, 45, []byte("prepare32"))
	require.Nil(t, reqErr)
	reqErr = env.RequestExternalData(1, 46, []byte("prepare32"))
	require.NotNil(t, reqErr)

	envErr := env.SaveRawDataRequests(ctx, keeper)
	require.Nil(t, envErr)
}

func TestGetExternalData(t *testing.T) {
	ctx, keeper := keep.CreateTestInput(t, false)
	keeper.SetRequest(ctx, 1, types.NewRequest(
		1, []byte("calldata"),
		[]sdk.ValAddress{sdk.ValAddress([]byte("val1")), sdk.ValAddress([]byte("val2"))},
		1, 0, 0, 100, 10000, "clientID",
	))

	keeper.SetRawDataReport(
		ctx,
		1,
		42,
		sdk.ValAddress([]byte("val1")),
		types.NewRawDataReport(42, []byte("data42")),
	)

	env, err := NewExecutionEnvironment(ctx, keeper, 1)
	require.Nil(t, err)

	err = env.LoadRawDataReports(ctx, keeper)
	require.Nil(t, err)
	// Get report from reported validator
	report, statusCode, envErr := env.GetExternalData(42, 0)
	require.Nil(t, envErr)
	require.Equal(t, []byte("data42"), report)
	require.Equal(t, uint8(42), statusCode)

	// Get report from missing validator
	_, _, envErr = env.GetExternalData(42, 1)
	require.NotNil(t, envErr)
	require.EqualError(t, envErr, "Unable to find raw data report with request ID (1) external ID (42) from (bandvaloper1weskcvsfgndm9): ItemNotFound")

	// Get report from invalid validator index
	_, _, envErr = env.GetExternalData(42, 2)
	require.NotNil(t, envErr, "validator out of range")

	// Get report from invalid validator index
	_, _, envErr = env.GetExternalData(42, -2)
	require.NotNil(t, envErr, "validator out of range")
}
