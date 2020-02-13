package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRawDataRequest is a function to save raw data request detail to the given request id and external id.
func (k Keeper) SetRawDataRequest(ctx sdk.Context, requestID, externalID, scriptID int64, calldata []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.RawDataRequestStoreKey(requestID, externalID),
		k.cdc.MustMarshalBinaryBare(types.NewRawDataRequest(scriptID, calldata)),
	)
}

// GetRawDataRequest is a function to get raw data request detail by the given request id and external id.
func (k Keeper) GetRawDataRequest(
	ctx sdk.Context, requestID, externalID int64,
) (types.RawDataRequest, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.RawDataRequestStoreKey(requestID, externalID)) {
		return types.RawDataRequest{}, types.ErrRawDataRequestNotFound(types.DefaultCodespace)
	}

	bz := store.Get(types.RawDataRequestStoreKey(requestID, externalID))
	var requestDetail types.RawDataRequest
	k.cdc.MustUnmarshalBinaryBare(bz, &requestDetail)
	return requestDetail, nil
}
