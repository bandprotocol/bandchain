package main

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
)

var (
	keybase  keyring.Keyring
	homePath string
	chainID  string
	nodeURI  string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "oracled",
		Short: "🔮 BandChain oracle daemon to subscribe and response to oracle requests",
	}

	rootCmd.AddCommand(
		configCmd(),
		keysCmd(),
		runCmd(),
	)

	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		var err error
		keybase, err = keyring.New("band", "test", homePath, nil)
		if err != nil {
			return err
		}
		return nil
	}

	rootCmd.PersistentFlags().StringVar(
		&homePath, flags.FlagHome, os.ExpandEnv("$HOME/.oracled"), "home directory",
	)
	rootCmd.PersistentFlags().StringVar(
		&chainID, flags.FlagChainID, "bandchain-dev", "chain ID of BandChain network",
	)
	rootCmd.PersistentFlags().StringVar(
		&nodeURI, flags.FlagNode, "tcp://localhost:26657", "RPC url to BandChain node",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
