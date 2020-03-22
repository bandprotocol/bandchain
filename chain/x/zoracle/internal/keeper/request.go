package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SetRequest is a function to save request to the given ID.
func (k Keeper) SetRequest(ctx sdk.Context, id types.RequestID, request types.Request) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RequestStoreKey(id), k.cdc.MustMarshalBinaryBare(request))
}

// GetRequest returns the entire Request metadata struct.
func (k Keeper) GetRequest(ctx sdk.Context, id types.RequestID) (types.Request, error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckRequestExists(ctx, id) {
		return types.Request{}, sdkerrors.Wrapf(types.ErrItemNotFound, "GetRequest: Unknown request ID %d.", id)
	}

	bz := store.Get(types.RequestStoreKey(id))
	var request types.Request
	k.cdc.MustUnmarshalBinaryBare(bz, &request)
	return request, nil
}

// AddRequest attempts to create a new request. An error is returned if some conditions failed.
func (k Keeper) AddRequest(
	ctx sdk.Context, oracleScriptID types.OracleScriptID, calldata []byte,
	requestedValidatorCount, sufficientValidatorCount, expiration int64, executeGas uint64,
) (types.RequestID, error) {
	if !k.CheckOracleScriptExists(ctx, oracleScriptID) {
		return 0, sdkerrors.Wrapf(types.ErrItemNotFound,
			"AddRequest: Unknown oracle script ID %d.",
			oracleScriptID,
		)
	}

	if int64(len(calldata)) > k.MaxCalldataSize(ctx) {
		return 0, sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddRequest: Calldata size (%d) exceeds the maximum size (%d).",
			len(calldata),
			int(k.MaxCalldataSize(ctx)),
		)
	}

	validatorsByPower := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	if int64(len(validatorsByPower)) < requestedValidatorCount {
		return 0, sdkerrors.Wrapf(types.ErrBadDataValue,
			"AddRequest: Requested validator count (%d) exceeds the current number of validators (%d).",
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
		return 0, sdkerrors.Wrapf(types.ErrBadDataValue,
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

// ValidateDataSourceCount validates that the number of raw data requests is
// not greater than `MaxDataSourceCountPerRequest`
func (k Keeper) ValidateDataSourceCount(ctx sdk.Context, id types.RequestID) error {
	dataSourceCount := k.GetRawDataRequestCount(ctx, id)
	if dataSourceCount > k.MaxDataSourceCountPerRequest(ctx) {
		return sdkerrors.Wrapf(types.ErrBadDataValue,
			"ValidateDataSourceCount: Data source count (%d) exceeds the limit (%d).",
			dataSourceCount,
			k.MaxDataSourceCountPerRequest(ctx),
		)
	}

	return nil
}

// PayDataSourceFees sends fees to the owners of the requested data sources.
func (k Keeper) PayDataSourceFees(ctx sdk.Context, id types.RequestID, sender sdk.AccAddress) error {
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

func (k Keeper) SetResolve(ctx sdk.Context, id types.RequestID, resolveStatus types.ResolveStatus) error {
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
// pending resolve list, which will be resolved during the EndBlock call. The move will happen exactly when
// the request receives sufficient raw reports from the validators.
func (k Keeper) ShouldBecomePendingResolve(ctx sdk.Context, id types.RequestID) bool {
	request, err := k.GetRequest(ctx, id)
	if err != nil {
		return false
	}
	return int64(len(request.ReceivedValidators)) == request.SufficientValidatorCount
}

// AddPendingRequest checks and append new request id to list if id already existed in list, it will return error.
func (k Keeper) AddPendingRequest(ctx sdk.Context, requestID types.RequestID) error {
	pendingList := k.GetPendingResolveList(ctx)
	for _, entry := range pendingList {
		if requestID == entry {
			return sdkerrors.Wrapf(types.ErrItemDuplication,
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

// GetPendingResolveList returns the list of pending request.
func (k Keeper) GetPendingResolveList(ctx sdk.Context) []types.RequestID {
	store := ctx.KVStore(k.storeKey)
	reqIDsBytes := store.Get(types.PendingResolveListStoreKey)

	// If the state is empty
	if len(reqIDsBytes) == 0 {
		return []types.RequestID{}
	}

	var reqIDs []types.RequestID
	k.cdc.MustUnmarshalBinaryBare(reqIDsBytes, &reqIDs)

	return reqIDs
}
