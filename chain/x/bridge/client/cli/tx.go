package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/bandprotocol/bandchain/chain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	tmmerkle "github.com/tendermint/tendermint/crypto/merkle"
)

const (
	flagName = "name"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	bridgeCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "bridge transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	bridgeCmd.AddCommand(flags.PostCommands(
		GetCmdRelayAndVerify(cdc),
	)...)

	return bridgeCmd
}

// GetCmdRelayAndVerify implements the relay and verify command handler.
func GetCmdRelayAndVerify(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relay-verify",
		Short: "Make a test",
		Args:  cobra.ExactArgs(0),
		Long: strings.TrimSpace(
			fmt.Sprintf(`test
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			proof := tmmerkle.Proof{
				Ops: []tmmerkle.ProofOp{
					{
						Type: "test",
						Key:  []byte("apple"),
						Data: []byte("appleData"),
					},
					{
						Type: "test2",
						Key:  []byte("kiwi"),
						Data: []byte("kiwiData"),
					},
				},
			}

			// fmt.Println(proof.)

			msg := types.NewMsgRelayAndVerify(
				proof,
				cliCtx.GetFromAddress(),
			)

			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
