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
		GetCmdStoreCode(cdc),
		GetCmdDeleteCode(cdc),
	)...)

	return zoracleCmd
}

// GetCmdRequest implements the request command handler
func GetCmdRequest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "request [report period] [code]",
		Short: "request open api data",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			var code []byte

			reportPeriod, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			code, err = hex.DecodeString(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgRequest(code, reportPeriod, cliCtx.GetFromAddress())
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

			var data []byte

			requestID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			data, err = hex.DecodeString(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgReport(requestID, data, sdk.ValAddress(cliCtx.GetFromAddress().Bytes()))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdStoreCode implements the store code command handler
func GetCmdStoreCode(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "store [code]",
		Short: "store wasm code to chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			code, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgStoreCode(code, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdDeleteCode implements the delete code command handler
func GetCmdDeleteCode(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [code hash]",
		Short: "delete wasm code by codehash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			codeHash, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteCode(codeHash, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
