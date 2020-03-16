package zoracle

import (
	"fmt"
	"math"

	"github.com/bandprotocol/d3n/chain/owasm"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates handler of this module
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
		default:
			errMsg := fmt.Sprintf("unrecognized zoracle message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgCreateDataSource is a function to handle MsgCreateDataSource.
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

// handleMsgEditDataSource is a function to handle MsgEditDataSource.
func handleMsgEditDataSource(ctx sdk.Context, keeper Keeper, msg MsgEditDataSource) sdk.Result {
	dataSource, err := keeper.GetDataSource(ctx, msg.DataSourceID)
	if err != nil {
		return err.Result()
	}

	if !dataSource.Owner.Equals(msg.Sender) {
		return types.ErrUnauthorizedPermission(
			"handleMsgEditDataSource: Sender (%s) is not data source owner (%s).",
			msg.Sender.String(),
			dataSource.Owner.String(),
		).Result()
	}

	err = keeper.EditDataSource(ctx, msg.DataSourceID, msg.Owner, msg.Name, msg.Description, msg.Fee, msg.Executable)
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

// handleMsgCreateOracleScript is a function to handle MsgCreateOracleScript.
func handleMsgCreateOracleScript(ctx sdk.Context, keeper Keeper, msg MsgCreateOracleScript) sdk.Result {
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

// handleMsgEditOracleScript is a function to handle MsgEditOracleScript.
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

	err = keeper.EditOracleScript(ctx, msg.OracleScriptID, msg.Owner, msg.Name, msg.Description, msg.Code)
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

func addUint64Overflow(a, b uint64) (uint64, bool) {
	if math.MaxUint64-a < b {
		return 0, true
	}

	return a + b, false
}

func handleEndBlock(ctx sdk.Context, keeper Keeper) sdk.Result {
	pendingList := keeper.GetPendingResolveList(ctx)
	endBlockExecuteGasLimit := keeper.EndBlockExecuteGasLimit(ctx)
	gasConsumed := uint64(0)
	firstUnresolvedRequestIndex := len(pendingList)

	for i, requestID := range pendingList {
		request, err := keeper.GetRequest(ctx, requestID)
		if err != nil { // should never happen
			keeper.SetResolve(ctx, requestID, types.Failure)
			continue
		}

		// Discard the request if execute gas is greater than EndBlockExecuteGasLimit.
		if request.ExecuteGas > endBlockExecuteGasLimit {
			keeper.SetResolve(ctx, requestID, types.Failure)
			continue
		}

		estimatedGasConsumed, overflow := addUint64Overflow(gasConsumed, request.ExecuteGas)
		if overflow || estimatedGasConsumed > endBlockExecuteGasLimit {
			firstUnresolvedRequestIndex = i
			break
		}

		env, err := NewExecutionEnvironment(ctx, keeper, requestID)
		if err != nil { // should never happen
			keeper.SetResolve(ctx, requestID, types.Failure)
			continue
		}

		script, err := keeper.GetOracleScript(ctx, request.OracleScriptID)
		if err != nil { // should never happen
			keeper.SetResolve(ctx, requestID, types.Failure)
			continue
		}

		result, gasUsed, errOwasm := owasm.Execute(
			&env, script.Code, "execute", request.Calldata, request.ExecuteGas,
		)

		if gasUsed > request.ExecuteGas {
			gasUsed = request.ExecuteGas
		}

		gasConsumed, overflow = addUint64Overflow(gasConsumed, gasUsed)
		// Must never overflow because we already checked for overflow above with
		// gasConsumed + request.ExecuteGas (which is >= gasUsed).
		if overflow {
			panic(sdk.ErrorGasOverflow{Descriptor: "ExecuteRequest"})
		}

		if errOwasm != nil {
			keeper.SetResolve(ctx, requestID, types.Failure)
			continue
		}

		errResult := keeper.AddResult(ctx, requestID, request.OracleScriptID, request.Calldata, result)
		if errResult != nil {
			keeper.SetResolve(ctx, requestID, types.Failure)
			continue
		}

		keeper.SetResolve(ctx, requestID, types.Success)
	}

	keeper.SetPendingResolveList(ctx, pendingList[firstUnresolvedRequestIndex:])

	// TODO: Emit event
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgRequestData(ctx sdk.Context, keeper Keeper, msg MsgRequestData) sdk.Result {
	id, err := keeper.AddRequest(
		ctx,
		msg.OracleScriptID,
		msg.Calldata,
		msg.RequestedValidatorCount,
		msg.SufficientValidatorCount,
		msg.Expiration,
		msg.ExecuteGas,
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

	ctx.GasMeter().ConsumeGas(msg.PrepareGas, "PrepareRequest")
	_, _, errOwasm := owasm.Execute(&env, script.Code, "prepare", msg.Calldata, msg.PrepareGas)
	if errOwasm != nil {
		return types.ErrBadWasmExecution(
			"handleMsgRequestData: An error occured while running Owasm prepare.",
		).Result()
	}

	err = keeper.ValidateDataSourceCount(ctx, id)
	if err != nil {
		return err.Result()
	}

	// Emit request event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRequest,
			sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgReportData(ctx sdk.Context, keeper Keeper, msg MsgReportData) sdk.Result {
	// Save new report to store
	err := keeper.AddReport(ctx, msg.RequestID, msg.DataSet, msg.Sender)
	if err != nil {
		return err.Result()
	}

	// Emit report event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeReport,
			sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", msg.RequestID)),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Sender.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
