package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	KeyMintAir = []byte("MintAir")
)

// mint parameters
type Params struct {
	mint.Params
	MintAir bool `json:"mint_air" yaml:"mint_air"` // should the coins be produced from nowhere
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string, inflationRateChange, inflationMax, inflationMin, goalBonded sdk.Dec, blocksPerYear uint64,
) Params {

	return Params{
		Params: mint.Params{
			MintDenom:           mintDenom,
			InflationRateChange: inflationRateChange,
			InflationMax:        inflationMax,
			InflationMin:        inflationMin,
			GoalBonded:          goalBonded,
			BlocksPerYear:       blocksPerYear,
		},
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		Params:  mint.DefaultParams(),
		MintAir: false,
	}
}

// validate params
func (p Params) Validate() error {
	if err := p.Params.Validate(); err != nil {
		return err
	}
	return nil

}

func (p Params) String() string {
	return fmt.Sprintf(`Minting Params:
  Mint Denom:             %s
  Inflation Rate Change:  %s
  Inflation Max:          %s
  Inflation Min:          %s
  Goal Bonded:            %s
  Blocks Per Year:        %d
  MintAir:                %v
`,
		p.MintDenom, p.InflationRateChange, p.InflationMax,
		p.InflationMin, p.GoalBonded, p.BlocksPerYear, p.MintAir,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return append(p.Params.ParamSetPairs(), params.NewParamSetPair(KeyMintAir, &p.MintAir, validateMintAir))
}

func validateMintAir(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
