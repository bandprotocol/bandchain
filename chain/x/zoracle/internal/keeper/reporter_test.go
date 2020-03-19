package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAddReporterSuccess(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	validatorAddress1 := sdk.ValAddress([]byte("validator1"))
	reporterAddress1 := sdk.AccAddress([]byte("reporter1"))

	err := keeper.AddReporter(ctx, validatorAddress1, reporterAddress1)
	require.Nil(t, err)

	err = keeper.AddReporter(ctx, validatorAddress1, reporterAddress1)
	require.NotNil(t, err)

	reporter := keeper.CheckReporter(ctx, validatorAddress1, reporterAddress1)
	require.True(t, reporter)
}

func TestRemoveReporterSuccess(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	validatorAddress1 := sdk.ValAddress([]byte("validator1"))
	reporterAddress1 := sdk.AccAddress([]byte("reporter1"))

	err := keeper.AddReporter(ctx, validatorAddress1, reporterAddress1)
	require.Nil(t, err)

	err = keeper.RemoveReporter(ctx, validatorAddress1, reporterAddress1)
	require.Nil(t, err)

	reporter := keeper.CheckReporter(ctx, validatorAddress1, reporterAddress1)
	require.False(t, reporter)
}

func TestRemoveReporterFail(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	validatorAddress1 := sdk.ValAddress([]byte("validator1"))
	reporterAddress1 := sdk.AccAddress([]byte("reporter1"))

	err := keeper.RemoveReporter(ctx, validatorAddress1, reporterAddress1)
	require.NotNil(t, err)
}
