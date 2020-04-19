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
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateDataSource:
			return handleMsgCreateDataSource(ctx, k, msg)
		case MsgEditDataSource:
			return handleMsgEditDataSource(ctx, k, msg)
		case MsgCreateOracleScript:
			return handleMsgCreateOracleScript(ctx, k, msg)
		case MsgEditOracleScript:
			return handleMsgEditOracleScript(ctx, k, msg)
		case MsgRequestData:
			return handleMsgRequestData(ctx, k, msg)
		case MsgReportData:
			return handleMsgReportData(ctx, k, msg)
		case MsgAddOracleAddress:
			return handleMsgAddOracleAddress(ctx, k, msg)
		case MsgRemoveOracleAddress:
			return handleMsgRemoveOracleAddress(ctx, k, msg)
		case channeltypes.MsgPacket:
			var requestData OracleRequestPacketData
			if err := types.ModuleCdc.UnmarshalJSON(msg.GetData(), &requestData); err == nil {
				calldata, err := hex.DecodeString(requestData.Calldata)
				if err != nil {
					return nil, err
				}
				newMsg := NewMsgRequestData(
					requestData.OracleScriptID, calldata, requestData.AskCount,
					requestData.MinCount, requestData.ClientID,
					msg.Signer,
				)
				return handleMsgRequestData(
					ctx, k, newMsg, msg.GetDestPort(), msg.GetDestChannel(),
				)
			}
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal oracle packet data")
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgCreateDataSource(ctx sdk.Context, k Keeper, msg MsgCreateDataSource) (*sdk.Result, error) {
	id, err := k.AddDataSource(
		ctx, msg.Owner, msg.Name, msg.Description, msg.Fee, msg.Executable,
	)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeCreateDataSource,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgEditDataSource(ctx sdk.Context, k Keeper, msg MsgEditDataSource) (*sdk.Result, error) {
	dataSource, err := k.GetDataSource(ctx, msg.DataSourceID)
	if err != nil {
		return nil, err
	}
	if !dataSource.Owner.Equals(msg.Sender) {
		return nil, types.ErrEditorNotAuthorized
	}
	err = k.EditDataSource(
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

func handleMsgCreateOracleScript(ctx sdk.Context, k Keeper, msg MsgCreateOracleScript) (*sdk.Result, error) {
	id, err := k.AddOracleScript(ctx, msg.Owner, msg.Name, msg.Description, msg.Code, msg.Schema, msg.SourceCodeURL)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeCreateOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgEditOracleScript(ctx sdk.Context, k Keeper, msg MsgEditOracleScript) (*sdk.Result, error) {
	oracleScript, err := k.GetOracleScript(ctx, msg.OracleScriptID)
	if err != nil {
		return nil, err
	}
	if !oracleScript.Owner.Equals(msg.Sender) {
		return nil, types.ErrEditorNotAuthorized
	}
	err = k.EditOracleScript(ctx, msg.OracleScriptID, msg.Owner, msg.Name, msg.Description, msg.Code, msg.Schema, msg.SourceCodeURL)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeEditOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", msg.OracleScriptID)),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgRequestData(ctx sdk.Context, k Keeper, msg MsgRequestData, ibcData ...string) (*sdk.Result, error) {

	validators, err := k.GetRandomValidators(ctx, int(msg.RequestedValidatorCount))
	if err != nil {
		return nil, err
	}
	// req.RequestedValidators = validators
	expirationHeight := ctx.BlockHeight() + int64(k.GetParam(ctx, types.KeyExpirationBlockCount))

	id, err := k.AddRequest(ctx, types.NewRequest(
		msg.OracleScriptID, msg.Calldata, validators, msg.SufficientValidatorCount,
		ctx.BlockHeight(), ctx.BlockTime().Unix(), expirationHeight, msg.ClientID,
	))

	// TODO: HACK AREA!
	if len(ibcData) == 2 {
		request, _ := k.GetRequest(ctx, id)
		request.SourcePort = ibcData[0]
		request.SourceChannel = ibcData[1]
		k.SetRequest(ctx, id, request)
	}
	// END HACK AREA!

	if err != nil {
		return nil, err
	}

	env, err := NewExecutionEnvironment(ctx, k, id, true, 0)
	if err != nil {
		return nil, err
	}

	script := k.MustGetOracleScript(ctx, msg.OracleScriptID)

	gasPrepare := k.GetParam(ctx, types.KeyPrepareGas)
	ctx.GasMeter().ConsumeGas(gasPrepare, "PrepareRequest")
	_, _, errOwasm := owasm.Execute(&env, script.Code, "prepare", msg.Calldata, 100000)

	if errOwasm != nil {
		return nil, sdkerrors.Wrapf(types.ErrBadWasmExecution,
			"handleMsgRequestData: An error occurred while running Owasm prepare.",
		)
	}

	err = env.SaveRawDataRequests(ctx, k)
	if err != nil {
		return nil, err
	}

	err = k.ValidateDataSourceCount(ctx, id)
	if err != nil {
		return nil, err
	}
	err = k.PayDataSourceFees(ctx, id, msg.Sender)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeRequest,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgReportData(ctx sdk.Context, k Keeper, msg MsgReportData) (*sdk.Result, error) {
	if !k.IsReporter(ctx, msg.Validator, msg.Reporter) {
		return nil, sdkerrors.Wrapf(types.ErrReporterNotAuthorized,
			"val: %s, addr: %s", msg.Reporter.String(), msg.Validator.String())
	}
	err := k.AddReport(ctx, msg.RequestID, types.NewReport(msg.DataSet, msg.Validator))
	if err != nil {
		return nil, err
	}
	if k.ShouldBecomePendingResolve(ctx, msg.RequestID) {
		k.AddPendingRequest(ctx, msg.RequestID)
	}
	ctx.EventManager().EmitEvents(sdk.Events{sdk.NewEvent(
		types.EventTypeReport,
		sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", msg.RequestID)),
		sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
	)})
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgAddOracleAddress(ctx sdk.Context, k Keeper, msg MsgAddOracleAddress) (*sdk.Result, error) {
	err := k.AddReporter(ctx, msg.Validator, msg.Reporter)
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

func handleMsgRemoveOracleAddress(ctx sdk.Context, k Keeper, msg MsgRemoveOracleAddress) (*sdk.Result, error) {
	err := k.RemoveReporter(ctx, msg.Validator, msg.Reporter)
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
