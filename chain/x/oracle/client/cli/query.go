package cli

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	oracleCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the oracle module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	oracleCmd.AddCommand(flags.GetCommands(
		GetCmdReadRequest(storeKey, cdc),
		GetCmdPendingRequest(storeKey, cdc),
	)...)

	return oracleCmd
}

// GetCmdReadRequest queries request by reqID
func GetCmdReadRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "request",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			reqID := args[0]

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/request/%s", queryRoute, reqID),
				nil,
			)
			if err != nil {
				fmt.Printf("read request fail - %s \n", reqID)
				return nil
			}
			out := types.RawBytes(res)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdPendingRequest queries request in pending state
func GetCmdPendingRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "pending_request",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/pending_request", queryRoute),
				nil,
			)
			if err != nil {
				fmt.Printf("get pending request fail \n")
				return nil
			}

			var out types.U64Array
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
