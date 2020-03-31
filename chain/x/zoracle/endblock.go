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

func handleEndBlock(ctx sdk.Context, keeper Keeper) sdk.Result {
	pendingList := keeper.GetPendingResolveList(ctx)
	endBlockExecuteGasLimit := keeper.GetParam(ctx, types.KeyEndBlockExecuteGasLimit)
	gasConsumed := uint64(0)
	firstUnresolvedRequestIndex := len(pendingList)
	newAttributes := []sdk.Attribute{}

	for i, requestID := range pendingList {
		request, err := keeper.GetRequest(ctx, requestID)
		if err != nil { // should never happen
			keeper.SetResolve(ctx, requestID, types.Failure)
			newAttributes = append(newAttributes, sdk.NewAttribute(fmt.Sprint(requestID), "Failure"))
			continue
		}

		// Discard the request if execute gas is greater than EndBlockExecuteGasLimit.
		if request.ExecuteGas > endBlockExecuteGasLimit {
			keeper.SetResolve(ctx, requestID, types.Failure)
			newAttributes = append(newAttributes, sdk.NewAttribute(fmt.Sprint(requestID), "Failure"))
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
			newAttributes = append(newAttributes, sdk.NewAttribute(fmt.Sprint(requestID), "Failure"))
			continue
		}

		script, err := keeper.GetOracleScript(ctx, request.OracleScriptID)
		if err != nil { // should never happen
			keeper.SetResolve(ctx, requestID, types.Failure)
			newAttributes = append(newAttributes, sdk.NewAttribute(fmt.Sprint(requestID), "Failure"))
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
			newAttributes = append(newAttributes, sdk.NewAttribute(fmt.Sprint(requestID), "Failure"))
			continue
		}

		errResult := keeper.AddResult(ctx, requestID, request.OracleScriptID, request.Calldata, result)
		if errResult != nil {
			keeper.SetResolve(ctx, requestID, types.Failure)
			newAttributes = append(newAttributes, sdk.NewAttribute(fmt.Sprint(requestID), "Failure"))
			continue
		}

		keeper.SetResolve(ctx, requestID, types.Success)
		newAttributes = append(newAttributes, sdk.NewAttribute(fmt.Sprint(requestID), "Success"))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEndBlock,
			newAttributes...,
		),
	})
	keeper.SetPendingResolveList(ctx, pendingList[firstUnresolvedRequestIndex:])

	return sdk.Result{Events: ctx.EventManager().Events()}
}
