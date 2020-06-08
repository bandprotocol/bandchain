package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

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
		GetQueryCmdRequests(storeKey, cdc),
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

// GetQueryCmdRequests implements the search requests command.
func GetQueryCmdRequests(route string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "requests [oracle-script-id]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			events := []string{
				fmt.Sprintf("%s.%s=%s", types.EventTypeRequest, types.AttributeKeyOracleScriptID, args[0]),
			}
			searchResult, err := authclient.QueryTxsByEvents(cliCtx, events, 1, 10, "")
			if err != nil {
				return err
			}
			out := []types.QueryRequestResult{}
			for _, tx := range searchResult.Txs {
				for _, log := range tx.Logs {
					for _, ev := range log.Events {
						if ev.Type != types.EventTypeRequest {
							continue
						}
						for _, attr := range ev.Attributes {
							if attr.Key != types.AttributeKeyID {
								continue
							}
							res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryRequests, attr.Value), nil)
							if err != nil {
								return err
							}
							var req types.QueryRequestResult
							cdc.MustUnmarshalJSON(res, &req)
							out = append(out, req)
						}
					}
				}
			}
			return cliCtx.PrintOutput(out)
		},
	}

	cmd.Flags().Uint32(flags.FlagPage, 1, "Query a specific page of paginated results")
	cmd.Flags().Uint32(flags.FlagLimit, 10, "Query number of transactions results per page returned")
	return cmd
}
