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

func (app *dbBandApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	response := app.BaseApp.DeliverTx(req)
	if response.Code == 0 {
		app.Logger().Info("Happy Tx")
		for _, event := range response.Events {
			app.Logger().Info(event.Type, event.Attributes)
			if event.Type != "message" {
				app.dbBand.HandleEvent(event.Type)
			}
		}
	} else {
		app.Logger().Error("Failed Tx")
	}
	return response
}
