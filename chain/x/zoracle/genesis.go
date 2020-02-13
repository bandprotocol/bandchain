package zoracle

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GenesisState is the zoracle state that must be provided at genesis.
type GenesisState struct {
	// Scripts []types.StoredCode `json:"scripts"`
	Params types.Params `json:"params" yaml:"params"` // module level parameters for zoracle
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params types.Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

func ValidateGenesis(data GenesisState) error {
	// for _, record := range data.RequestRecords {
	// if record.Owner == nil {
	// 	return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Owner", record.Value)
	// }
	// if record.Value == "" {
	// 	return fmt.Errorf("invalid WhoisRecord: Owner: %s. Error: Missing Value", record.Owner)
	// }
	// if record.Price == nil {
	// 	return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Price", record.Value)
	// }
	// }
	return nil
}

// DefaultGenesisState returns the default genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	k.SetMaxDataSourceExecutableSize(ctx, data.Params.MaxDataSourceExecutableSize)
	k.SetMaxOracleScriptCodeSize(ctx, data.Params.MaxOracleScriptCodeSize)
	k.SetMaxCalldataSize(ctx, data.Params.MaxCalldataSize)
	k.SetMaxDataSourceCountPerRequest(ctx, data.Params.MaxDataSourceCountPerRequest)
	k.SetMaxRawDataReportSize(ctx, data.Params.MaxRawDataReportSize)

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return GenesisState{
		Params: k.GetParams(ctx),
	}
}
