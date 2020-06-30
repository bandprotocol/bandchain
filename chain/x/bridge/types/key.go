package types

const (
	// ModuleName is the name of the module.
	ModuleName = "bridge"
	// StoreKey to be used when creating the KVStore.
	StoreKey = ModuleName
)

var (
	// GlobalStoreKeyPrefix is the prefix for global primitive state variables.
	GlobalStoreKeyPrefix = []byte{0x00}
)
