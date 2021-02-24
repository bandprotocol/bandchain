package mint

import (
	"encoding/json"
	"fmt"
	"github.com/GeoDB-Limited/odincore/chain/x/mint/client/cli"
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/keeper"
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//____________________________________________________________________________

type AppModuleBasic struct {
	mint.AppModuleBasic
}

// AppModule implements an application module for the odinmint module.
type AppModule struct {
	AppModuleBasic
	mint.AppModule

	keeper Keeper
}

// GetQueryCmd returns the root query command for the mint module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the mint
// module.
func (amb AppModuleBasic) DefaultGenesis() json.RawMessage {
	return mint.ModuleCdc.MustMarshalJSON(utils.ToGenesisFormat(DefaultGenesisState()))
}

// Just call the specific underlying method
// GetTxCmd returns no root tx command for the mint module.
func (amb AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return amb.AppModuleBasic.GetTxCmd(cdc)
}

// RegisterCodec registers the mint module's types for the given codec.
func (amb AppModuleBasic) RegisterCodec(cdc *codec.Codec) { amb.AppModuleBasic.RegisterCodec(cdc) }

// RegisterRESTRoutes registers the REST routes for the mint module.
func (amb AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	amb.AppModuleBasic.RegisterRESTRoutes(ctx, rtr)
}

// ValidateGenesis performs genesis state validation for the mint module.
func (amb AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var tempData utils.GenesisFormat
	if err := mint.ModuleCdc.UnmarshalJSON(bz, &tempData); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", mint.ModuleName, err)
	}

	return ValidateGenesis(utils.FromGenesisFormat(tempData))
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, appModule mint.AppModule) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		AppModule:      appModule,
		keeper:         keeper,
	}
}

// Name returns the mint module's name.
func (AppModule) Name() string {
	return mint.ModuleName
}

// NewQuerierHandler returns the mint module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the mint module. It returns
// no validator updates.
// todo
// ModuleCdc does not allow to unmarshal embedded struct correctly
// that is why, we need this intermediate struct, in order to read in it
// we should later come up with better solution
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisFormat utils.GenesisFormat
	mint.ModuleCdc.MustUnmarshalJSON(data, &genesisFormat)
	InitGenesis(ctx, am.keeper, utils.FromGenesisFormat(genesisFormat))
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the mint
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return mint.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the mint module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}
