package main

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	genaccscli "github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/bandprotocol/bandchain/chain/app"
	banddb "github.com/bandprotocol/bandchain/chain/db"
)

const (
	flagInvCheckPeriod         = "inv-check-period"
	flagWithDB                 = "with-db"
	flagUptimeLookBackDuration = "uptime-look-back"
)

var (
	invCheckPeriod uint
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(config)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "bandd",
		Short:             "BandChain App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		InitCmd(ctx, cdc, app.NewDefaultGenesisState(), app.GetDefaultDataSourcesAndOracleScripts, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, app.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(ctx, cdc),
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
			genaccounts.AppModuleBasic{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		genaccscli.AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
		client.NewCompletionCmd(rootCmd, true),
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "BAND", app.DefaultNodeHome)
	rootCmd.PersistentFlags().String(flagWithDB, "", "[Experimental] Flush blockchain state to SQL database")
	rootCmd.PersistentFlags().Int64(flagUptimeLookBackDuration, 1000, "[Experimental] Historical node uptime lookback duration")
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	if viper.IsSet(flagWithDB) {
		dbSplit := strings.SplitN(viper.GetString(flagWithDB), ":", 2)
		if len(dbSplit) != 2 {
			panic("Invalid DB string format")
		}
		metadata := map[string]string{
			banddb.KeyUptimeLookBackDuration: viper.GetString(flagUptimeLookBackDuration),
		}
		bandDB, err := banddb.NewDB(dbSplit[0], dbSplit[1], metadata)
		if err != nil {
			panic(err)
		}
		return app.NewDBBandApp(
			logger, db, traceStore, true, invCheckPeriod, bandDB,
			baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
			baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
			baseapp.SetHaltHeight(uint64(viper.GetInt(server.FlagHaltHeight))))
	} else {
		return app.NewBandApp(
			logger, db, traceStore, true, invCheckPeriod,
			baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
			baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
			baseapp.SetHaltHeight(uint64(viper.GetInt(server.FlagHaltHeight))))
	}
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		bandApp := app.NewBandApp(logger, db, traceStore, false, uint(1))
		err := bandApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return bandApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	bandApp := app.NewBandApp(logger, db, traceStore, true, uint(1))
	return bandApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
