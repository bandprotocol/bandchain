package keeper

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// AddResult validates the result's size and saves it to the store.
func (k Keeper) AddResult(
	ctx sdk.Context, requestID types.RequestID,
	requestPacket types.OracleRequestPacketData,
	responsePacket types.OracleResponsePacketData,
) error {
	result, err := hex.DecodeString(responsePacket.Result)
	if err != nil {
		return err
	}
	if uint64(len(result)) > k.GetParam(ctx, types.KeyMaxResultSize) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddResult: Result size (%d) exceeds the maximum size (%d).",
			len(result), k.GetParam(ctx, types.KeyMaxResultSize),
		)
	}
	store := ctx.KVStore(k.storeKey)

	h := sha256.New()
	h.Write(k.cdc.MustMarshalBinaryBare(requestPacket))
	reqPacketHash := h.Sum(nil)

	h = sha256.New()
	h.Write(k.cdc.MustMarshalBinaryBare(responsePacket))
	resPacketHash := h.Sum(nil)

	h = sha256.New()
	h.Write(append(reqPacketHash, resPacketHash...))
	store.Set(
		types.ResultStoreKey(requestID),
		h.Sum(nil),
	)
	return nil
}

// GetResult returns the result bytes in the store.
func (k Keeper) GetResult(
	ctx sdk.Context, requestID types.RequestID,
) ([]byte, error) {
	if !k.HasResult(ctx, requestID) {
		return nil, sdkerrors.Wrapf(types.ErrItemNotFound,
			"GetResult: Result for request ID %d is not available.", requestID,
		)
	}
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.ResultStoreKey(requestID)), nil
}

// HasResult checks whether the result at this request id exists in the store.
func (k Keeper) HasResult(
	ctx sdk.Context, requestID types.RequestID,
) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ResultStoreKey(requestID))
}
