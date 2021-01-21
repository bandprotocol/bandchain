package oracle_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/testapp"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

func fromHex(hexStr string) []byte {
	res, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return res
}

func TestRollingSeedCorrect(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false)
	// Initially rolling seed should be all zeros.
	require.Equal(t, fromHex("0000000000000000000000000000000000000000000000000000000000000000"), k.GetRollingSeed(ctx))
	// Every begin block, the rolling seed should get updated.
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash: fromHex("0100000000000000000000000000000000000000000000000000000000000000"),
	})
	require.Equal(t, fromHex("0000000000000000000000000000000000000000000000000000000000000001"), k.GetRollingSeed(ctx))
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash: fromHex("0200000000000000000000000000000000000000000000000000000000000000"),
	})
	require.Equal(t, fromHex("0000000000000000000000000000000000000000000000000000000000000102"), k.GetRollingSeed(ctx))
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash: fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
	})
	require.Equal(t, fromHex("00000000000000000000000000000000000000000000000000000000000102ff"), k.GetRollingSeed(ctx))
}

func TestAllocateTokensCalledOnBeginBlock(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false)
	votes := []abci.VoteInfo{{
		Validator:       abci.Validator{Address: testapp.Validator1.PubKey.Address(), Power: 70},
		SignedLastBlock: true,
	}, {
		Validator:       abci.Validator{Address: testapp.Validator2.PubKey.Address(), Power: 30},
		SignedLastBlock: true,
	}}
	// Set collected fee to 100uband + 70% oracle reward proportion + disable minting inflation.
	// NOTE: we intentionally keep ctx.BlockHeight = 0, so distr's AllocateTokens doesn't get called.
	feeCollector := app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	feeCollector.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("uband", 100)))
	app.AccountKeeper.SetAccount(ctx, feeCollector)
	mintParams := app.MintKeeper.GetParams(ctx)
	mintParams.InflationMin = sdk.ZeroDec()
	mintParams.InflationMax = sdk.ZeroDec()
	app.MintKeeper.SetParams(ctx, mintParams)
	k.SetParam(ctx, types.KeyOracleRewardPercentage, 70)
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 100)), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	// If there are no validators active, Calling begin block should be no-op.
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash:           fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		LastCommitInfo: abci.LastCommitInfo{Votes: votes},
	})
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 100)), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	// 1 validator active, begin block should take 70% of the fee. 2% of that goes to comm pool.
	k.Activate(ctx, testapp.Validator2.ValAddress)
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash:           fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		LastCommitInfo: abci.LastCommitInfo{Votes: votes},
	})
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 30)), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 70)), app.SupplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins())
	// 100*70%*2% = 1.4uband
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(14, 1)}}, app.DistrKeeper.GetFeePool(ctx).CommunityPool)
	// 0uband
	require.Equal(t, sdk.DecCoins(nil), app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator1.ValAddress))
	// 100*70%*98% = 68.6uband
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(686, 1)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator2.ValAddress))
	// 2 validators active now. 70% of the remaining fee pool will be split 3 ways (comm pool + val1 + val2).
	k.Activate(ctx, testapp.Validator1.ValAddress)
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash:           fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		LastCommitInfo: abci.LastCommitInfo{Votes: votes},
	})
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 9)), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 91)), app.SupplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins())
	// 1.4uband + 30*70%*2% = 1.82uband
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(182, 2)}}, app.DistrKeeper.GetFeePool(ctx).CommunityPool)
	// 30*70%*98%*70% = 14.406uband
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(14406, 3)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator1.ValAddress))
	// 68.6uband + 30*70%*98%*30% = 74.774uband
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(74774, 3)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator2.ValAddress))
	// 1 validator becomes in active, and will not get reward this time.
	k.MissReport(ctx, testapp.Validator2.ValAddress, testapp.ParseTime(100))
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash:           fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		LastCommitInfo: abci.LastCommitInfo{Votes: votes},
	})
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 3)), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 97)), app.SupplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins())
	// 1.82uband + 6*2% = 1.82uband
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(194, 2)}}, app.DistrKeeper.GetFeePool(ctx).CommunityPool)
	// 14.406uband + 6*98% = 20.286uband
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(20286, 3)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator1.ValAddress))
	// 74.774uband
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(74774, 3)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator2.ValAddress))
}

func TestAllocateTokensWithDistrAllocateTokens(t *testing.T) {
	app, ctx, k := testapp.CreateTestInput(false)
	ctx = ctx.WithBlockHeight(10) // Set block height to ensure distr's AllocateTokens gets called.
	votes := []abci.VoteInfo{{
		Validator:       abci.Validator{Address: testapp.Validator1.PubKey.Address(), Power: 70},
		SignedLastBlock: true,
	}, {
		Validator:       abci.Validator{Address: testapp.Validator2.PubKey.Address(), Power: 30},
		SignedLastBlock: true,
	}}
	// Set collected fee to 100uband + 70% oracle reward proportion + disable minting inflation.
	feeCollector := app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
	feeCollector.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("uband", 50)))
	app.AccountKeeper.SetAccount(ctx, feeCollector)
	mintParams := app.MintKeeper.GetParams(ctx)
	mintParams.InflationMin = sdk.ZeroDec()
	mintParams.InflationMax = sdk.ZeroDec()
	app.MintKeeper.SetParams(ctx, mintParams)
	k.SetParam(ctx, types.KeyOracleRewardPercentage, 70)
	// Set block proposer to Validator2, who will receive 5% bonus.
	app.DistrKeeper.SetPreviousProposerConsAddr(ctx, testapp.Validator2.Address.Bytes())
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 50)), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	// Only validator1 active. After we call begin block:
	//   35uband = 70% go to oracle pool
	//     0.7uband (2%) go to community pool
	//     34.3uband go to validator1 (active)
	//   15uband = 30% go to distr pool
	//     0.3uband (2%) go to community pool
	//     2.25uband (15%) go to validator2 (proposer)
	//     12.45uband split among voters
	//        8.715uband (70%) go to validator1
	//        3.735uband (30%) go to validator2
	// In summary
	//   Community pool: 0.7 + 0.3 = 1
	//   Validator1: 34.3 + 8.715 = 43.015
	//   Validator2: 2.25 + 3.735 = 5.985
	k.Activate(ctx, testapp.Validator1.ValAddress)
	app.BeginBlocker(ctx, abci.RequestBeginBlock{
		Hash:           fromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		LastCommitInfo: abci.LastCommitInfo{Votes: votes},
	})
	require.Equal(t, sdk.Coins(nil), app.SupplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName).GetCoins())
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 50)), app.SupplyKeeper.GetModuleAccount(ctx, distribution.ModuleName).GetCoins())
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDec(1)}}, app.DistrKeeper.GetFeePool(ctx).CommunityPool)
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(43015, 3)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator1.ValAddress))
	require.Equal(t, sdk.DecCoins{{Denom: "uband", Amount: sdk.NewDecWithPrec(5985, 3)}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, testapp.Validator2.ValAddress))
}
