package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
)

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator stakingexported.ValidatorI) (stop bool))
	Validator(ctx sdk.Context, address sdk.ValAddress) stakingexported.ValidatorI
	Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec)
	Jail(ctx sdk.Context, consAddr sdk.ConsAddress)
}

// ScopedKeeper defines the expected scoped keeper.
type ScopedKeeper interface {
	GetCapability(ctx sdk.Context, name string) (*capability.Capability, bool)
}

// PortKeeper defines the expected IBC port keeper.
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capability.Capability
}

// ChannelKeeper defines the expected IBC channel keeper.
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channel.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, channelCap *capability.Capability, packet channelexported.PacketI) error
	PacketExecuted(ctx sdk.Context, chanCap *capability.Capability, packet channelexported.PacketI, acknowledgement []byte) error
	ChanCloseInit(ctx sdk.Context, portID, channelID string, chanCap *capability.Capability) error
	TimeoutExecuted(ctx sdk.Context, chanCap *capability.Capability, packet channelexported.PacketI) error
}
