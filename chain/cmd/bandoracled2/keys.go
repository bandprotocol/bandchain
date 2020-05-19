package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
)

func keysCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"k"},
		Short:   "Manage key held by the oracle process",
	}
	cmd.AddCommand(keysAddCmd(c))
	cmd.AddCommand(keysListCmd(c))
	return cmd
}

func keysAddCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [name]",
		Aliases: []string{"a"},
		Short:   "Add a new key to the keychain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Allow mnemonic import
			seed, err := bip39.NewEntropy(256)
			if err != nil {
				return err
			}

			mnemonic, err := bip39.NewMnemonic(seed)
			if err != nil {
				return err
			}

			info, err := keybase.NewAccount(
				args[0], mnemonic, "", hd.CreateHDPath(494, 0, 0).String(), hd.Secp256k1,
			)
			if err != nil {
				return err
			}

			fmt.Printf("Mnemonic: %s\n", mnemonic)
			fmt.Printf("Address: %s", info.GetAddress().String())
			return nil
		},
	}
	return cmd
}

func keysListCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List all the keys in the keychain",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			keys, err := keybase.List()
			if err != nil {
				return err
			}
			for _, key := range keys {
				fmt.Printf("%s => %s\n", key.GetName(), key.GetAddress().String())
			}
			return nil
		},
	}
	return cmd
}
