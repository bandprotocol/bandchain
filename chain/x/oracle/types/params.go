package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// nolint
const (
	DefaultParamspace             = ModuleName
	DefaultMaxDataSourceCount     = uint64(16)
	DefaultMaxAskCount            = uint64(16)
	DefaultExpirationBlockCount   = uint64(20)
	DefaultMaxConsecutiveMisses   = uint64(10)
	DefaultBaseRequestGas         = uint64(150000)
	DefaultPerValidatorRequestGas = uint64(30000)
)

// nolint
var (
	KeyMaxDataSourceCount     = []byte("MaxDataSourceCount")
	KeyMaxAskCount            = []byte("MaxAskCount")
	KeyExpirationBlockCount   = []byte("ExpirationBlockCount")
	KeyMaxConsecutiveMisses   = []byte("MaxConsecutiveMisses")
	KeyBaseRequestGas         = []byte("BaseRequestGas")
	KeyPerValidatorRequestGas = []byte("PerValidatorRequestGas")
)

// String implements the stringer interface for Params.
func (p Params) String() string {
	return fmt.Sprintf(`oracle Params:
  MaxDataSourceCount:     %d
  MaxAskCount:            %d
  ExpirationBlockCount:   %d
  MaxConsecutiveMisses:   %d
  BaseRequestGas          %d
  PerValidatorRequestGas: %d
`,
		p.MaxDataSourceCount,
		p.MaxAskCount,
		p.ExpirationBlockCount,
		p.MaxConsecutiveMisses,
		p.BaseRequestGas,
		p.PerValidatorRequestGas,
	)
}

// ParamSetPairs implements the params.ParamSet interface for Params.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxDataSourceCount, &p.MaxDataSourceCount, validateUint64("max data source count", true)),
		paramtypes.NewParamSetPair(KeyMaxAskCount, &p.MaxAskCount, validateUint64("max ask count", true)),
		paramtypes.NewParamSetPair(KeyExpirationBlockCount, &p.ExpirationBlockCount, validateUint64("expiration block count", true)),
		paramtypes.NewParamSetPair(KeyMaxConsecutiveMisses, &p.MaxConsecutiveMisses, validateUint64("max consecutive misses", false)),
		paramtypes.NewParamSetPair(KeyBaseRequestGas, &p.BaseRequestGas, validateUint64("base request gas", false)),
		paramtypes.NewParamSetPair(KeyPerValidatorRequestGas, &p.PerValidatorRequestGas, validateUint64("per validator request gas", false)),
	}
}

// DefaultParams defines the default parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxDataSourceCount,
		DefaultMaxAskCount,
		DefaultExpirationBlockCount,
		DefaultMaxConsecutiveMisses,
		DefaultBaseRequestGas,
		DefaultPerValidatorRequestGas,
	)
}

func validateUint64(name string, positiveOnly bool) func(interface{}) error {
	return func(i interface{}) error {
		v, ok := i.(uint64)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}
		if v <= 0 && positiveOnly {
			return fmt.Errorf("%s must be positive: %d", name, v)
		}
		return nil
	}
}
