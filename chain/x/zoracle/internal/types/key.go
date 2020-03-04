package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

	// PendingResolveListStoreKey is a key that help getting pending request
	PendingResolveListStoreKey = append(GlobalStoreKeyPrefix, []byte("PendingList")...)

	// DataSourceCountStoreKey is a key that keeps the current data source count state variable.
	DataSourceCountStoreKey = append(GlobalStoreKeyPrefix, []byte("DataSourceCount")...)

	// OracleScriptCountStoreKey is a key that keeps the current oracle script count state variable.
	OracleScriptCountStoreKey = append(GlobalStoreKeyPrefix, []byte("OracleScriptCount")...)

	// ========================================================================

	// RequestStoreKeyPrefix is a prefix for request store
	RequestStoreKeyPrefix = []byte{0x01}

	// ResultStoreKeyPrefix is a prefix for storing result
	ResultStoreKeyPrefix = []byte{0xff}

	// RawDataRequestStoreKeyPrefix is a prefix for storing raw data request detail.
	RawDataRequestStoreKeyPrefix = []byte{0x02}

	// RawDataReportStoreKeyPrefix is a prefix for report store
	RawDataReportStoreKeyPrefix = []byte{0x03}

	// DataSourceStoreKeyPrefix is a prefix for data source store.
	DataSourceStoreKeyPrefix = []byte{0x04}

	// OracleScriptStoreKeyPrefix is a prefix for oracle script store.
	OracleScriptStoreKeyPrefix = []byte{0x05}
)

// RequestStoreKey is a function to generate key for each request in store
func RequestStoreKey(requestID RequestID) []byte {
	return append(RequestStoreKeyPrefix, int64ToBytes(int64(requestID))...)
}

// ResultStoreKey is a function to generate key for each result in store
func ResultStoreKey(requestID RequestID, oracleScriptID OracleScriptID, calldata []byte) []byte {
	buf := append(ResultStoreKeyPrefix, int64ToBytes(int64(requestID))...)
	buf = append(buf, int64ToBytes(int64(oracleScriptID))...)
	buf = append(buf, calldata...)
	return buf
}

// RawDataRequestStoreKey is a function to generate key for each raw data request in store
func RawDataRequestStoreKey(requestID RequestID, externalID int64) []byte {
	buf := append(RawDataRequestStoreKeyPrefix, int64ToBytes(int64(requestID))...)
	buf = append(buf, int64ToBytes(externalID)...)
	return buf
}

// RawDataReportStoreKey is a function to generate key for each raw data report in store.
func RawDataReportStoreKey(requestID RequestID, externalID int64, validatorAddress sdk.ValAddress) []byte {
	buf := append(RawDataReportStoreKeyPrefix, int64ToBytes(int64(requestID))...)
	buf = append(buf, int64ToBytes(externalID)...)
	buf = append(buf, validatorAddress.Bytes()...)
	return buf
}

// DataSourceStoreKey is a function to generate key for each data source in store.
func DataSourceStoreKey(dataSourceID int64) []byte {
	return append(DataSourceStoreKeyPrefix, int64ToBytes(int64(dataSourceID))...)
}

// OracleScriptStoreKey is a function to generate key for each oracle script in store.
func OracleScriptStoreKey(oracleScriptID OracleScriptID) []byte {
	return append(OracleScriptStoreKeyPrefix, int64ToBytes(int64(oracleScriptID))...)
}

// GetIteratorPrefix is a function to get specific prefix
func GetIteratorPrefix(prefix []byte, requestID RequestID) []byte {
	return append(prefix, int64ToBytes(int64(requestID))...)
}

// GetExternalIDFromRawDataRequestKey is a function to get external id from raw data request key.
func GetExternalIDFromRawDataRequestKey(key []byte) int64 {
	prefixLength := len(RawDataRequestStoreKeyPrefix)
	externalIDBytes := key[prefixLength+8 : prefixLength+16]
	return int64(binary.BigEndian.Uint64(externalIDBytes))
}

// GetValidatorAddressAndExternalID is a function to get validator address and external id from raw data report key.
func GetValidatorAddressAndExternalID(
	key []byte, requestID RequestID,
) (sdk.ValAddress, int64) {
	prefixLength := len(RawDataReportStoreKeyPrefix)
	externalIDBytes := key[prefixLength+8 : prefixLength+16]
	externalID := int64(binary.BigEndian.Uint64(externalIDBytes))
	return key[prefixLength+16:], externalID
}
