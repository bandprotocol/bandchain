package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module.
	ModuleName = "oracle"
	// StoreKey to be used when creating the KVStore.
	StoreKey = ModuleName
	// Default PortID that oracle module binds to.
	PortID = ModuleName
)

var (
	// GlobalStoreKeyPrefix is the prefix for global primitive state variables.
	GlobalStoreKeyPrefix = []byte{0x00}
	// RollingSeedStoreKey is the key that keeps the seed based on the first 8-bit of the most recent 32 block hashes.
	RollingSeedStoreKey = append(GlobalStoreKeyPrefix, []byte("RollingSeed")...)
	// RequestCountStoreKey is the key that keeps the total request count.
	RequestCountStoreKey = append(GlobalStoreKeyPrefix, []byte("RequestCount")...)
	// RequestLastExpiredStoreKey is the key that keeps the ID of the last expired request, or 0 if none.
	RequestLastExpiredStoreKey = append(GlobalStoreKeyPrefix, []byte("RequestLastExpired")...)
	// PendingResolveListStoreKey is the key that keeps the list of pending-resolve requests.
	PendingResolveListStoreKey = append(GlobalStoreKeyPrefix, []byte("PendingList")...)
	// DataSourceCountStoreKey is the key that keeps the total data source count.
	DataSourceCountStoreKey = append(GlobalStoreKeyPrefix, []byte("DataSourceCount")...)
	// OracleScriptCountStoreKey is the key that keeps the total oracle sciprt count.
	OracleScriptCountStoreKey = append(GlobalStoreKeyPrefix, []byte("OracleScriptCount")...)

	// RequestStoreKeyPrefix is the prefix for request store.
	RequestStoreKeyPrefix = []byte{0x01}
	// ReportStoreKeyPrefix is the prefix for report store.
	ReportStoreKeyPrefix = []byte{0x02}
	// DataSourceStoreKeyPrefix is the prefix for data source store.
	DataSourceStoreKeyPrefix = []byte{0x03}
	// OracleScriptStoreKeyPrefix is the prefix for oracle script store.
	OracleScriptStoreKeyPrefix = []byte{0x04}
	// ReporterStoreKeyPrefix is the prefix for reporter store.
	ReporterStoreKeyPrefix = []byte{0x05}
	// ValidatorStatusKeyPrefix is the prefix for validator status store.
	ValidatorStatusKeyPrefix = []byte{0x06}
	// ResultStoreKeyPrefix is the prefix for request result store.
	ResultStoreKeyPrefix = []byte{0xff}
)

// RequestStoreKey returns the key to retrieve a specfic request from the store.
func RequestStoreKey(requestID RequestID) []byte {
	return append(RequestStoreKeyPrefix, sdk.Uint64ToBigEndian(uint64(requestID))...)
}

// ReportStoreKey returns the key to retrieve all data reports for a request.
func ReportStoreKey(requestID RequestID) []byte {
	return append(ReportStoreKeyPrefix, sdk.Uint64ToBigEndian(uint64(requestID))...)
}

// ReportStoreKeyPerValidator returns the key to retrieve the data report from a validator to a request.
func ReportStoreKeyPerValidator(reqID RequestID, val sdk.ValAddress) []byte {
	buf := append(ReportStoreKeyPrefix, sdk.Uint64ToBigEndian(uint64(reqID))...)
	buf = append(buf, val.Bytes()...)
	return buf
}

// DataSourceStoreKey returns the key to retrieve a specific data source from the store.
func DataSourceStoreKey(dataSourceID DataSourceID) []byte {
	return append(DataSourceStoreKeyPrefix, sdk.Uint64ToBigEndian(uint64(dataSourceID))...)
}

// DataSourceStoreKey returns the key to retrieve a specific oracle script from the store.
func OracleScriptStoreKey(oracleScriptID OracleScriptID) []byte {
	return append(OracleScriptStoreKeyPrefix, sdk.Uint64ToBigEndian(uint64(oracleScriptID))...)
}

// ReporterStoreKey returns the key to check whether an address is a reporter of a validator.
func ReporterStoreKey(validatorAddress sdk.ValAddress, reporterAddress sdk.AccAddress) []byte {
	buf := append(ReporterStoreKeyPrefix, []byte(validatorAddress)...)
	buf = append(buf, []byte(reporterAddress)...)
	return buf
}

// ValidatorStatusStoreKey returns the key to a validator's status.
func ValidatorStatusStoreKey(v sdk.ValAddress) []byte {
	return append(ValidatorStatusKeyPrefix, v.Bytes()...)
}

// ResultStoreKey returns the key to a request result in the store.
func ResultStoreKey(requestID RequestID) []byte {
	return append(ResultStoreKeyPrefix, sdk.Uint64ToBigEndian(uint64(requestID))...)
}

// ValidatorReporterPrefixKey returns the key to a validator's reporters.
func ValidatorReporterPrefixKey(val sdk.ValAddress) []byte {
	return append(ReporterStoreKeyPrefix, val.Bytes()...)
}
