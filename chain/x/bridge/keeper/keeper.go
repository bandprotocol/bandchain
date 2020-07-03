package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/merkle"
	tmmerkle "github.com/tendermint/tendermint/crypto/merkle"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"

	"github.com/bandprotocol/bandchain/chain/x/bridge/types"
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

// SetLatestRelayBlockHeight sets the latest block height that relay block to the store.
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

// SetAppHash sets the app hash to the given height to the store
func (k Keeper) SetAppHash(ctx sdk.Context, height int64, appHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AppHashStoreKey(height), appHash)
}

// GetAppHash returns the app hash of the given height
func (k Keeper) GetAppHash(ctx sdk.Context, height int64) []byte {
	return ctx.KVStore(k.storeKey).Get(types.AppHashStoreKey(height))
}

// SetLatestResponse sets the latest response of the given request packet to the store
func (k Keeper) SetLatestResponse(ctx sdk.Context, requestPacket otypes.OracleRequestPacketData, responsePacket otypes.OracleResponsePacketData) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastestResponseStoreKey(requestPacket), k.cdc.MustMarshalBinaryBare(responsePacket))
}

// GetLatestResponse returns the lastest response of the given request packet
func (k Keeper) GetLatestResponse(ctx sdk.Context, requestPacket otypes.OracleRequestPacketData) otypes.OracleResponsePacketData {
	var responsePacket otypes.OracleResponsePacketData
	bz := ctx.KVStore(k.storeKey).Get(types.LastestResponseStoreKey(requestPacket))
	if len(bz) == 0 {
		return otypes.OracleResponsePacketData{}
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &responsePacket)
	return responsePacket
}

// UpdateValidators sets the new validators to the store, the previous validators set will be delete.
func (k Keeper) UpdateValidators(ctx sdk.Context, validators []tmtypes.Validator) error {
	// Delete previous set of validators on the store
	k.deleteValidators(ctx)

	// Set the new set of validators
	for _, val := range validators {
		k.setValidator(ctx, val)
	}
	return nil
}

// setValidator sets the validator to the store
func (k Keeper) setValidator(ctx sdk.Context, validator tmtypes.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ValidatorStoreKey(validator), k.cdc.MustMarshalBinaryBare(validator))
}

// deleteValidators deletes all validators on the store
func (k Keeper) deleteValidators(ctx sdk.Context) {
	var keys [][]byte
	iterator := k.GetValidatorsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, iterator.Key())
	}
	for _, key := range keys {
		ctx.KVStore(k.storeKey).Delete(key)
	}
}

// GetValidatorsIterator returns the iterator for all validators on validator set
func (k Keeper) GetValidatorsIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.ValidatorStoreKeyPrefix)
}

// GetValidators returns the current validators set from the store
func (k Keeper) GetValidators(ctx sdk.Context) (validators []tmtypes.Validator) {
	iterator := k.GetValidatorsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var validator tmtypes.Validator
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &validator)
		validators = append(validators, validator)
	}
	return validators
}

// GetTotalValidatorsVotingPower returns the total voting power of all validators
func (k Keeper) GetTotalValidatorsVotingPower(ctx sdk.Context) (totalPower int64) {
	iterator := k.GetValidatorsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var validator tmtypes.Validator
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &validator)
		totalPower += validator.VotingPower
	}
	return totalPower
}

// Relay - relay the block on BandChain, set app hash to the given height if relay block is valid, or return err if relay block is invalid.
func (k Keeper) Relay(ctx sdk.Context, signedHeader tmtypes.SignedHeader) error {
	// Check relay block height, returns err if block height is older than the latest update validators set block height
	latestUpdateValidatorBlockHeight := k.GetLatestValidatorsUpdateBlockHeight(ctx)
	if signedHeader.Height < latestUpdateValidatorBlockHeight {
		return sdkerrors.Wrapf(types.ErrRelayTooOldBlockHeight, "min block height: %d", latestUpdateValidatorBlockHeight)
	}

	chainID := k.GetChainID(ctx)
	validators := k.GetValidators(ctx)
	totalVotingPower := k.GetTotalValidatorsVotingPower(ctx)
	sumVotingPower := int64(0)
	valsMap := make(map[string]tmtypes.Validator)
	for _, val := range validators {
		valsMap[val.Address.String()] = val
	}

	for idx, commitSign := range signedHeader.Commit.Signatures {
		if val, ok := valsMap[commitSign.ValidatorAddress.String()]; ok {
			msg := signedHeader.Commit.VoteSignBytes(chainID, idx)
			if val.PubKey.VerifyBytes(msg, commitSign.Signature) {
				sumVotingPower += val.VotingPower
				delete(valsMap, val.Address.String())
			}
		}
	}

	if 3*sumVotingPower > 2*totalVotingPower {
		k.SetAppHash(ctx, signedHeader.Height, signedHeader.AppHash)
		return nil
	}
	return sdkerrors.Wrapf(types.ErrRelayBlock, "sum voting power: %d, total voting power: %d", sumVotingPower, totalVotingPower)
}

// VerifyProof - verify the given proof that the request and reponse is from the BandChain
func (k Keeper) VerifyProof(ctx sdk.Context, height int64, proof tmmerkle.Proof, requestPacket otypes.OracleRequestPacketData, responsePacket otypes.OracleResponsePacketData) error {
	// Verify given proof
	prt := rootmulti.DefaultProofRuntime()
	appHash := k.GetAppHash(ctx, height)
	if appHash == nil {
		return sdkerrors.Wrapf(types.ErrAppHashNotFound, "height: %d", height)
	}

	kp := merkle.KeyPath{}
	kp = kp.AppendKey([]byte("oracle"), tmmerkle.KeyEncodingURL)
	kp = kp.AppendKey(proof.Ops[0].Key, tmmerkle.KeyEncodingURL)

	value := obi.MustEncode(requestPacket, responsePacket)
	err := prt.VerifyValue(&proof, appHash, kp.String(), value)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrVerifyProofFail, "err: %s", err.Error())
	}

	// Check the lastest response packet is newer than the previous packet
	previousPacket := k.GetLatestResponse(ctx, requestPacket)
	if responsePacket.ResolveTime <= previousPacket.ResolveTime {
		return sdkerrors.Wrapf(types.ErrResponsePacketOutdated, "lastest packet resolve time: %d", previousPacket.ResolveTime)
	}

	k.SetLatestResponse(ctx, requestPacket, responsePacket)
	return nil
}
