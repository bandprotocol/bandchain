package oracle

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/bandchain/chain/x/oracle/client/cli"
	"github.com/bandprotocol/bandchain/chain/x/oracle/client/rest"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ port.IBCModule        = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is Band Oracle's module basic object.
type AppModuleBasic struct{}

// Name implements AppModuleBasic interface.
func (AppModuleBasic) Name() string { return ModuleName }

// RegisterCodec implements AppModuleBasic interface.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { RegisterCodec(cdc) }

// DefaultGenesis returns the default genesis state as raw bytes.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the oracle module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, bz json.RawMessage) error {
	var data GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	// Once json successfully marshalled, passes along to genesis.go
	return ValidateGenesis(data)
}

// RegisterRESTRoutes implements AppModuleBasic interface.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, StoreKey)
}

// GetQueryCmd implements AppModuleBasic interface.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

// GetTxCmd implements AppModuleBasic interface.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

// AppModule represents the AppModule for this module.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(k Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
}

// RegisterInvariants implements the AppModule interface.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route implements the AppModule interface.
func (am AppModule) Route() string { return RouterKey }

// NewHandler implements the AppModule interface.
func (am AppModule) NewHandler() sdk.Handler { return NewHandler(am.keeper) }

// QuerierRoute implements the AppModule interface.
func (am AppModule) QuerierRoute() string { return ModuleName }

// NewQuerierHandler implements the AppModule interface.
func (am AppModule) NewQuerierHandler() sdk.Querier { return NewQuerier(am.keeper) }

// BeginBlock implements the AppModule interface.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock implements the AppModule interface.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	handleEndBlock(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

// InitGenesis performs genesis initialization for the oracle module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis returns the current state as genesis raw bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// OnChanOpenInit implements ics-26 IBCModule interface.
func (am AppModule) OnChanOpenInit(
	ctx sdk.Context, order channelexported.Order, connectionHops []string,
	portID string, channelID string, chanCap *capability.Capability,
	counterparty channeltypes.Counterparty, version string,
) error {
	// Claim channel capability passed back by IBC module.
	capPath := ibctypes.ChannelCapabilityPath(portID, channelID)
	if err := am.keeper.ScopedKeeper.ClaimCapability(ctx, chanCap, capPath); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
	}
	return nil
}

// OnChanOpenTry implements ics-26 IBCModule interface.
func (am AppModule) OnChanOpenTry(
	ctx sdk.Context, order channelexported.Order, connectionHops []string, portID string,
	channelID string, chanCap *capability.Capability, counterparty channeltypes.Counterparty,
	version string, counterpartyVersion string,
) error {
	// Claim channel capability passed back by IBC module.
	capPath := ibctypes.ChannelCapabilityPath(portID, channelID)
	if err := am.keeper.ScopedKeeper.ClaimCapability(ctx, chanCap, capPath); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
	}
	return nil
}

// OnChanOpenAck implements ics-26 IBCModule interface.
func (am AppModule) OnChanOpenAck(ctx sdk.Context, portID string, channelID string, counterpartyVersion string) error {
	return nil
}

// OnChanOpenConfirm implements ics-26 IBCModule interface.
func (am AppModule) OnChanOpenConfirm(ctx sdk.Context, portID string, channelID string) error {
	return nil
}

// OnChanCloseInit implements ics-26 IBCModule interface.
func (am AppModule) OnChanCloseInit(ctx sdk.Context, portID string, channelID string) error {
	// Disallow user-initiated channel closing for transfer channels
	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements ics-26 IBCModule interface.
func (am AppModule) OnChanCloseConfirm(ctx sdk.Context, portID string, channelID string) error {
	return nil
}

// OnRecvPacket implements ics-26 IBCModule interface.
func (am AppModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, error) {
	var req OracleRequestPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &req); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal request packet data: %s", err.Error())
	}

	// TODO: Mock data source fee payer
	newMsg := types.NewMsgRequestData(
		req.OracleScriptID, req.Calldata, req.AskCount, req.MinCount, req.ClientID,
		sdk.AccAddress([]byte("Unknown")),
	)
	return handleMsgRequestDataIBC(
		ctx, am.keeper, newMsg, packet.GetDestPort(), packet.GetDestChannel(),
	)
}

// OnAcknowledgementPacket implements ics-26 IBCModule interface.
func (am AppModule) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte) (*sdk.Result, error) {
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

// OnTimeoutPacket implements ics-26 IBCModule interface.
func (am AppModule) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, error) {
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}
