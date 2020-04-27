package consuming

import (
	"encoding/hex"
	"fmt"

	"github.com/bandprotocol/band-consumer/x/consuming/types"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
)

// NewHandler creates the msg handler of this module, as required by Cosmos-SDK standard.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgRequestData:
			sourceChannelEnd, found := keeper.ChannelKeeper.GetChannel(ctx, "consuming", msg.SourceChannel)
			if !found {
				return nil, sdkerrors.Wrapf(
					sdkerrors.ErrUnknownRequest,
					"unknown channel %s port consuming",
					msg.SourceChannel,
				)
			}
			destinationPort := sourceChannelEnd.Counterparty.PortID
			destinationChannel := sourceChannelEnd.Counterparty.ChannelID
			sequence, found := keeper.ChannelKeeper.GetNextSequenceSend(
				ctx, "consuming", msg.SourceChannel,
			)
			if !found {
				return nil, sdkerrors.Wrapf(
					sdkerrors.ErrUnknownRequest,
					"unknown sequence number for channel %s port oracle",
					msg.SourceChannel,
				)
			}
			packet := oracle.NewOracleRequestPacketData(
				msg.ClientID, msg.OracleScriptID, hex.EncodeToString(msg.Calldata),
				msg.AskCount, msg.MinCount,
			)
			err := keeper.ChannelKeeper.SendPacket(ctx, channel.NewPacket(packet.GetBytes(),
				sequence, "consuming", msg.SourceChannel, destinationPort, destinationChannel,
				1000000000, // Arbitrarily high timeout for now
			))
			if err != nil {
				return nil, err
			}
			return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
		case channeltypes.MsgPacket:
			var responseData oracle.OracleResponsePacketData
			if err := types.ModuleCdc.UnmarshalJSON(msg.GetData(), &responseData); err == nil {
				fmt.Println("I GOT DATA", responseData.Result, responseData.ResolveTime)
				return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
			}
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal oracle packet data")
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}
