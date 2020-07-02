package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is a bridge Keeper instance.
type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
}

// NewKeeper creates a new bridge Keeper instance.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: key,
	}
}

// SetChainID sets the chainID for relay and verify proof.
func (k Keeper) SetChainID(ctx sdk.Context, chainID string) {
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(chainID)
	ctx.KVStore(k.storeKey).Set(types.ChainIDStoreKey, bz)
}

// GetChainID returns the chain ID for relay and verify proof.
func (k Keeper) GetChainID(ctx sdk.Context) string {
	var chainID string
	bz := ctx.KVStore(k.storeKey).Get(types.ChainIDStoreKey)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &chainID)
	return chainID
}

// SetLatestRelayBlockHeight sets the latest block height that relay block.
func (k Keeper) SetLatestRelayBlockHeight(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LatestRelayBlockHeightStoreKey, k.cdc.MustMarshalBinaryLengthPrefixed(height))
}

// GetLatestRelayBlockHeight returns the latest block height that relay block.
func (k Keeper) GetLatestRelayBlockHeight(ctx sdk.Context) int64 {
	var height int64
	bz := ctx.KVStore(k.storeKey).Get(types.LatestRelayBlockHeightStoreKey)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &height)
	return height
}

// SetLatestValidatorsUpdateBlockHeight sets the lastest block height that validator set is updated.
func (k Keeper) SetLatestValidatorsUpdateBlockHeight(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LatestValidatorsUpdateBlockHeightStoreKey, k.cdc.MustMarshalBinaryLengthPrefixed(height))
}

// GetLatestValidatorsUpdateBlockHeight returns the latest block height that validator set is updated.
func (k Keeper) GetLatestValidatorsUpdateBlockHeight(ctx sdk.Context) int64 {
	var height int64
	bz := ctx.KVStore(k.storeKey).Get(types.LatestValidatorsUpdateBlockHeightStoreKey)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &height)
	return height
}

// SetAppHash sets the app hash to the given height
func (k Keeper) SetAppHash(ctx sdk.Context, height int64, appHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AppHashStoreKey(height), k.cdc.MustMarshalBinaryBare(appHash))
}

// GetAppHash returns the app hash of the given height
func (k Keeper) GetAppHash(ctx sdk.Context, height int64) []byte {
	var hash []byte
	bz := ctx.KVStore(k.storeKey).Get(types.AppHashStoreKey(height))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &hash)
	return hash
}
