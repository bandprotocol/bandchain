package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	CoinKeeper    bank.Keeper
	StakingKeeper staking.Keeper
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, coinKeeper bank.Keeper, stakingKeeper staking.Keeper) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		CoinKeeper:    coinKeeper,
		StakingKeeper: stakingKeeper,
	}
}

// GetRequestCount returns current count of requests
func (k Keeper) GetRequestCount(ctx sdk.Context) uint64 {
	var requestNumber uint64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RequestsCountStoreKey)
	if bz == nil {
		return 0
	}
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &requestNumber)
	if err != nil {
		panic(err)
	}
	return requestNumber
}

// GetNextRequestID returns and increments the current requests.
// If the global request count is not set, it initializes it with value 0.
func (k Keeper) GetNextRequestID(ctx sdk.Context) uint64 {
	requestNumber := k.GetRequestCount(ctx)
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(requestNumber + 1)
	store.Set(types.RequestsCountStoreKey, bz)
	return requestNumber + 1
}
