package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// // AddResult validates the result's size and saves it to the store.
// func (k Keeper) AddResult(
// 	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID,
// 	calldata []byte, result []byte,
// ) error {
// 	if uint64(len(result)) > k.GetParam(ctx, types.KeyMaxResultSize) {
// 		return sdkerrors.Wrapf(types.ErrBadDataValue,
// 			"AddResult: Result size (%d) exceeds the maximum size (%d).",
// 			len(result), k.GetParam(ctx, types.KeyMaxResultSize),
// 		)
// 	}

// 	request, err := k.GetRequest(ctx, requestID)
// 	if err != nil {
// 		return err
// 	}

// 	k.SetResult(ctx, requestID, oracleScriptID, calldata, types.NewResult(
// 		request.RequestTime, ctx.BlockTime().Unix(), int64(len(request.RequestedValidators)),
// 		request.SufficientValidatorCount, int64(len(request.ReceivedValidators)), result,
// 	))
// 	return nil
// }

// // SetResult saves the given result of code execution to the store without performing validation.
// func (k Keeper) SetResult(
// 	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID,
// 	calldata []byte, result types.Result,
// ) {
// 	store := ctx.KVStore(k.storeKey)
// 	store.Set(
// 		types.ResultStoreKey(requestID, oracleScriptID, calldata),
// 		result.Hash(),
// 	)
// }

// // GetResult returns the result bytes in the store.
// func (k Keeper) GetResult(
// 	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID,
// 	calldata []byte,
// ) ([]byte, error) {
// 	if !k.HasResult(ctx, requestID, oracleScriptID, calldata) {
// 		return nil, sdkerrors.Wrapf(types.ErrItemNotFound,
// 			"GetResult: Result for request ID %d is not available.", requestID,
// 		)
// 	}
// 	store := ctx.KVStore(k.storeKey)
// 	return store.Get(types.ResultStoreKey(requestID, oracleScriptID, calldata)), nil
// }

// // HasResult checks whether the result at this request id exists in the store.
// func (k Keeper) HasResult(
// 	ctx sdk.Context, requestID types.RequestID, oracleScriptID types.OracleScriptID,
// 	calldata []byte,
// ) bool {
// 	store := ctx.KVStore(k.storeKey)
// 	return store.Has(types.ResultStoreKey(requestID, oracleScriptID, calldata))
// }

func (k Keeper) AddPacketHash(
	ctx sdk.Context, requestID types.RequestID,
	requestPacket types.OracleRequestPacketData,
	responsePacket types.OracleResponsePacketData,
) error {

	if uint64(len([]byte(responsePacket.Result))) > k.GetParam(ctx, types.KeyMaxResultSize) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddResult: Result size (%d) exceeds the maximum size (%d).",
			len([]byte(responsePacket.Result)), k.GetParam(ctx, types.KeyMaxResultSize),
		)
	}
	store := ctx.KVStore(k.storeKey)

	hashData := k.cdc.MustMarshalBinaryBare(requestPacket)
	hashData = append(hashData, k.cdc.MustMarshalBinaryBare(responsePacket)...)
	store.Set(
		types.PacketHashStoreKey(requestID),
		k.cdc.MustMarshalBinaryBare(hashData),
	)
	return nil
}

func (k Keeper) HasPacketHash(
	ctx sdk.Context, requestID types.RequestID,
) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.PacketHashStoreKey(requestID))
}

func (k Keeper) GetPacketHash(
	ctx sdk.Context, requestID types.RequestID,
) ([]byte, error) {
	if !k.HasPacketHash(ctx, requestID) {
		return nil, sdkerrors.Wrapf(types.ErrItemNotFound,
			"GetResult: Result for request ID %d is not available.", requestID,
		)
	}
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.PacketHashStoreKey(requestID)), nil
}
