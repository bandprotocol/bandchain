package main

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/bandchain/chain/app"
)

func main() {
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(config)
	config.Seal()

	ctx := &Context{}
	rootCmd := &cobra.Command{
		Use:   "oracled",
		Short: "BandChain oracle daemon to subscribe and response to oracle requests",
	}

	rootCmd.AddCommand(
		configCmd(),
		keysCmd(ctx),
		runCmd(ctx),
	)

	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		var err error
		ctx.keybase, err = keyring.New("band", "test", ctx.homePath, nil)
		if err != nil {
			return err
		}
		return nil
	}

	rootCmd.PersistentFlags().StringVar(
		&ctx.homePath, flags.FlagHome, os.ExpandEnv("$HOME/.oracled"), "home directory",
	)
	rootCmd.PersistentFlags().StringVar(
		&ctx.chainID, flags.FlagChainID, "bandchain-dev", "chain ID of BandChain network",
	)
	rootCmd.PersistentFlags().StringVar(
		&ctx.nodeURI, flags.FlagNode, "tcp://localhost:26657", "RPC url to BandChain node",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
