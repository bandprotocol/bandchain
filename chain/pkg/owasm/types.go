package owasm

// ExecEnv encapsulates the operations that an Owasm script
// can call to interact with the external world. An operation can fail and
// when that occurs, the whole function call will fail.
type ExecEnv interface {
	// GetAskCount returns the number of validators that the current
	// data request specifies for the oracle query.
	GetAskCount() int64

	// GetMinCount returns the number of validators
	// that is enough to push this data request into the aggregation phase.
	GetMinCount() int64

	// GetAnsCount returns the number of validators among the
	// requested ones that replied with raw data reports. Return zero during the
	// *preparation* phase.
	GetAnsCount() int64

	// GetPrepareBlockTime returns the time at which the *preparation* phase of
	// this data request was being run.
	GetPrepareBlockTime() int64

	// GetAggregateBlockTime returns the time at which the *aggregation* phase of
	// this data request was being run. Return zero during the *preparation* phase.
	GetAggregateBlockTime() int64

	// GetValidatorPublic returns the 20-byte address of the block validator
	// at the specified index.
	GetValidatorAddress(validatorIndex int64) ([]byte, error)

	// GetMaximumResultSize returns the maxixmum size of result data that returns from
	// execute function.
	GetMaximumResultSize() int64

	// GetMaximumCalldataOfDataSourceSize returns the maximum size of call data using in
	// data source execution.
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
