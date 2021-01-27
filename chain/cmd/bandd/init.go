package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/libs/tempfile"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
)

const (
	flagTimeoutCommit = "timeout-commit"
	flagOverwrite     = "overwrite"
)

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(cdc *codec.Codec, info printInfo) error {
	out, err := codec.MarshalJSONIndent(cdc, info)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stderr, "%s\n", string(sdk.MustSortJSON(out)))
	return err
}

func genFilePVIfNotExists(cdc *codec.Codec, keyFilePath, stateFilePath string) error {
	if !tmos.FileExists(keyFilePath) {
		privKey := secp256k1.GenPrivKey()
		pv := &privval.FilePV{
			Key: privval.FilePVKey{
				Address: privKey.PubKey().Address(),
				PubKey:  privKey.PubKey(),
				PrivKey: privKey,
			},
			LastSignState: privval.FilePVLastSignState{
				Step: 0,
			},
		}
		jsonBytes, err := cdc.MarshalJSONIndent(pv.Key, "", "  ")
		if err != nil {
			return err
		}
		if err = tempfile.WriteFileAtomic(keyFilePath, jsonBytes, 0600); err != nil {
			return err
		}
		jsonBytes, err = cdc.MarshalJSONIndent(pv.LastSignState, "", "  ")
		if err != nil {
			return err
		}
		if err = tempfile.WriteFileAtomic(stateFilePath, jsonBytes, 0600); err != nil {
			return err
		}
	}
	return nil
}

// InitCmd returns a command that initializes all files needed for Tendermint and BandChain app.
func InitCmd(ctx *server.Context, cdc *codec.Codec, customAppState map[string]json.RawMessage, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			config.Consensus.TimeoutCommit = time.Duration(viper.GetInt(flagTimeoutCommit)) * time.Millisecond
			if err := genFilePVIfNotExists(cdc, config.PrivValidatorKeyFile(), config.PrivValidatorStateFile()); err != nil {
				return err
			}
			chainID := viper.GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", tmrand.Str(6))
			}
			nodeID, _, err := genutil.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}
			config.Moniker = args[0]
			genFile := config.GenesisFile()
			if !viper.GetBool(flagOverwrite) && tmos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}
			appState, err := codec.MarshalJSONIndent(cdc, customAppState)
			if err != nil {
				return err
			}
			genDoc := &types.GenesisDoc{}
			genDoc.ChainID = chainID
			genDoc.Validators = nil
			genDoc.AppState = appState
			genDoc.ConsensusParams = types.DefaultConsensusParams()
			genDoc.ConsensusParams.Block.MaxBytes = 1000000 // 1M bytes
			genDoc.ConsensusParams.Block.MaxGas = 5000000   // 5M gas
			genDoc.ConsensusParams.Block.TimeIotaMs = 1000  // 1 second
			genDoc.ConsensusParams.Validator.PubKeyTypes = []string{types.ABCIPubKeyTypeSecp256k1}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return err
				}
			}
			if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
				return err
			}
			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)
			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(cdc, toPrint)
		},
	}
	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Int(flagTimeoutCommit, 1500, "timeout commit of the node in ms")
	return cmd
}
