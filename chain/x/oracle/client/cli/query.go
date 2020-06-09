package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	clientcmn "github.com/bandprotocol/bandchain/chain/x/oracle/client/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	oracleCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the oracle module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	oracleCmd.AddCommand(flags.GetCommands(
		GetQueryCmdParams(storeKey, cdc),
		GetQueryCmdCounts(storeKey, cdc),
		GetQueryCmdDataSource(storeKey, cdc),
		GetQueryCmdOracleScript(storeKey, cdc),
		GetQueryCmdRequest(storeKey, cdc),
		GetQueryCmdRequestSearch(storeKey, cdc),
	)...)
	return oracleCmd
}

// GetQueryCmdParams implements the query parameters command.
func GetQueryCmdParams(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "params",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", route, types.QueryParams), nil)
			if err != nil {
				return err
			}
			var out types.Params
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetQueryCmdCounts implements the query counts command.
func GetQueryCmdCounts(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "counts",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", route, types.QueryCounts), nil)
			if err != nil {
				return err
			}
			var out types.QueryCountsResult
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetQueryCmdDataSource implements the query data source command.
func GetQueryCmdDataSource(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "data-source [id]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryDataSources, args[0]), nil)
			if err != nil {
				return err
			}
			var out types.DataSource
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetQueryCmdOracleScript implements the query oracle script command.
func GetQueryCmdOracleScript(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "oracle-script [id]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryOracleScripts, args[0]), nil)
			if err != nil {
				return err
			}
			var out types.OracleScript
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetQueryCmdRequest implements the query request command.
func GetQueryCmdRequest(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "request [id]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryRequests, args[0]), nil)
			if err != nil {
				return err
			}
			var out types.QueryRequestResult
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetQueryCmdRequestSearch implements the search request command.
func GetQueryCmdRequestSearch(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "request-search [oracle-script-id] [calldata] [ask-count] [min-count]",
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := clientcmn.QuerySearchLatestRequest(route, cliCtx, args[0], args[1], args[2], args[3])
			if err != nil {
				return err
			}
			var out types.QueryRequestResult
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
