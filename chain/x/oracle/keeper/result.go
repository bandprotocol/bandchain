package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AddResult validates the result's size and saves it to the store.
func (k Keeper) AddResult(
	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID,
	calldata []byte, result []byte,
) error {
	if uint64(len(result)) > k.GetParam(ctx, types.KeyMaxResultSize) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddResult: Result size (%d) exceeds the maximum size (%d).",
			len(result), k.GetParam(ctx, types.KeyMaxResultSize),
		)
	}

	request := k.MustGetRequest(ctx, requestID)

	k.SetResult(ctx, requestID, oracleScriptID, calldata, types.NewResult(
		request.RequestTime, ctx.BlockTime().Unix(), int64(len(request.RequestedValidators)),
		request.SufficientValidatorCount, k.GetReportCount(ctx, requestID), result,
	))
	return nil
}

// SetResult saves the given result of code execution to the store without performing validation.
func (k Keeper) SetResult(
	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID,
	calldata []byte, result types.Result,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.ResultStoreKey(requestID, oracleScriptID, calldata),
		result.EncodeResult(),
	)
}

// GetResult returns the result bytes in the store.
func (k Keeper) GetResult(
	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID,
	calldata []byte,
) (types.Result, error) {
	if !k.HasResult(ctx, requestID, oracleScriptID, calldata) {
		return types.Result{}, sdkerrors.Wrapf(types.ErrItemNotFound,
			"GetResult: Result for request ID %d is not available.", requestID,
		)
	}
	store := ctx.KVStore(k.storeKey)
	return types.MustDecodeResult(
		store.Get(types.ResultStoreKey(requestID, oracleScriptID, calldata)),
	), nil
}

// HasResult checks whether the result at this request id exists in the store.
func (k Keeper) HasResult(
	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID,
	calldata []byte,
) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ResultStoreKey(requestID, oracleScriptID, calldata))
}
