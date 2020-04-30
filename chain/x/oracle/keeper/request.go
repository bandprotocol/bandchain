package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

// HasRequest checks if the request of this ID exists in the storage.
func (k Keeper) HasRequest(ctx sdk.Context, id types.RID) bool {
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
	if err := k.EnsureLength(ctx, types.KeyMaxCalldataSize, len(req.Calldata)); err != nil {
		return 0, err
	}
	id := k.GetNextRequestID(ctx)
	k.SetRequest(ctx, id, req)
	return id, nil
}

func (k Keeper) resolveRequest(
	ctx sdk.Context, reqID types.RequestID, resolveStatus types.ResolveStatus, result []byte,
) types.OracleResponsePacketData {
	request := k.MustGetRequest(ctx, reqID)
	reqPacketData := types.NewOracleRequestPacketData(
		request.ClientID, request.OracleScriptID,
		hex.EncodeToString(request.Calldata), request.SufficientValidatorCount,
		int64(len(request.RequestedValidators)),
	)
	resPacketData := types.NewOracleResponsePacketData(
		request.ClientID,
		reqID,
		int64(k.GetReportCount(ctx, reqID)),
		request.RequestTime, ctx.BlockTime().Unix(),
		types.Success, hex.EncodeToString(result),
	)

	if resolveStatus != types.Success {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRequestExecute,
				sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", reqID)),
				sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", resolveStatus)),
			))
		return resPacketData
	}

	resultHash, err := k.AddResult(ctx, reqID, reqPacketData, resPacketData)
	if err != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRequestExecute,
				sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", reqID)),
				sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", types.Failure)),
			))
		return resPacketData
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRequestExecute,
			sdk.NewAttribute(types.AttributeKeyClientID, reqPacketData.ClientID),
			sdk.NewAttribute(types.AttributeKeyOracleScriptID, fmt.Sprintf("%d", reqPacketData.OracleScriptID)),
			sdk.NewAttribute(types.AttributeKeyCalldata, reqPacketData.Calldata),
			sdk.NewAttribute(types.AttributeKeyAskCount, fmt.Sprintf("%d", reqPacketData.AskCount)),
			sdk.NewAttribute(types.AttributeKeyMinCount, fmt.Sprintf("%d", reqPacketData.MinCount)),
			sdk.NewAttribute(types.AttributeKeyRequestID, fmt.Sprintf("%d", resPacketData.RequestID)),
			sdk.NewAttribute(types.AttributeKeyResolveStatus, fmt.Sprintf("%d", types.Success)),
			sdk.NewAttribute(types.AttributeKeyAnsCount, fmt.Sprintf("%d", resPacketData.AnsCount)),
			sdk.NewAttribute(types.AttributeKeyRequestTime, fmt.Sprintf("%d", request.RequestTime)),
			sdk.NewAttribute(types.AttributeKeyResolveTime, fmt.Sprintf("%d", resPacketData.ResolveTime)),
			sdk.NewAttribute(types.AttributeKeyResult, resPacketData.Result),
			sdk.NewAttribute(types.AttributeKeyResultHash, hex.EncodeToString(resultHash)),
		))
	return resPacketData
}

// ProcessOracleResponse takes a
func (k Keeper) ProcessOracleResponse(
	ctx sdk.Context, reqID types.RequestID, resolveStatus types.ResolveStatus, result []byte,
) {
	resPacketData := k.resolveRequest(ctx, reqID, resolveStatus, result)
	if resolveStatus == types.Expired {
		return
	}

	request := k.MustGetRequest(ctx, reqID)

	if request.RequestIBC == nil {
		return
	}

	sourceChannelEnd, found := k.ChannelKeeper.GetChannel(ctx, request.RequestIBC.SourcePort, request.RequestIBC.SourceChannel)
	if !found {
		fmt.Println("SOURCE NOT FOUND", request.RequestIBC.SourcePort, request.RequestIBC.SourceChannel)
		return
	}

	destinationPort := sourceChannelEnd.Counterparty.PortID
	destinationChannel := sourceChannelEnd.Counterparty.ChannelID

	sequence, found := k.ChannelKeeper.GetNextSequenceSend(ctx, request.RequestIBC.SourcePort, request.RequestIBC.SourceChannel)
	if !found {
		fmt.Println("SEQUENCE NOT FOUND", request.RequestIBC.SourcePort, request.RequestIBC.SourceChannel)
		return
	}

	channelCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(destinationPort, destinationChannel))
	if !ok {
		fmt.Println("GET CAPABILITY ERROR", request.RequestIBC.SourcePort, request.RequestIBC.SourceChannel)
		return
	}
	err := k.ChannelKeeper.SendPacket(ctx, channelCap, channel.NewPacket(resPacketData.GetBytes(),
		sequence, request.RequestIBC.SourcePort, request.RequestIBC.SourceChannel, destinationPort, destinationChannel,
		1000000000, 1000000000, // Arbitrarily height and timestamp timeout for now
	))

	if err != nil {
		fmt.Println("SEND PACKET ERROR", err)
		return
	}
}

// ProcessExpiredRequests removes all expired data requests from the store, and
// sends oracle response packets for the ones that have never been resolved.
func (k Keeper) ProcessExpiredRequests(ctx sdk.Context) {
	currentReqID := k.GetRequestBeginID(ctx)
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
		if k.GetReportCount(ctx, currentReqID) < request.SufficientValidatorCount {
			k.ProcessOracleResponse(ctx, currentReqID, types.Expired, nil)
		}
		// We are done with this request. Remove it and its dependencies from the store.
		k.DeleteRequest(ctx, currentReqID)
		k.DeleteRawRequests(ctx, currentReqID)
		k.DeleteReports(ctx, currentReqID)
	}
	// Lastly, we update RequestBeginID to reflect the most up-to-date ID for open requests.
	k.SetRequestBeginID(ctx, currentReqID)
}

// AddPendingRequest adds the request to the pending list. DO NOT add same request more than once.
func (k Keeper) AddPendingRequest(ctx sdk.Context, requestID types.RequestID) {
	pendingList := k.GetPendingResolveList(ctx)
	pendingList = append(pendingList, requestID)
	k.SetPendingResolveList(ctx, pendingList)
}

// SetPendingResolveList saves the list of pending request that will be resolved at end block.
func (k Keeper) SetPendingResolveList(ctx sdk.Context, reqIDs []types.RequestID) {
	bz := k.cdc.MustMarshalBinaryBare(reqIDs)
	if bz == nil {
		bz = []byte{}
	}
	ctx.KVStore(k.storeKey).Set(types.PendingResolveListStoreKey, bz)
}

// GetPendingResolveList returns the list of pending requests to be executed during EndBlock.
func (k Keeper) GetPendingResolveList(ctx sdk.Context) (reqIDs []types.RequestID) {
	bz := ctx.KVStore(k.storeKey).Get(types.PendingResolveListStoreKey)
	if len(bz) == 0 { // Return an empty list if the key does not exist in the store.
		return []types.RequestID{}
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &reqIDs)
	return reqIDs
}
