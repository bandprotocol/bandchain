package oracle

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState is the oracle state that must be provided at genesis.
type GenesisState struct {
	Params        types.Params                `json:"params" yaml:"params"`
	DataSources   []types.DataSource          `json:"data_sources"  yaml:"data_sources"`
	OracleScripts []types.OracleScript        `json:"oracle_scripts"  yaml:"oracle_scripts"`
	ReportInfos   []types.ValidatorReportInfo `json:"report_infos" yaml:"report_infos"`
	Results       [][]byte                    `json:"results" yaml:"results"`
}

// DefaultGenesisState returns the default oracle genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:        types.DefaultParams(),
		DataSources:   []types.DataSource{},
		OracleScripts: []types.OracleScript{},
		ReportInfos:   []types.ValidatorReportInfo{},
		Results:       [][]byte{},
	}
}

// InitGenesis performs genesis initialization for the oracle module.
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	k.SetParam(ctx, types.KeyMaxRawRequestCount, data.Params.MaxRawRequestCount)
	k.SetParam(ctx, types.KeyMaxAskCount, data.Params.MaxAskCount)
	k.SetParam(ctx, types.KeyExpirationBlockCount, data.Params.ExpirationBlockCount)
	k.SetParam(ctx, types.KeyMaxConsecutiveMisses, data.Params.MaxConsecutiveMisses)
	k.SetParam(ctx, types.KeyBaseRequestGas, data.Params.BaseRequestGas)
	k.SetParam(ctx, types.KeyPerValidatorRequestGas, data.Params.PerValidatorRequestGas)
	for _, dataSource := range data.DataSources {
		_ = k.AddDataSource(ctx, dataSource)
	}
	for _, oracleScript := range data.OracleScripts {
		_ = k.AddOracleScript(ctx, oracleScript)
	}
	for _, info := range data.ReportInfos {
		k.SetValidatorReportInfo(ctx, info.Validator, info)
	}
	for idx, result := range data.Results {
		if result != nil {
			k.SetResult(ctx, types.RequestID(idx+1), result)
		}
	}
	k.SetRequestCount(ctx, int64(len(data.Results)))
	err := k.BindPort(ctx, PortID)
	if err != nil {
		panic(fmt.Sprintf("could not claim port capability: %v", err))
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Params:        k.GetParams(ctx),
		DataSources:   k.GetAllDataSources(ctx),
		OracleScripts: k.GetAllOracleScripts(ctx),
		ReportInfos:   k.GetAllValidatorReportInfos(ctx),
		Results:       k.GetAllResults(ctx),
	}
}
