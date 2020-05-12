package main

import (
	"github.com/spf13/cobra"
)

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure oracle environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			// logger.Info("TODO ~~~")
			return nil
		},
	}
	return cmd
}
