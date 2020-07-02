package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"

	"github.com/bandprotocol/bandchain/chain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	tmmerkle "github.com/tendermint/tendermint/crypto/merkle"
)

const (
	flagProof = "proof"
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
		GetCmdVerifyProof(cdc),
	)...)

	return bridgeCmd
}

// GetCmdVerifyProof implements the verify proof command handler.
func GetCmdVerifyProof(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify a proof",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// Read file from path
			proofPath, err := cmd.Flags().GetString(flagProof)
			if err != nil {
				return err
			}

			proofData, err := ioutil.ReadFile(proofPath)
			if err != nil {
				return err
			}

			var proof tmmerkle.Proof
			// aminoCdc := makeCodec()
			_ = cdc.UnmarshalJSON([]byte(proofData), &proof)

			fmt.Println("proof", proof)

			msg := types.NewMsgVerifyProof(
				proof,
				cliCtx.GetFromAddress(),
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagProof, "", "path to proof")
	return cmd
}

// func makeCodec() *codec.Codec {
// 	var cdc = codec.New()
// 	sdk.RegisterCodec(cdc)
// 	codec.RegisterCrypto(cdc)
// 	authtypes.RegisterCodec(cdc)
// 	// cdc.RegisterConcrete(sdk.StdTx{}, "cosmos-sdk/StdTx", nil)
// 	return cdc
// }
