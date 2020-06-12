package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCheckSelfReporter(t *testing.T) {
	_, ctx, k := createTestInput()
	// Owner must always be a reporter of himself.
	require.True(t, k.IsReporter(ctx, Owner.ValAddress, Owner.Address))
}

func TestAddReporter(t *testing.T) {
	_, ctx, k := createTestInput()
	// Before we do anything, Bob and Carol must not be a reporter of Alice.
	require.False(t, k.IsReporter(ctx, Alice.ValAddress, Bob.Address))
	// Adds Bob as a reporter Alice. IsReporter should return true for Bob, false for Carol.
	err := k.AddReporter(ctx, Alice.ValAddress, Bob.Address)
	require.Nil(t, err)
	require.True(t, k.IsReporter(ctx, Alice.ValAddress, Bob.Address))
	require.False(t, k.IsReporter(ctx, Alice.ValAddress, Carol.Address))
	// We should get an error if we try to add Bob again.
	err = k.AddReporter(ctx, Alice.ValAddress, Bob.Address)
	require.NotNil(t, err)
}

func TestRemoveReporter(t *testing.T) {
	_, ctx, k := createTestInput()
	// Removing Bob from Alice's reporter list should error as Bob is not a reporter of Alice.
	err := k.RemoveReporter(ctx, Alice.ValAddress, Bob.Address)
	require.NotNil(t, err)
	// Adds Bob as the reporter of Alice. We now should be able to remove Bob, but not Carol.
	err = k.AddReporter(ctx, Alice.ValAddress, Bob.Address)
	require.Nil(t, err)
	err = k.RemoveReporter(ctx, Alice.ValAddress, Bob.Address)
	require.Nil(t, err)
	err = k.RemoveReporter(ctx, Alice.ValAddress, Carol.Address)
	require.NotNil(t, err)
	// By the end of everything, no one should be a reporter of Alice.
	require.False(t, k.IsReporter(ctx, Alice.ValAddress, Bob.Address))
	require.False(t, k.IsReporter(ctx, Alice.ValAddress, Carol.Address))
}

func TestGetReporters(t *testing.T) {
	_, ctx, k := createTestInput()
	aliceReporters := []sdk.AccAddress{
		Alice.Address, //self reporter of validator
		Bob.Address,
		Carol.Address,
	}

	// Adds Alice validator
	err := k.AddReporter(ctx, Alice.ValAddress, aliceReporters[1])
	require.NoError(t, err)
	err = k.AddReporter(ctx, Alice.ValAddress, aliceReporters[2])
	require.NoError(t, err)

	err = k.AddReporter(ctx, Bob.ValAddress, Alice.Address)
	require.NoError(t, err)

	reporters := k.GetReporters(ctx, Alice.ValAddress)
	require.Equal(t, len(reporters), len(aliceReporters))
}
