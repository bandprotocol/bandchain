package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
)

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channel.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, packet channelexported.PacketI) error
	PacketExecuted(ctx sdk.Context, packet channelexported.PacketI, acknowledgement channelexported.PacketAcknowledgementI) error
	ChanCloseInit(ctx sdk.Context, portID, channelID string) error
	TimeoutExecuted(ctx sdk.Context, packet channelexported.PacketI) error
}
