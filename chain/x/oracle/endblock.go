package oracle

import (
	"github.com/bandprotocol/bandchain/chain/owasm"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// resolveRequest resolves the given request, sends response packet out (if applicable),
// and saves result hash to the store. Assumes that the given request is in a resolvable state.
func resolveRequest(ctx sdk.Context, keeper Keeper, reqID types.RequestID) {
	request := keeper.MustGetRequest(ctx, reqID)
	env, err := NewExecutionEnvironment(ctx, keeper, reqID, false, keeper.GetReportCount(ctx, reqID))
	if err != nil {
		panic(err)
	}
	err = env.LoadRawDataReports(ctx, keeper)
	if err != nil {
		panic(err)
	}
	script, err := keeper.GetOracleScript(ctx, request.OracleScriptID)
	if err != nil {
		panic(err)
	}

	// TODO: Refactor this code. For now we hardcode execute gas to 100k
	executeGas := keeper.GetParam(ctx, KeyExecuteGas)
	result, _, err := owasm.Execute(&env, script.Code, "execute", request.Calldata, executeGas)

	var resolveStatus types.ResolveStatus
	if err != nil {
		resolveStatus = types.Failure
	} else {
		resolveStatus = types.Success
	}
	keeper.ProcessOracleResponse(ctx, reqID, resolveStatus, result)
}

// handleEndBlock cleans up the state during end block. See comment in the implementation!
func handleEndBlock(ctx sdk.Context, keeper Keeper) {
	// Loops through all requests in the resolvable list to resolve all of them!
	for _, reqID := range keeper.GetPendingResolveList(ctx) {
		resolveRequest(ctx, keeper, reqID)
	}
	// Once all the requests are resolved, we can clear the list.
	keeper.SetPendingResolveList(ctx, []types.RequestID{})
	// Lastly, we clean up old data requests from the primary storage.
	keeper.ProcessExpiredRequests(ctx)
}
