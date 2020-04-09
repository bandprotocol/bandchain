package oracle

import (
	"encoding/hex"
	"fmt"

	"github.com/bandprotocol/bandchain/chain/owasm"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
)

// NewHandler creates the msg handler of this module, as required by Cosmos-SDK standard.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateDataSource:
			return handleMsgCreateDataSource(ctx, keeper, msg)
		case MsgEditDataSource:
			return handleMsgEditDataSource(ctx, keeper, msg)
		case MsgCreateOracleScript:
			return handleMsgCreateOracleScript(ctx, keeper, msg)
		case MsgEditOracleScript:
			return handleMsgEditOracleScript(ctx, keeper, msg)
		case MsgRequestData:
			// TODO: Remove this hack!!!
			// Here we assume that call data contains "sourceChannel + data"
			// sourceChannel is always 10 characters
			// sourceChannel := string(msg.Calldata[:10])
			// calldata := hex.EncodeToString(msg.Calldata[10:])
			// sourceChannelEnd, found := keeper.ChannelKeeper.GetChannel(ctx, "oracle", sourceChannel)
			// if !found {
			// 	return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown channel %s port oracle", sourceChannel)
			// }
			// destinationPort := sourceChannelEnd.Counterparty.PortID
			// destinationChannel := sourceChannelEnd.Counterparty.ChannelID
			// sequence, found := keeper.ChannelKeeper.GetNextSequenceSend(ctx, "oracle", sourceChannel)
			// if !found {
			// 	return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown sequence number for channel %s port oracle", sourceChannel)
			// }

			// fmt.Println(msg)
			// packet := NewOracleRequestPacketData(
			// 	msg.OracleScriptID, calldata, msg.RequestedValidatorCount,
			// 	msg.SufficientValidatorCount, msg.Expiration, msg.PrepareGas,
			// 	msg.ExecuteGas,
			// )
			// fmt.Println(packet.GetBytes())
			// err := keeper.ChannelKeeper.SendPacket(ctx, channel.NewPacket(packet.GetBytes(),
			// 	sequence, "oracle", sourceChannel, destinationPort, destinationChannel,
			// 	1000000000, // Arbitrarily high timeout for now
			// ))
			// if err != nil {
			// 	return nil, err
			// }

			// return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
			return handleMsgRequestData(ctx, keeper, msg)
		case MsgReportData:
			return handleMsgReportData(ctx, keeper, msg)
		case MsgAddOracleAddress:
			return handleMsgAddOracleAddress(ctx, keeper, msg)
		case MsgRemoveOracleAddress:
			return handleMsgRemoveOracleAddress(ctx, keeper, msg)
		case channeltypes.MsgPacket:
			var requestData OracleRequestPacketData
			if err := types.ModuleCdc.UnmarshalJSON(msg.GetData(), &requestData); err == nil {
				calldata, err := hex.DecodeString(requestData.Calldata)
				if err != nil {
					return nil, err
				}
				newMsg := NewMsgRequestData(
					requestData.OracleScriptID, calldata, requestData.RequestedValidatorCount,
					requestData.SufficientValidatorCount, requestData.Expiration,
					requestData.PrepareGas, requestData.ExecuteGas,
					sdk.AccAddress([]byte("NOT_IMPORTANT")),
				)
				return handleMsgRequestData(
					ctx, keeper, newMsg, msg.GetDestPort(), msg.GetDestChannel(),
				)
			}
			var responseData OracleResponsePacketData
			if err := types.ModuleCdc.UnmarshalJSON(msg.GetData(), &responseData); err == nil {
				fmt.Println("I GOT DATA", responseData.Result)
				return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
			}
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal oracle packet data")
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgCreateDataSource(
	ctx sdk.Context, keeper Keeper, msg MsgCreateDataSource,
) (*sdk.Result, error) {
	dataSourceID, err := keeper.AddDataSource(
		ctx, msg.Owner, msg.Name, msg.Description, msg.Fee, msg.Executable,
	)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeCreateDataSource,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", dataSourceID)),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgEditDataSource(
	ctx sdk.Context, keeper Keeper, msg MsgEditDataSource,
) (*sdk.Result, error) {
	dataSource, err := keeper.GetDataSource(ctx, msg.DataSourceID)
	if err != nil {
		return nil, err
	}
	if !dataSource.Owner.Equals(msg.Sender) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorizedPermission,
			"%s is not authorized to edit this data source", msg.Sender.String(),
		)
	}
	err = keeper.EditDataSource(
		ctx, msg.DataSourceID, msg.Owner, msg.Name, msg.Description,
		msg.Fee, msg.Executable,
	)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeEditDataSource,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", msg.DataSourceID)),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgCreateOracleScript(
	ctx sdk.Context, keeper Keeper, msg MsgCreateOracleScript,
) (*sdk.Result, error) {
	oracleScriptID, err := keeper.AddOracleScript(
		ctx, msg.Owner, msg.Name, msg.Description, msg.Code,
	)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeCreateOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", oracleScriptID)),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgEditOracleScript(
	ctx sdk.Context, keeper Keeper, msg MsgEditOracleScript,
) (*sdk.Result, error) {
	oracleScript, err := keeper.GetOracleScript(ctx, msg.OracleScriptID)
	if err != nil {
		return nil, err
	}
	if !oracleScript.Owner.Equals(msg.Sender) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorizedPermission,
			"%s is not authorized to edit this oracle script", msg.Sender.String(),
		)
	}
	err = keeper.EditOracleScript(
		ctx, msg.OracleScriptID, msg.Owner, msg.Name, msg.Description, msg.Code,
	)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeEditOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", msg.OracleScriptID)),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgRequestData(
	ctx sdk.Context, keeper Keeper, msg MsgRequestData, ibcData ...string,
) (*sdk.Result, error) {
	id, err := keeper.AddRequest(
		ctx, msg.OracleScriptID, msg.Calldata, msg.RequestedValidatorCount,
		msg.SufficientValidatorCount, msg.Expiration, msg.ExecuteGas,
	)
	// TODO: HACK AREA!
	if len(ibcData) == 2 {
		request, _ := keeper.GetRequest(ctx, id)
		request.SourcePort = ibcData[0]
		request.SourceChannel = ibcData[1]
		keeper.SetRequest(ctx, id, request)
	}
	// END HACK AREA!

	if err != nil {
		return nil, err
	}

	env, err := NewExecutionEnvironment(ctx, keeper, id)
	if err != nil {
		return nil, err
	}

	script, err := keeper.GetOracleScript(ctx, msg.OracleScriptID)
	if err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(msg.PrepareGas, "PrepareRequest")
	_, _, errOwasm := owasm.Execute(&env, script.Code, "prepare", msg.Calldata, msg.PrepareGas)
	if errOwasm != nil {
		return nil, sdkerrors.Wrapf(types.ErrBadWasmExecution,
			"handleMsgRequestData: An error occurred while running Owasm prepare.",
		)
	}

	err = env.SaveRawDataRequests(ctx, keeper)
	if err != nil {
		return nil, err
	}

	err = keeper.ValidateDataSourceCount(ctx, id)
	if err != nil {
		return nil, err
	}
	err = keeper.PayDataSourceFees(ctx, id, msg.Sender)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeRequest,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgReportData(
	ctx sdk.Context, keeper Keeper, msg MsgReportData,
) (*sdk.Result, error) {
	err := keeper.AddReport(ctx, msg.RequestID, msg.DataSet, msg.Validator, msg.Reporter)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeReport,
		sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", msg.RequestID)),
		sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgAddOracleAddress(
	ctx sdk.Context, keeper Keeper, msg MsgAddOracleAddress,
) (*sdk.Result, error) {
	err := keeper.AddReporter(ctx, msg.Validator, msg.Reporter)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeAddOracleAddress,
		sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
		sdk.NewAttribute(types.AttributeKeyReporter, msg.Reporter.String()),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgRemoveOracleAddress(
	ctx sdk.Context, keeper Keeper, msg MsgRemoveOracleAddress,
) (*sdk.Result, error) {
	err := keeper.RemoveReporter(ctx, msg.Validator, msg.Reporter)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeRemoveOracleAddress,
		sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
		sdk.NewAttribute(types.AttributeKeyReporter, msg.Reporter.String()),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}
