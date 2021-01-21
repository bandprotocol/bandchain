package main

import (
	"fmt"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/GeoDB-Limited/odincore/chain/app"
)

const (
	flagPort   = "port"
	flagAmount = "amount"
)

// Config data structure for faucet server.
type Config struct {
	ChainID   string `mapstructure:"chain-id"`   // ChainID of the target chain
	NodeURI   string `mapstructure:"node"`       // Remote RPC URI of BandChain node to connect to
	GasPrices string `mapstructure:"gas-prices"` // Gas prices of the transaction
	Port      string `mapstructure:"port"`       // Port of faucet service
	Amount    int64  `mapstructure:"amount"`     // Amount of BAND for each request
}

// Global instances.
var (
	cfg     Config
	keybase keys.Keybase
)

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(flags.FlagHome)
	if err != nil {
		return err
	}
	viper.SetConfigFile(path.Join(home, "config.yaml"))
	_ = viper.ReadInConfig() // If we fail to read config file, we'll just rely on cmd flags.
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
		Use:   "faucet",
		Short: "Faucet server for devnet",
	}

	rootCmd.AddCommand(configCmd(), keysCmd(ctx), runCmd(ctx))
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		home, err := rootCmd.PersistentFlags().GetString(flags.FlagHome)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(home, os.ModePerm); err != nil {
			return err
		}
		keybase, err = keys.NewKeyring("band", "test", home, nil)
		if err != nil {
			return err
		}
		return initConfig(rootCmd)
	}
	rootCmd.PersistentFlags().String(flags.FlagHome, os.ExpandEnv("$HOME/.faucet"), "home directory")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
