package zoracle

import (
	"fmt"

	"github.com/bandprotocol/d3n/chain/owasm"
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates handler of this module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgRequestData:
			return handleMsgRequest(ctx, keeper, msg)
		// case MsgReport:
		// 	return handleMsgReport(ctx, keeper, msg)
		// case MsgStoreCode:
		// 	return handleMsgStoreCode(ctx, keeper, msg)
		// case MsgDeleteCode:
		// 	return handleMsgDeleteCode(ctx, keeper, msg)
		case MsgCreateDataSource:
			return handleMsgCreateDataSource(ctx, keeper, msg)
		case MsgEditDataSource:
			return handleMsgEditDataSource(ctx, keeper, msg)
		case MsgCreateOracleScript:
			return handleMsgCreateOracleScript(ctx, keeper, msg)
		case MsgEditOracleScript:
			return handleMsgEditOracleScript(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized zoracle message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgRequest(ctx sdk.Context, keeper Keeper, msg MsgRequestData) sdk.Result {
	id, err := keeper.AddRequest(
		ctx,
		msg.OracleScriptID,
		msg.Calldata,
		msg.RequestedValidatorCount,
		msg.SufficientValidatorCount,
		msg.Expiration,
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
	_, _, errOwasm := owasm.Execute(&env, script.Code, "prepare", msg.Calldata, 100000)
	if errOwasm != nil {
		// TODO: error
		return sdk.ErrUnknownRequest(errOwasm.Error()).Result()
	}

	err = keeper.ValidateDataSourceCount(ctx, id)
	if err != nil {
		return err.Result()
	}

	// Emit request event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRequest,
			sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", id)),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

// func handleMsgReport(ctx sdk.Context, keeper Keeper, msg MsgReport) sdk.Result {
// 	// check request id is valid.
// 	request, err := keeper.GetRequest(ctx, msg.RequestID)
// 	if err != nil {
// 		return err.Result()
// 	}

// 	storedCode, err := keeper.GetCode(ctx, request.CodeHash)
// 	if err != nil {
// 		return err.Result()
// 	}

// 	// check request is in period of reporting
// 	if uint64(ctx.BlockHeight()) > request.ReportEndAt {
// 		return types.ErrOutOfReportPeriod(types.DefaultCodespace).Result()
// 	}

// 	// Validate sender
// 	validators := keeper.StakingKeeper.GetLastValidators(ctx)

// 	isFound := false
// 	for _, validator := range validators {
// 		if msg.Validator.Equals(validator.GetOperator()) {
// 			isFound = true
// 			break
// 		}
// 	}
// 	if !isFound {
// 		return types.ErrInvalidValidator(types.DefaultCodespace).Result()
// 	}

// 	keeper.SetReport(ctx, msg.RequestID, msg.Validator, msg.Data)
// 	// Emit report event
// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		sdk.NewEvent(
// 			types.EventTypeReport,
// 			sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", msg.RequestID)),
// 			sdk.NewAttribute(types.AttributeKeyCodeName, storedCode.Name),
// 			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
// 		),
// 	})
// 	return sdk.Result{Events: ctx.EventManager().Events()}
// }

// func handleMsgStoreCode(ctx sdk.Context, keeper Keeper, msg MsgStoreCode) sdk.Result {
// 	sc := types.NewStoredCode(msg.Code, msg.Name, msg.Owner)
// 	codeHash := sc.GetCodeHash()
// 	if keeper.CheckCodeHashExists(ctx, codeHash) {
// 		return types.ErrCodeAlreadyExisted(types.DefaultCodespace).Result()
// 	}
// 	keeper.SetCode(ctx, msg.Code, msg.Name, msg.Owner)

// 	// Emit store code event
// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		sdk.NewEvent(
// 			types.EventTypeStoreCode,
// 			sdk.NewAttribute(types.AttributeKeyCodeHash, hex.EncodeToString(codeHash)),
// 			sdk.NewAttribute(types.AttributeKeyCodeName, msg.Name),
// 		),
// 	})
// 	return sdk.Result{Events: ctx.EventManager().Events()}
// }

// func handleMsgDeleteCode(ctx sdk.Context, keeper Keeper, msg MsgDeleteCode) sdk.Result {
// 	storedCode, sdkErr := keeper.GetCode(ctx, msg.CodeHash)
// 	if sdkErr != nil {
// 		return types.ErrCodeHashNotFound(types.DefaultCodespace).Result()
// 	}
// 	if !storedCode.Owner.Equals(msg.Owner) {
// 		return types.ErrInvalidOwner(types.DefaultCodespace).Result()
// 	}

// 	keeper.DeleteCode(ctx, msg.CodeHash)
// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		sdk.NewEvent(
// 			types.EventTypeDeleteCode,
// 			sdk.NewAttribute(types.AttributeKeyCodeHash, hex.EncodeToString(msg.CodeHash)),
// 			sdk.NewAttribute(types.AttributeKeyCodeName, storedCode.Name),
// 		),
// 	})
// 	return sdk.Result{Events: ctx.EventManager().Events()}
// }

// handleMsgCreateDataSource is a function to handle MsgCreateDataSource.
func handleMsgCreateDataSource(ctx sdk.Context, keeper Keeper, msg MsgCreateDataSource) sdk.Result {
	err := keeper.AddDataSource(ctx, msg.Owner, msg.Name, msg.Fee, msg.Executable)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgEditDataSource is a function to handle MsgEditDataSource.
func handleMsgEditDataSource(ctx sdk.Context, keeper Keeper, msg MsgEditDataSource) sdk.Result {
	dataSource, err := keeper.GetDataSource(ctx, msg.DataSourceID)
	if err != nil {
		return err.Result()
	}

	if !dataSource.Owner.Equals(msg.Sender) {
		// TODO: change it later.
		return types.ErrInvalidOwner(types.DefaultCodespace).Result()
	}

	err = keeper.EditDataSource(ctx, msg.DataSourceID, msg.Owner, msg.Name, msg.Fee, msg.Executable)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgCreateOracleScript is a function to handle MsgCreateOracleScript.
func handleMsgCreateOracleScript(ctx sdk.Context, keeper Keeper, msg MsgCreateOracleScript) sdk.Result {
	err := keeper.AddOracleScript(ctx, msg.Owner, msg.Name, msg.Code)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgEditOracleScript is a function to handle MsgEditOracleScript.
func handleMsgEditOracleScript(ctx sdk.Context, keeper Keeper, msg MsgEditOracleScript) sdk.Result {
	oracleScript, err := keeper.GetOracleScript(ctx, msg.OracleScriptID)
	if err != nil {
		return err.Result()
	}

	if !oracleScript.Owner.Equals(msg.Sender) {
		// TODO: change it later.
		return types.ErrInvalidOwner(types.DefaultCodespace).Result()
	}

	err = keeper.EditOracleScript(ctx, msg.OracleScriptID, msg.Owner, msg.Name, msg.Code)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleEndBlock(ctx sdk.Context, keeper Keeper) sdk.Result {
	pendingList := keeper.GetPendingResolveList(ctx)

	for _, requestID := range pendingList {
		request, err := keeper.GetRequest(ctx, requestID)
		if err != nil {
			// Don't expect to happen
			continue
		}

		env, err := NewExecutionEnvironment(ctx, keeper, requestID)
		if err != nil {
			continue
		}

		script, err := keeper.GetOracleScript(ctx, request.OracleScriptID)
		if err != nil {
			continue
		}

		result, _, errOwasm := owasm.Execute(&env, script.Code, "execute", request.Calldata, 100000)
		// TODO: Handle error if happen
		if errOwasm == nil {
			keeper.SetResult(ctx, requestID, result)
			keeper.SetResolve(ctx, requestID, true)
		}
	}

	keeper.SetPendingResolveList(ctx, []int64{})

	// TODO: Emit event
	return sdk.Result{Events: ctx.EventManager().Events()}
}
