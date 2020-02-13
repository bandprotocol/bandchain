package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName

	// the maximum size of data source executable size, in bytes.
	// default value is set to 10000 bytes or 10 kb
	DefaultMaxDataSourceExecutableSize = int64(10000)

	// the maximum size of Owasm code, in bytes.
	// default value is set to 500000 bytes or 500 kb
	DefaultMaxOracleScriptCodeSize = int64(500000)

	// the maximum size of calldata when invoking for oracle scripts or data sources.
	// default value is set to 1000 bytes or 1 kb
	DefaultMaxCalldataSize = int64(1000)

	// the maximum number of data sources a request can make.
	// default value is set to 16.
	DefaultMaxDataSourceCountPerRequest = int64(16)

	// the maximum size of raw data report per data source.
	// default value is set to 1000 bytes or 1 kb
	DefaultMaxRawDataReportSize = int64(1000)
)

// Parameter store keys
var (
	KeyMaxDataSourceExecutableSize  = []byte("MaxDataSourceExecutableSize")
	KeyMaxOracleScriptCodeSize      = []byte("MaxOracleScriptCodeSize")
	KeyMaxCalldataSize              = []byte("MaxCalldataSize")
	KeyMaxDataSourceCountPerRequest = []byte("MaxDataSourceCountPerRequest")
	KeyMaxRawDataReportSize         = []byte("MaxRawDataReportSize")
)

// Params - used for initializing default parameter for zoracle at genesis
type Params struct {
	MaxDataSourceExecutableSize  int64 `json:"max_data_source_executable_size" yaml:"max_data_source_executable_size"`
	MaxOracleScriptCodeSize      int64 `json:"max_oracle_script_code_size" yaml:"max_oracle_script_code_size"`
	MaxCalldataSize              int64 `json:"max_calldata_size" yaml:"max_calldata_size"`
	MaxDataSourceCountPerRequest int64 `json:"max_data_source_count_per_request" yaml:"max_data_source_count_per_request"`
	MaxRawDataReportSize         int64 `json:"max_raw_data_report_size" yaml:"max_raw_data_report_size"`
}

// ParamKeyTable for slashing module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	maxDataSourceExecutableSize int64,
	maxOracleScriptCodeSize int64,
	maxCalldataSize int64,
	maxDataSourceCountPerRequest int64,
	maxRawDataReportSize int64,
) Params {
	return Params{
		MaxDataSourceExecutableSize:  maxDataSourceExecutableSize,
		MaxOracleScriptCodeSize:      maxOracleScriptCodeSize,
		MaxCalldataSize:              maxCalldataSize,
		MaxDataSourceCountPerRequest: maxDataSourceCountPerRequest,
		MaxRawDataReportSize:         maxRawDataReportSize,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Slashing Params:
		MaxDataSourceExecutableSize:          %d
		MaxOracleScriptCodeSize:      %d
		MaxCalldataSize:      %d
		MaxDataSourceCountPerRequest:    %d
		MaxRawDataReportSize: %d`, p.MaxDataSourceExecutableSize,
		p.MaxOracleScriptCodeSize, p.MaxCalldataSize,
		p.MaxDataSourceCountPerRequest, p.MaxRawDataReportSize)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMaxDataSourceExecutableSize, Value: &p.MaxDataSourceExecutableSize},
		{Key: KeyMaxOracleScriptCodeSize, Value: &p.MaxOracleScriptCodeSize},
		{Key: KeyMaxCalldataSize, Value: &p.MaxCalldataSize},
		{Key: KeyMaxDataSourceCountPerRequest, Value: &p.MaxDataSourceCountPerRequest},
		{Key: KeyMaxRawDataReportSize, Value: &p.MaxRawDataReportSize},
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(
		DefaultMaxDataSourceCountPerRequest,
		DefaultMaxOracleScriptCodeSize,
		DefaultMaxCalldataSize,
		DefaultMaxDataSourceCountPerRequest,
		DefaultMaxRawDataReportSize,
	)
}
