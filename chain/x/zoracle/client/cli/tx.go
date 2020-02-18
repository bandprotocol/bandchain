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
	flagName                = "name"
	flagScript              = "script"
	flagCallFee             = "call-fee"
	flagOwner               = "owner"
	flagCalldata            = "calldata"
	flagRequireValidator    = "require-validator"
	flagSufficientValidator = "sufficient-validator"
	flagExpiration          = "expiration"
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
		GetCmdRequest(cdc),
		GetCmdReport(cdc),
		GetCmdCreateDataSource(cdc),
		GetCmdEditDataSource(cdc),
	)...)

	return zoracleCmd
}

// GetCmdRequest implements the request command handler
func GetCmdRequest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request [oracleScriptID] (--calldata [calldata]) (--require-validator [requestedValidatorCount]) (--sufficient-validator [sufficientValidatorCount]) (--expiration [Expiration])",
		Short: "Request a new request from an existed oracle script",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create new request from an existed oracle script with configure flags.
Example:
$ %s tx zoracle request 3 --calldata 03455448 --require-validator 4 --sufficient-validator 3 --expiration 20 --from mykey
`,
				version.ClientName,
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

			requestedValidatorCount, err := cmd.Flags().GetInt64(flagRequireValidator)
			if err != nil {
				return err
			}

			sufficientValidatorCount, err := cmd.Flags().GetInt64(flagSufficientValidator)
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

	cmd.Flags().BytesHex(flagCalldata, nil, "calldata used in calling oracle script")
	cmd.Flags().Int64(flagRequireValidator, 0, "the number of top validators that need to report for this request")
	cmd.MarkFlagRequired(flagRequireValidator)
	cmd.Flags().Int64(flagSufficientValidator, 0, "the minimum number of reports that require to execute the script")
	cmd.MarkFlagRequired(flagSufficientValidator)
	cmd.Flags().Int64(flagExpiration, 0, "report period before determined as expired request")
	cmd.MarkFlagRequired(flagExpiration)

	return cmd
}

// GetCmdReport implements the report command handler
func GetCmdReport(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "report [requestid] ([data]...)",
		Short: "Report to given request id",
		Args:  cobra.MinimumNArgs(2),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Report to unresolved request, need to report data equal to the number of raw data request
Example:
$ %s tx zoracle report 1 1:172.5 2:{"price":197.6} --from mykey
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

			// Sort report by external ID
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
		Use:   "create-data-source (--name [name]) (--script [path_to_script]) (--call-fee [fee]) (--owner [owner])",
		Short: "Create a new data source",
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

// GetCmdEditDataSource implements the edit data source command handler.
func GetCmdEditDataSource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-data-source [id] (--name [name]) (--script [path_to_script]) (--call-fee [fee]) (--owner [owner])",
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
	cmd.Flags().String(flagName, "", "name of data source")
	cmd.Flags().String(flagScript, "", "path to data source script")
	cmd.Flags().String(flagCallFee, "", "fee for query this data source")
	cmd.Flags().String(flagOwner, "", "owner of this data source")

	return cmd
}
