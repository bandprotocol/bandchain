package consuming

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/band-consumer/x/consuming/client/cli"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

// AppModule Basics object
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(DefaultGenesisState())
}

// Validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, bz json.RawMessage) error {
	var data GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	// Once json successfully marshalled, passes along to genesis.go
	return ValidateGenesis(data)
}

// Register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	// rest.RegisterRoutes(ctx, rtr, StoreKey)
}

// Get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

// Get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(k Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) Route() string {
	return RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}
func (am AppModule) QuerierRoute() string {
	return ModuleName
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// Implement IBCModule callbacks
func (am AppModule) OnChanOpenInit(
	ctx sdk.Context,
	order channelexported.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capability.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	// // TODO: Enforce ordering, currently relayers use ORDERED channels

	// if version != types.Version {
	// 	return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, "ics20-1")
	// }

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ScopedKeeper.ClaimCapability(ctx, chanCap, ibctypes.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
	}

	// TODO: escrow
	return nil
}

func (am AppModule) OnChanOpenTry(
	ctx sdk.Context,
	order channelexported.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capability.Capability,
	counterparty channeltypes.Counterparty,
	version,
	counterpartyVersion string,
) error {
	// TODO: Enforce ordering, currently relayers use ORDERED channels

	// if version != types.Version {
	// 	return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, "ics20-1")
	// }

	// if counterpartyVersion != types.Version {
	// 	return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, "ics20-1")
	// }

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ScopedKeeper.ClaimCapability(ctx, chanCap, ibctypes.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
	}

	// TODO: escrow
	return nil
}

func (am AppModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyVersion string,
) error {
	// if counterpartyVersion != types.Version {
	// 	return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, "ics20-1")
	// }
	return nil
}

func (am AppModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (am AppModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for transfer channels
	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

func (am AppModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (am AppModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, error) {
	var responseData oracle.OracleResponsePacketData
	if err := oracle.ModuleCdc.UnmarshalJSON(packet.GetData(), &responseData); err == nil {
		fmt.Println("I GOT DATA", responseData.Result, responseData.ResolveTime)
	}
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil

	// acknowledgement := FungibleTokenPacketAcknowledgement{
	// 	Success: true,
	// 	Error:   "",
	// }
	// if err := am.keeper.OnRecvPacket(ctx, packet, data); err != nil {
	// 	acknowledgement = FungibleTokenPacketAcknowledgement{
	// 		Success: false,
	// 		Error:   err.Error(),
	// 	}
	// }

	// if err := am.keeper.ChannelKeeper.PacketExecuted(ctx, ,packet, acknowledgement.GetBytes()); err != nil {
	// 	return nil, err
	// }

	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		types.EventTypePacket,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	// 		// sdk.NewAttribute(types.AttributeKeyReceiver, data.Receiver),
	// 		// sdk.NewAttribute(types.AttributeKeyValue, data.Amount.String()),
	// 	),
	// )

	// return &sdk.Result{
	// 	Events: ctx.EventManager().Events().ToABCIEvents(),
	// }, nil
}

func (am AppModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) (*sdk.Result, error) {
	// 	var ack FungibleTokenPacketAcknowledgement
	// 	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
	// 		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	// 	}
	// 	var data FungibleTokenPacketData
	// 	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
	// 		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	// 	}

	// 	if err := am.keeper.OnAcknowledgementPacket(ctx, packet, data, ack); err != nil {
	// 		return nil, err
	// 	}

	// 	ctx.EventManager().EmitEvent(
	// 		sdk.NewEvent(
	// 			EventTypePacket,
	// 			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
	// 			sdk.NewAttribute(AttributeKeyReceiver, data.Receiver),
	// 			sdk.NewAttribute(AttributeKeyValue, data.Amount.String()),
	// 			sdk.NewAttribute(AttributeKeyAckSuccess, fmt.Sprintf("%t", ack.Success)),
	// 		),
	// 	)

	// 	if !ack.Success {
	// 		ctx.EventManager().EmitEvent(
	// 			sdk.NewEvent(
	// 				EventTypePacket,
	// 				sdk.NewAttribute(AttributeKeyAckError, ack.Error),
	// 			),
	// 		)
	// 	}

	// 	return &sdk.Result{
	// 		Events: ctx.EventManager().Events().ToABCIEvents(),
	// 	}, nil
	return nil, nil
}

func (am AppModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, error) {
	// 	var data FungibleTokenPacketData
	// 	if err := types.ModuleCdc.UnmarshalBinaryBare(packet.GetData(), &data); err != nil {
	// 		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	// 	}
	// 	// refund tokens
	// 	if err := am.keeper.OnTimeoutPacket(ctx, packet, data); err != nil {
	// 		return nil, err
	// 	}

	// 	ctx.EventManager().EmitEvent(
	// 		sdk.NewEvent(
	// 			EventTypeTimeout,
	// 			sdk.NewAttribute(AttributeKeyRefundReceiver, data.Sender),
	// 			sdk.NewAttribute(AttributeKeyRefundValue, data.Amount.String()),
	// 			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
	// 		),
	// 	)

	// 	return &sdk.Result{
	// 		Events: ctx.EventManager().Events().ToABCIEvents(),
	// 	}, nil
	return nil, nil
}
