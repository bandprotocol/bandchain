package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/GeoDB-Limited/odincore/chain/app"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	extypes "github.com/cosmos/cosmos-sdk/x/genutil"
	v038 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v0_38"
	v039 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v0_39"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	flagGenesisTime = "genesis-time"
	flagChainID     = "chain-id"
)

// GenesisDocFromFile reads JSON data from a file and unmarshalls it into a GenesisDoc.
func GenesisDocFromFile(genDocFile string, cdc *codec.Codec) (*tmtypes.GenesisDoc, error) {
	jsonBlob, err := ioutil.ReadFile(genDocFile)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read GenesisDoc file")
	}

	var genDoc tmtypes.GenesisDoc
	err = cdc.UnmarshalJSON(jsonBlob, &genDoc)
	if err != nil {
		return nil, err
	}

	// Set up Tendermint consensus params with default value.
	genDoc.ConsensusParams.Evidence = tmtypes.DefaultEvidenceParams()
	genDoc.ConsensusParams.Block.MaxBytes = 1000000 // 1M bytes
	genDoc.ConsensusParams.Block.MaxGas = 5000000   // 5M gas
	genDoc.ConsensusParams.Block.TimeIotaMs = 1000  // 1 second

	if err := genDoc.ValidateAndComplete(); err != nil {
		return nil, err
	}

	return &genDoc, nil
}

// MigrateGenesisCmd returns a command to execute genesis state migration.
// nolint: funlen
func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [genesis-file]",
		Short: "Migrate Wenchang genesis to Guanyu version",
		Long: fmt.Sprintf(`Migrate the Wenchang genesis into the Guanyu version and print to STDOUT.

Example:
$ %s migrate /path/to/genesis.json --chain-id=band-guanyu --genesis-time=2020-08-11T17:00:00Z
`, version.ServerName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			importGenesis := args[0]
			genDoc, err := GenesisDocFromFile(importGenesis, cdc)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var state extypes.AppMap
			if err := cdc.UnmarshalJSON(genDoc.AppState, &state); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			// Migrate from Wenchang (0.36 like genesis file) to cosmos-sdk v0.39
			state = v039.Migrate(v038.Migrate(state))

			// Add genesis state of the new modules get state from `app` default genesis.
			defaultGuanyu := app.NewDefaultGenesisState()
			state[oracle.ModuleName] = cdc.MustMarshalJSON(defaultGuanyu[oracle.ModuleName])
			state[evidence.ModuleName] = cdc.MustMarshalJSON(defaultGuanyu[evidence.ModuleName])
			state[upgrade.ModuleName] = cdc.MustMarshalJSON(defaultGuanyu[upgrade.ModuleName])

			genDoc.AppState, err = cdc.MarshalJSON(state)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}
			if err := genDoc.ValidateAndComplete(); err != nil {
				return err
			}

			genesisTime := cmd.Flag(flagGenesisTime).Value.String()
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal genesis time")
				}

				genDoc.GenesisTime = t
			}

			chainID := cmd.Flag(flagChainID).Value.String()
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			bz, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}

			fmt.Println(string(sortedBz))
			return nil
		},
	}
	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "band-guanyu", "override chain_id with this flag")
	return cmd
}
