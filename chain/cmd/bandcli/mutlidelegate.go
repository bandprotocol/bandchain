package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

// multiDelegateCommand sends one transaction with multiple delegation messages.
//
//    [{
//      "to": "bandvaloper1ggmufk3tfrrctr44tg9red3f8hps7nge68z75z",
//      "amount": "100uband"
//    }, {
//      "to": "bandvaloper1asec2q0fyd30kwx6zj7hc5336shmegw0mll724",
//      "amount": "10uband"
//    }]
func multiDelegateCommand(cdc *amino.Codec) *cobra.Command {
	return client.PostCommands(&cobra.Command{
		Use:   "multidelegate [config.json]",
		Short: "Submit a transaction with multiple delegation messages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			type Delegation struct {
				To     string `json:"to"`
				Amount string `json:"amount"`
			}

			var delegations []Delegation
			content, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			err = json.Unmarshal(content, &delegations)
			if err != nil {
				return err
			}

			delAddr := cliCtx.GetFromAddress()
			msgs := []sdk.Msg{}

			for _, delegation := range delegations {
				amount, err := sdk.ParseCoin(delegation.Amount)
				if err != nil {
					return err
				}

				valAddr, err := sdk.ValAddressFromBech32(delegation.To)
				if err != nil {
					return err
				}

				msg := types.NewMsgDelegate(delAddr, valAddr, amount)
				err = msg.ValidateBasic()
				if err != nil {
					return err
				}

				msgs = append(msgs, msg)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
		},
	})[0]
}
