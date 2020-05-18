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

	// The maximum consecutive misses allowance
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

// ParamSetPairs implements the params.ParamSet interface for Params.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxRawRequestCount, &p.MaxRawRequestCount, validateMaxRawRequestCount),
		paramtypes.NewParamSetPair(KeyMaxResultSize, &p.MaxResultSize, validateMaxResultSize),
		paramtypes.NewParamSetPair(KeyGasPerRawDataRequestPerValidator, &p.GasPerRawDataRequestPerValidator, validateGasPerRawDataRequestPerValidator),
		paramtypes.NewParamSetPair(KeyExpirationBlockCount, &p.ExpirationBlockCount, validateExpirationBlockCount),
		paramtypes.NewParamSetPair(KeyMaxConsecutiveMisses, &p.MaxConsecutiveMisses, validateMaxConsecutiveMisses),
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

func validateMaxRawRequestCount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("max raw request count must be positive: %d", v)
	}
	return nil
}

func validateMaxResultSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("max result size must be positive: %d", v)
	}
	return nil
}

func validateGasPerRawDataRequestPerValidator(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("gas per raw data request per validator must be positive: %d", v)
	}
	return nil
}

func validateExpirationBlockCount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("expiration block count must be positive: %d", v)
	}
	return nil
}

func validateMaxConsecutiveMisses(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("max consecutive misses must be positive: %d", v)
	}
	return nil
}
