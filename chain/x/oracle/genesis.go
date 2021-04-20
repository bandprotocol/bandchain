package oracle

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// GenesisState is the oracle state that must be provided at genesis.
type GenesisState struct {
	Params        types.Params                  `json:"params" yaml:"params"`
	DataSources   []types.DataSource            `json:"data_sources"  yaml:"data_sources"`
	OracleScripts []types.OracleScript          `json:"oracle_scripts"  yaml:"oracle_scripts"`
	Reporters     []types.ReportersPerValidator `json:"reporters" yaml:"reporters"`
}

// DefaultGenesisState returns the default oracle genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:        types.DefaultParams(),
		DataSources:   []types.DataSource{},
		OracleScripts: []types.OracleScript{},
		Reporters:     []types.ReportersPerValidator{},
	}
}

// InitGenesis performs genesis initialization for the oracle module.
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	k.SetParam(ctx, types.KeyMaxRawRequestCount, data.Params.MaxRawRequestCount)
	k.SetParam(ctx, types.KeyMaxAskCount, data.Params.MaxAskCount)
	k.SetParam(ctx, types.KeyExpirationBlockCount, data.Params.ExpirationBlockCount)
	k.SetParam(ctx, types.KeyBaseRequestGas, data.Params.BaseRequestGas)
	k.SetParam(ctx, types.KeyPerValidatorRequestGas, data.Params.PerValidatorRequestGas)
	k.SetParam(ctx, types.KeySamplingTryCount, data.Params.SamplingTryCount)
	k.SetParam(ctx, types.KeyOracleRewardPercentage, data.Params.OracleRewardPercentage)
	k.SetParam(ctx, types.KeyInactivePenaltyDuration, data.Params.InactivePenaltyDuration)
	k.SetDataSourceCount(ctx, 0)
	k.SetOracleScriptCount(ctx, 0)
	k.SetRequestCount(ctx, 0)
	k.SetRequestLastExpired(ctx, 0)
	k.SetRollingSeed(ctx, make([]byte, types.RollingSeedSizeInBytes))
	for _, dataSource := range data.DataSources {
		_ = k.AddDataSource(ctx, dataSource)
	}
	for _, oracleScript := range data.OracleScripts {
		_ = k.AddOracleScript(ctx, oracleScript)
	}
	for _, reportersPerValidator := range data.Reporters {
		for _, reporter := range reportersPerValidator.Reporters {
			k.AddReporter(ctx, reportersPerValidator.Validator, reporter)
		}
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Params:        k.GetParams(ctx),
		DataSources:   k.GetAllDataSources(ctx),
		OracleScripts: k.GetAllOracleScripts(ctx),
		Reporters:     k.GetAllReporters(ctx),
	}
}

// GetGenesisStateFromAppState returns x/oracle GenesisState given raw application genesis state.
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
