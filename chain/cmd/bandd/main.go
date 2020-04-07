package main

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
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

var invCheckPeriod uint

func main() {
	cdc := codecstd.MakeCodec(app.ModuleBasics)
	appCodec := codecstd.NewAppCodec(cdc)

	config := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(config)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "bandd",
		Short:             "BandChain Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(InitCmd(ctx, cdc, app.NewDefaultGenesisState(), app.GetDefaultDataSourcesAndOracleScripts, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, bank.GenesisBalancesIterator{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.MigrateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
			bank.GenesisBalancesIterator{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
	)
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc, appCodec, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(flags.NewCompletionCmd(rootCmd, true))
	// rootCmd.AddCommand(testnetCmd(ctx, cdc, app.ModuleBasics, bank.GenesisBalancesIterator{}))
	// rootCmd.AddCommand(replayCmd())
	rootCmd.AddCommand(debug.Cmd(cdc))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "BAND", app.DefaultNodeHome)

	rootCmd.PersistentFlags().UintVar(
		&invCheckPeriod, flagInvCheckPeriod, 0, "Assert registered invariants every N blocks",
	)
	rootCmd.PersistentFlags().String(
		flagWithDB, "", "[Experimental] Flush blockchain state to SQL database",
	)
	rootCmd.PersistentFlags().Int64(
		flagUptimeLookBackDuration, 1000, "[Experimental] Historical node uptime lookback duration",
	)

	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	var cache sdk.MultiStorePersistentCache

	if viper.GetBool(server.FlagInterBlockCache) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range viper.GetIntSlice(server.FlagUnsafeSkipUpgrades) {
		skipUpgradeHeights[int64(h)] = true
	}

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
			logger, db, traceStore, true, invCheckPeriod, skipUpgradeHeights,
			viper.GetString(flags.FlagHome), bandDB,
			baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
			baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
			baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
			baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
			baseapp.SetInterBlockCache(cache),
		)
	} else {
		return app.NewBandApp(
			logger, db, traceStore, true, invCheckPeriod, skipUpgradeHeights,
			viper.GetString(flags.FlagHome),
			baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
			baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
			baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
			baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
			baseapp.SetInterBlockCache(cache),
		)
	}
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		bandApp := app.NewBandApp(logger, db, traceStore, false, uint(1), map[int64]bool{}, "")
		err := bandApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}

		return bandApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	bandApp := app.NewBandApp(logger, db, traceStore, true, uint(1), map[int64]bool{}, "")
	return bandApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
