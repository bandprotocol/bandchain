package oracle

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
