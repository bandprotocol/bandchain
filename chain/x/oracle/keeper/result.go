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
func (k Keeper) SetResult(ctx sdk.Context, reqID types.RequestID, result []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ResultStoreKey(reqID), result)
}

// GetResult returns the result for the given request ID or error if not exists.
func (k Keeper) GetResult(ctx sdk.Context, id types.RequestID) (types.Result, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.ResultStoreKey(id))
	if bz == nil {
		return types.Result{}, sdkerrors.Wrapf(types.ErrResultNotFound, "id: %d", id)
	}
	var result types.Result
	err := obi.Decode(bz, &result)
	if err != nil {
		return types.Result{}, types.ErrOBIDecode
	}
	return result, nil
}

// SaveResult saves the result packets for the given request and returns back the response packet.
func (k Keeper) SaveResult(ctx sdk.Context, id types.RequestID, status types.ResolveStatus, result []byte) {
	r := k.MustGetRequest(ctx, id)
	req := types.NewOracleRequestPacketData(
		r.ClientID, r.OracleScriptID, r.Calldata, uint64(len(r.RequestedValidators)), r.MinCount,
	)
	res := types.NewOracleResponsePacketData(
		r.ClientID, id, k.GetReportCount(ctx, id), r.RequestTime.Unix(),
		ctx.BlockTime().Unix(), status, result,
	)
	k.SetResult(ctx, id, obi.MustEncode(req, res))
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRequestExecute,
		sdk.NewAttribute(types.AttributeKeyClientID, req.ClientID),
		sdk.NewAttribute(types.AttributeKeyOracleScriptID, fmt.Sprintf("%d", req.OracleScriptID)),
		sdk.NewAttribute(types.AttributeKeyCalldata, string(req.Calldata)),
		sdk.NewAttribute(types.AttributeKeyAskCount, fmt.Sprintf("%d", req.AskCount)),
		sdk.NewAttribute(types.AttributeKeyMinCount, fmt.Sprintf("%d", req.MinCount)),
		sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", res.RequestID)),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", status)),
		sdk.NewAttribute(types.AttributeKeyAnsCount, fmt.Sprintf("%d", res.AnsCount)),
		sdk.NewAttribute(types.AttributeKeyRequestTime, fmt.Sprintf("%d", res.RequestTime)),
		sdk.NewAttribute(types.AttributeKeyResolveTime, fmt.Sprintf("%d", res.ResolveTime)),
		sdk.NewAttribute(types.AttributeKeyResult, string(res.Result)),
	))
}

// MustGetResult returns the result for the given request ID. Panics on error.
func (k Keeper) MustGetResult(ctx sdk.Context, id types.RequestID) types.Result {
	result, err := k.GetResult(ctx, id)
	if err != nil {
		panic(err)
	}
	return result
}
