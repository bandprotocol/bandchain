package app

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

// GenesisState defines a type alias for the Band genesis application state.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	cdc := MakeCodec()
	denom := "uband"
	// Get default genesis states of the modules we are to override.
	authGenesis := auth.DefaultGenesisState()
	stakingGenesis := staking.DefaultGenesisState()
	distrGenesis := distr.DefaultGenesisState()
	mintGenesis := mint.DefaultGenesisState()
	govGenesis := gov.DefaultGenesisState()
	crisisGenesis := crisis.DefaultGenesisState()
	slashingGenesis := slashing.DefaultGenesisState()
	// Override the genesis parameters.
	authGenesis.Params.TxSizeCostPerByte = 5
	stakingGenesis.Params.BondDenom = denom
	stakingGenesis.Params.HistoricalEntries = 1000
	distrGenesis.Params.BaseProposerReward = sdk.NewDecWithPrec(3, 2)   // 3%
	distrGenesis.Params.BonusProposerReward = sdk.NewDecWithPrec(12, 2) // 12%
	mintGenesis.Params.BlocksPerYear = 10519200                         // target 3-second block time
	mintGenesis.Params.MintDenom = denom
	govGenesis.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(denom, sdk.TokensFromConsensusPower(1000)))
	crisisGenesis.ConstantFee = sdk.NewCoin(denom, sdk.TokensFromConsensusPower(10000))
	// slashingGenesis.Params.SignedBlocksWindow = 30000                         // approximately 1 day
	// slashingGenesis.Params.MinSignedPerWindow = sdk.NewDecWithPrec(5, 2)      // 5%
	// slashingGenesis.Params.DowntimeJailDuration = 60 * 10 * time.Second       // 10 minutes
	slashingGenesis.Params.SignedBlocksWindow = 50
	slashingGenesis.Params.MinSignedPerWindow = sdk.NewDecWithPrec(5, 2) // 5%
	slashingGenesis.Params.DowntimeJailDuration = 1 * time.Minute
	slashingGenesis.Params.SlashFractionDoubleSign = sdk.NewDecWithPrec(5, 2) // 5%
	slashingGenesis.Params.SlashFractionDowntime = sdk.NewDecWithPrec(1, 4)   // 0.01%
	stakingGenesis.Params.UnbondingTime = 2 * time.Minute
	govGenesis.DepositParams.MaxDepositPeriod = 5 * time.Minute
	govGenesis.VotingParams.VotingPeriod = 5 * time.Minute
	return GenesisState{
		genutil.ModuleName:  genutil.AppModuleBasic{}.DefaultGenesis(),
		auth.ModuleName:     cdc.MustMarshalJSON(authGenesis),
		bank.ModuleName:     bank.AppModuleBasic{}.DefaultGenesis(),
		supply.ModuleName:   supply.AppModuleBasic{}.DefaultGenesis(),
		staking.ModuleName:  cdc.MustMarshalJSON(stakingGenesis),
		mint.ModuleName:     cdc.MustMarshalJSON(mintGenesis),
		distr.ModuleName:    cdc.MustMarshalJSON(distrGenesis),
		gov.ModuleName:      cdc.MustMarshalJSON(govGenesis),
		crisis.ModuleName:   cdc.MustMarshalJSON(crisisGenesis),
		slashing.ModuleName: cdc.MustMarshalJSON(slashingGenesis),
		upgrade.ModuleName:  upgrade.AppModuleBasic{}.DefaultGenesis(),
		evidence.ModuleName: evidence.AppModuleBasic{}.DefaultGenesis(),
		oracle.ModuleName:   oracle.AppModuleBasic{}.DefaultGenesis(),
	}
}
