package app

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
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
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
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

func GetDefaultDataSourcesAndOracleScripts(owner sdk.AccAddress) json.RawMessage {
	state := oracle.DefaultGenesisState()
	dataSources := []struct {
		name        string
		description string
		path        string
	}{
		{
			"Coingecko script",
			"The script that queries crypto price from https://coingecko.com",
			"./datasources/coingecko_price.py",
		},
		{
			"Crypto compare script",
			"The script that queries crypto price from https://cryptocompare.com",
			"./datasources/crypto_compare_price.sh",
		},
		{
			"Binance price",
			"The script that queries crypto price from https://www.binance.com/en",
			"./datasources/binance_price.sh",
		},
		{
			"Open weather map",
			"The script that queries weather information from https://api.openweathermap.org",
			"./datasources/open_weather_map.sh",
		},
		{
			"Gold price",
			"The script that queries current gold price",
			"./datasources/gold_price.sh",
		},
		{
			"Alphavantage",
			"The script that queries stock price from Alphavantage",
			"./datasources/alphavantage.sh",
		},
		{
			"Bitcoin block count",
			"The script that queries latest block height of Bitcoin",
			"./datasources/bitcoin_count.sh",
		},
		{
			"Bitcoin block hash",
			"The script that queries block hash of Bitcoin",
			"./datasources/bitcoin_hash.sh",
		},
		{
			"Coingecko volume script",
			"The script that queries crypto volume from Coingecko",
			"./datasources/coingecko_volume.sh",
		},
		{
			"Crypto compare volume script",
			"The script that queries crypto volume from Crypto compare",
			"./datasources/crypto_compare_volume.sh",
		},
		{
			"ETH gas station",
			"The script that queries current Ethereum gas price https://ethgasstation.info",
			"./datasources/ethgasstation.sh",
		},
		{
			"Open sky network",
			"The script that queries flight information from https://opensky-network.org",
			"./datasources/open_sky_network.sh",
		},
		{
			"Quantum random numbers",
			"The script that queries array of random number from https://qrng.anu.edu.au",
			"./datasources/qrng_anu.sh",
		},
		{
			"Yahoo finance",
			"The script that queries stock price from https://finance.yahoo.com",
			"./datasources/yahoo_finance.py",
		},
	}

	// TODO: Find a better way to specify path to data sources
	state.DataSources = make([]otypes.DataSource, len(dataSources))
	for i, dataSource := range dataSources {
		script, err := ioutil.ReadFile(dataSource.path)
		if err != nil {
			panic(err)
		}
		f := filecache.New(filepath.Join(viper.GetString(cli.HomeFlag), "files"))
		filename := f.AddFile(script)
		state.DataSources[i] = otypes.NewDataSource(
			owner,
			dataSource.name,
			dataSource.description,
			filename,
		)
	}

	// TODO: Find a better way to specify path to oracle scripts
	oracleScripts := []struct {
		name          string
		description   string
		path          string
		schema        string
		sourceCodeURL string
	}{
		{
			"Crypto price script",
			"Oracle script for getting the current an average cryptocurrency price from various sources.",
			"./pkg/owasm/res/crypto_price.wasm",
			`{symbol:string,multiplier:string}/{px:u64}`,
			"https://ipfs.io/ipfs/QmUbdfoRR9ge6P39EoqDjBhQoDeaT6gJu76Ce9kKsz938N",
		},
		{
			"Gold price script",
			"Oracle script for getting the current average gold price in ATOMs",
			"./pkg/owasm/res/gold_price.wasm",
			`{symbol:string,multiplier:string}/{px:u64}`,
			"https://ipfs.io/ipfs/Qmbcdr3UZXMrJeoRtHzTtHHepnzjyX1gWNhewWe6BXgmPm",
		},
		{
			"Alpha Vantage stock price script",
			"Oracle script for getting the current price of a stock from Alpha Vantage",
			"./pkg/owasm/res/alphavantage.wasm",
			`{symbol:string,api_key:string,multiplier:string}/{px:u64}`,
			"https://ipfs.io/ipfs/QmPsSmJ9gEdBoeQqwtk6bJykyFtqSpztCeGb9J1VFW65av",
		},
		{
			"Bitcoin block count",
			"Oracle script for getting the latest Bitcoin block height",
			"./pkg/owasm/res/bitcoin_block_count.wasm",
			`{_unused:u8}/{block_count:u64}`,
			"https://ipfs.io/ipfs/QmUkpTCvdKMEFxwgeTpjP9hszZ11e5ioXZAS7XLpQLbV2k",
		},
		{
			"Bitcoin block hash",
			"Oracle script for getting the Bitcoin block hash at the given block height",
			"./pkg/owasm/res/bitcoin_block_hash.wasm",
			`{block_height:u64}/{block_hash:string}`,
			"https://ipfs.io/ipfs/QmXu5NyUrtbcdPxut4WhVsRT4KjsPjy2NwJzEts7rjuEDf",
		},
		{
			"CoinGecko crypto volume",
			"Oracle script for getting a cryptocurrency's average trading volume for the past day from Coingecko",
			"./pkg/owasm/res/coingecko_volume.wasm",
			`{symbol:string,multiplier:string}/{volume:u64}`,
			"https://ipfs.io/ipfs/QmVuYP5cSujNSv33ZNMiFbcSMoRtBa7WMzG2q55j21Vhxj",
		},
		{
			"CryptoCompare crypto volume",
			"Oracle script for getting a cryptocurrency's average trading volume for the past day from CryptoCompare",
			"./pkg/owasm/res/crypto_compare_volume.wasm",
			`{symbol:string,multiplier:string}/{volume:u64}`,
			"https://ipfs.io/ipfs/Qmf2e5VF3uscGzBMwfQRZbaYLWxAZPWAgwv3m3iSv6BoGE",
		},
		{
			"Ethereum gas price",
			"Oracle script for getting the current Ethereum gas price from ETH gas station",
			"./pkg/owasm/res/eth_gas_station.wasm",
			`{gas_option:string}/{gweix10:u64}`,
			"https://ipfs.io/ipfs/QmP1i61XdPnfKSewh7vyh3xLgjxT42Gqpiv7CYLFK6V3Mg",
		},
		{
			"Open sky network",
			"Oracle script for checking whether a given flight number exists during the given time period",
			"./pkg/owasm/res/open_sky_network.wasm",
			`{flight_op:string,icao24:string,begin:string,end:string}/{flight_existence:u8}`,
			"https://ipfs.io/ipfs/QmST4us1xAXmfXZFqBRZqjsDpNhTMY1CRsgjNccwmB4FTX",
		},
		{
			"Open weather map",
			"Oracle script for getting the current weather data of a location",
			"./pkg/owasm/res/open_weather_map.wasm",
			`{country:string,main_field:string,sub_field:string,multiplier:u64}/{value:u64}`,
			"https://ipfs.io/ipfs/QmNWvYfqZztrMNKjKyKLTATvnVPVUMPCdLJFkV5HfHBQoo",
		},
		{
			"Quantum random number generator",
			"Oracle script for getting a big random number from quantum computer",
			"./pkg/owasm/res/qrng.wasm",
			`{size:u64}/{random_bytes:string}`,
			"https://ipfs.io/ipfs/QmZ62dxgAmCtDnt5XcAs2zP4UjGzozDzdKiHoR1Wo9MVeV",
		},
		{
			"Yahoo stock price",
			"Oracle script for getting the current price of a stock from Yahoo Finance",
			"./pkg/owasm/res/yahoo_price.wasm",
			`{symbol:string,multiplier:u64}/{px:u64}`,
			"https://ipfs.io/ipfs/QmfEUKFoX9PY3LHnT7Deixwb8qRrgWvdf5v8MzTinTYXLu",
		},
		{
			"Fair price from 3 sources",
			"Oracle script that query prices from many markets and then aggregate them together",
			"./pkg/owasm/res/fair_crypto_market_price.wasm",
			`{base_symbol:string,quote_symbol:string,aggregation_method:string,multiplier:u64}/{px:u64}`,
			"https://ipfs.io/ipfs/QmbnRei1WG8gdstsVuU7Qqq4PwqED9LuHFvjDgS5asShoM",
		},
	}
	state.OracleScripts = make([]otypes.OracleScript, len(oracleScripts))
	for i, oracleScript := range oracleScripts {
		code, err := ioutil.ReadFile(oracleScript.path)
		if err != nil {
			panic(err)
		}
		f := filecache.New(filepath.Join(viper.GetString(cli.HomeFlag), "files"))
		filename := f.AddFile(code)
		state.OracleScripts[i] = otypes.NewOracleScript(
			owner,
			oracleScript.name,
			oracleScript.description,
			filename,
			oracleScript.schema,
			oracleScript.sourceCodeURL,
		)
	}
	return oracle.ModuleCdc.MustMarshalJSON(state)
}
