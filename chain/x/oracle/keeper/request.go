package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

// ProcessOracleResponse takes a
func (k Keeper) ProcessOracleResponse(
	ctx sdk.Context, reqID types.RequestID, resolveStatus types.ResolveStatus, result []byte,
) {
	request := k.MustGetRequest(ctx, reqID)

	// TODO: Send IBC packets + save data to result tree
	reqPacketData := types.OracleRequestPacketData{}
	resPacketData := types.OracleResponsePacketData{}

	_ = request
	_ = reqPacketData
	_ = resPacketData

	// SOME OLD CODE FOR YOU!
	// 	event, packet := handleResolveRequest(ctx, keeper, requestID)
	// 	// TODO: Refactor this packet code
	// 	request, err := keeper.GetRequest(ctx, requestID)
	// 	events = append(events, event)
	// 	sourceChannelEnd, found := keeper.ChannelKeeper.GetChannel(ctx, request.SourcePort, request.SourceChannel)
	// 	if !found {
	// 		fmt.Println("SOURCE NOT FOUND", request.SourcePort, request.SourceChannel)
	// 		continue
	// 	}

	// 	destinationPort := sourceChannelEnd.Counterparty.PortID
	// 	destinationChannel := sourceChannelEnd.Counterparty.ChannelID

	// 	// get the next sequence
	// 	sequence, found := keeper.ChannelKeeper.GetNextSequenceSend(ctx, request.SourcePort, request.SourceChannel)
	// 	if !found {
	// 		fmt.Println("SEQUENCE NOT FOUND", request.SourcePort, request.SourceChannel)
	// 		continue
	// 	}

	// 	err = keeper.ChannelKeeper.SendPacket(ctx, channel.NewPacket(packet.GetBytes(),
	// 		sequence, request.SourcePort, request.SourceChannel, destinationPort, destinationChannel,
	// 		1000000000, // Arbitrarily high timeout for now
	// 	))

	// 	if err != nil {
	// 		fmt.Println("SEND PACKET ERROR", err)
	// 	}
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

// ValidateDataSourceCount returns whether the number of raw data requests exceeds the maximum
// allowed value, as specified by `MaxDataSourceCountPerRequest` parameter.
func (k Keeper) ValidateDataSourceCount(ctx sdk.Context, id types.RequestID) error {
	dataSourceCount := k.GetRawRequestCount(ctx, id)
	if uint64(dataSourceCount) > k.GetParam(ctx, types.KeyMaxDataSourceCountPerRequest) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"ValidateDataSourceCount: Data source count (%d) exceeds the limit (%d).",
			dataSourceCount, k.GetParam(ctx, types.KeyMaxDataSourceCountPerRequest),
		)
	}
	return nil
}

// PayDataSourceFees sends fees from the sender to the owner of the requested data source.
func (k Keeper) PayDataSourceFees(ctx sdk.Context, id types.RID, sender sdk.AccAddress) error {
	rawDataRequests := k.GetRawRequests(ctx, id)
	for _, rawDataRequest := range rawDataRequests {
		dataSource := k.MustGetDataSource(ctx, rawDataRequest.DataSourceID)
		if dataSource.Owner.Equals(sender) {
			continue
		}
		if dataSource.Fee.IsZero() {
			continue
		}
		err := k.CoinKeeper.SendCoins(ctx, sender, dataSource.Owner, dataSource.Fee)
		if err != nil {
			return err
		}
	}
	return nil
}

// ShouldBecomePendingResolve checks and returns whether the given request should be moved to the
// pending resolve list, which will be resolved during the EndBlock call. The move will happen
// exactly once will the request receives sufficient raw reports from the validators.
func (k Keeper) ShouldBecomePendingResolve(ctx sdk.Context, id types.RequestID) bool {
	request := k.MustGetRequest(ctx, id)
	return k.GetReportCount(ctx, id) == request.SufficientValidatorCount
}

// AddPendingRequest adds the request to the pending list. DO NOT add same request more than once.
func (k Keeper) AddPendingRequest(ctx sdk.Context, requestID types.RequestID) {
	pendingList := k.GetPendingResolveList(ctx)
	pendingList = append(pendingList, requestID)
	k.SetPendingResolveList(ctx, pendingList)
}

// SetPendingResolveList saves the list of pending request that will be resolved at end block.
func (k Keeper) SetPendingResolveList(ctx sdk.Context, reqIDs []types.RequestID) {
	store := ctx.KVStore(k.storeKey)
	encoded := k.cdc.MustMarshalBinaryBare(reqIDs)
	if encoded == nil {
		encoded = []byte{}
	}
	store.Set(types.PendingResolveListStoreKey, encoded)
}

// GetPendingResolveList returns the list of pending requests to be executed during EndBlock.
func (k Keeper) GetPendingResolveList(ctx sdk.Context) []types.RequestID {
	store := ctx.KVStore(k.storeKey)
	reqIDsBytes := store.Get(types.PendingResolveListStoreKey)
	if len(reqIDsBytes) == 0 {
		// Return an empty list if the key does not exist in the store.
		return []types.RequestID{}
	}
	var reqIDs []types.RequestID
	k.cdc.MustUnmarshalBinaryBare(reqIDsBytes, &reqIDs)
	return reqIDs
}
