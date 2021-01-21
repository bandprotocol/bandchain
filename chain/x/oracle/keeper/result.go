package keeper

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/GeoDB-Limited/odincore/chain/pkg/obi"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
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

// ResolveSuccess resolves the given request as success with the given result.
func (k Keeper) ResolveSuccess(ctx sdk.Context, id types.RequestID, result []byte, gasUsed uint32) {
	k.SaveResult(ctx, id, types.ResolveStatus_Success, result)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", types.ResolveStatus_Success)),
		sdk.NewAttribute(types.AttributeKeyResult, hex.EncodeToString(result)),
		sdk.NewAttribute(types.AttributeKeyGasUsed, fmt.Sprintf("%d", gasUsed)),
	))
}

// ResolveFailure resolves the given request as failure with the given reason.
func (k Keeper) ResolveFailure(ctx sdk.Context, id types.RequestID, reason string) {
	k.SaveResult(ctx, id, types.ResolveStatus_Failure, []byte{})
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", types.ResolveStatus_Failure)),
		sdk.NewAttribute(types.AttributeKeyReason, reason),
	))
}

// ResolveExpired resolves the given request as expired.
func (k Keeper) ResolveExpired(ctx sdk.Context, id types.RequestID) {
	k.SaveResult(ctx, id, types.ResolveStatus_Expired, []byte{})
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeResolve,
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", id)),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", types.ResolveStatus_Expired)),
	))
}

// SaveResult saves the result packets for the request with the given resolve status and result.
func (k Keeper) SaveResult(
	ctx sdk.Context, id types.RequestID, status types.ResolveStatus, result []byte,
) {
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
}
