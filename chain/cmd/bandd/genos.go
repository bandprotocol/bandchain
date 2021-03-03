package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/bandprotocol/go-owasm/api"
)

// AddGenesisOracleScriptCmd returns add-oracle-script cobra Command.
func AddGenesisOracleScriptCmd(ctx *server.Context, cdc *codec.Codec, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-oracle-script [name] [description] [schema] [url] [owner] [filepath]",
		Short: "Add a data source to genesis.json",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			f := filecache.New(filepath.Join(viper.GetString(cli.HomeFlag), "files"))
			data, err := ioutil.ReadFile(args[5])
			if err != nil {
				return err
			}
			vm, err := api.NewVm(0) // Compile didn't use cache
			if err != nil {
				return err
			}
			compiledData, err := vm.Compile(data, types.MaxCompiledWasmCodeSize)
			if err != nil {
				return err
			}
			filename := f.AddFile(compiledData)
			owner, err := sdk.AccAddressFromBech32(args[4])
			if err != nil {
				return err
			}
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}
			oracleGenState := oracle.GetGenesisStateFromAppState(cdc, appState)
			oracleGenState.OracleScripts = append(oracleGenState.OracleScripts, types.NewOracleScript(
				owner, args[0], args[1], filename, args[2], args[3],
			))
			appState[oracle.ModuleName] = cdc.MustMarshalJSON(oracleGenState)
			appStateJSON := cdc.MustMarshalJSON(appState)
			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}
	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	return cmd
}
