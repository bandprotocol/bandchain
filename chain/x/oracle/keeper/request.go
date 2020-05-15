package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// HasRequest checks if the request of this ID exists in the storage.
func (k Keeper) HasRequest(ctx sdk.Context, id types.RequestID) bool {
	return ctx.KVStore(k.storeKey).Has(types.RequestStoreKey(id))
}

// GetRequest returns the request struct for the given ID or error if not exists.
func (k Keeper) GetRequest(ctx sdk.Context, id types.RequestID) (types.Request, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.RequestStoreKey(id))
	if bz == nil {
		return types.Request{}, sdkerrors.Wrapf(types.ErrRequestNotFound, "id: %d", id)
	}
	var request types.Request
	k.cdc.MustUnmarshalBinaryBare(bz, &request)
	return request, nil
}

// MustGetRequest returns the request struct for the given ID. Panics error if not exists.
func (k Keeper) MustGetRequest(ctx sdk.Context, id types.RequestID) types.Request {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		panic(err)
	}
	return request
}

// SetRequest saves the given data request to the store without performing any validation.
func (k Keeper) SetRequest(ctx sdk.Context, id types.RequestID, request types.Request) {
	ctx.KVStore(k.storeKey).Set(types.RequestStoreKey(id), k.cdc.MustMarshalBinaryBare(request))
}

// DeleteRequest removes the given data request from the store.
func (k Keeper) DeleteRequest(ctx sdk.Context, id types.RequestID) {
	ctx.KVStore(k.storeKey).Delete(types.RequestStoreKey(id))
}

// GetRandomValidators returns a pseudorandom list of active validators. Each validator has
// chance of getting selected directly proportional to the amount of voting power it has.
func (k Keeper) GetRandomValidators(ctx sdk.Context, size int) ([]sdk.ValAddress, error) {
	// TODO: Make this function actually return random validators.
	validatorsByPower := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	if len(validatorsByPower) < size {
		return nil, sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddRequest: Requested validator count (%d) exceeds the number of validators (%d).",
			size, len(validatorsByPower),
		)
	}

	validators := make([]sdk.ValAddress, size)
	for i := 0; i < size; i++ {
		validators[i] = validatorsByPower[i].GetOperator()
	}
	return validators, nil
}

// AddRequest attempts to create and save a new request. Returns error if some conditions failed.
func (k Keeper) AddRequest(ctx sdk.Context, req types.Request) (types.RequestID, error) {
	if !k.HasOracleScript(ctx, req.OracleScriptID) {
		return 0, sdkerrors.Wrapf(types.ErrOracleScriptNotFound, "id: %d", req.OracleScriptID)
	}
	id := k.GetNextRequestID(ctx)
	k.SetRequest(ctx, id, req)
	return id, nil
}

// ResolveRequest updates the request with resolve status and result, and saves the commitment
// pair of oracle request/response packets to the store. Returns back the response packet.
func (k Keeper) ResolveRequest(
	ctx sdk.Context, id types.RequestID, status types.ResolveStatus, result []byte,
) types.OracleResponsePacketData {

	request := k.MustGetRequest(ctx, id)
	req := types.NewOracleRequestPacketData(
		request.ClientID, request.OracleScriptID,
		hex.EncodeToString(request.Calldata), request.MinCount,
		int64(len(request.RequestedValidators)),
	)
	res := types.NewOracleResponsePacketData(
		request.ClientID, id, int64(k.GetReportCount(ctx, id)), request.RequestTime,
		ctx.BlockTime().Unix(), types.ResolveStatus_Success, hex.EncodeToString(result),
	)

	if status != types.ResolveStatus_Success {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRequestExecute,
				sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", id)),
				sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", status)),
			))
		return res
	}

	resultHash, err := k.AddResult(ctx, id, req, res)
	if err != nil {
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeRequestExecute,
			sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", id)),
			sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", types.ResolveStatus_Failure)),
		))
		return res
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRequestExecute,
		sdk.NewAttribute(types.AttributeKeyClientID, req.ClientID),
		sdk.NewAttribute(types.AttributeKeyOracleScriptID, fmt.Sprintf("%d", req.OracleScriptID)),
		sdk.NewAttribute(types.AttributeKeyCalldata, req.Calldata),
		sdk.NewAttribute(types.AttributeKeyAskCount, fmt.Sprintf("%d", req.AskCount)),
		sdk.NewAttribute(types.AttributeKeyMinCount, fmt.Sprintf("%d", req.MinCount)),
		sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", res.RequestID)),
		sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", types.ResolveStatus_Success)),
		sdk.NewAttribute(types.AttributeKeyAnsCount, fmt.Sprintf("%d", res.AnsCount)),
		sdk.NewAttribute(types.AttributeKeyRequestTime, fmt.Sprintf("%d", request.RequestTime)),
		sdk.NewAttribute(types.AttributeKeyResolveTime, fmt.Sprintf("%d", res.ResolveTime)),
		sdk.NewAttribute(types.AttributeKeyResult, res.Result),
		sdk.NewAttribute(types.AttributeKeyResultHash, hex.EncodeToString(resultHash)),
	))
	return res
}

// ProcessExpiredRequests removes all expired data requests from the store, and sends oracle
// response packets for the ones that have never been resolved.
func (k Keeper) ProcessExpiredRequests(ctx sdk.Context) {
	iter := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.RequestStoreKeyPrefix)
	if !iter.Valid() { // No request currently in the store.
		return
	}
	currentReqID := types.RequestID(sdk.BigEndianToUint64(iter.Key()[1:])) // First available request ID
	lastReqID := types.RequestID(k.GetRequestCount(ctx))
	expirationBlockCount := int64(k.GetParam(ctx, types.KeyExpirationBlockCount))
	// Loop through all data requests in chronological order. If a request reaches its
	// expiration time, it will be removed from the storage. Note that we will need to
	// send oracle response packets with status EXPIRED for those that are not yet resolved.
	for ; currentReqID <= lastReqID; currentReqID++ {
		request := k.MustGetRequest(ctx, currentReqID)
		// This request is not yet expired, so there's nothing to do here. Ditto for
		// all other requests that come after this. Thus we can just break the loop.
		if request.RequestHeight+expirationBlockCount > ctx.BlockHeight() {
			break
		}
		// If the number of reports still doesn't reach the minimum, that means this request
		// is never resolved. Here we process the response as EXPIRED.
		if k.GetReportCount(ctx, currentReqID) < request.MinCount {
			res := k.ResolveRequest(ctx, currentReqID, types.ResolveStatus_Expired, nil)
			if request.IBC != nil {
				k.SendOracleResponse(ctx, request.IBC.SourcePort, request.IBC.SourceChannel, res)
			}
		}
		// We are done with this request. Remove it and its dependencies from the store.
		k.DeleteRequest(ctx, currentReqID)
		k.DeleteRawRequests(ctx, currentReqID)
		k.DeleteReports(ctx, currentReqID)
	}
}

// AddPendingRequest adds the request to the pending list. DO NOT add same request more than once.
func (k Keeper) AddPendingRequest(ctx sdk.Context, id types.RequestID) {
	pendingList := k.GetPendingResolveList(ctx)
	pendingList = append(pendingList, id)
	k.SetPendingResolveList(ctx, pendingList)
}

// SetPendingResolveList saves the list of pending request that will be resolved at end block.
func (k Keeper) SetPendingResolveList(ctx sdk.Context, ids []types.RequestID) {
	bz := k.cdc.MustMarshalBinaryBare(ids)
	if bz == nil {
		bz = []byte{}
	}
	ctx.KVStore(k.storeKey).Set(types.PendingResolveListStoreKey, bz)
}

// GetPendingResolveList returns the list of pending requests to be executed during EndBlock.
func (k Keeper) GetPendingResolveList(ctx sdk.Context) (ids []types.RequestID) {
	bz := ctx.KVStore(k.storeKey).Get(types.PendingResolveListStoreKey)
	if len(bz) == 0 { // Return an empty list if the key does not exist in the store.
		return []types.RequestID{}
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &ids)
	return ids
}
