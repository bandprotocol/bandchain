package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint
const (
	DefaultParamspace              = ModuleName
	DefaultMaxRawRequestCount      = uint64(12)
	DefaultMaxAskCount             = uint64(16)
	DefaultExpirationBlockCount    = uint64(100)
	DefaultBaseRequestGas          = uint64(150000)
	DefaultPerValidatorRequestGas  = uint64(30000)
	DefaultSamplingTryCount        = uint64(3)
	DefaultOracleRewardPercentage  = uint64(70)
	DefaultInactivePenaltyDuration = uint64(10 * time.Minute)
)

// nolint
var (
	KeyMaxRawRequestCount      = []byte("MaxRawRequestCount")
	KeyMaxAskCount             = []byte("MaxAskCount")
	KeyExpirationBlockCount    = []byte("ExpirationBlockCount")
	KeyBaseRequestGas          = []byte("BaseRequestGas")
	KeyPerValidatorRequestGas  = []byte("PerValidatorRequestGas")
	KeySamplingTryCount        = []byte("SamplingTryCount")
	KeyOracleRewardPercentage  = []byte("OracleRewardPercentage")
	KeyInactivePenaltyDuration = []byte("InactivePenaltyDuration")
)

// String implements the stringer interface for Params.
func (p Params) String() string {
	return fmt.Sprintf(`oracle Params:
  MaxRawRequestCount:      %d
  MaxAskCount:             %d
  ExpirationBlockCount:    %d
  BaseRequestGas           %d
  PerValidatorRequestGas:  %d
  SamplingTryCount:        %d
  OracleRewardPercentage:  %d
  InactivePenaltyDuration: %d
`,
		p.MaxRawRequestCount,
		p.MaxAskCount,
		p.ExpirationBlockCount,
		p.BaseRequestGas,
		p.PerValidatorRequestGas,
		p.SamplingTryCount,
		p.OracleRewardPercentage,
		p.InactivePenaltyDuration,
	)
}

// ParamSetPairs implements the params.ParamSet interface for Params.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMaxRawRequestCount, &p.MaxRawRequestCount, validateUint64("max data source count", true)),
		params.NewParamSetPair(KeyMaxAskCount, &p.MaxAskCount, validateUint64("max ask count", true)),
		params.NewParamSetPair(KeyExpirationBlockCount, &p.ExpirationBlockCount, validateUint64("expiration block count", true)),
		params.NewParamSetPair(KeyBaseRequestGas, &p.BaseRequestGas, validateUint64("base request gas", false)),
		params.NewParamSetPair(KeyPerValidatorRequestGas, &p.PerValidatorRequestGas, validateUint64("per validator request gas", false)),
		params.NewParamSetPair(KeySamplingTryCount, &p.SamplingTryCount, validateUint64("sampling try count", true)),
		params.NewParamSetPair(KeyOracleRewardPercentage, &p.OracleRewardPercentage, validateUint64("oracle reward percentage", false)),
		params.NewParamSetPair(KeyInactivePenaltyDuration, &p.InactivePenaltyDuration, validateUint64("inactive penalty duration", false)),
	}
}

// DefaultParams defines the default parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxRawRequestCount,
		DefaultMaxAskCount,
		DefaultExpirationBlockCount,
		DefaultBaseRequestGas,
		DefaultPerValidatorRequestGas,
		DefaultSamplingTryCount,
		DefaultOracleRewardPercentage,
		DefaultInactivePenaltyDuration,
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
