package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace.
const (
	DefaultParamspace = ModuleName

	// The maximum size of data source executable size, in bytes.
	// Default value is set to 10 kB.
	DefaultMaxDataSourceExecutableSize = uint64(10 * 1024)

	// The maximum size of Owasm code, in bytes.
	// Default value is set to 500 kB.
	DefaultMaxOracleScriptCodeSize = uint64(500 * 1024)

	// The maximum size of calldata when invoking for oracle scripts or data sources.
	// Default value is set 1 kB.
	DefaultMaxCalldataSize = uint64(1 * 1024)

	// The maximum number of raw requests that a request can make.
	// Default value is set to 16.
	DefaultMaxRawRequestCount = uint64(16)

	// The maximum size of raw data report per data source.
	// Default value is set to 1 kB.
	DefaultMaxRawDataReportSize = uint64(1 * 1024)

	// The maximum size of result after execution.
	// Default value is set 1 kB.
	DefaultMaxResultSize = uint64(1 * 1024)

	// The maximum size of name length.
	// Default value is 280
	DefaultMaxNameLength = uint64(280)

	// The maximum size of description length.
	// Default value 4096
	DefaultDescriptionLength = uint64(4096)

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
	KeyMaxExecutableSize                = []byte("MaxExecutableSize")
	KeyMaxOracleScriptCodeSize          = []byte("MaxOracleScriptCodeSize")
	KeyMaxCalldataSize                  = []byte("MaxCalldataSize")
	KeyMaxRawRequestCount               = []byte("MaxRawRequestCount")
	KeyMaxRawDataReportSize             = []byte("MaxRawDataReportSize")
	KeyMaxResultSize                    = []byte("MaxResultSize")
	KeyMaxNameLength                    = []byte("MaxNameLength")
	KeyMaxDescriptionLength             = []byte("MaxDescriptionLength")
	KeyGasPerRawDataRequestPerValidator = []byte("GasPerRawDataRequestPerValidator")
	KeyExpirationBlockCount             = []byte("ExpirationBlockCount")
	KeyExecuteGas                       = []byte("ExecuteGas")
	KeyPrepareGas                       = []byte("PrepareGas")
)

// Params - used for initializing default parameter for oracle at genesis.
type Params struct {
	MaxDataSourceExecutableSize      uint64 `json:"max_data_source_executable_size" yaml:"max_data_source_executable_size"`
	MaxOracleScriptCodeSize          uint64 `json:"max_oracle_script_code_size" yaml:"max_oracle_script_code_size"`
	MaxCalldataSize                  uint64 `json:"max_calldata_size" yaml:"max_calldata_size"`
	MaxRawRequestCount               uint64 `json:"max_raw_request_count" yaml:"max_raw_request_count%"`
	MaxRawDataReportSize             uint64 `json:"max_raw_data_report_size" yaml:"max_raw_data_report_size"`
	MaxResultSize                    uint64 `json:"max_result_size" yaml:"max_result_size"`
	MaxNameLength                    uint64 `json:"max_name_length" yaml:"max_name_length"`
	MaxDescriptionLength             uint64 `json:"max_description_length" yaml:"max_description_length"`
	GasPerRawDataRequestPerValidator uint64 `json:"gas_per_raw_data_request" yaml:"gas_per_raw_data_request"`
	ExpirationBlockCount             uint64 `json:"expiration_block_count"`
	ExecuteGas                       uint64 `json:"execute_gas"`
	PrepareGas                       uint64 `json:"prepare_gas"`
}

// NewParams creates a new Params object.
func NewParams(
	maxDataSourceExecutableSize uint64,
	maxOracleScriptCodeSize uint64,
	maxCalldataSize uint64,
	maxDataSourceCountPerRequest uint64,
	maxRawDataReportSize uint64,
	maxResultSize uint64,
	maxNameLength uint64,
	maxDescriptionLength uint64,
	gasPerRawDataRequestPerValidator uint64,
	expirationBlockCount uint64,
	executeGas uint64,
	prepareGas uint64,
) Params {
	return Params{
		MaxDataSourceExecutableSize:      maxDataSourceExecutableSize,
		MaxOracleScriptCodeSize:          maxOracleScriptCodeSize,
		MaxCalldataSize:                  maxCalldataSize,
		MaxRawRequestCount:               maxDataSourceCountPerRequest,
		MaxRawDataReportSize:             maxRawDataReportSize,
		MaxResultSize:                    maxResultSize,
		MaxNameLength:                    maxNameLength,
		MaxDescriptionLength:             maxDescriptionLength,
		GasPerRawDataRequestPerValidator: gasPerRawDataRequestPerValidator,
		ExpirationBlockCount:             expirationBlockCount,
		ExecuteGas:                       executeGas,
		PrepareGas:                       prepareGas,
	}
}

// String implements the stringer interface for Params.
func (p Params) String() string {
	return fmt.Sprintf(`oracle Params:
  MaxDataSourceExecutableSize:      %d
  MaxOracleScriptCodeSize:          %d
  MaxCalldataSize:                  %d
  MaxRawRequestCount:               %d
  MaxRawDataReportSize:             %d
  MaxResultSize:                    %d
  MaxNameLength:                    %d
  MaxDescriptionLength:             %d
  GasPerRawDataRequestPerValidator: %d
  ExpirationBlockCount:             %d
  ExecuteGas:                       %d
  PrepareGas:                       %d
`, p.MaxDataSourceExecutableSize,
		p.MaxOracleScriptCodeSize,
		p.MaxCalldataSize,
		p.MaxRawRequestCount,
		p.MaxRawDataReportSize,
		p.MaxResultSize,
		p.MaxNameLength,
		p.MaxDescriptionLength,
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
		paramtypes.NewParamSetPair(KeyMaxExecutableSize, &p.MaxDataSourceExecutableSize, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxOracleScriptCodeSize, &p.MaxOracleScriptCodeSize, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxCalldataSize, &p.MaxCalldataSize, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxRawRequestCount, &p.MaxRawRequestCount, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxRawDataReportSize, &p.MaxRawDataReportSize, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxResultSize, &p.MaxResultSize, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxNameLength, &p.MaxNameLength, validateNoOp),
		paramtypes.NewParamSetPair(KeyMaxDescriptionLength, &p.MaxDescriptionLength, validateNoOp),
		paramtypes.NewParamSetPair(KeyGasPerRawDataRequestPerValidator, &p.GasPerRawDataRequestPerValidator, validateNoOp),
		paramtypes.NewParamSetPair(KeyExpirationBlockCount, &p.ExpirationBlockCount, validateNoOp),
		paramtypes.NewParamSetPair(KeyExecuteGas, &p.ExecuteGas, validateNoOp),
		paramtypes.NewParamSetPair(KeyPrepareGas, &p.PrepareGas, validateNoOp),
	}
}

// DefaultParams defines the default parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxDataSourceExecutableSize,
		DefaultMaxOracleScriptCodeSize,
		DefaultMaxCalldataSize,
		DefaultMaxRawRequestCount,
		DefaultMaxRawDataReportSize,
		DefaultMaxResultSize,
		DefaultMaxNameLength,
		DefaultDescriptionLength,
		DefaultGasPerRawDataRequestPerValidator,
		DefaultExpirationBlockCount,
		DefaultExecuteGas,
		DefaultPrepareGas,
	)
}
