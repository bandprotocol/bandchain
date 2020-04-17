package oracle

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState is the oracle state that must be provided at genesis.
type GenesisState struct {
	Params        types.Params         `json:"params" yaml:"params"` // module level parameters for oracle
	DataSources   []types.DataSource   `json:"data_sources"  yaml:"data_sources"`
	OracleScripts []types.OracleScript `json:"oracle_scripts"  yaml:"oracle_scripts"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	params types.Params, dataSources []types.DataSource, oracleScripts []types.OracleScript,
) GenesisState {
	return GenesisState{
		Params:        params,
		DataSources:   dataSources,
		OracleScripts: oracleScripts,
	}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}

// DefaultGenesisState returns the default genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:        DefaultParams(),
		DataSources:   []types.DataSource{},
		OracleScripts: []types.OracleScript{},
	}
}

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	k.SetParam(ctx, KeyMaxExecutableSize, data.Params.MaxDataSourceExecutableSize)
	k.SetParam(ctx, KeyMaxOracleScriptCodeSize, data.Params.MaxOracleScriptCodeSize)
	k.SetParam(ctx, KeyMaxCalldataSize, data.Params.MaxCalldataSize)
	k.SetParam(ctx, KeyMaxDataSourceCountPerRequest, data.Params.MaxDataSourceCountPerRequest)
	k.SetParam(ctx, KeyMaxRawDataReportSize, data.Params.MaxRawDataReportSize)
	k.SetParam(ctx, KeyMaxResultSize, data.Params.MaxResultSize)
	k.SetParam(ctx, KeyEndBlockExecuteGasLimit, data.Params.EndBlockExecuteGasLimit)
	k.SetParam(ctx, KeyMaxNameLength, data.Params.MaxNameLength)
	k.SetParam(ctx, KeyMaxDescriptionLength, data.Params.MaxDescriptionLength)
	k.SetParam(ctx, KeyGasPerRawDataRequestPerValidator, data.Params.GasPerRawDataRequestPerValidator)
	k.SetParam(ctx, KeyExpirationBlockCount, data.Params.ExpirationBlockCount)
	k.SetParam(ctx, KeyExecuteGas, data.Params.ExecuteGas)
	k.SetParam(ctx, KeyPrepareGas, data.Params.PrepareGas)

	for _, dataSource := range data.DataSources {
		_, err := k.AddDataSource(
			ctx,
			dataSource.Owner,
			dataSource.Name,
			dataSource.Description,
			dataSource.Fee,
			dataSource.Executable,
		)
		if err != nil {
			panic(err)
		}
	}

	for _, oracleScript := range data.OracleScripts {
		_, err := k.AddOracleScript(
			ctx, oracleScript.Owner, oracleScript.Name, oracleScript.Description, oracleScript.Code, oracleScript.Schema, oracleScript.SourceCodeURL)
		if err != nil {
			panic(err)
		}
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Params:        k.GetParams(ctx),
		DataSources:   k.GetAllDataSources(ctx),
		OracleScripts: k.GetAllOracleScripts(ctx),
	}
}
