package oracle

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// resolveRequest resolves the given request, sends response packet out (if applicable),
// and saves result hash to the store. Assumes that the given request is in a resolvable state.
func resolveRequest(ctx sdk.Context, k Keeper, reqID types.RequestID) {
	req := k.MustGetRequest(ctx, reqID)
	env := NewExecEnv(ctx, k, req)
	env.SetReports(k.GetReports(ctx, reqID))
	script := k.MustGetOracleScript(ctx, req.OracleScriptID)
	result, _, err := k.OwasmExecute(env, script.Code, "execute", req.Calldata, types.WasmExecuteGas)

	var res types.OracleResponsePacketData
	if err != nil {
		k.Logger(ctx).Info(fmt.Sprintf(
			"failed to execute request id: %d with error: %s", reqID, err.Error(),
		))
		res = k.ResolveRequest(ctx, reqID, types.ResolveStatus_Failure, nil)
	} else {
		res = k.ResolveRequest(ctx, reqID, types.ResolveStatus_Success, result)
	}
	if req.IBC != nil {
		k.SendOracleResponse(ctx, req.IBC.SourcePort, req.IBC.SourceChannel, res)
	}
}

// handleEndBlock cleans up the state during end block. See comment in the implementation!
func handleEndBlock(ctx sdk.Context, k Keeper) {
	// Loops through all requests in the resolvable list to resolve all of them!
	for _, reqID := range k.GetPendingResolveList(ctx) {
		resolveRequest(ctx, k, reqID)
	}
	// Once all the requests are resolved, we can clear the list.
	k.SetPendingResolveList(ctx, []types.RequestID{})
	// Lastly, we clean up old data requests from the primary storage.
	k.ProcessExpiredRequests(ctx)
}
