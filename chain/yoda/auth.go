package yoda

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/yoda/executor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	JWTSecretKey = "jwt-secret-key"
)

func authorizerCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "authorizer",
		Aliases: []string{"auth"},
		Short:   "Manage key held by the oracle process",
	}
	cmd.AddCommand(singedTokenCmd(c))
	cmd.AddCommand(addSecretKey(c))
	return cmd
}

func singedTokenCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signed-token",
		Short: "Generate singed token",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(executor.GetSingedToken(cfg.JWTSecretKey))
			return nil
		},
	}

	return cmd
}

func addSecretKey(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-key",
		Short: "Add secret key",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.Set(JWTSecretKey, args[0])
			return viper.WriteConfig()
		},
	}
	return cmd
}
