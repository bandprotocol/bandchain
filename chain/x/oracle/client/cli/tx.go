package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

const (
	flagName          = "name"
	flagDescription   = "description"
	flagScript        = "script"
	flagOwner         = "owner"
	flagCalldata      = "calldata"
	flagClientID      = "client-id"
	flagSchema        = "schema"
	flagSourceCodeURL = "url"
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
		GetCmdActivate(cdc),
		GetCmdAddReporters(cdc),
		GetCmdRemoveReporter(cdc),
	)...)

	return oracleCmd
}

// GetCmdRequest implements the request command handler.
func GetCmdRequest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request [oracle-script-id] [ask-count] [min-count] (-c [calldata]) (-m [client-id])",
		Short: "Make a new data request via an existing oracle script",
		Args:  cobra.ExactArgs(3),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Make a new request via an existing oracle script with the configuration flags.
Example:
$ %s tx oracle request 1 4 3 -c 1234abcdef -x 20 -m client-id --from mykey
$ %s tx oracle request 1 4 3 --calldata 1234abcdef --client-id cliend-id --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			int64OracleScriptID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			oracleScriptID := types.OracleScriptID(int64OracleScriptID)

			askCount, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			minCount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			calldata, err := cmd.Flags().GetBytesHex(flagCalldata)
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
				askCount,
				minCount,
				clientID,
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
	cmd.Flags().StringP(flagClientID, "m", "", "Requester can match up the request with response by clientID")

	return cmd
}

// GetCmdCreateDataSource implements the create data source command handler.
func GetCmdCreateDataSource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-data-source (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner])",
		Short: "Create a new data source",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new data source that will be used by oracle scripts.
Example:
$ %s tx oracle create-data-source --name coingecko-price --description "The script that queries crypto price from cryptocompare" --script ../price.sh --owner band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

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
	cmd.Flags().String(flagName, "", "Name of this data source")
	cmd.Flags().String(flagDescription, "", "Description of this data source")
	cmd.Flags().String(flagScript, "", "Path to this data source script")
	cmd.Flags().String(flagOwner, "", "Owner of this data source")

	return cmd
}

// GetCmdEditDataSource implements the edit data source command handler.
func GetCmdEditDataSource(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-data-source [id] (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner])",
		Short: "Edit data source",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit an existing data source. The caller must be the current data source's owner.
Example:
$ %s tx oracle edit-data-source 1 --name coingecko-price --description The script that queries crypto price from cryptocompare --script ../price.sh --owner band15d4apf20449ajvwycq8ruaypt7v6d345n9fpt9 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

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
			execBytes := types.DoNotModifyBytes
			if scriptPath != types.DoNotModify {
				execBytes, err = ioutil.ReadFile(scriptPath)
				if err != nil {
					return err
				}
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
	cmd.Flags().String(flagName, types.DoNotModify, "Name of this data source")
	cmd.Flags().String(flagDescription, types.DoNotModify, "Description of this data source")
	cmd.Flags().String(flagScript, types.DoNotModify, "Path to this data source script")
	cmd.Flags().String(flagOwner, "", "Owner of this data source")

	return cmd
}

// GetCmdCreateOracleScript implements the create oracle script command handler.
func GetCmdCreateOracleScript(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-oracle-script (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner]) (--schema [schema]) (--url [source-code-url])",
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
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

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

			schema, err := cmd.Flags().GetString(flagSchema)
			if err != nil {
				return err
			}

			sourceCodeURL, err := cmd.Flags().GetString(flagSourceCodeURL)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateOracleScript(
				owner,
				name,
				description,
				scriptCode,
				schema,
				sourceCodeURL,
				cliCtx.GetFromAddress(),
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagName, "", "Name of this oracle script")
	cmd.Flags().String(flagDescription, "", "Description of this oracle script")
	cmd.Flags().String(flagScript, "", "Path to this oracle script")
	cmd.Flags().String(flagOwner, "", "Owner of this oracle script")
	cmd.Flags().String(flagSchema, "", "Schema of this oracle script")
	cmd.Flags().String(flagSourceCodeURL, "", "URL for the source code of this oracle script")

	return cmd
}

// GetCmdEditOracleScript implements the editing of oracle script command handler.
func GetCmdEditOracleScript(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-oracle-script [id] (--name [name]) (--description [description]) (--script [path-to-script]) (--owner [owner]) (--schema [schema]) (--url [source-code-url])",
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
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

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
			scriptCode := types.DoNotModifyBytes
			if scriptPath != types.DoNotModify {
				scriptCode, err = ioutil.ReadFile(scriptPath)
				if err != nil {
					return err
				}
			}

			ownerStr, err := cmd.Flags().GetString(flagOwner)
			if err != nil {
				return err
			}
			owner, err := sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				return err
			}

			schema, err := cmd.Flags().GetString(flagSchema)
			if err != nil {
				return err
			}

			sourceCodeURL, err := cmd.Flags().GetString(flagSourceCodeURL)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditOracleScript(
				oracleScriptID,
				owner,
				name,
				description,
				scriptCode,
				schema,
				sourceCodeURL,
				cliCtx.GetFromAddress(),
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagName, types.DoNotModify, "Name of this oracle script")
	cmd.Flags().String(flagDescription, types.DoNotModify, "Description of this oracle script")
	cmd.Flags().String(flagScript, types.DoNotModify, "Path to this oracle script")
	cmd.Flags().String(flagOwner, "", "Owner of this oracle script")
	cmd.Flags().String(flagSchema, types.DoNotModify, "Schema of this oracle script")
	cmd.Flags().String(flagSourceCodeURL, types.DoNotModify, "URL for the source code of this oracle script")

	return cmd
}

// GetCmdActivate implements the activate command handler.
func GetCmdActivate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Activate myself to become an oracle validator.",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Activate myself to become an oracle validator.
Example:
$ %s tx oracle activate --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			validator := sdk.ValAddress(cliCtx.GetFromAddress())
			msg := types.NewMsgActivate(validator)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdAddReporters implements the add reporters command handler.
func GetCmdAddReporters(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-reporters [reporter1] [reporter2] ...",
		Short: "Add agents authorized to submit report transactions.",
		Args:  cobra.MinimumNArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add agents authorized to submit report transactions.
Example:
$ %s tx oracle add-reporters band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			validator := sdk.ValAddress(cliCtx.GetFromAddress())
			msgs := make([]sdk.Msg, len(args))
			for i, raw := range args {
				reporter, err := sdk.AccAddressFromBech32(raw)
				if err != nil {
					return err
				}
				msgs[i] = types.NewMsgAddReporter(
					validator,
					reporter,
				)
				err = msgs[i].ValidateBasic()
				if err != nil {
					return err
				}
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
		},
	}

	return cmd
}

// GetCmdRemoveReporter implements the remove reporter command handler.
func GetCmdRemoveReporter(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-reporter [reporter]",
		Short: "Remove an agent from the list of authorized reporters.",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Remove an agent from the list of authorized reporters.
Example:
$ %s tx oracle remove-reporter band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			validator := sdk.ValAddress(cliCtx.GetFromAddress())
			reporter, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgRemoveReporter(
				validator,
				reporter,
			)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
