package oracle

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// handleEndBlock cleans up the state during end block. See comment in the implementation!
func handleEndBlock(ctx sdk.Context, k Keeper) {
	// Loops through all requests in the resolvable list to resolve all of them!
	for _, reqID := range k.GetPendingResolveList(ctx) {
		k.ResolveRequest(ctx, reqID)
	}
	// Once all the requests are resolved, we can clear the list.
	k.SetPendingResolveList(ctx, []types.RequestID{})
	// Lastly, we clean up data requests that are supposed to be expired.
	k.ProcessExpiredRequests(ctx)
	// TODO: We can remove old requests from storage to optimize state space, using `k.DeleteRequest`
	// and `k.DeleteReports`. We don't do that for now as it is premature optimization at this state.
}
