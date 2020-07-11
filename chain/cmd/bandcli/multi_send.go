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

func mulCoins(coins sdk.Coins, multiplier int64) sdk.Coins {
	var newCoins sdk.Coins
	for _, coin := range coins {
		newCoins = append(newCoins, sdk.NewCoin(coin.Denom, coin.Amount.MulRaw(multiplier)))
	}
	return newCoins
}

// MultiSendTxCmd will create a multi-send tx and sign it with the given key.
func MultiSendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-send [amount] [to_address1] [to_address2] ....",
		Short: "Create and sign a multi-send tx",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}
			msg := bank.NewMsgMultiSend(
				[]bank.Input{bank.NewInput(cliCtx.GetFromAddress(), mulCoins(coins, int64(len(args[1:]))))},
				[]bank.Output{},
			)
			for _, arg := range args[1:] {
				to, err := sdk.AccAddressFromBech32(arg)
				if err != nil {
					return err
				}
				msg.Outputs = append(msg.Outputs, bank.NewOutput(to, coins))
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd = flags.PostCommands(cmd)[0]
	return cmd
}
