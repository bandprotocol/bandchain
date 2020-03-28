package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRequest saves the given data request to the store without performing any validation.
func (k Keeper) SetRequest(ctx sdk.Context, id types.RequestID, request types.Request) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RequestStoreKey(id), k.cdc.MustMarshalBinaryBare(request))
}

// GetRequest returns the entire Request metadata struct from the store.
func (k Keeper) GetRequest(ctx sdk.Context, id types.RequestID) (types.Request, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckRequestExists(ctx, id) {
		return types.Request{}, types.ErrItemNotFound("GetRequest: Unknown request ID %d.", id)
	}

	bz := store.Get(types.RequestStoreKey(id))
	var request types.Request
	k.cdc.MustUnmarshalBinaryBare(bz, &request)
	return request, nil
}

// AddRequest attempts to create and save a new request. Returns error some conditions failed.
func (k Keeper) AddRequest(
	ctx sdk.Context, oracleScriptID types.OracleScriptID, calldata []byte,
	requestedValidatorCount, sufficientValidatorCount, expiration int64, executeGas uint64,
) (types.RequestID, sdk.Error) {
	if !k.CheckOracleScriptExists(ctx, oracleScriptID) {
		return 0, types.ErrItemNotFound(
			"AddRequest: Unknown oracle script ID %d.",
			oracleScriptID,
		)
	}

	if int64(len(calldata)) > k.MaxCalldataSize(ctx) {
		return 0, types.ErrBadDataValue(
			"AddRequest: Calldata size (%d) exceeds the maximum size (%d).",
			len(calldata),
			int(k.MaxCalldataSize(ctx)),
		)
	}

	validatorsByPower := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	if int64(len(validatorsByPower)) < requestedValidatorCount {
		return 0, types.ErrBadDataValue(
			"AddRequest: Requested validator count (%d) exceeds the number of validators (%d).",
			requestedValidatorCount,
			len(validatorsByPower),
		)
	}

	validators := make([]sdk.ValAddress, requestedValidatorCount)
	for i := int64(0); i < requestedValidatorCount; i++ {
		validators[i] = validatorsByPower[i].GetOperator()
	}

	ctx.GasMeter().ConsumeGas(executeGas, "ExecuteGas")
	if executeGas > k.EndBlockExecuteGasLimit(ctx) {
		return 0, types.ErrBadDataValue(
			"AddRequest: Execute gas (%d) exceeds the maximum limit (%d).",
			executeGas,
			k.EndBlockExecuteGasLimit(ctx),
		)
	}

	requestID := k.GetNextRequestID(ctx)
	k.SetRequest(ctx, requestID, types.NewRequest(
		oracleScriptID,
		calldata,
		validators,
		sufficientValidatorCount,
		ctx.BlockHeight(),
		ctx.BlockTime().Unix(),
		ctx.BlockHeight()+expiration,
		executeGas,
	))

	return requestID, nil
}

// ValidateDataSourceCount returns whether the number of raw data requests exceeds the maximum
// allowed value, as specified by `MaxDataSourceCountPerRequest` parameter.
func (k Keeper) ValidateDataSourceCount(ctx sdk.Context, id types.RequestID) sdk.Error {
	dataSourceCount := k.GetRawDataRequestCount(ctx, id)
	if dataSourceCount > k.MaxDataSourceCountPerRequest(ctx) {
		return types.ErrBadDataValue(
			"ValidateDataSourceCount: Data source count (%d) exceeds the limit (%d).",
			dataSourceCount,
			k.MaxDataSourceCountPerRequest(ctx),
		)
	}
	return nil
}

// PayDataSourceFees sends fees from the sender to the owner of the requested data source.
func (k Keeper) PayDataSourceFees(
	ctx sdk.Context, id types.RequestID, sender sdk.AccAddress,
) sdk.Error {
	rawDataRequests := k.GetRawDataRequests(ctx, id)
	for _, rawDataRequest := range rawDataRequests {
		dataSource, err := k.GetDataSource(ctx, rawDataRequest.DataSourceID)
		if err != nil {
			return err
		}
		if dataSource.Owner.Equals(sender) {
			continue
		}
		if dataSource.Fee.IsZero() {
			continue
		}
		err = k.CoinKeeper.SendCoins(ctx, sender, dataSource.Owner, dataSource.Fee)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetResolve updates the resolve status of the given request as specified by id.
func (k Keeper) SetResolve(
	ctx sdk.Context, id types.RequestID, resolveStatus types.ResolveStatus,
) sdk.Error {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		return err
	}

	request.ResolveStatus = resolveStatus
	k.SetRequest(ctx, id, request)
	return nil
}

// CheckRequestExists checks if the request at this id is present in the store or not.
func (k Keeper) CheckRequestExists(ctx sdk.Context, id types.RequestID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.RequestStoreKey(id))
}

// ShouldBecomePendingResolve checks and returns whether the given request should be moved to the
// pending resolve list, which will be resolved during the EndBlock call. The move will happen
// exactly once will the request receives sufficient raw reports from the validators.
func (k Keeper) ShouldBecomePendingResolve(ctx sdk.Context, id types.RequestID) bool {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		return false
	}
	return int64(len(request.ReceivedValidators)) == request.SufficientValidatorCount
}

// AddPendingRequest appends the given request to the pending list. Returns error if the request
// already exists in the list.
func (k Keeper) AddPendingRequest(ctx sdk.Context, requestID types.RequestID) sdk.Error {
	pendingList := k.GetPendingResolveList(ctx)
	for _, entry := range pendingList {
		if requestID == entry {
			return types.ErrItemDuplication(
				"AddPendingRequest: Request ID %d already exists in the pending list",
				requestID,
			)
		}
	}
	pendingList = append(pendingList, requestID)
	k.SetPendingResolveList(ctx, pendingList)
	return nil
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
