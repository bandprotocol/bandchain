package main

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
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
	"github.com/bandprotocol/bandchain/chain/hooks/emitter"
	"github.com/bandprotocol/bandchain/chain/hooks/price"
	"github.com/bandprotocol/bandchain/chain/hooks/replay"
	"github.com/bandprotocol/bandchain/chain/hooks/request"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

const (
	flagInvCheckPeriod        = "inv-check-period"
	flagWithEmitter           = "with-emitter"
	flagDisableFeelessReports = "disable-feeless-reports"
	flagEnableFastSync        = "enable-fast-sync"
	flagWithPricer            = "with-pricer"
	flagWithRequestSearch     = "with-request-search"
	flagReplayMode            = "replay-mode"
)

var invCheckPeriod uint

func main() {
	cdc := app.MakeCodec()
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
	// Add subcommands to bandd root cmd.
	rootCmd.AddCommand(InitCmd(ctx, cdc, app.NewDefaultGenesisState(), app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome))
	rootCmd.AddCommand(MigrateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(genutilcli.GenTxCmd(ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{}, auth.GenesisAccountIterator{}, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(AddGenesisDataSourceCmd(ctx, cdc, app.DefaultNodeHome))
	rootCmd.AddCommand(AddGenesisOracleScriptCmd(ctx, cdc, app.DefaultNodeHome))
	rootCmd.AddCommand(flags.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(debug.Cmd(cdc))
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)
	// Prepare and add persistent flags.
	executor := cli.PrepareBaseCmd(rootCmd, "BAND", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod, 0, "Assert registered invariants every N blocks")
	rootCmd.PersistentFlags().Bool(flagDisableFeelessReports, false, "[Experimental] Disable allowance of feeless reports")
	rootCmd.PersistentFlags().String(flagWithEmitter, "", "[Experimental] Use Kafka emitter")
	rootCmd.PersistentFlags().Bool(flagEnableFastSync, false, "[Experimental] Enable fast sync mode")
	rootCmd.PersistentFlags().String(flagWithRequestSearch, "", "[Experimental] Enable mode to save request in sql database")
	rootCmd.PersistentFlags().String(flagWithPricer, "", "[Experimental] Enable mode to save price in level db")
	rootCmd.PersistentFlags().String(flagReplayMode, "", "[Experimental] Use emitter replay mode")
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
	pruningOpts, err := server.GetPruningOptionsFromFlags()
	if err != nil {
		panic(err)
	}

	bandApp := app.NewBandApp(
		logger, db, traceStore, true, invCheckPeriod, skipUpgradeHeights,
		viper.GetString(flags.FlagHome),
		viper.GetBool(flagDisableFeelessReports),
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
		baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
		baseapp.SetInterBlockCache(cache),
	)

	// Add hooks based on flag
	if viper.IsSet(flagWithEmitter) {
		bandApp.AddHook(emitter.NewHook(
			bandApp.Codec(), bandApp.AccountKeeper, bandApp.BankKeeper, bandApp.SupplyKeeper,
			bandApp.StakingKeeper, bandApp.MintKeeper, bandApp.DistrKeeper, bandApp.GovKeeper,
			bandApp.OracleKeeper, viper.GetString(flagWithEmitter), viper.GetBool(flagEnableFastSync)))
	}
	if viper.IsSet(flagWithRequestSearch) {
		bandApp.AddHook(request.NewHook(
			bandApp.Codec(), bandApp.OracleKeeper,
			viper.GetString(flagWithRequestSearch)))
	}
	if viper.IsSet(flagWithPricer) {
		rawOids := strings.Split(viper.GetString(flagWithPricer), ",")
		oids := make([]types.OracleScriptID, len(rawOids))
		for idx, rawOid := range rawOids {
			oid, err := strconv.ParseInt(rawOid, 10, 64)
			if err != nil {
				panic(err)
			}
			oids[idx] = types.OracleScriptID(oid)
		}
		bandApp.AddHook(price.NewHook(bandApp.Codec(), bandApp.OracleKeeper, oids, filepath.Join(viper.GetString(cli.HomeFlag), "prices")))
	}
	if viper.IsSet(flagReplayMode) {
		bandApp.AddHook(replay.NewHook(
			bandApp.Codec(), bandApp.AccountKeeper, bandApp.BankKeeper, bandApp.SupplyKeeper,
			bandApp.StakingKeeper, bandApp.MintKeeper, bandApp.DistrKeeper, bandApp.GovKeeper,
			bandApp.OracleKeeper, viper.GetString(flagReplayMode)))
	}
	return bandApp
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		bandApp := app.NewBandApp(logger, db, traceStore, false, uint(1), map[int64]bool{}, "", viper.GetBool(flagDisableFeelessReports))
		err := bandApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}

		return bandApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	bandApp := app.NewBandApp(logger, db, traceStore, true, uint(1), map[int64]bool{}, "", viper.GetBool(flagDisableFeelessReports))
	return bandApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
