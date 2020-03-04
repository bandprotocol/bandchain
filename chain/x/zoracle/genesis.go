package zoracle

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState is the zoracle state that must be provided at genesis.
type GenesisState struct {
	Params        types.Params         `json:"params" yaml:"params"` // module level parameters for zoracle
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
	k.SetMaxDataSourceExecutableSize(ctx, data.Params.MaxDataSourceExecutableSize)
	k.SetMaxOracleScriptCodeSize(ctx, data.Params.MaxOracleScriptCodeSize)
	k.SetMaxCalldataSize(ctx, data.Params.MaxCalldataSize)
	k.SetMaxDataSourceCountPerRequest(ctx, data.Params.MaxDataSourceCountPerRequest)
	k.SetMaxRawDataReportSize(ctx, data.Params.MaxRawDataReportSize)
	k.SetMaxResultSize(ctx, data.Params.MaxResultSize)
	k.SetEndBlockExecuteGasLimit(ctx, data.Params.EndBlockExecuteGasLimit)
	k.SetMaxNameLength(ctx, data.Params.MaxNameLength)
	k.SetMaxDescriptionLength(ctx, data.Params.MaxDescriptionLength)

	for _, dataSource := range data.DataSources {
		err := k.AddDataSource(
			ctx, dataSource.Owner, dataSource.Name, dataSource.Fee, dataSource.Executable,
		)
		if err != nil {
			panic(err)
		}
	}

	for _, oracleScript := range data.OracleScripts {
		err := k.AddOracleScript(
			ctx, oracleScript.Owner, oracleScript.Name, oracleScript.Code,
		)
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
