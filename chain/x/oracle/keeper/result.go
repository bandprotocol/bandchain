package keeper

import (
	"crypto/sha256"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasReport checks if the result of this request ID exists in the storage.
func (k Keeper) HasResult(ctx sdk.Context, id types.RequestID) bool {
	return ctx.KVStore(k.storeKey).Has(types.ResultStoreKey(id))
}

// GetDataSource returns the result bytes for the given request ID or error if not exists.
func (k Keeper) GetResult(ctx sdk.Context, id types.RequestID) ([]byte, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.ResultStoreKey(id))
	if bz == nil {
		return nil, sdkerrors.Wrapf(types.ErrResultNotFound, "id: %d", id)
	}
	return bz, nil
}

// AddResult validates the result's size and saves it to the store.
func (k Keeper) AddResult(
	ctx sdk.Context, id types.RequestID,
	req types.OracleRequestPacketData, res types.OracleResponsePacketData,
) ([]byte, error) {
	if uint64(len(res.Result)) > k.GetParam(ctx, types.KeyMaxResultSize) {
		return nil, sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddResult: Result size (%d) exceeds the maximum size (%d).",
			len(res.Result), k.GetParam(ctx, types.KeyMaxResultSize),
		)
	}
	store := ctx.KVStore(k.storeKey)

	h := sha256.New()
	h.Write(k.cdc.MustMarshalBinaryBare(req))
	reqPacketHash := h.Sum(nil)

	h = sha256.New()
	h.Write(k.cdc.MustMarshalBinaryBare(res))
	resPacketHash := h.Sum(nil)

	h = sha256.New()
	h.Write(append(reqPacketHash, resPacketHash...))
	resultHash := h.Sum(nil)

	store.Set(types.ResultStoreKey(id), resultHash)
	return resultHash, nil
}

// SetResult set result to the store.
func (k Keeper) SetResult(ctx sdk.Context, regID types.RequestID, result []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ResultStoreKey(regID), result)
}

// GetAllResults return the list of all results in the store, or nil if there is none
func (k Keeper) GetAllResults(ctx sdk.Context) (results [][]byte) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ResultStoreKeyPrefix)
	var previousReqID types.RequestID
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		currentReqID := types.RequestID(sdk.BigEndianToUint64(iterator.Key()[1:]))
		diffReqIDCount := currentReqID - previousReqID

		// Insert nil for the request that result is empty
		for i := int64(0); i < int64(diffReqIDCount)-1; i++ {
			results = append(results, nil)
		}

		results = append(results, iterator.Value())
		previousReqID = currentReqID
	}
	return results
}
