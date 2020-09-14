package yoda

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/yoda/executor"
	"github.com/spf13/cobra"
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
