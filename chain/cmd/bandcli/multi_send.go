package main

import (
	"bufio"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/spf13/cobra"
)

// MultiSendTxCmd creates a multi-send tx and signs it with the given key.
func MultiSendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-send [amount] [to_address1] [to_address2] ....",
		Short: "Create and sign a multi-send tx",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// Parse the coins we are trying to send
			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}
			accounts := args[1:]
			inputCoins := sdk.NewCoins()
			outputs := make([]bank.Output, 0, len(accounts))
			for _, acc := range accounts {
				to, err := sdk.AccAddressFromBech32(acc)
				if err != nil {
					return err
				}
				outputs = append(outputs, bank.NewOutput(to, coins))
				inputCoins = inputCoins.Add(coins...)
			}
			msg := bank.NewMsgMultiSend(
				[]bank.Input{bank.NewInput(cliCtx.GetFromAddress(), inputCoins)},
				outputs,
			)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd = flags.PostCommands(cmd)[0]
	return cmd
}
