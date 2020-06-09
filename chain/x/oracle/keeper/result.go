package keeper

import (
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

// GetResult returns the result bytes for the given request ID or error if not exists.
func (k Keeper) GetResult(ctx sdk.Context, id types.RequestID) (types.OracleRequestPacketData, types.OracleResponsePacketData, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.ResultStoreKey(id))
	if bz == nil {
		return types.OracleRequestPacketData{}, types.OracleResponsePacketData{}, sdkerrors.Wrapf(types.ErrResultNotFound, "id: %d", id)
	}
	var result types.Result
	err := obi.Decode(bz, &result)
	if err != nil {
		return types.OracleRequestPacketData{}, types.OracleResponsePacketData{}, types.ErrOBIDecode
	}

	return result.RequestPacketData, result.ResponsePacketData, nil
}

// MustGetResult returns the result bytes for the given request ID. Panics on error.
func (k Keeper) MustGetResult(ctx sdk.Context, id types.RequestID) (types.OracleRequestPacketData, types.OracleResponsePacketData) {
	bz := ctx.KVStore(k.storeKey).Get(types.ResultStoreKey(id))
	var result types.Result
	obi.MustDecode(bz, &result)
	return result.RequestPacketData, result.ResponsePacketData
}

// GetAllResults returns the list of all results in the store. Nil will be added for skipped results.
func (k Keeper) GetAllResults(ctx sdk.Context) (results [][]byte) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ResultStoreKeyPrefix)
	var previousReqID types.RequestID
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		currentReqID := types.RequestID(sdk.BigEndianToUint64(iterator.Key()[1:]))
		diffReqIDCount := int(currentReqID - previousReqID)

		// Insert nil for each request without result.
		for i := 0; i < diffReqIDCount-1; i++ {
			results = append(results, nil)
		}

		results = append(results, iterator.Value())
		previousReqID = currentReqID
	}
	return results
}
