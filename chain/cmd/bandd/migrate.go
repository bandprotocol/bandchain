package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
	genutil "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	extypes "github.com/cosmos/cosmos-sdk/x/genutil"
	v038 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v0_38"
	v039 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v0_39"
)

const (
	flagGenesisTime = "genesis-time"
	flagChainID     = "chain-id"
)

// Allow applications to extend and modify the migration process.
//
// Ref: https://github.com/cosmos/cosmos-sdk/issues/5041
var migrationMap = extypes.MigrationMap{
	"v0.38": v038.Migrate, // NOTE: v0.37 and v0.38 are genesis compatible
	"v0.39": v039.Migrate,
}

func Guanyu(appState genutil.AppMap) genutil.AppMap {
	oracleCodec := codec.New()
	codec.RegisterCrypto(oracleCodec)
	// migrate distribution state
	// var oracleGenState oracle.GenesisState
	// fmt.Println(oracle.DefaultGenesisState())
	// panic("yo")
	appState[oracle.ModuleName] = oracleCodec.MustMarshalJSON(oracle.DefaultGenesisState())
	return appState
}

// GetMigrationCallback returns a MigrationCallback for a given version.
func GetMigrationCallback(version string) extypes.MigrationCallback {
	return migrationMap[version]
}

// GetMigrationVersions get all migration version in a sorted slice.
func GetMigrationVersions() []string {
	versions := make([]string, len(migrationMap))

	var i int
	for version := range migrationMap {
		versions[i] = version
		i++
	}

	return versions
}

func GenesisDocFromJSON(jsonBlob []byte, cdc *codec.Codec) (*tmtypes.GenesisDoc, error) {
	genDoc := tmtypes.GenesisDoc{}
	// cdc.
	err := cdc.UnmarshalJSON(jsonBlob, &genDoc)
	if err != nil {
		return nil, err
	}

	genDoc.ConsensusParams.Evidence.MaxAgeNumBlocks = 100000
	genDoc.ConsensusParams.Evidence.MaxAgeDuration = 172800000000000

	if err := genDoc.ValidateAndComplete(); err != nil {
		return nil, err
	}

	return &genDoc, err
}

// GenesisDocFromFile reads JSON data from a file and unmarshalls it into a GenesisDoc.
func GenesisDocFromFile(genDocFile string, cdc *codec.Codec) (*tmtypes.GenesisDoc, error) {
	jsonBlob, err := ioutil.ReadFile(genDocFile)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read GenesisDoc file")
	}
	genDoc, err := GenesisDocFromJSON(jsonBlob, cdc)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error reading GenesisDoc at %v", genDocFile))
	}
	return genDoc, nil
}

// MigrateGenesisCmd returns a command to execute genesis state migration.
// nolint: funlen
func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Example:
$ %s migrate v0.36 /path/to/genesis.json --chain-id=cosmoshub-3 --genesis-time=2019-04-22T17:00:00Z
`, version.ServerName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			target := args[0]
			importGenesis := args[1]
			genDoc, err := GenesisDocFromFile(importGenesis, cdc)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState extypes.AppMap
			if err := cdc.UnmarshalJSON(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}
			_ = GetMigrationCallback(target)
			// if migrationFunc == nil {
			// 	return fmt.Errorf("unknown migration function for version: %s", target)
			// }

			// TODO: handler error from migrationFunc call
			initialState = v038.Migrate(initialState)
			initialState = v039.Migrate(initialState)
			newGenState := Guanyu(initialState)

			genDoc.AppState, err = cdc.MarshalJSON(newGenState)

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
	cmd.Flags().String(flagChainID, "", "override chain_id with this flag")

	return cmd
}
