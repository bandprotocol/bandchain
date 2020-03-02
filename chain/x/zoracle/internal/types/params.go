package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace.
const (
	DefaultParamspace = ModuleName

	// The maximum size of data source executable size, in bytes.
	// Default value is set to 10 kB.
	DefaultMaxDataSourceExecutableSize = int64(10 * 1024)

	// The maximum size of Owasm code, in bytes.
	// Default value is set to 500 kB.
	DefaultMaxOracleScriptCodeSize = int64(500 * 1024)

	// The maximum size of calldata when invoking for oracle scripts or data sources.
	// Default value is set 1 kB.
	DefaultMaxCalldataSize = int64(1 * 1024)

	// The maximum number of data sources a request can make.
	// Default value is set to 16.
	DefaultMaxDataSourceCountPerRequest = int64(16)

	// The maximum size of raw data report per data source.
	// Default value is set to 1 kB.
	DefaultMaxRawDataReportSize = int64(1 * 1024)

	// The maximum size of result after execution.
	// Default value is set 1 kB.
	DefaultMaxResultSize = int64(1 * 1024)

	// The maximum gas that can be used to resolve requests at endblock time
	// Default value is 1000000
	DefaultEndBlockExecuteGasLimit = uint64(1000000)

	// The maximum size of name length.
	// Default value is 280
	DefaultMaxNameLength = int64(280)
)

// Parameter store keys.
var (
	KeyMaxDataSourceExecutableSize  = []byte("MaxDataSourceExecutableSize")
	KeyMaxOracleScriptCodeSize      = []byte("MaxOracleScriptCodeSize")
	KeyMaxCalldataSize              = []byte("MaxCalldataSize")
	KeyMaxDataSourceCountPerRequest = []byte("MaxDataSourceCountPerRequest")
	KeyMaxRawDataReportSize         = []byte("MaxRawDataReportSize")
	KeyMaxResultSize                = []byte("MaxResultSize")
	KeyEndBlockExecuteGasLimit      = []byte("EndBlockExecuteGasLimit")
	KeyMaxNameLength                = []byte("MaxNameLength")
)

// Params - used for initializing default parameter for zoracle at genesis.
type Params struct {
	MaxDataSourceExecutableSize  int64  `json:"max_data_source_executable_size" yaml:"max_data_source_executable_size"`
	MaxOracleScriptCodeSize      int64  `json:"max_oracle_script_code_size" yaml:"max_oracle_script_code_size"`
	MaxCalldataSize              int64  `json:"max_calldata_size" yaml:"max_calldata_size"`
	MaxDataSourceCountPerRequest int64  `json:"max_data_source_count_per_request" yaml:"max_data_source_count_per_request"`
	MaxRawDataReportSize         int64  `json:"max_raw_data_report_size" yaml:"max_raw_data_report_size"`
	MaxResultSize                int64  `json:"max_result_size" yaml:"max_result_size"`
	EndBlockExecuteGasLimit      uint64 `json:"end_block_execute_gas_limit" yaml:"end_block_execute_gas_limit"`
	MaxNameLength                int64  `json:"max_name_length" yaml:"max_name_length"`
}

// NewParams creates a new Params object.
func NewParams(
	maxDataSourceExecutableSize int64,
	maxOracleScriptCodeSize int64,
	maxCalldataSize int64,
	maxDataSourceCountPerRequest int64,
	maxRawDataReportSize int64,
	maxResultSize int64,
	endBlockExecuteGasLimit uint64,
	maxNameLength int64,
) Params {
	return Params{
		MaxDataSourceExecutableSize:  maxDataSourceExecutableSize,
		MaxOracleScriptCodeSize:      maxOracleScriptCodeSize,
		MaxCalldataSize:              maxCalldataSize,
		MaxDataSourceCountPerRequest: maxDataSourceCountPerRequest,
		MaxRawDataReportSize:         maxRawDataReportSize,
		MaxResultSize:                maxResultSize,
		EndBlockExecuteGasLimit:      endBlockExecuteGasLimit,
		MaxNameLength:                maxNameLength,
	}
}

// String implements the stringer interface for Params.
func (p Params) String() string {
	return fmt.Sprintf(`Zoracle Params:
  MaxDataSourceExecutableSize:  %d
  MaxOracleScriptCodeSize:      %d
  MaxCalldataSize:              %d
  MaxDataSourceCountPerRequest: %d
  MaxRawDataReportSize:         %d
  MaxResultSize:                %d
  EndBlockExecuteGasLimit:      %d
  MaxNameLength					%d
`, p.MaxDataSourceExecutableSize,
		p.MaxOracleScriptCodeSize,
		p.MaxCalldataSize,
		p.MaxDataSourceCountPerRequest,
		p.MaxRawDataReportSize,
		p.MaxResultSize,
		p.EndBlockExecuteGasLimit,
		p.MaxNameLength,
	)
}

// ParamSetPairs implements the params.ParamSet interface for Params.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMaxDataSourceExecutableSize, Value: &p.MaxDataSourceExecutableSize},
		{Key: KeyMaxOracleScriptCodeSize, Value: &p.MaxOracleScriptCodeSize},
		{Key: KeyMaxCalldataSize, Value: &p.MaxCalldataSize},
		{Key: KeyMaxDataSourceCountPerRequest, Value: &p.MaxDataSourceCountPerRequest},
		{Key: KeyMaxRawDataReportSize, Value: &p.MaxRawDataReportSize},
		{Key: KeyMaxResultSize, Value: &p.MaxResultSize},
		{Key: KeyEndBlockExecuteGasLimit, Value: &p.EndBlockExecuteGasLimit},
		{Key: KeyMaxNameLength, Value: &p.MaxNameLength},
	}
}

// DefaultParams defines the default parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxDataSourceExecutableSize,
		DefaultMaxOracleScriptCodeSize,
		DefaultMaxCalldataSize,
		DefaultMaxDataSourceCountPerRequest,
		DefaultMaxRawDataReportSize,
		DefaultMaxResultSize,
		DefaultEndBlockExecuteGasLimit,
		DefaultMaxNameLength,
	)
}
