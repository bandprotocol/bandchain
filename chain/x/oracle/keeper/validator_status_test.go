package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/testapp"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

func defaultVotes() []abci.VoteInfo {
	return []abci.VoteInfo{{
		Validator: abci.Validator{
			Address: testapp.Validator1.PubKey.Address(),
			Power:   70,
		},
		SignedLastBlock: true,
	}, {
		Validator: abci.Validator{
			Address: testapp.Validator2.PubKey.Address(),
			Power:   20,
		},
		SignedLastBlock: true,
	}, {
		Validator: abci.Validator{
			Address: testapp.Validator3.PubKey.Address(),
			Power:   10,
		},
		SignedLastBlock: true,
	}}
}

func TestAllocateTokenNoActiveValidators(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false)
	// Set collected fee to 1000000uband and 70% oracle reward proportion.
	feeCollector := app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	feeCollector.SetCoins(Coins1000000uband)
	app.AccountKeeper.SetAccount(ctx, feeCollector)
	k.SetParam(ctx, types.KeyOracleRewardPercentage, 70)
	require.Equal(t, Coins1000000uband, app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	// No active oracle validators so nothing should happen.
	k.AllocateTokens(ctx, defaultVotes())
	require.Equal(t, Coins1000000uband, app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	require.Equal(t, sdk.Coins(nil), app.SupplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins())
}

func TestAllocateTokensOneActive(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false)
	// Set collected fee to 1000000uband + 70% oracle reward proportion.
	feeCollector := app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	feeCollector.SetCoins(Coins1000000uband)
	app.AccountKeeper.SetAccount(ctx, feeCollector)
	k.SetParam(ctx, types.KeyOracleRewardPercentage, 70)
	require.Equal(t, Coins1000000uband, app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	// From 70% of fee, 2% should go to community pool, the rest goes to the only active validator.
	k.Activate(ctx, testapp.Validator2.ValAddress)
	k.AllocateTokens(ctx, defaultVotes())
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 300000)), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 700000)), app.SupplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins())
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(14000)}}, app.DistrKeeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, sdk.DecCoins(nil), app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator1.ValAddress))
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(686000)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator2.ValAddress))
	require.Equal(t, sdk.DecCoins(nil), app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator3.ValAddress))
}

func TestAllocateTokensAllActive(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(true)
	// Set collected fee to 1000000uband + 70% oracle reward proportion.
	feeCollector := app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	feeCollector.SetCoins(Coins1000000uband)
	app.AccountKeeper.SetAccount(ctx, feeCollector)
	k.SetParam(ctx, types.KeyOracleRewardPercentage, 70)
	require.Equal(t, Coins1000000uband, app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	// From 70% of fee, 2% should go to community pool, the rest get split to validators.
	k.AllocateTokens(ctx, defaultVotes())
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 300000)), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 700000)), app.SupplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins())
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(14000)}}, app.DistrKeeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(480200)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator1.ValAddress))
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(137200)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator2.ValAddress))
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(68600)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator3.ValAddress))
}

func TestGetDefaultValidatorStatus(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	vs := k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress)
	require.Equal(t, types.NewValidatorStatus(false, time.Time{}), vs)
}

func TestGetSetValidatorStatus(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	// After setting status of the 1st validator, we should be able to get it back.
	k.SetValidatorStatus(ctx, testapp.Validator1.ValAddress, types.NewValidatorStatus(true, now))
	vs := k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress)
	require.Equal(t, types.NewValidatorStatus(true, now), vs)
	vs = k.GetValidatorStatus(ctx, testapp.Validator2.ValAddress)
	require.Equal(t, types.NewValidatorStatus(false, time.Time{}), vs)
}

func TestActivateValidatorOK(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	ctx = ctx.WithBlockTime(now)
	err := k.Activate(ctx, testapp.Validator1.ValAddress)
	require.NoError(t, err)
	vs := k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress)
	require.Equal(t, types.NewValidatorStatus(true, now), vs)
	vs = k.GetValidatorStatus(ctx, testapp.Validator2.ValAddress)
	require.Equal(t, types.NewValidatorStatus(false, time.Time{}), vs)
}

func TestFailActivateAlreadyActive(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	ctx = ctx.WithBlockTime(now)
	err := k.Activate(ctx, testapp.Validator1.ValAddress)
	require.NoError(t, err)
	err = k.Activate(ctx, testapp.Validator1.ValAddress)
	require.Error(t, err)
}

func TestFailActivateTooSoon(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	// Set validator to be inactive just now.
	k.SetValidatorStatus(ctx, testapp.Validator1.ValAddress, types.NewValidatorStatus(false, now))
	// You can't activate until it's been at least InactivePenaltyDuration nanosec.
	penaltyDuration := k.GetParam(ctx, types.KeyInactivePenaltyDuration)
	require.Error(t, k.Activate(ctx.WithBlockTime(now), testapp.Validator1.ValAddress))
	require.Error(t, k.Activate(ctx.WithBlockTime(now.Add(time.Duration(penaltyDuration/2))), testapp.Validator1.ValAddress))
	// So far there must be no changes to the validator's status.
	vs := k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress)
	require.Equal(t, types.NewValidatorStatus(false, now), vs)
	// Now the time has come.
	require.NoError(t, k.Activate(ctx.WithBlockTime(now.Add(time.Duration(penaltyDuration))), testapp.Validator1.ValAddress))
	vs = k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress)
	require.Equal(t, types.NewValidatorStatus(true, now.Add(time.Duration(penaltyDuration))), vs)
}

func TestMissReportSuccess(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	next := now.Add(time.Duration(10))
	k.SetValidatorStatus(ctx, testapp.Validator1.ValAddress, types.NewValidatorStatus(true, now))
	k.MissReport(ctx.WithBlockTime(next), testapp.Validator1.ValAddress, next)
	vs := k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress)
	require.Equal(t, types.NewValidatorStatus(false, next), vs)
}

func TestMissReportTooSoonNoop(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	prev := time.Now().UTC()
	now := prev.Add(time.Duration(10))
	k.SetValidatorStatus(ctx, testapp.Validator1.ValAddress, types.NewValidatorStatus(true, now))
	k.MissReport(ctx.WithBlockTime(prev), testapp.Validator1.ValAddress, prev)
	vs := k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress)
	require.Equal(t, types.NewValidatorStatus(true, now), vs)
}

func TestMissReportAlreadyInactiveNoop(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(false)
	now := time.Now().UTC()
	next := now.Add(time.Duration(10))
	k.SetValidatorStatus(ctx, testapp.Validator1.ValAddress, types.NewValidatorStatus(false, now))
	k.MissReport(ctx.WithBlockTime(next), testapp.Validator1.ValAddress, next)
	vs := k.GetValidatorStatus(ctx, testapp.Validator1.ValAddress)
	require.Equal(t, types.NewValidatorStatus(false, now), vs)
}
