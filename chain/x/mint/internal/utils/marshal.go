package utils

import (
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
)

type InnerFormatParams struct {
	MintDenom           string  `json:"mint_denom" yaml:"mint_denom"`                       // type of coin to mint
	InflationRateChange sdk.Dec `json:"inflation_rate_change" yaml:"inflation_rate_change"` // maximum annual change in inflation rate
	InflationMax        sdk.Dec `json:"inflation_max" yaml:"inflation_max"`                 // maximum inflation rate
	InflationMin        sdk.Dec `json:"inflation_min" yaml:"inflation_min"`                 // minimum inflation rate
	GoalBonded          sdk.Dec `json:"goal_bonded" yaml:"goal_bonded"`                     // goal of percent bonded atoms
	BlocksPerYear       uint64  `json:"blocks_per_year" yaml:"blocks_per_year"`             // expected blocks per year
	MintAir             bool    `json:"mint_air" yaml:"mint_air"`                           // should the coins be produced from nowhere
}

type FormatParams struct {
	Params InnerFormatParams `json:"params" yaml:"params"`
}

type GenesisFormat struct {
	Params InnerFormatParams `json:"params" yaml:"params"`
	Minter mint.Minter       `json:"minter" yaml:"minter"` // minter object
}

func fromInnerFormat(params InnerFormatParams) types.Params {
	return types.Params{
		Params: mint.Params{
			MintDenom:           params.MintDenom,
			InflationRateChange: params.InflationRateChange,
			InflationMax:        params.InflationMax,
			InflationMin:        params.InflationMin,
			GoalBonded:          params.GoalBonded,
			BlocksPerYear:       params.BlocksPerYear,
		},
		MintAir: params.MintAir,
	}
}

func toInnerFormat(params types.Params) InnerFormatParams {
	return InnerFormatParams{
		MintDenom:           params.MintDenom,
		InflationRateChange: params.InflationRateChange,
		InflationMax:        params.InflationMax,
		InflationMin:        params.InflationMin,
		GoalBonded:          params.GoalBonded,
		BlocksPerYear:       params.BlocksPerYear,
		MintAir:             params.MintAir,
	}
}

func ToFormat(params types.Params) FormatParams {
	return FormatParams{
		Params: toInnerFormat(params),
	}
}

func FromFormat(params FormatParams) types.Params {
	return fromInnerFormat(params.Params)
}

func ToGenesisFormat(genesisState types.GenesisState) GenesisFormat {
	return GenesisFormat{
		Params: toInnerFormat(genesisState.Params),
		Minter: mint.Minter{
			Inflation:        genesisState.Minter.Inflation,
			AnnualProvisions: genesisState.Minter.AnnualProvisions,
		},
	}
}

func FromGenesisFormat(genesisFormat GenesisFormat) types.GenesisState {
	return types.GenesisState{
		Params: fromInnerFormat(genesisFormat.Params),
	}
}
