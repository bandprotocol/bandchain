package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure oracle environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("TODO ~~~")
			return nil
		},
	}
	return cmd
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "oracled",
		Short: "PLACEHOLDER FOR ORACLED SHORT EXPLANATION",
	}

	rootCmd.AddCommand(
		configCmd(),
		runCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
