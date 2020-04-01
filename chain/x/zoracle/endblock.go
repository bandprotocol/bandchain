package zoracle

import (
	"fmt"
	"math"

	"github.com/bandprotocol/bandchain/chain/owasm"
	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func addUint64Overflow(a, b uint64) (uint64, bool) {
	if math.MaxUint64-a < b {
		return 0, true
	}

	return a + b, false
}

func setResolveAndAddNewAttribute(ctx sdk.Context, keeper Keeper, requestID RequestID, resolveStatus types.ResolveStatus, newAttributes *[]sdk.Attribute) {
	keeper.SetResolve(ctx, requestID, resolveStatus)
	if resolveStatus == types.Failure {
		*newAttributes = append(*newAttributes, sdk.NewAttribute(fmt.Sprint(requestID), "Failure"))
	} else if resolveStatus == types.Success {
		*newAttributes = append(*newAttributes, sdk.NewAttribute(fmt.Sprint(requestID), "Success"))
	}
}

func handleEndBlock(ctx sdk.Context, keeper Keeper) sdk.Result {
	pendingList := keeper.GetPendingResolveList(ctx)
	endBlockExecuteGasLimit := keeper.GetParam(ctx, types.KeyEndBlockExecuteGasLimit)
	gasConsumed := uint64(0)
	firstUnresolvedRequestIndex := len(pendingList)
	newAttributes := []sdk.Attribute{}

	for i, requestID := range pendingList {
		request, err := keeper.GetRequest(ctx, requestID)
		if err != nil { // should never happen
			setResolveAndAddNewAttribute(ctx, keeper, requestID, types.Failure, &newAttributes)
			continue
		}

		// Discard the request if execute gas is greater than EndBlockExecuteGasLimit.
		if request.ExecuteGas > endBlockExecuteGasLimit {
			setResolveAndAddNewAttribute(ctx, keeper, requestID, types.Failure, &newAttributes)
			continue
		}

		estimatedGasConsumed, overflow := addUint64Overflow(gasConsumed, request.ExecuteGas)
		if overflow || estimatedGasConsumed > endBlockExecuteGasLimit {
			firstUnresolvedRequestIndex = i
			break
		}

		env, err := NewExecutionEnvironment(ctx, keeper, requestID)
		if err != nil { // should never happen
			setResolveAndAddNewAttribute(ctx, keeper, requestID, types.Failure, &newAttributes)
			continue
		}

		script, err := keeper.GetOracleScript(ctx, request.OracleScriptID)
		if err != nil { // should never happen
			setResolveAndAddNewAttribute(ctx, keeper, requestID, types.Failure, &newAttributes)
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
			setResolveAndAddNewAttribute(ctx, keeper, requestID, types.Failure, &newAttributes)
			continue
		}

		errResult := keeper.AddResult(ctx, requestID, request.OracleScriptID, request.Calldata, result)
		if errResult != nil {
			setResolveAndAddNewAttribute(ctx, keeper, requestID, types.Failure, &newAttributes)
			continue
		}

		setResolveAndAddNewAttribute(ctx, keeper, requestID, types.Success, &newAttributes)

	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRequestExecute,
			newAttributes...,
		),
	})
	keeper.SetPendingResolveList(ctx, pendingList[firstUnresolvedRequestIndex:])

	return sdk.Result{Events: ctx.EventManager().Events()}
}
