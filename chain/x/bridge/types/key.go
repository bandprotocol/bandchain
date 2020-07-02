package types

import (
	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	// ModuleName is the name of the module.
	ModuleName = "bridge"
	// StoreKey to be used when creating the KVStore.
	StoreKey = ModuleName
)

var (
	// GlobalStoreKeyPrefix is the prefix for global primitive state variables.
	GlobalStoreKeyPrefix = []byte{0x00}
	// ChainIDStoreKey is the key that keeps the chain ID.
	ChainIDStoreKey = append(GlobalStoreKeyPrefix, []byte("ChainID")...)
	// LatestValidatorsUpdateBlockHeightStoreKey is the key that keeps the latest block height that validators are updated.
	LatestValidatorsUpdateBlockHeightStoreKey = append(GlobalStoreKeyPrefix, []byte("LatestUpdatedValidatorsBlockHeight")...)
	//LatestRelayBlockHeightStoreKey is the key that keeps the latest block height that relay.
	LatestRelayBlockHeightStoreKey = append(GlobalStoreKeyPrefix, []byte("LatestRelayBlockHeight")...)

	// ValidatorStoreKeyPrefix is the prefix for validator store.
	ValidatorStoreKeyPrefix = []byte{0x01}
	// AppHashStoreKeyPrefix is the prefix for app hash store.
	AppHashStoreKeyPrefix = []byte{0x02}
	// LatestResponseStoreKeyPrefix is the prefix for latest response store.
	LatestResponseStoreKeyPrefix = []byte{0x03}
)

// AppHashStoreKey returns the key to retrieve the app hash from block height.
func AppHashStoreKey(height int64) []byte {
	return append(AppHashStoreKeyPrefix, sdk.Uint64ToBigEndian(uint64(height))...)
}

// LastestResponseStoreKey returns the key to retrieve the latest response packet from the given request packet
func LastestResponseStoreKey(requestPacket otypes.OracleRequestPacketData) []byte {
	return append(LatestResponseStoreKeyPrefix, requestPacket.GetBytes()...)
}

// ValidatorStoreKey returns the key to retrieve the validator from the validator address
func ValidatorStoreKey(validator tmtypes.Validator) []byte {
	return append(ValidatorStoreKeyPrefix, []byte(validator.Address)...)
}
