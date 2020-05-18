package owasm

// ExecEnv encapsulates the operations that an Owasm script can interact with BandChain.
type ExecEnv interface {
	// GetAskCount returns the number of validators asked to work for this oracle query.
	GetAskCount() int64
	// GetMinCount returns the minimum number of validators to move to aggregation phase.
	GetMinCount() int64
	// GetAnsCount returns the number of validators that submit answers. Zero during preparation.
	GetAnsCount() int64
	// GetPrepareBlockTime returns the time of *preparation* phase was run.
	GetPrepareBlockTime() int64
	// GetAggregateBlockTime returns the time of *aggregation* phase. Zero during preparation.
	GetAggregateBlockTime() int64
	// GetValidatorAddress returns the 20-byte validator address at the specified index.
	GetValidatorAddress(validatorIndex int64) ([]byte, error)
	// GetMaximumResultSize returns the maxixmum size of aggregation result.
	GetMaximumResultSize() int64
	// GetMaximumCalldataOfDataSourceSize returns the maximum size of data source's calldata.
	GetMaximumCalldataOfDataSourceSize() int64

	// RequestExternalData performs a request to the specified data source
	// with and assigns the request with the external data ID. The function must
	// only be called during the *preparation* phase of an oracle script.
	RequestExternalData(
		dataSourceID int64,
		externalID int64,
		calldata []byte,
	) error

	// GetExternalData reads from the execution environment state for a raw
	// data report for the specified external data ID from the specified validator.
	// The function must only be called during the *aggregation* phase.
	GetExternalData(
		externalID int64,
		validatorIndex int64,
	) ([]byte, uint32, error)
}
