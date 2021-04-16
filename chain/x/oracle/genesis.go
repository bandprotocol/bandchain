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
	Params            types.Params              `json:"params" yaml:"params"`
	DataSources       []types.DataSource        `json:"data_sources"  yaml:"data_sources"`
	OracleScripts     []types.OracleScript      `json:"oracle_scripts"  yaml:"oracle_scripts"`
	DataSourceFiles   map[string][]byte         `json:"data_source_files" yaml:"data_source_files"`
	OracleScriptFiles map[string][]byte         `json:"oracle_script_files" yaml:"oracle_script_files"`
	Reporters         map[string]sdk.ValAddress `json:"reporters" yaml:"reporters"`
}

// DefaultGenesisState returns the default oracle genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:            types.DefaultParams(),
		DataSources:       []types.DataSource{},
		OracleScripts:     []types.OracleScript{},
		DataSourceFiles:   make(map[string][]byte),
		OracleScriptFiles: make(map[string][]byte),
		Reporters:         make(map[string]sdk.ValAddress),
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
		k.AddExecutableFile(data.DataSourceFiles[dataSource.Filename])
		_ = k.AddDataSource(ctx, dataSource)
	}
	for _, oracleScript := range data.OracleScripts {
		k.AddOracleScriptFile(data.OracleScriptFiles[oracleScript.Filename])
		_ = k.AddOracleScript(ctx, oracleScript)
	}
	for reporterAddrBech32, valAddr := range data.Reporters {
		reporterAddr, _ := sdk.AccAddressFromBech32(reporterAddrBech32)
		k.AddReporter(ctx, valAddr, reporterAddr)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	dataSources := k.GetAllDataSources(ctx)
	oracleScripts := k.GetAllOracleScripts(ctx)

	dataSourceFiles := make(map[string][]byte)
	for _, dataSource := range dataSources {
		dataSourceFile := k.GetFile(dataSource.Filename)
		if len(dataSourceFile) > 0 {
			dataSourceFiles[dataSource.Filename] = dataSourceFile
		}
	}

	oracleScriptFiles := make(map[string][]byte)
	for _, oracleScript := range oracleScripts {
		oracleScriptFile := k.GetFile(oracleScript.Filename)
		if len(oracleScriptFile) > 0 {
			oracleScriptFiles[oracleScript.Filename] = oracleScriptFile
		}
	}

	reporters := k.GetAllReporters(ctx)

	return GenesisState{
		Params:            params,
		DataSources:       dataSources,
		OracleScripts:     oracleScripts,
		DataSourceFiles:   dataSourceFiles,
		OracleScriptFiles: oracleScriptFiles,
		Reporters:         reporters,
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
