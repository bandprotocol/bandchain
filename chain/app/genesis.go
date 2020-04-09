package app

import (
	"encoding/json"
	"io/ioutil"
	"time"

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
		genutil.ModuleName:  genutil.AppModuleBasic{}.DefaultGenesis(cdc),
		auth.ModuleName:     auth.AppModuleBasic{}.DefaultGenesis(cdc),
		bank.ModuleName:     bank.AppModuleBasic{}.DefaultGenesis(cdc),
		staking.ModuleName:  cdc.MustMarshalJSON(stakingGenesis),
		mint.ModuleName:     cdc.MustMarshalJSON(mintGenesis),
		distr.ModuleName:    distr.AppModuleBasic{}.DefaultGenesis(cdc),
		gov.ModuleName:      cdc.MustMarshalJSON(govGenesis),
		crisis.ModuleName:   cdc.MustMarshalJSON(crisisGenesis),
		slashing.ModuleName: cdc.MustMarshalJSON(slashingGenesis),
		supply.ModuleName:   supply.AppModuleBasic{}.DefaultGenesis(cdc),
		ibc.ModuleName:      ibc.AppModuleBasic{}.DefaultGenesis(cdc),
		upgrade.ModuleName:  upgrade.AppModuleBasic{}.DefaultGenesis(cdc),
		evidence.ModuleName: evidence.AppModuleBasic{}.DefaultGenesis(cdc),
		transfer.ModuleName: transfer.AppModuleBasic{}.DefaultGenesis(cdc),
		oracle.ModuleName:  oracle.AppModuleBasic{}.DefaultGenesis(cdc),
	}
}

func GetDefaultDataSourcesAndOracleScripts(owner sdk.AccAddress) json.RawMessage {
	state := oracle.DefaultGenesisState()
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
	state.DataSources = make([]oracle.DataSource, len(dataSources))
	for i, dataSource := range dataSources {
		script, err := ioutil.ReadFile(dataSource.path)
		if err != nil {
			panic(err)
		}
		state.DataSources[i] = oracle.NewDataSource(
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
	state.OracleScripts = make([]oracle.OracleScript, len(oracleScripts))
	for i, oracleScript := range oracleScripts {
		code, err := ioutil.ReadFile(oracleScript.path)
		if err != nil {
			panic(err)
		}
		state.OracleScripts[i] = oracle.NewOracleScript(
			owner,
			oracleScript.name,
			oracleScript.description,
			code,
		)
	}
	return oracle.ModuleCdc.MustMarshalJSON(state)
}
