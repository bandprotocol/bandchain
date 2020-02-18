package cli

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

const (
	flagName                     = "name"
	flagScript                   = "script"
	flagCallFee                  = "call-fee"
	flagOwner                    = "owner"
	flagCalldata                 = "calldata"
	flagRequestedValidatorCount  = "requested-validator-count"
	flagSufficientValidatorCount = "sufficient-validator-count"
	flagExpiration               = "expiration"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	zoracleCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "zoracle transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	zoracleCmd.AddCommand(client.PostCommands(
		GetCmdCreateDataSource(cdc),
		GetCmdEditDataSource(cdc),
<<<<<<< HEAD
		GetCmdRequest(cdc),
		GetCmdReport(cdc),
=======
		GetCmdCreateOracleScript(cdc),
		GetCmdEditOracleScript(cdc),
>>>>>>> add function GetCmdCreateOracleScript
	)...)

	return zoracleCmd
}

// GetCmdRequest implements the request command handler.
func GetCmdRequest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request [oracle-script-id] (-c [calldata]) (-r [requested-validator-count]) (-v [sufficient-validator-count]) (-x [expiration])",
		Short: "Make a new data request via an existing oracle script",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Make a new request via an existing oracle script with the configuration flags.
Example:
$ %s tx zoracle request 1 -c 1234abcdef -r 4 -v 3 -x 20 --from mykey
$ %s tx zoracle request 1 --calldata 1234abcdef --requested-validator-count 4 --sufficient-validator-count 3 --expiration 20 --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			oracleScriptID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			calldata, err := cmd.Flags().GetBytesHex(flagCalldata)
			if err != nil {
				return err
			}

			requestedValidatorCount, err := cmd.Flags().GetInt64(flagRequestedValidatorCount)
			if err != nil {
				return err
			}

			sufficientValidatorCount, err := cmd.Flags().GetInt64(flagSufficientValidatorCount)
			if err != nil {
				return err
			}

			expiration, err := cmd.Flags().GetInt64(flagExpiration)
			if err != nil {
				return err
			}

			msg := types.NewMsgRequestData(
				oracleScriptID,
				calldata,
				requestedValidatorCount,
				sufficientValidatorCount,
				expiration,
				cliCtx.GetFromAddress(),
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().BytesHexP(flagCalldata, "c", nil, "Calldata used in calling the oracle script")
	cmd.Flags().Int64P(flagRequestedValidatorCount, "r", 0, "Number of top validators that need to report data for this request")
	cmd.MarkFlagRequired(flagRequestedValidatorCount)
	cmd.Flags().Int64P(flagSufficientValidatorCount, "v", 0, "Minimum number of reports sufficient to conclude the request's result")
	cmd.MarkFlagRequired(flagSufficientValidatorCount)
	cmd.Flags().Int64P(flagExpiration, "x", 0, "Maximum block count before the data request is considered expired")
	cmd.MarkFlagRequired(flagExpiration)

	return cmd
}

// GetCmdReport implements the report command handler.
func GetCmdReport(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "report [request-id] ([data]...)",
		Short: "Report raw data for the given request ID",
		Args:  cobra.MinimumNArgs(2),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Report raw data for an unresolved request. All raw data requests must be reported at once.
Example:
$ %s tx zoracle report 1 1:172.5 2:HELLOWORLD --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			requestID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			var dataset []types.RawDataReport
			for _, arg := range args[1:] {
				reportRaw := strings.SplitN(arg, ":", 2)
				if len(reportRaw) != 2 {
					return fmt.Errorf("Invalid report format: %s", reportRaw[0])
				}
				externalID, err := strconv.ParseInt(reportRaw[0], 10, 64)
				if err != nil {
					return err
				}
				dataset = append(dataset, types.NewRawDataReport(externalID, []byte(reportRaw[1])))
			}

			// Sort data reports by external ID
			sort.Slice(dataset, func(i, j int) bool {
				return dataset[i].ExternalDataID < dataset[j].ExternalDataID
			})

			msg := types.NewMsgReportData(requestID, dataset, sdk.ValAddress(cliCtx.GetFromAddress()))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdCreateDataSource implements the create data source command handler.
func GetCmdCreateDataSource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-data-source (--name [name]) (--script [path-to-script]) (--call-fee [fee]) (--owner [owner])",
		Short: "Create a new data source",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new data source that will be used by oracle scripts.
Example:
$ %s tx zoracle create-data-source --name coingecko-price --script ../price.sh --call-fee 100uband --owner band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}
			scriptPath, err := cmd.Flags().GetString(flagScript)
			if err != nil {
				return err
			}
			execBytes, err := ioutil.ReadFile(scriptPath)
			if err != nil {
				return err
			}

			feeStr, err := cmd.Flags().GetString(flagCallFee)
			if err != nil {
				return err
			}

			fee, err := sdk.ParseCoins(feeStr)
			if err != nil {
				return err
			}

			ownerStr, err := cmd.Flags().GetString(flagOwner)
			if err != nil {
				return err
			}
			owner, err := sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateDataSource(
				owner,
				name,
				fee,
				execBytes,
				cliCtx.GetFromAddress(),
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagName, "", "Name of data source")
	cmd.Flags().String(flagScript, "", "Path to data source script")
	cmd.Flags().String(flagCallFee, "", "Fee for querying this data source")
	cmd.Flags().String(flagOwner, "", "Owner of this data source")

	return cmd
}

// GetCmdEditDataSource implements the edit data source command handler.
func GetCmdEditDataSource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-data-source [id] (--name [name]) (--script [path-to-script]) (--call-fee [fee]) (--owner [owner])",
		Short: "Edit data source",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit an existing data source. The caller must be the current data source's owner.
Example:
$ %s tx zoracle edit-data-source 1 --name coingecko-price --script ../price.sh --call-fee 100uband --owner band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}
			scriptPath, err := cmd.Flags().GetString(flagScript)
			if err != nil {
				return err
			}
			execBytes, err := ioutil.ReadFile(scriptPath)
			if err != nil {
				return err
			}

			feeStr, err := cmd.Flags().GetString(flagCallFee)
			if err != nil {
				return err
			}

			fee, err := sdk.ParseCoins(feeStr)
			if err != nil {
				return err
			}

			ownerStr, err := cmd.Flags().GetString(flagOwner)
			if err != nil {
				return err
			}
			owner, err := sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditDataSource(
				id,
				owner,
				name,
				fee,
				execBytes,
				cliCtx.GetFromAddress(),
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagName, "", "Name of data source")
	cmd.Flags().String(flagScript, "", "Path to data source script")
	cmd.Flags().String(flagCallFee, "", "Fee for querying this data source")
	cmd.Flags().String(flagOwner, "", "Owner of this data source")

	return cmd
}

// GetCmdCreateOracleScript implements the create oracle script command handler.
func GetCmdCreateOracleScript(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-oracle-script (--name [name]) (--script [path_to_script]) (--owner [owner])",
		Short: "Create a new oracle script",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create new data source that will be used by oracle scripts.
Example:
$ %s tx zoracle create-data-source --name coingecko-price --script ../price.sh --call-fee 100uband --owner band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}
			scriptPath, err := cmd.Flags().GetString(flagScript)
			if err != nil {
				return err
			}
			execBytes, err := ioutil.ReadFile(scriptPath)
			if err != nil {
				return err
			}

			feeStr, err := cmd.Flags().GetString(flagCallFee)
			if err != nil {
				return err
			}

			fee, err := sdk.ParseCoins(feeStr)
			if err != nil {
				return err
			}

			ownerStr, err := cmd.Flags().GetString(flagOwner)
			if err != nil {
				return err
			}
			owner, err := sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateDataSource(
				owner,
				name,
				fee,
				execBytes,
				cliCtx.GetFromAddress(),
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagName, "", "name of data source")
	cmd.Flags().String(flagScript, "", "path to data source script")
	cmd.Flags().String(flagCallFee, "", "fee for query this data source")
	cmd.Flags().String(flagOwner, "", "owner of this data source")

	return cmd
}
