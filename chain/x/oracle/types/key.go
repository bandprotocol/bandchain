package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "oracle"
	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	// GlobalStoreKeyPrefix is a prefix for global primitive state variable
	GlobalStoreKeyPrefix = []byte{0x00}

	// RequestBeginStoreKey TODO
	RequestBeginStoreKey = append(GlobalStoreKeyPrefix, []byte("RequestBeginStoreKey")...)

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

	// RawRequestStoreKeyPrefix is a prefix for storing raw data request detail.
	RawRequestStoreKeyPrefix = []byte{0x02}

	// RawDataReportStoreKeyPrefix is a prefix for report store
	RawDataReportStoreKeyPrefix = []byte{0x03}

	// DataSourceStoreKeyPrefix is a prefix for data source store.
	DataSourceStoreKeyPrefix = []byte{0x04}

	// OracleScriptStoreKeyPrefix is a prefix for oracle script store.
	OracleScriptStoreKeyPrefix = []byte{0x05}

	// ReporterStoreKeyPrefix is a prefix for reporter store.
	ReporterStoreKeyPrefix = []byte{0x06}
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

// RawRequestStoreKey is a function to generate key for each raw data request in store
func RawRequestStoreKey(requestID RequestID, externalID ExternalID) []byte {
	buf := append(RawRequestStoreKeyPrefix, int64ToBytes(int64(requestID))...)
	buf = append(buf, int64ToBytes(int64(externalID))...)
	return buf
}

// RawDataReportStoreKey is a function to generate key for each raw data report in store.
func RawDataReportStoreKey(requestID RequestID, validatorAddress sdk.ValAddress) []byte {
	buf := append(RawDataReportStoreKeyPrefix, int64ToBytes(int64(requestID))...)
	buf = append(buf, validatorAddress.Bytes()...)
	return buf
}

// RawDataReportStoreKeyUnique is a function to generate key for each raw data report in store.
func RawDataReportStoreKeyUnique(requestID RequestID, externalID ExternalID, validatorAddress sdk.ValAddress) []byte {
	buf := append(RawDataReportStoreKeyPrefix, int64ToBytes(int64(requestID))...)
	buf = append(buf, int64ToBytes(int64(externalID))...)
	buf = append(buf, validatorAddress.Bytes()...)
	return buf
}

// DataSourceStoreKey is a function to generate key for each data source in store.
func DataSourceStoreKey(dataSourceID DataSourceID) []byte {
	return append(DataSourceStoreKeyPrefix, int64ToBytes(int64(dataSourceID))...)
}

// OracleScriptStoreKey is a function to generate key for each oracle script in store.
func OracleScriptStoreKey(oracleScriptID OracleScriptID) []byte {
	return append(OracleScriptStoreKeyPrefix, int64ToBytes(int64(oracleScriptID))...)
}

// ReporterStoreKey is a function to generate key for each validator-reporter pair in store.
func ReporterStoreKey(validatorAddress sdk.ValAddress, reporterAddress sdk.AccAddress) []byte {
	buff := append(ReporterStoreKeyPrefix, []byte(validatorAddress)...)
	buff = append(buff, []byte(reporterAddress)...)
	return buff
}

// GetIteratorPrefix is a function to get specific prefix
func GetIteratorPrefix(prefix []byte, requestID RequestID) []byte {
	return append(prefix, int64ToBytes(int64(requestID))...)
}

// GetValidatorAddressAndExternalID is a function to get validator address and external id from raw data report key.
func GetValidatorAddressAndExternalID(
	key []byte, requestID RequestID,
) (sdk.ValAddress, ExternalID) {
	prefixLength := len(RawDataReportStoreKeyPrefix)
	externalIDBytes := key[prefixLength+8 : prefixLength+16]
	externalID := ExternalID(binary.BigEndian.Uint64(externalIDBytes))
	return key[prefixLength+16:], externalID
}
