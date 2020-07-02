package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/testapp"
)

func TestCheckSelfReporter(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	// Owner must always be a reporter of himself.
	require.True(t, k.IsReporter(ctx, testapp.Owner.ValAddress, testapp.Owner.Address))
}

func TestAddReporter(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	// Before we do anything, Bob and Carol must not be a reporter of Alice.
	require.False(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address))
	// Adds Bob as a reporter Alice. IsReporter should return true for Bob, false for Carol.
	err := k.AddReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.Nil(t, err)
	require.True(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address))
	require.False(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Carol.Address))
	// We should get an error if we try to add Bob again.
	err = k.AddReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.NotNil(t, err)
}

func TestRemoveReporter(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	// Removing Bob from Alice's reporter list should error as Bob is not a reporter of Alice.
	err := k.RemoveReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.NotNil(t, err)
	// Adds Bob as the reporter of Alice. We now should be able to remove Bob, but not Carol.
	err = k.AddReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.Nil(t, err)
	err = k.RemoveReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address)
	require.Nil(t, err)
	err = k.RemoveReporter(ctx, testapp.Alice.ValAddress, testapp.Carol.Address)
	require.NotNil(t, err)
	// By the end of everything, no one should be a reporter of Alice.
	require.False(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Bob.Address))
	require.False(t, k.IsReporter(ctx, testapp.Alice.ValAddress, testapp.Carol.Address))
}

func TestGetReporters(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput()
	aliceReporters := []sdk.AccAddress{
		testapp.Alice.Address, //self reporter of validator
		testapp.Bob.Address,
		testapp.Carol.Address,
	}

	// Adds Alice validator
	err := k.AddReporter(ctx, testapp.Alice.ValAddress, aliceReporters[1])
	require.NoError(t, err)
	err = k.AddReporter(ctx, testapp.Alice.ValAddress, aliceReporters[2])
	require.NoError(t, err)

	err = k.AddReporter(ctx, testapp.Bob.ValAddress, testapp.Alice.Address)
	require.NoError(t, err)

	reporters := k.GetReporters(ctx, testapp.Alice.ValAddress)
	require.Equal(t, len(reporters), len(aliceReporters))
}
