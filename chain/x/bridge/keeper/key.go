package types

const (
	// ModuleName is the name of the module.
	ModuleName = "bridge"
	// StoreKey to be used when creating the KVStore.
	StoreKey = ModuleName
	// Default PortID that oracle module binds to.
	PortID = ModuleName
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