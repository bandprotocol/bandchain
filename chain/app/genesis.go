package app

import (
	"encoding/json"
	"io/ioutil"

	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
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
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/bandprotocol/bandchain/chain/x/zoracle"
)

// GenesisState defines a type alias for the Band genesis application state.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	cdc := codecstd.MakeCodec(ModuleBasics)
	return GenesisState{
		genutil.ModuleName:  genutil.AppModuleBasic{}.DefaultGenesis(cdc),
		auth.ModuleName:     auth.AppModuleBasic{}.DefaultGenesis(cdc),
		bank.ModuleName:     bank.AppModuleBasic{}.DefaultGenesis(cdc),
		staking.ModuleName:  staking.AppModuleBasic{}.DefaultGenesis(cdc),
		mint.ModuleName:     mint.AppModuleBasic{}.DefaultGenesis(cdc),
		distr.ModuleName:    distr.AppModuleBasic{}.DefaultGenesis(cdc),
		gov.ModuleName:      gov.AppModuleBasic{}.DefaultGenesis(cdc),
		crisis.ModuleName:   crisis.AppModuleBasic{}.DefaultGenesis(cdc),
		slashing.ModuleName: slashing.AppModuleBasic{}.DefaultGenesis(cdc),
		supply.ModuleName:   supply.AppModuleBasic{}.DefaultGenesis(cdc),
		ibc.ModuleName:      ibc.AppModuleBasic{}.DefaultGenesis(cdc),
		upgrade.ModuleName:  upgrade.AppModuleBasic{}.DefaultGenesis(cdc),
		evidence.ModuleName: evidence.AppModuleBasic{}.DefaultGenesis(cdc),
		transfer.ModuleName: transfer.AppModuleBasic{}.DefaultGenesis(cdc),
		zoracle.ModuleName:  zoracle.AppModuleBasic{}.DefaultGenesis(cdc),
		// staking.ModuleName: staking.ModuleCdc.MustMarshalJSON(staking.GenesisState{
		// 	Params: staking.Params{
		// 		UnbondingTime: time.Hour * 24 * 7 * 3, // 3 weeks
		// 		BondDenom:     "uband",
		// 		MaxEntries:    7,
		// 		MaxValidators: 100,
		// 	},
		// }),
		// mint.ModuleName: mint.ModuleCdc.MustMarshalJSON(mint.GenesisState{
		// 	Minter: mint.Minter{
		// 		AnnualProvisions: sdk.NewDecWithPrec(0, 0),
		// 		Inflation:        sdk.NewDecWithPrec(135, 3), // 13.5%
		// 	},
		// 	Params: mint.Params{
		// 		BlocksPerYear:       10519200,                  // 3 second  block times
		// 		GoalBonded:          sdk.NewDecWithPrec(67, 2), // 67%
		// 		InflationMax:        sdk.NewDecWithPrec(20, 2), // 20%
		// 		InflationMin:        sdk.NewDecWithPrec(7, 2),  // 7%
		// 		InflationRateChange: sdk.NewDecWithPrec(13, 2), // 13%
		// 		MintDenom:           "uband",
		// 	},
		// }),
		// distr.ModuleName: distr.ModuleCdc.MustMarshalJSON(distr.GenesisState{
		// 	FeePool:                         distr.InitialFeePool(),
		// 	CommunityTax:                    sdk.NewDecWithPrec(2, 2), // 2%
		// 	BaseProposerReward:              sdk.NewDecWithPrec(1, 2), // 1%
		// 	BonusProposerReward:             sdk.NewDecWithPrec(4, 2), // 4%
		// 	WithdrawAddrEnabled:             true,
		// 	DelegatorWithdrawInfos:          []distr.DelegatorWithdrawInfo{},
		// 	PreviousProposer:                nil,
		// 	OutstandingRewards:              []distr.ValidatorOutstandingRewardsRecord{},
		// 	ValidatorAccumulatedCommissions: []distr.ValidatorAccumulatedCommissionRecord{},
		// 	ValidatorHistoricalRewards:      []distr.ValidatorHistoricalRewardsRecord{},
		// 	ValidatorCurrentRewards:         []distr.ValidatorCurrentRewardsRecord{},
		// 	DelegatorStartingInfos:          []distr.DelegatorStartingInfoRecord{},
		// 	ValidatorSlashEvents:            []distr.ValidatorSlashEventRecord{},
		// }),
		// gov.ModuleName: gov.ModuleCdc.MustMarshalJSON(gov.GenesisState{
		// 	StartingProposalID: 1,
		// 	DepositParams: gov.DepositParams{
		// 		MinDeposit:       sdk.NewCoins(sdk.NewCoin("uband", sdk.TokensFromConsensusPower(1000))),
		// 		MaxDepositPeriod: 86400 * 14 * time.Second, // 14 days
		// 	},
		// 	VotingParams: gov.VotingParams{
		// 		VotingPeriod: 86400 * 14 * time.Second, // 14 days
		// 	},
		// 	TallyParams: gov.TallyParams{
		// 		Quorum:    sdk.NewDecWithPrec(4, 1),   // 40%
		// 		Threshold: sdk.NewDecWithPrec(5, 1),   // 50%
		// 		Veto:      sdk.NewDecWithPrec(334, 3), // 33.4%
		// 	},
		// }),
		// crisis.ModuleName: crisis.ModuleCdc.MustMarshalJSON(crisis.GenesisState{
		// 	ConstantFee: sdk.NewCoin("uband", sdk.TokensFromConsensusPower(10000)),
		// }),
		// slashing.ModuleName: slashing.ModuleCdc.MustMarshalJSON(slashing.GenesisState{
		// 	Params: slashing.Params{
		// 		MaxEvidenceAge:          60 * 30240 * time.Second, // 3 weeks
		// 		SignedBlocksWindow:      int64(30000),
		// 		MinSignedPerWindow:      sdk.NewDecWithPrec(5, 2), // 5%
		// 		DowntimeJailDuration:    60 * 10 * time.Second,    // 10 minutes
		// 		SlashFractionDoubleSign: sdk.NewDecWithPrec(5, 2), // 5%
		// 		SlashFractionDowntime:   sdk.NewDecWithPrec(1, 4), // 0.01%
		// 	},
		// 	SigningInfos: make(map[string]slashing.ValidatorSigningInfo),
		// 	MissedBlocks: make(map[string][]slashing.MissedBlock),
		// }),
		// supply.ModuleName: supply.ModuleCdc.MustMarshalJSON(supply.GenesisState{
		// 	Supply: sdk.NewCoins(),
		// }),
		// zoracle.ModuleName: zoracle.ModuleCdc.MustMarshalJSON(zoracle.DefaultGenesisState()),
	}
}

func GetDefaultDataSourcesAndOracleScripts(owner sdk.AccAddress) json.RawMessage {
	state := zoracle.DefaultGenesisState()
	dataSources := []struct {
		name        string
		description string
		path        string
	}{
		{
			"Coingecko script",
			"The Script that queries crypto price from https://coingecko.com",
			"./datasources/coingecko_price.sh",
		},
		{
			"Crypto compare script",
			"The Script that queries crypto price from https://cryptocompare.com",
			"./datasources/crypto_compare_price.sh",
		},
		{
			"Binance price",
			"The Script that queries crypto price from https://www.binance.com/en",
			"./datasources/binance_price.sh",
		},
		{
			"Open weather",
			"The script that queries current weather",
			"./datasources/open_weather_map.sh",
		},
	}

	// TODO: Find a better way to specify path to data sources
	state.DataSources = make([]zoracle.DataSource, len(dataSources))
	for i, dataSource := range dataSources {
		script, err := ioutil.ReadFile(dataSource.path)
		if err != nil {
			panic(err)
		}
		state.DataSources[i] = zoracle.NewDataSource(
			owner,
			dataSource.name,
			dataSource.description,
			sdk.Coins{},
			script,
		)
	}

	// TODO: Find a better way to specify path to oracle scripts
	oracleScripts := []struct {
		name        string
		description string
		path        string
	}{
		{
			"Crypto price script",
			"Oracle script for getting an average crypto price from many sources.",
			"./owasm/res/crypto_price.wasm",
		},
		{
			"Crypto price script (Borsh version)",
			"Oracle script for getting an average crypto price from many sources encoding parameter by borsh.",
			"./owasm/res/crypto_price_borsh.wasm",
		},
	}
	state.OracleScripts = make([]zoracle.OracleScript, len(oracleScripts))
	for i, oracleScript := range oracleScripts {
		code, err := ioutil.ReadFile(oracleScript.path)
		if err != nil {
			panic(err)
		}
		state.OracleScripts[i] = zoracle.NewOracleScript(
			owner,
			oracleScript.name,
			oracleScript.description,
			code,
		)
	}
	return zoracle.ModuleCdc.MustMarshalJSON(state)
}
