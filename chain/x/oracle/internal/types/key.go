package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "zoracle"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	// GlobalStoreKeyPrefix is a prefix for global primitive state variable
	GlobalStoreKeyPrefix = []byte{0x00}

	// RequestsCountStoreKey is a key that help getting to current requests count state variable
	RequestsCountStoreKey = append(GlobalStoreKeyPrefix, []byte("RequestsCount")...)

	// PendingListStoreKey is a key that help getting pending request
	PendingListStoreKey = append(GlobalStoreKeyPrefix, []byte("PendingList")...)

	// ========================================================================

	// DataPointStoreKeyPrefix is a prefix for datapoint store
	DataPointStoreKeyPrefix = []byte{0x01}

	// CodeHashKeyPrefix is a prefix for code store
	CodeHashKeyPrefix = []byte{0x02}

	// ReportKeyPrefix is a prefix for report store
	ReportKeyPrefix = []byte{0x03}
)

// DataPointStoreKey is a function to generate key for each datapoint in store
func DataPointStoreKey(requestID uint64) []byte {
	buf := uint64ToBytes(requestID)
	return append(DataPointStoreKeyPrefix, buf...)
}

// CodeHashStoreKey is a function to generate key for codehash to actual code in store
func CodeHashStoreKey(codeHash []byte) []byte {
	return append(CodeHashKeyPrefix, codeHash...)
}

// ReportStoreKey is a function to generate key for each report from
// validator calculate from validator address and request id
func ReportStoreKey(requestID uint64, validatorAddress sdk.ValAddress) []byte {
	buf := append(ReportKeyPrefix, uint64ToBytes(requestID)...)
	return append(buf, validatorAddress.Bytes()...)
}

// GetIteratorPrefix is a function to get specific prefix
func GetIteratorPrefix(prefix []byte, requestID uint64) []byte {
	buf := uint64ToBytes(requestID)
	return append(prefix, buf...)
}

// GetValidatorAddress is a function to get validator address from key
func GetValidatorAddress(key []byte, prefix []byte, requestID uint64) sdk.ValAddress {
	lenRequest := len(uint64ToBytes(requestID))
	return key[len(prefix)+lenRequest:]
}
