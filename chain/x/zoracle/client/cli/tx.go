package cli

import (
	"encoding/hex"
	"strconv"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
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
	)...)

	return zoracleCmd
}

// GetCmdRequest implements the request command handler
func GetCmdRequest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "request [oracleScriptID] [calldata] [requestedCount] [sufficientValidator] [Expiration]",
		Short: "request open api data",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			oracleScriptID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			calldata, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}

			requestedValidatorCount, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			sufficientValidatorCount, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			expiration, err := strconv.ParseInt(args[4], 10, 64)
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
}

// GetCmdReport implements the report command handler
func GetCmdReport(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "report [requestid] [data]",
		Short: "report data",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			requestID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			var dataset []types.RawDataReport
			err = cdc.UnmarshalJSON([]byte(args[1]), &dataset)
			if err != nil {
				return err
			}

			msg := types.NewMsgReportData(requestID, dataset, sdk.ValAddress(cliCtx.GetFromAddress()))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
