package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// HasResult checks if the result of this request ID exists in the storage.
func (k Keeper) HasResult(ctx sdk.Context, id types.RequestID) bool {
	return ctx.KVStore(k.storeKey).Has(types.ResultStoreKey(id))
}

// SetResult sets result to the store.
func (k Keeper) SetResult(ctx sdk.Context, reqID types.RequestID, result types.Result) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ResultStoreKey(reqID), obi.MustEncode(result))
}

// GetResult returns the result for the given request ID or error if not exists.
func (k Keeper) GetResult(ctx sdk.Context, id types.RequestID) (types.Result, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.ResultStoreKey(id))
	if bz == nil {
		return types.Result{}, sdkerrors.Wrapf(types.ErrResultNotFound, "id: %d", id)
	}
	var result types.Result
	obi.MustDecode(bz, &result)
	return result, nil
}

// MustGetResult returns the result for the given request ID. Panics on error.
func (k Keeper) MustGetResult(ctx sdk.Context, id types.RequestID) types.Result {
	result, err := k.GetResult(ctx, id)
	if err != nil {
		panic(err)
	}
	return result
}

// Resolve saves the result packets for the given request and emits the resolve event.
func (k Keeper) Resolve(ctx sdk.Context, id types.RequestID, status types.ResolveStatus, result []byte) {
	r := k.MustGetRequest(ctx, id)
	reqPacket := types.NewOracleRequestPacketData(
		r.ClientID,                         // ClientID
		r.OracleScriptID,                   // OracleScriptID
		r.Calldata,                         // Calldata
		uint64(len(r.RequestedValidators)), // AskCount
		r.MinCount,                         // Mincount
	)
	resPacket := types.NewOracleResponsePacketData(
		r.ClientID,                // ClientID
		id,                        // RequestID
		k.GetReportCount(ctx, id), // AnsCount
		r.RequestTime.Unix(),      // RequestTime
		ctx.BlockTime().Unix(),    // ResolveTime
		status,                    // ResolveStatus
		result,                    // Result
	)
	k.SetResult(ctx, id, types.NewResult(reqPacket, resPacket))
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", status)),
	))
}
