package oracle_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/simapp"
	me "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
)

var (
	Owner      = simapp.Owner
	Alice      = simapp.Alice
	Bob        = simapp.Bob
	Carol      = simapp.Carol
	Validator1 = simapp.Validator1
	Validator2 = simapp.Validator2
)

func createTestInput() (*bandapp.BandApp, sdk.Context, me.Keeper) {
	app := simapp.NewSimApp()
	ctx := app.BaseApp.NewContext(false, abci.Header{})
	return app, ctx, app.OracleKeeper
}
