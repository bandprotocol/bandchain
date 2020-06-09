package app

import (
	"encoding/json"
	"time"

	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

// GenesisState defines a type alias for the Band genesis application state.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	cdc := codecstd.MakeCodec(ModuleBasics)
	denom := "uband"

	stakingGenesis := staking.DefaultGenesisState()
	mintGenesis := mint.DefaultGenesisState()
	govGenesis := gov.DefaultGenesisState()
	crisisGenesis := crisis.DefaultGenesisState()
	slashingGenesis := slashing.DefaultGenesisState()

	stakingGenesis.Params.BondDenom = denom
	stakingGenesis.Params.HistoricalEntries = 1000
	mintGenesis.Params.BlocksPerYear = 10519200 // target 3-second block time
	mintGenesis.Params.MintDenom = denom
	govGenesis.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(denom, sdk.TokensFromConsensusPower(1000)))
	crisisGenesis.ConstantFee = sdk.NewCoin(denom, sdk.TokensFromConsensusPower(10000))
	slashingGenesis.Params.SignedBlocksWindow = 30000                         // approximately 1 day
	slashingGenesis.Params.MinSignedPerWindow = sdk.NewDecWithPrec(5, 2)      // 5%
	slashingGenesis.Params.DowntimeJailDuration = 60 * 10 * time.Second       // 10 minutes
	slashingGenesis.Params.SlashFractionDoubleSign = sdk.NewDecWithPrec(5, 2) // 5%
	slashingGenesis.Params.SlashFractionDowntime = sdk.NewDecWithPrec(1, 4)   // 0.01%

	return GenesisState{
		genutil.ModuleName:    genutil.AppModuleBasic{}.DefaultGenesis(cdc),
		auth.ModuleName:       auth.AppModuleBasic{}.DefaultGenesis(cdc),
		bank.ModuleName:       bank.AppModuleBasic{}.DefaultGenesis(cdc),
		staking.ModuleName:    cdc.MustMarshalJSON(stakingGenesis),
		mint.ModuleName:       cdc.MustMarshalJSON(mintGenesis),
		distr.ModuleName:      distr.AppModuleBasic{}.DefaultGenesis(cdc),
		gov.ModuleName:        cdc.MustMarshalJSON(govGenesis),
		crisis.ModuleName:     cdc.MustMarshalJSON(crisisGenesis),
		slashing.ModuleName:   cdc.MustMarshalJSON(slashingGenesis),
		ibc.ModuleName:        ibc.AppModuleBasic{}.DefaultGenesis(cdc),
		capability.ModuleName: capability.AppModuleBasic{}.DefaultGenesis(cdc),
		upgrade.ModuleName:    upgrade.AppModuleBasic{}.DefaultGenesis(cdc),
		evidence.ModuleName:   evidence.AppModuleBasic{}.DefaultGenesis(cdc),
		transfer.ModuleName:   transfer.AppModuleBasic{}.DefaultGenesis(cdc),
		oracle.ModuleName:     oracle.AppModuleBasic{}.DefaultGenesis(cdc),
	}
}
