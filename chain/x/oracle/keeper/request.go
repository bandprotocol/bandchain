package keeper

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/pkg/bandrng"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
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
func (k Keeper) GetRandomValidators(ctx sdk.Context, size int, nextReqID int64) ([]sdk.ValAddress, error) {
	valOperators := []sdk.ValAddress{}
	valPowers := []uint64{}
	k.StakingKeeper.IterateBondedValidatorsByPower(ctx, func(idx int64, val exported.ValidatorI) (stop bool) {
		valOperators = append(valOperators, val.GetOperator())
		valPowers = append(valPowers, val.GetTokens().Uint64())
		return false
	})
	if len(valOperators) < size {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientValidators, "%d < %d", len(valOperators), size)
	}
	seed := fmt.Sprintf("%x:%d:%d",
		ctx.BlockHeader().LastBlockId.Hash, ctx.BlockHeader().Time.Nanosecond(), nextReqID,
	)
	luckyValidatorIndexes := bandrng.ChooseK(bandrng.NewRng(seed), valPowers, size)
	validators := make([]sdk.ValAddress, size)
	for i, idx := range luckyValidatorIndexes {
		validators[i] = valOperators[idx]
	}
	return validators, nil
}

// AddRequest attempts to create and save a new request.
func (k Keeper) AddRequest(ctx sdk.Context, req types.Request) types.RequestID {
	id := k.GetNextRequestID(ctx)
	k.SetRequest(ctx, id, req)
	return id
}

// SaveResult updates the request with resolve status and result, and saves the commitment
// pair of oracle request/response packets to the store. Returns back the response packet.
func (k Keeper) SaveResult(ctx sdk.Context, id types.RequestID, status types.ResolveStatus, result []byte) types.OracleResponsePacketData {
	r := k.MustGetRequest(ctx, id)
	req := types.NewOracleRequestPacketData(
		r.ClientID, r.OracleScriptID, r.Calldata, uint64(len(r.RequestedValidators)), r.MinCount,
	)
	res := types.NewOracleResponsePacketData(
		r.ClientID, id, k.GetReportCount(ctx, id), r.RequestTime,
		ctx.BlockTime().Unix(), status, result,
	)
	k.SetResult(ctx, id, types.CalculateEncodedResult(req, res))
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
	return res
}

// ProcessExpiredRequests resolves and sends response packets for all expired-but-unresolved requests.
func (k Keeper) ProcessExpiredRequests(ctx sdk.Context) {
	currentReqID := types.RequestID(k.GetRequestLastExpired(ctx) + 1)
	lastReqID := types.RequestID(k.GetRequestCount(ctx))
	expirationBlockCount := int64(k.GetParam(ctx, types.KeyExpirationBlockCount))
	// Loop through all data requests in chronological order. If a request reaches its
	// expiration time, it will be removed from the storage. Note that we will need to
	// send oracle response packets with status EXPIRED for those that are not yet resolved.
	for ; currentReqID <= lastReqID; currentReqID++ {
		req := k.MustGetRequest(ctx, currentReqID)
		// This request is not yet expired, so there's nothing to do here. Ditto for
		// all other requests that come after this. Thus we can just break the loop.
		if req.RequestHeight+expirationBlockCount > ctx.BlockHeight() {
			break
		}
		// If the number of reports still doesn't reach the minimum, that means this request
		// is never resolved. Here we process the response as EXPIRED.
		if k.GetReportCount(ctx, currentReqID) < req.MinCount {
			// res := k.SaveResult(ctx, currentReqID, types.ResolveStatus_Expired, nil)
			k.SaveResult(ctx, currentReqID, types.ResolveStatus_Expired, nil)
			// if req.IBCInfo != nil {
			// 	k.SendOracleResponse(ctx, req.IBCInfo.SourcePort, req.IBCInfo.SourceChannel, res)
			// }
		}
		// Update report info for requested validators.
		k.UpdateReportInfos(ctx, currentReqID)
		// Set last expired request ID to be this current request.
		k.SetRequestLastExpired(ctx, int64(currentReqID))
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
