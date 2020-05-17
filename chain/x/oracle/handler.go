package oracle

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
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
		case MsgAddReporter:
			return handleMsgAddReporter(ctx, k, msg)
		case MsgRemoveReporter:
			return handleMsgRemoveReporter(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgCreateDataSource(ctx sdk.Context, k Keeper, m MsgCreateDataSource) (*sdk.Result, error) {
	id, err := k.AddDataSource(ctx, types.NewDataSource(
		m.Owner, m.Name, m.Description, m.Executable,
	))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreateDataSource,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
	))
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgEditDataSource(ctx sdk.Context, k Keeper, m MsgEditDataSource) (*sdk.Result, error) {
	dataSource, err := k.GetDataSource(ctx, m.DataSourceID)
	if err != nil {
		return nil, err
	}
	if !dataSource.Owner.Equals(m.Sender) {
		return nil, types.ErrEditorNotAuthorized
	}
	err = k.EditDataSource(ctx, m.DataSourceID, types.NewDataSource(
		m.Owner, m.Name, m.Description, m.Executable,
	))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEditDataSource,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", m.DataSourceID)),
	))
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgCreateOracleScript(ctx sdk.Context, k Keeper, m MsgCreateOracleScript) (*sdk.Result, error) {
	id, err := k.AddOracleScript(ctx, types.NewOracleScript(
		m.Owner, m.Name, m.Description, m.Code, m.Schema, m.SourceCodeURL,
	))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreateOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
	))
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgEditOracleScript(ctx sdk.Context, k Keeper, m MsgEditOracleScript) (*sdk.Result, error) {
	oracleScript, err := k.GetOracleScript(ctx, m.OracleScriptID)
	if err != nil {
		return nil, err
	}
	if !oracleScript.Owner.Equals(m.Sender) {
		return nil, types.ErrEditorNotAuthorized
	}
	err = k.EditOracleScript(ctx, m.OracleScriptID, types.NewOracleScript(
		m.Owner, m.Name, m.Description, m.Code, m.Schema, m.SourceCodeURL,
	))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEditOracleScript,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", m.OracleScriptID)),
	))
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgRequestData(ctx sdk.Context, k Keeper, m MsgRequestData) (*sdk.Result, error) {
	validators, err := k.GetRandomValidators(ctx, int(m.AskCount))
	if err != nil {
		return nil, err
	}

	req := types.NewRequest(
		m.OracleScriptID, m.Calldata, validators, m.MinCount,
		ctx.BlockHeight(), ctx.BlockTime().Unix(), m.ClientID, nil,
	)
	return prepareRequest(ctx, k, m, req)
}

func handleMsgRequestDataIBC(ctx sdk.Context, k Keeper, m MsgRequestData, sourcePort string, sourceChannel string) (*sdk.Result, error) {
	validators, err := k.GetRandomValidators(ctx, int(m.AskCount))
	if err != nil {
		return nil, err
	}

	req := types.NewRequest(
		m.OracleScriptID, m.Calldata, validators, m.MinCount,
		ctx.BlockHeight(), ctx.BlockTime().Unix(), m.ClientID,
		&types.RequestIBC{sourcePort, sourceChannel},
	)

	return prepareRequest(ctx, k, m, req)
}

func prepareRequest(ctx sdk.Context, k Keeper, m MsgRequestData, req types.Request) (*sdk.Result, error) {
	env := NewExecEnv(ctx, k, req)
	script, err := k.GetOracleScript(ctx, m.OracleScriptID)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrOracleScriptNotFound, "id: %d", m.OracleScriptID)
	}
	// gasPrepare := k.GetParam(ctx, types.KeyPrepareGas)
	// ctx.GasMeter().ConsumeGas(gasPrepare, "PrepareRequest")
	// TODO: Consume gas for request fixed size portion

	// TODO: We will need to also validate call data size here
	_, _, err = k.OwasmExecute(env, script.Code, "prepare", m.Calldata, types.WasmPrepareGas)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrBadWasmExecution, err.Error())
	}

	id, err := k.AddRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	event := sdk.NewEvent(types.EventTypeRequest)
	event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)))
	for _, val := range req.RequestedValidators {
		event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyValidator, val.String()))
	}
	ctx.EventManager().EmitEvent(event)

	for _, raw := range env.GetRawRequests() {
		// TODO: Consume more gas if using more raw requests.
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, fmt.Sprintf("%d", raw.DataSourceID)),
			sdk.NewAttribute(types.AttributeKeyExternalID, fmt.Sprintf("%d", raw.ExternalID)),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(raw.Calldata)),
		))
		// TODO: Remove raw request keeper. Make cacher and bandoracled parse from events.
		err = k.AddRawRequest(ctx, id, raw)
		if err != nil {
			return nil, err
		}
	}

	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgReportData(ctx sdk.Context, k Keeper, m MsgReportData) (*sdk.Result, error) {
	if !k.IsReporter(ctx, m.Validator, m.Reporter) {
		return nil, types.ErrReporterNotAuthorized
	}
	err := k.AddReport(ctx, m.RequestID, types.NewReport(m.Validator, m.DataSet))
	if err != nil {
		return nil, err
	}
	req := k.MustGetRequest(ctx, m.RequestID)
	if k.GetReportCount(ctx, m.RequestID) == req.MinCount {
		// At the exact moment when the number of reports is sufficient, we add the request to
		// the pending resolve list. This can happen at most one time for any request.
		k.AddPendingRequest(ctx, m.RequestID)
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeReport,
		sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", m.RequestID)),
		sdk.NewAttribute(types.AttributeKeyValidator, m.Validator.String()),
	))
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgAddReporter(ctx sdk.Context, k Keeper, m MsgAddReporter) (*sdk.Result, error) {
	err := k.AddReporter(ctx, m.Validator, m.Reporter)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAddReporter,
		sdk.NewAttribute(types.AttributeKeyValidator, m.Validator.String()),
		sdk.NewAttribute(types.AttributeKeyReporter, m.Reporter.String()),
	))
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleMsgRemoveReporter(ctx sdk.Context, k Keeper, m MsgRemoveReporter) (*sdk.Result, error) {
	err := k.RemoveReporter(ctx, m.Validator, m.Reporter)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRemoveReporter,
		sdk.NewAttribute(types.AttributeKeyValidator, m.Validator.String()),
		sdk.NewAttribute(types.AttributeKeyReporter, m.Reporter.String()),
	))
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}
