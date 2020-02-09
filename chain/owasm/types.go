package owasm

// ExecutionEnvironment encapsulates the operations that an Owasm script
// can call to interact with the external world. An operation can fail and
// when that occurs, the whole function call will fail.
type ExecutionEnvironment interface {
	// GetCurrentRequestID returns the unique identifier that can is the reference
	// to the current data request.
	GetCurrentRequestID() (int64, error)

	// GetRequestedValidatorCount returns the number of validators that the current
	// data request specifies for the oracle query.
	GetRequestedValidatorCount() (int32, error)

	// GetReceivedValidatorCount returns the number of validators among the
	// requested ones that replied with raw data reports. Only available during
	// the *aggregation* phase.
	GetReceivedValidatorCount() (int32, error)

	// GetValidatorPublic returns the 33-byte compressed public key of the
	// block validator at the specified index.
	GetValidatorPubKey(validatorIndex int32) ([]byte, error)

	// GetPrepareBlockTime returns the time at which the *preparation* phase of
	// this data request was being run.
	GetPrepareBlockTime() (int64, error)

	// GetAggregateBlockTime returns the time at which the *aggregation* phase of
	// this data request was being run. Only available during the *aggregation* phase.
	GetAggregateBlockTime() (int64, error)

	// RequestExternalData performs a request to the specified data source
	// with and assigns the request with the external data ID. The function must
	// only be called during the *preparation* phase of an oracle script.
	RequestExternalData(
		dataSourceID int64,
		externalDataID int64,
		calldata []byte,
	) error

	// GetExternalData reads from the execution environment state for a raw
	// data report for the specified external data ID from the specified validator
	// The function must only be called during the *aggregation* phase.
	GetExternalData(
		externalDataID int64,
		validatorIndex int32,
	) ([]byte, error)
}
