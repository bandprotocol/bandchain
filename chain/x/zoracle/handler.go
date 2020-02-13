package zoracle

import (
	"fmt"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates handler of this module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		// case MsgRequest:
		// 	return handleMsgRequest(ctx, keeper, msg)
		// case MsgReport:
		// 	return handleMsgReport(ctx, keeper, msg)
		// case MsgStoreCode:
		// 	return handleMsgStoreCode(ctx, keeper, msg)
		// case MsgDeleteCode:
		// 	return handleMsgDeleteCode(ctx, keeper, msg)
		case MsgCreateDataSource:
			return handleMsgCreateDataSource(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized zoracle message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// func handleMsgRequest(ctx sdk.Context, keeper Keeper, msg MsgRequest) sdk.Result {
// 	// Get Code from code hash
// 	storedCode, sdkError := keeper.GetCode(ctx, msg.CodeHash)
// 	if sdkError != nil {
// 		return sdkError.Result()
// 	}

// 	newRequestID := keeper.GetNextRequestID(ctx)

// 	newRequest := types.NewRequest(
// 		msg.CodeHash,
// 		msg.Params,
// 		uint64(ctx.BlockHeight())+msg.ReportPeriod,
// 	)

// 	prepare, err := wasm.Prepare(storedCode.Code, msg.Params)
// 	if err != nil {
// 		return sdk.NewError(types.DefaultCodespace, types.WasmError, err.Error()).Result()
// 	}

// 	// Save Request to state
// 	keeper.SetRequest(ctx, newRequestID, newRequest)
// 	// Add new request to pending bucket
// 	pendingRequests := keeper.GetPendingRequests(ctx)
// 	pendingRequests = append(pendingRequests, newRequestID)
// 	keeper.SetPendingRequests(ctx, pendingRequests)
// 	// Emit request event
// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		sdk.NewEvent(
// 			types.EventTypeRequest,
// 			sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", newRequestID)),
// 			sdk.NewAttribute(types.AttributeKeyCodeHash, hex.EncodeToString(msg.CodeHash)),
// 			sdk.NewAttribute(types.AttributeKeyCodeName, storedCode.Name),
// 			sdk.NewAttribute(types.AttributeKeyPrepare, hex.EncodeToString(prepare)),
// 		),
// 	})
// 	return sdk.Result{Events: ctx.EventManager().Events()}
// }

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

func handleMsgCreateDataSource(ctx sdk.Context, keeper Keeper, msg MsgCreateDataSource) sdk.Result {
	newDataSourceID := keeper.GetNextDataSourceID(ctx)

	newDataSource := types.NewDataSource(
		msg.Owner,
		msg.Name,
		msg.Fee,
		msg.Executable,
	)

	// Save DataSource to state
	keeper.SetDataSource(ctx, newDataSourceID, newDataSource)
	return sdk.Result{}
}

func handleEndBlock(ctx sdk.Context, keeper Keeper) sdk.Result {
	// 	reqIDs := keeper.GetPendingRequests(ctx)
	// 	remainingReqIDs := reqIDs

	// 	for _, reqID := range reqIDs {
	// 		request, err := keeper.GetRequest(ctx, reqID)
	// 		if err != nil {
	// 			return err.Result()
	// 		}

	// 		// pack data from validator together
	// 		packedReport := keeper.GetDataReports(ctx, reqID)
	// 		validatorSize := len(keeper.StakingKeeper.GetLastValidators(ctx))

	// 		// skip this request because it's not end and some validators haven't sent report yet.
	// 		if uint64(ctx.BlockHeight()) < request.ReportEndAt && len(packedReport) < validatorSize {
	// 			continue
	// 		}

	// 		var packedData [][]byte
	// 		for _, report := range packedReport {
	// 			packedData = append(packedData, report.Data)
	// 		}

	// 		storedCode, err := keeper.GetCode(ctx, request.CodeHash)
	// 		if err != nil {
	// 			// remove reqID if can't get code
	// 			remainingReqIDs = remove(remainingReqIDs, reqID)
	// 			continue
	// 		}

	// 		result, errWasm := wasm.Execute(storedCode.Code, request.Params, packedData)
	// 		if errWasm == nil {
	// 			keeper.SetResult(ctx, reqID, request.CodeHash, request.Params, result)
	// 		}

	// 		// remove reqID when set result
	// 		remainingReqIDs = remove(remainingReqIDs, reqID)
	// 	}

	// 	keeper.SetPendingRequests(ctx, remainingReqIDs)

	// TODO: Emit event
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func remove(pending []uint64, removeElement uint64) (ret []uint64) {
	for _, s := range pending {
		if s != removeElement {
			ret = append(ret, s)
		}
	}
	return
}
