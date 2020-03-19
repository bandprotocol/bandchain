package app

import (
	"io"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/bandprotocol/d3n/chain/db"
)

type dbBandApp struct {
	*bandApp
	dbBand *db.BandDB
}

func NewDBBandApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, dbBand *db.BandDB, baseAppOptions ...func(*bam.BaseApp),
) *dbBandApp {
	app := NewBandApp(logger, db, traceStore, loadLatest, invCheckPeriod, baseAppOptions...)

	return &dbBandApp{bandApp: app, dbBand: dbBand}
}

func (app *dbBandApp) InitChain(req abci.RequestInitChain) abci.ResponseInitChain {
	app.dbBand.OpenTransaction()

	app.dbBand.SaveChainID(req.GetChainId())
	app.dbBand.SetBlockHeight(0)

	app.dbBand.Commit()

	return app.bandApp.InitChain(req)
}

func (app *dbBandApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	response := app.bandApp.DeliverTx(req)
	if response.IsOK() && app.dbBand.GetBlockHeight()+1 == app.DeliverContext.BlockHeight() {
		for _, event := range response.Events {
			kvMap := make(map[string]string)
			for _, kv := range event.Attributes {
				kvMap[string(kv.GetKey())] = string(kv.GetValue())
			}
			app.dbBand.HandleEvent(event.Type, kvMap)
		}
	}
	return response
}

func (app *dbBandApp) BeginBlock(req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	res := app.bandApp.BeginBlock(req)
	// Open transaction
	app.dbBand.OpenTransaction()
	if err := app.dbBand.ValidateChainID(app.DeliverContext.ChainID()); err != nil {
		panic(err)
	}

	return res
}

func (app *dbBandApp) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	res := app.bandApp.EndBlock(req)
	// Do other logic
	app.dbBand.SetBlockHeight(req.GetHeight())
	app.dbBand.Commit()
	return res
}
