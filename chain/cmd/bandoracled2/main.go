package main

import (
	"fmt"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bandprotocol/bandchain/chain/app"
)

const (
	flagValidator = "validator"
)

// Config data structure for bandoracled daemon.
type Config struct {
	ChainID   string `mapstructure:"chain-id"`  // ChainID of the target chain
	NodeURI   string `mapstructure:"node"`      // Remote RPC URI of BandChain node to connect to
	Validator string `mapstructure:"validator"` // The validator address that I'm responsible for
}

// Global instances.
var (
	cfg     Config
	keybase keyring.Keyring
)

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(flags.FlagHome)
	if err != nil {
		return err
	}
	viper.SetConfigFile(path.Join(home, "config", "config.yaml"))
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	return nil
}

func main() {
	appConfig := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(appConfig)
	appConfig.Seal()

	ctx := &Context{}
	rootCmd := &cobra.Command{
		Use:   "oracled",
		Short: "BandChain oracle daemon to subscribe and response to oracle requests",
	}

	rootCmd.AddCommand(configCmd(), keysCmd(ctx), runCmd(ctx))
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		home, err := rootCmd.PersistentFlags().GetString(flags.FlagHome)
		if err != nil {
			return err
		}
		keybase, err = keyring.New("band", "test", home, nil)
		if err != nil {
			return err
		}
		return initConfig(rootCmd)
	}
	rootCmd.PersistentFlags().String(flags.FlagHome, os.ExpandEnv("$HOME/.oracled"), "home directory")
	rootCmd.PersistentFlags().String(flags.FlagChainID, "bandchain-dev", "chain ID of BandChain network")
	rootCmd.PersistentFlags().String(flags.FlagNode, "tcp://localhost:26657", "RPC url to BandChain node")
	rootCmd.PersistentFlags().String(flagValidator, "", "validator address")
	viper.BindPFlag(flags.FlagChainID, rootCmd.PersistentFlags().Lookup(flags.FlagChainID))
	viper.BindPFlag(flags.FlagNode, rootCmd.PersistentFlags().Lookup(flags.FlagNode))
	viper.BindPFlag(flagValidator, rootCmd.PersistentFlags().Lookup(flagValidator))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
