package keeper

// import (
// 	"fmt"
// 	"time"

// 	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
// 	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
// )

// const (
// 	// DefaultPacketTimeoutHeight is set to 0 to disable block height based timeout.
// 	DefaultPacketTimeoutHeight = 0 // NOTE: in blocks

// 	// DefaultPacketTimeoutTimestampDuration is set to 10 minutes from BandChain's block time.
// 	DefaultPacketTimeoutTimestampDuration = uint64(600 * time.Second) // NOTE: in nanoseconds
// )

// // SendOracleResponse sends the given oracle response packet through the given IBC channel/port.
// func (k Keeper) SendOracleResponse(ctx sdk.Context, portID, channelID string, res types.OracleResponsePacketData) {
// 	logger := k.Logger(ctx)
// 	ch, ok := k.ChannelKeeper.GetChannel(ctx, portID, channelID)
// 	if !ok {
// 		logger.Error(fmt.Sprintf("failed to get channel with id: %s, port: %s", channelID, portID))
// 		return
// 	}

// 	seq, ok := k.ChannelKeeper.GetNextSequenceSend(ctx, portID, channelID)
// 	if !ok {
// 		logger.Error(fmt.Sprintf("failed to get next sequence for channel id: %s, port: %s", channelID, portID))
// 		return
// 	}

// 	chCap, ok := k.ScopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(portID, channelID))
// 	if !ok {
// 		logger.Error(fmt.Sprintf("failed to get channel capability for id: %s, port: %s", channelID, portID))
// 		return
// 	}

// 	packet := channel.NewPacket(
// 		res.GetBytes(), seq, portID, channelID, ch.Counterparty.PortID, ch.Counterparty.ChannelID,
// 		DefaultPacketTimeoutHeight, // TimeoutHeight
// 		uint64(ctx.BlockTime().UnixNano())+DefaultPacketTimeoutTimestampDuration, // TimeoutTimestamp
// 	)
// 	err := k.ChannelKeeper.SendPacket(ctx, chCap, packet)
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("failed to send IBC packet for channel id: %s, port: %s, err: %s", channelID, portID, err.Error()))
// 	}
// }
