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

	// The maximum size of raw data report per data source.
	// Default value is set to 1 kB.
	DefaultMaxRawDataReportSize = uint64(1 * 1024)

	// The maximum size of result after execution.
	// Default value is set 1 kB.
	DefaultMaxResultSize = uint64(1 * 1024)

	// Gas cost per validator for each raw data request.
	DefaultGasPerRawDataRequestPerValidator = uint64(25000)

	// Expiration block count value 20
	DefaultExpirationBlockCount = uint64(20)

	// Execute gas cost value is 100000
	DefaultExecuteGas = uint64(100000)

	// Prepare gas cost value is 100000
	DefaultPrepareGas = uint64(100000)
)

// Parameter store keys.
var (
	KeyMaxRawRequestCount               = []byte("MaxRawRequestCount")
	KeyMaxRawDataReportSize             = []byte("MaxRawDataReportSize")
	KeyMaxResultSize                    = []byte("MaxResultSize")
	KeyGasPerRawDataRequestPerValidator = []byte("GasPerRawDataRequestPerValidator")
	KeyExpirationBlockCount             = []byte("ExpirationBlockCount")
	KeyExecuteGas                       = []byte("ExecuteGas")
	KeyPrepareGas                       = []byte("PrepareGas")
)

// String implements the stringer interface for Params.
func (p Params) String() string {
	return fmt.Sprintf(`oracle Params:
  MaxRawRequestCount:               %d
  MaxRawDataReportSize:             %d
  MaxResultSize:                    %d
  GasPerRawDataRequestPerValidator: %d
  ExpirationBlockCount:             %d
  ExecuteGas:                       %d
  PrepareGas:                       %d
`,
		p.MaxRawRequestCount,
		p.MaxRawDataReportSize,
		p.MaxResultSize,
		p.GasPerRawDataRequestPerValidator,
		p.ExpirationBlockCount,
		p.ExecuteGas,
		p.PrepareGas,
	)
}

func validateNoOp(_ interface{}) error { return nil }

// ParamSetPairs implements the params.ParamSet interface for Params.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	// TODO: Make validation real. Not just noop
	return params.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxRawRequestCount, &p.MaxRawRequestCount, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxRawDataReportSize, &p.MaxRawDataReportSize, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxResultSize, &p.MaxResultSize, validateNoOp),
		paramtypes.NewParamSetPair(KeyGasPerRawDataRequestPerValidator, &p.GasPerRawDataRequestPerValidator, validateNoOp),
		paramtypes.NewParamSetPair(KeyExpirationBlockCount, &p.ExpirationBlockCount, validateNoOp),
		paramtypes.NewParamSetPair(KeyExecuteGas, &p.ExecuteGas, validateNoOp),
		paramtypes.NewParamSetPair(KeyPrepareGas, &p.PrepareGas, validateNoOp),
	}
}

// DefaultParams defines the default parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxRawRequestCount,
		DefaultMaxRawDataReportSize,
		DefaultMaxResultSize,
		DefaultGasPerRawDataRequestPerValidator,
		DefaultExpirationBlockCount,
		DefaultExecuteGas,
		DefaultPrepareGas,
	)
}
