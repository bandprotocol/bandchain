package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
)

func createTestInput() (*bandapp.BandApp, sdk.Context) {
	db := dbm.NewMemDB()
	app := bandapp.NewBandApp(log.NewNopLogger(), db, nil, true, 0, map[int64]bool{}, "")
	app.InitChain(abci.RequestInitChain{
		Validators:    []abci.ValidatorUpdate{},
		AppStateBytes: codec.MustMarshalJSONIndent(app.Codec(), bandapp.NewDefaultGenesisState()),
	})
	ctx := app.BaseApp.NewContext(false, abci.Header{})
	return app, ctx
}
