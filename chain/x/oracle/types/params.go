package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace.
const (
	DefaultParamspace = ModuleName

	// The maximum number of raw requests that a request can make.
	// Default value is set to 16.
	DefaultMaxRawRequestCount = uint64(16)

	// The maximum size of result after execution.
	// Default value is set 1 kB.
	DefaultMaxResultSize = uint64(1 * 1024)

	// Gas cost per validator for each raw data request.
	DefaultGasPerRawDataRequestPerValidator = uint64(25000)

	// Expiration block count value 20
	DefaultExpirationBlockCount = uint64(20)

	// Reported window size is 100
	DefaulMaxConsecutiveMisses = uint64(10)
)

// Parameter store keys.
var (
	KeyMaxRawRequestCount               = []byte("MaxRawRequestCount")
	KeyMaxResultSize                    = []byte("MaxResultSize")
	KeyGasPerRawDataRequestPerValidator = []byte("GasPerRawDataRequestPerValidator")
	KeyExpirationBlockCount             = []byte("ExpirationBlockCount")
	KeyMaxConsecutiveMisses             = []byte("MaxConsecutiveMisses")
)

// String implements the stringer interface for Params.
func (p Params) String() string {
	return fmt.Sprintf(`oracle Params:
  MaxRawRequestCount:               %d
  MaxResultSize:                    %d
  GasPerRawDataRequestPerValidator: %d
  ExpirationBlockCount:             %d
  MaxConsecutiveMisses:             %d
`,
		p.MaxRawRequestCount,
		p.MaxResultSize,
		p.GasPerRawDataRequestPerValidator,
		p.ExpirationBlockCount,
		p.MaxConsecutiveMisses,
	)
}

func validateNoOp(_ interface{}) error { return nil }

// ParamSetPairs implements the params.ParamSet interface for Params.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	// TODO: Make validation real. Not just noop
	return params.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxRawRequestCount, &p.MaxRawRequestCount, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxResultSize, &p.MaxResultSize, validateNoOp),
		paramtypes.NewParamSetPair(KeyGasPerRawDataRequestPerValidator, &p.GasPerRawDataRequestPerValidator, validateNoOp),
		paramtypes.NewParamSetPair(KeyExpirationBlockCount, &p.ExpirationBlockCount, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxConsecutiveMisses, &p.MaxConsecutiveMisses, validateNoOp),
	}
}

// DefaultParams defines the default parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxRawRequestCount,
		DefaultMaxResultSize,
		DefaultGasPerRawDataRequestPerValidator,
		DefaultExpirationBlockCount,
		DefaulMaxConsecutiveMisses,
	)
}
