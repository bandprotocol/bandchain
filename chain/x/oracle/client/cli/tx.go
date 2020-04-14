package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/spf13/cobra"
)

const (
	flagName                     = "name"
	flagDescription              = "description"
	flagScript                   = "script"
	flagCallFee                  = "call-fee"
	flagOwner                    = "owner"
	flagCalldata                 = "calldata"
	flagRequestedValidatorCount  = "requested-validator-count"
	flagSufficientValidatorCount = "sufficient-validator-count"
	flagClientID                 = "client-id"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	oracleCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "oracle transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	oracleCmd.AddCommand(flags.PostCommands(
		GetCmdCreateDataSource(cdc),
		GetCmdEditDataSource(cdc),
		GetCmdCreateOracleScript(cdc),
		GetCmdEditOracleScript(cdc),
		GetCmdRequest(cdc),
		GetCmdReport(cdc),
		GetCmdAddOracleAddress(cdc),
		GetCmdRemoveOracleAddress(cdc),
	)...)

	return oracleCmd
}

// GetCmdRequest implements the request command handler.
func GetCmdRequest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request [oracle-script-id] (-c [calldata]) (-r [requested-validator-count]) (-v [sufficient-validator-count]) (-m [client-id])",
		Short: "Make a new data request via an existing oracle script",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Make a new request via an existing oracle script with the configuration flags.
Example:
$ %s tx oracle request 1 -c 1234abcdef -r 4 -v 3 -x 20 -m client-id --from mykey
$ %s tx oracle request 1 --calldata 1234abcdef --requested-validator-count 4 --sufficient-validator-count 3 --client-id cliend-id --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			int64OracleScriptID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			oracleScriptID := types.OracleScriptID(int64OracleScriptID)

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

			clientID, err := cmd.Flags().GetString(flagClientID)
			if err != nil {
				return err
			}

			msg := types.NewMsgRequestData(
				oracleScriptID,
				calldata,
				requestedValidatorCount,
				sufficientValidatorCount,
				clientID,
				cliCtx.GetFromAddress(),
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().BytesHexP(flagCalldata, "c", nil, "Calldata used in calling the oracle script")
	cmd.Flags().Int64P(flagRequestedValidatorCount, "r", 0, "Number of top validators that need to report data for this request")
	cmd.MarkFlagRequired(flagRequestedValidatorCount)
	cmd.Flags().Int64P(flagSufficientValidatorCount, "v", 0, "Minimum number of reports sufficient to conclude the request's result")
	cmd.MarkFlagRequired(flagSufficientValidatorCount)
	cmd.Flags().StringP(flagClientID, "m", "", "Requester can match up the request with response by clientID")

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
$ %s tx oracle report 1 1:172.5 2:HELLOWORLD --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			int64RequestID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			requestID := types.RequestID(int64RequestID)

			if err != nil {
				return err
			}

			var dataset []types.RawDataReportWithID
			for _, arg := range args[1:] {
				reportRaw := strings.SplitN(arg, ":", 2)
				if len(reportRaw) != 2 {
					return fmt.Errorf("Invalid report format: %s", reportRaw[0])
				}
				int64ExternalID, err := strconv.ParseInt(reportRaw[0], 10, 64)
				if err != nil {
					return err
				}
				externalID := types.ExternalID(int64ExternalID)

				// TODO: Do not hardcode exit code
				dataset = append(dataset, types.NewRawDataReportWithID(externalID, 0, []byte(reportRaw[1])))
			}

			// Sort data reports by external ID
			sort.Slice(dataset, func(i, j int) bool {
				return dataset[i].ExternalDataID < dataset[j].ExternalDataID
			})

			msg := types.NewMsgReportData(requestID, dataset, sdk.ValAddress(cliCtx.GetFromAddress()), cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdCreateDataSource implements the create data source command handler.
func GetCmdCreateDataSource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-data-source (--name [name]) (--description [description]) (--script [path-to-script]) (--call-fee [fee]) (--owner [owner])",
		Short: "Create a new data source",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new data source that will be used by oracle scripts.
Example:
$ %s tx oracle create-data-source --name coingecko-price --description "The script that queries crypto price from cryptocompare" --script ../price.sh --call-fee 100uband --owner band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(flagDescription)
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
				description,
				fee,
				execBytes,
				cliCtx.GetFromAddress(),
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagName, "", "Name of this data source")
	cmd.Flags().String(flagDescription, "", "Description of this data source")
	cmd.Flags().String(flagScript, "", "Path to this data source script")
	cmd.Flags().String(flagCallFee, "", "Fee for querying this data source")
	cmd.Flags().String(flagOwner, "", "Owner of this data source")

	return cmd
}

// GetCmdEditDataSource implements the edit data source command handler.
func GetCmdEditDataSource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-data-source [id] (--name [name]) (--description [description])(--script [path-to-script]) (--call-fee [fee]) (--owner [owner])",
		Short: "Edit data source",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit an existing data source. The caller must be the current data source's owner.
Example:
$ %s tx oracle edit-data-source 1 --name coingecko-price --description The script that queries crypto price from cryptocompare --script ../price.sh --call-fee 100uband --owner band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			int64ID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			dataSourceID := types.DataSourceID(int64ID)
			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(flagDescription)
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
				dataSourceID,
				owner,
				name,
				description,
				fee,
				execBytes,
				cliCtx.GetFromAddress(),
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagName, "", "Name of this data source")
	cmd.Flags().String(flagDescription, "", "Description of this data source")
	cmd.Flags().String(flagScript, "", "Path to this data source script")
	cmd.Flags().String(flagCallFee, "", "Fee for querying this data source")
	cmd.Flags().String(flagOwner, "", "Owner of this data source")

	return cmd
}

// GetCmdCreateOracleScript implements the create oracle script command handler.
func GetCmdCreateOracleScript(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-oracle-script (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner])",
		Short: "Create a new oracle script that will be used by data requests.",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new oracle script that will be used by data requests.
Example:
$ %s tx oracle create-oracle-script --name eth-price --description "Oracle script for getting Ethereum price" --script ../eth_price.wasm --owner band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}
			description, err := cmd.Flags().GetString(flagDescription)
			if err != nil {
				return err
			}

			scriptPath, err := cmd.Flags().GetString(flagScript)
			if err != nil {
				return err
			}
			scriptCode, err := ioutil.ReadFile(scriptPath)
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

			msg := types.NewMsgCreateOracleScript(
				owner,
				name,
				description,
				scriptCode,
				cliCtx.GetFromAddress(),
				"placeholder schema",
				"placeholder url",
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagName, "", "Name of this oracle script")
	cmd.Flags().String(flagDescription, "", "Description of this oracle script")
	cmd.Flags().String(flagScript, "", "Path to this oracle script")
	cmd.Flags().String(flagOwner, "", "Owner of this oracle script")

	return cmd
}

// GetCmdEditOracleScript implements the editing of oracle script command handler.
func GetCmdEditOracleScript(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-oracle-script [id] (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner])",
		Short: "Edit an existing oracle script that will be used by data requests.",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit an existing oracle script that will be used by data requests.
Example:
$ %s tx oracle edit-oracle-script 1 --name eth-price --description "Oracle script for getting Ethereum price" --script ../eth_price.wasm --owner band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			oracleScriptID := types.OracleScriptID(id)
			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(flagDescription)
			if err != nil {
				return err
			}

			scriptPath, err := cmd.Flags().GetString(flagScript)
			if err != nil {
				return err
			}
			scriptCode, err := ioutil.ReadFile(scriptPath)
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

			msg := types.NewMsgEditOracleScript(
				oracleScriptID,
				owner,
				name,
				description,
				scriptCode,
				cliCtx.GetFromAddress(),
				"placeholder schema",
				"placeholder url",
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagName, "", "Name of this oracle script")
	cmd.Flags().String(flagDescription, "", "Description of this oracle script")
	cmd.Flags().String(flagScript, "", "Path to this oracle script")
	cmd.Flags().String(flagOwner, "", "Owner of this oracle script")

	return cmd
}

// GetCmdAddOracleAddress implements the adding of oracle address command handler.
func GetCmdAddOracleAddress(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-oracle-address [reporter]",
		Short: "Add an agent authorized to submit report transactions.",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add an agent authorized to submit report transactions.
Example:
$ %s tx oracle add-oracle-address band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			validator := sdk.ValAddress(cliCtx.GetFromAddress())
			reporter, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgAddOracleAddress(
				validator,
				reporter,
			)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdRemoveOracleAddress implements the Removing of oracle address command handler.
func GetCmdRemoveOracleAddress(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-oracle-address [reporter]",
		Short: "Remove an agent from the list of authorized reporters.",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Remove an agent from the list of authorized reporters.
Example:
$ %s tx oracle remove-oracle-address band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			validator := sdk.ValAddress(cliCtx.GetFromAddress())
			reporter, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgRemoveOracleAddress(
				validator,
				reporter,
			)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
