package zoracle

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/owasm"
	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates the msg handler of this module, as required by Cosmos-SDK standard.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
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
			return handleMsgRequestData(ctx, keeper, msg)
		case MsgReportData:
			return handleMsgReportData(ctx, keeper, msg)
		case MsgAddOracleAddress:
			return handleMsgAddOracleAddress(ctx, keeper, msg)
		case MsgRemoveOracleAddress:
			return handleMsgRemoveOracleAddress(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized zoracle message type: %T.", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateDataSource(ctx sdk.Context, keeper Keeper, msg MsgCreateDataSource) sdk.Result {
	dataSourceID, err := keeper.AddDataSource(
		ctx, msg.Owner, msg.Name, msg.Description, msg.Fee, msg.Executable,
	)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateDataSource,
			sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", dataSourceID)),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgEditDataSource(ctx sdk.Context, keeper Keeper, msg MsgEditDataSource) sdk.Result {
	dataSource, err := keeper.GetDataSource(ctx, msg.DataSourceID)
	if err != nil {
		return err.Result()
	}
	if !dataSource.Owner.Equals(msg.Sender) {
		return types.ErrUnauthorizedPermission(
			"handleMsgEditDataSource: Sender (%s) is not data source owner (%s).",
			msg.Sender.String(), dataSource.Owner.String(),
		).Result()
	}
	err = keeper.EditDataSource(
		ctx, msg.DataSourceID, msg.Owner, msg.Name, msg.Description, msg.Fee, msg.Executable,
	)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditDataSource,
			sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", msg.DataSourceID)),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgCreateOracleScript(
	ctx sdk.Context, keeper Keeper, msg MsgCreateOracleScript,
) sdk.Result {
	oracleScriptID, err := keeper.AddOracleScript(
		ctx, msg.Owner, msg.Name, msg.Description, msg.Code,
	)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateOracleScript,
			sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", oracleScriptID)),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgEditOracleScript(ctx sdk.Context, keeper Keeper, msg MsgEditOracleScript) sdk.Result {
	oracleScript, err := keeper.GetOracleScript(ctx, msg.OracleScriptID)
	if err != nil {
		return err.Result()
	}
	if !oracleScript.Owner.Equals(msg.Sender) {
		return types.ErrUnauthorizedPermission(
			"handleMsgEditOracleScript: Sender (%s) is not oracle owner (%s).",
			msg.Sender.String(),
			oracleScript.Owner.String(),
		).Result()
	}
	err = keeper.EditOracleScript(
		ctx, msg.OracleScriptID, msg.Owner, msg.Name, msg.Description, msg.Code,
	)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditOracleScript,
			sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", msg.OracleScriptID)),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgRequestData(ctx sdk.Context, keeper Keeper, msg MsgRequestData) sdk.Result {
	id, err := keeper.AddRequest(
		ctx, msg.OracleScriptID, msg.Calldata, msg.RequestedValidatorCount,
		msg.SufficientValidatorCount, msg.Expiration, msg.ExecuteGas,
	)
	if err != nil {
		return err.Result()
	}

	env, err := NewExecutionEnvironment(ctx, keeper, id)
	if err != nil {
		return err.Result()
	}

	script, err := keeper.GetOracleScript(ctx, msg.OracleScriptID)
	if err != nil {
		return err.Result()
	}

	_, gasUsed, errOwasm := owasm.Execute(&env,
		script.Code, "prepare", msg.Calldata,
		ctx.GasMeter().Limit()-ctx.GasMeter().GasConsumedToLimit())
	if errOwasm != nil {
		return types.ErrBadWasmExecution(
			"handleMsgRequestData: An error occurred while running Owasm prepare.",
		).Result()
	}
	ctx.GasMeter().ConsumeGas(gasUsed, "PrepareRequest")
	err = env.SaveRawDataRequests(ctx, keeper)
	if err != nil {
		return err.Result()
	}

	err = keeper.ValidateDataSourceCount(ctx, id)
	if err != nil {
		return err.Result()
	}
	err = keeper.PayDataSourceFees(ctx, id, msg.Sender)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRequest,
			sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgReportData(ctx sdk.Context, keeper Keeper, msg MsgReportData) sdk.Result {
	err := keeper.AddReport(ctx, msg.RequestID, msg.DataSet, msg.Validator, msg.Reporter)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReport,
			sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", msg.RequestID)),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgAddOracleAddress(ctx sdk.Context, keeper Keeper, msg MsgAddOracleAddress) sdk.Result {
	err := keeper.AddReporter(ctx, msg.Validator, msg.Reporter)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddOracleAddress,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyReporter, msg.Reporter.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgRemoveOracleAddress(
	ctx sdk.Context, keeper Keeper, msg MsgRemoveOracleAddress,
) sdk.Result {
	err := keeper.RemoveReporter(ctx, msg.Validator, msg.Reporter)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveOracleAddress,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyReporter, msg.Reporter.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
