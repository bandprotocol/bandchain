package oracle

import (
	"fmt"

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
	k.SetParam(ctx, KeyMaxRawRequestCount, data.Params.MaxRawRequestCount)
	k.SetParam(ctx, KeyMaxResultSize, data.Params.MaxResultSize)
	k.SetParam(ctx, KeyGasPerRawDataRequestPerValidator, data.Params.GasPerRawDataRequestPerValidator)
	k.SetParam(ctx, KeyExpirationBlockCount, data.Params.ExpirationBlockCount)

	for _, dataSource := range data.DataSources {
		_, err := k.AddDataSource(ctx, types.NewDataSource(
			dataSource.Owner, dataSource.Name, dataSource.Description, dataSource.Executable,
		))
		if err != nil {
			panic(err)
		}
	}

	for _, oracleScript := range data.OracleScripts {
		_, err := k.AddOracleScript(ctx, types.NewOracleScript(
			oracleScript.Owner, oracleScript.Name, oracleScript.Description,
			oracleScript.Code, oracleScript.Schema, oracleScript.SourceCodeURL,
		))
		if err != nil {
			panic(err)
		}
	}
	err := k.BindPort(ctx, PortID)
	if err != nil {
		panic(fmt.Sprintf("could not claim port capability: %v", err))
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
