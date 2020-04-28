package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"

	"github.com/bandprotocol/band-consumer/x/consuming/types"
)

type Keeper struct {
	storeKey        sdk.StoreKey
	cdc             codec.Marshaler
	ChannelKeeper   transfertypes.ChannelKeeper
	ScopedIBCKeeper capability.ScopedKeeper
}

// NewKeeper creates a new band consumer Keeper instance.
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey, channelKeeper transfertypes.ChannelKeeper, scopedIBCKeeper capability.ScopedKeeper) Keeper {
	return Keeper{
		storeKey:        key,
		cdc:             cdc,
		ChannelKeeper:   channelKeeper,
		ScopedIBCKeeper: scopedIBCKeeper,
	}
}

func (k Keeper) SetResult(ctx sdk.Context, requestID oracle.RequestID, result []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ResultStoreKey(requestID), result)
}

func (k Keeper) GetResult(ctx sdk.Context, requestID oracle.RequestID) ([]byte, error) {
	if !k.HasResult(ctx, requestID) {
		return nil, sdkerrors.Wrapf(types.ErrItemNotFound,
			"GetResult: Result for request ID %d is not available.", requestID,
		)
	}
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.ResultStoreKey(requestID)), nil
}

func (k Keeper) HasResult(ctx sdk.Context, requestID oracle.RequestID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ResultStoreKey(requestID))
}
