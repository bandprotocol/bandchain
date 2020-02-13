package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRawDataRequest is a function to save raw data request detail to the given request id and external id.
func (k Keeper) SetRawDataRequest(
	ctx sdk.Context, requestID, externalID, dataSourceID int64, calldata []byte,
) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	// TODO: Check calldata size

	// Check request exist
	if !k.CheckRequestExists(ctx, requestID) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}
	// Check data source exist
	if !k.CheckDataSourceExists(ctx, dataSourceID) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}
	if k.CheckRawDataRequestExists(ctx, requestID, externalID) {
		// TODO: fix error later
		return types.ErrRequestNotFound(types.DefaultCodespace)
	}
	store.Set(
		types.RawDataRequestStoreKey(requestID, externalID),
		k.cdc.MustMarshalBinaryBare(types.NewRawDataRequest(dataSourceID, calldata)),
	)
	return nil
}

// GetRawDataRequest is a function to get raw data request detail by the given request id and external id.
func (k Keeper) GetRawDataRequest(
	ctx sdk.Context, requestID, externalID int64,
) (types.RawDataRequest, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if !k.CheckRawDataRequestExists(ctx, requestID, externalID) {
		return types.RawDataRequest{}, types.ErrRawDataRequestNotFound(types.DefaultCodespace)
	}

	bz := store.Get(types.RawDataRequestStoreKey(requestID, externalID))
	var requestDetail types.RawDataRequest
	k.cdc.MustUnmarshalBinaryBare(bz, &requestDetail)
	return requestDetail, nil
}

func (k Keeper) CheckRawDataRequestExists(ctx sdk.Context, requestID, externalID int64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.RawDataRequestStoreKey(requestID, externalID))
}
