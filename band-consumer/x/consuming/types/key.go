package types

import (
	"encoding/binary"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

const (
	// ModuleName defines the IBC transfer name
	ModuleName = "consuming"

	// Version defines the current version the IBC tranfer
	// module supports
	Version = "ics20-1"

	// Default PortID that transfer module binds to
	PortID = "consuming"

	// StoreKey is the store key string for IBC transfer
	StoreKey = ModuleName

	// RouterKey is the message route for IBC transfer
	RouterKey = ModuleName

	// Key to store portID in our store
	PortKey = "portID"

	// QuerierRoute is the querier route for IBC transfer
	QuerierRoute = ModuleName
)

var (
	// ResultStoreKeyPrefix is a prefix for storing result
	ResultStoreKeyPrefix = []byte{0xff}
)

// ResultStoreKey is a function to generate key for each result in store
func ResultStoreKey(requestID oracle.RequestID) []byte {
	return append(ResultStoreKeyPrefix, int64ToBytes(int64(requestID))...)
}

func int64ToBytes(num int64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(num))
	return result
}
