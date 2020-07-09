package yoda

import (
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/input"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
)

const (
	flagAccount  = "account"
	flagIndex    = "index"
	flagCoinType = "coin-type"
	flagRecover  = "recover"
)

func keysCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"k"},
		Short:   "Manage key held by the oracle process",
	}
	cmd.AddCommand(keysAddCmd(c))
	cmd.AddCommand(keysDeleteCmd(c))
	cmd.AddCommand(keysListCmd(c))
	cmd.AddCommand(keysShowCmd(c))
	return cmd
}

func keysAddCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [name]",
		Aliases: []string{"a"},
		Short:   "Add a new key to the keychain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var mnemonic string
			recover, err := cmd.Flags().GetBool(flagRecover)
			if err != nil {
				return err
			}
			if recover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				var err error
				mnemonic, err = input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}
			} else {
				seed, err := bip39.NewEntropy(256)
				if err != nil {
					return err
				}
				mnemonic, err = bip39.NewMnemonic(seed)
				if err != nil {
					return err
				}
				fmt.Printf("Mnemonic: %s\n", mnemonic)
			}

			if err != nil {
				return err
			}
			account, err := cmd.Flags().GetUint32(flagAccount)
			if err != nil {
				return err
			}
			index, err := cmd.Flags().GetUint32(flagIndex)
			if err != nil {
				return err
			}
			hdPath := keys.CreateHDPath(account, index)
			info, err := keybase.CreateAccount(args[0], mnemonic, "", ckeys.DefaultKeyPass, hdPath.String(), keys.Secp256k1)
			if err != nil {
				return err
			}
			fmt.Printf("Address: %s\n", info.GetAddress().String())
			return nil
		},
	}
	cmd.Flags().Bool(flagRecover, false, "Provide seed phrase to recover existing key instead of creating")
	cmd.Flags().Uint32(flagAccount, 0, "Account number for HD derivation")
	cmd.Flags().Uint32(flagIndex, 0, "Address index number for HD derivation")
	return cmd
}

func keysDeleteCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [name]",
		Aliases: []string{"d"},
		Short:   "Delete a key from the keychain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			_, err := keybase.Get(name)
			if err != nil {
				return err
			}

			inBuf := bufio.NewReader(cmd.InOrStdin())
			confirmInput, err := input.GetString("Key will be deleted. Continue?[y/N]", inBuf)
			if err != nil {
				return err
			}

			if confirmInput != "y" {
				fmt.Println("Cancel")
				return nil
			}

			if err := keybase.Delete(name, "", true); err != nil {
				return err
			}

			fmt.Printf("Deleted key: %s\n", name)
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

func keysShowCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [name]",
		Aliases: []string{"s"},
		Short:   "Show address from name in the keychain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			key, err := keybase.Get(name)
			if err != nil {
				return err
			}
			fmt.Println(key.GetAddress().String())
			return nil
		},
	}
	return cmd
}
