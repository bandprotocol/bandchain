package cli

import (
	"bufio"
	"io/ioutil"

	"github.com/bandprotocol/bandchain/chain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
)

const (
	flagProof      = "proof"
	flagRelay      = "relay"
	flagValidators = "validators"
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
		GetCmdUpdateChainID(cdc),
		GetCmdUpdateValidators(cdc),
		GetCmdVerifyProof(cdc),
		GetCmdRelay(cdc),
	)...)

	return bridgeCmd
}

// GetCmdUpdateChainID implements the update chain id command handler.
func GetCmdUpdateChainID(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-chain-id [chain-id]",
		Short: "Update chain id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgUpdateChainID(
				args[0],
				cliCtx.GetFromAddress(),
			)

			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagRelay, "", "path to relay")
	return cmd
}

// GetCmdUpdateValidators implements the update validators command handler.
func GetCmdUpdateValidators(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-validators",
		Short: "Update validators on BandChain",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// Read file from path
			validatorsPath, err := cmd.Flags().GetString(flagValidators)
			if err != nil {
				return err
			}

			validatorsData, err := ioutil.ReadFile(validatorsPath)
			if err != nil {
				return err
			}

			var msg types.MsgUpdateValidators
			aminoCdc := makeCodec()
			_ = aminoCdc.UnmarshalJSON([]byte(validatorsData), &msg)

			msg.Sender = cliCtx.GetFromAddress()

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagValidators, "", "path to validators")
	return cmd
}

// GetCmdRelay implements the relay block command handler.
func GetCmdRelay(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relay",
		Short: "Relay a block on BandChain",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// Read file from path
			relayPath, err := cmd.Flags().GetString(flagRelay)
			if err != nil {
				return err
			}

			relayData, err := ioutil.ReadFile(relayPath)
			if err != nil {
				return err
			}

			var msg types.MsgRelay
			aminoCdc := makeCodec()
			_ = aminoCdc.UnmarshalJSON([]byte(relayData), &msg)
			msg.Sender = cliCtx.GetFromAddress()

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagRelay, "", "path to relay")
	return cmd
}

// GetCmdVerifyProof implements the verify proof command handler.
func GetCmdVerifyProof(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify a proof from json file",
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

			var msg types.MsgVerifyProof
			aminoCdc := makeCodec()
			_ = aminoCdc.UnmarshalJSON([]byte(proofData), &msg)
			msg.Sender = cliCtx.GetFromAddress()

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

func makeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	authtypes.RegisterCodec(cdc)
	// cdc.RegisterConcrete(sdk.StdTx{}, "cosmos-sdk/StdTx", nil)
	return cdc
}
