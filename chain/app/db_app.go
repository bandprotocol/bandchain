package app

import (
	"io"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/bandprotocol/bandchain/chain/db"
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
	app.dbBand.BeginTransaction()

	err := app.dbBand.SaveChainID(req.GetChainId())
	if err != nil {
		panic(err)
	}
	err = app.dbBand.SetLastProcessedHeight(0)
	if err != nil {
		panic(err)
	}

	app.dbBand.Commit()

	return app.bandApp.InitChain(req)
}

func (app *dbBandApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	res = app.bandApp.DeliverTx(req)
	lastProcessHeight, err := app.dbBand.GetLastProcessedHeight()
	if err != nil {
		panic(err)
	}
	if lastProcessHeight+1 == app.DeliverContext.BlockHeight() {
		return res
	}
	if !res.IsOK() {
		// TODO: We should not completely ignore failed transactions.
		return res
	}
	logs, err := sdk.ParseABCILogs(res.Log)
	if err != nil {
		panic(err)
	}
	tx, err := app.TxDecoder(req.Tx)
	if err != nil {
		panic(err)
	}
	if stdTx, ok := tx.(auth.StdTx); ok {
		txHash := tmhash.Sum(req.Tx)
		app.dbBand.HandleTransaction(stdTx, txHash, logs)
	}
	return res
}

func (app *dbBandApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {
	res = app.bandApp.BeginBlock(req)
	// Begin transaction
	app.dbBand.BeginTransaction()
	err := app.dbBand.ValidateChainID(app.DeliverContext.ChainID())
	if err != nil {
		panic(err)
	}

	return res
}

func (app *dbBandApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
	res = app.bandApp.EndBlock(req)
	err := app.dbBand.SetLastProcessedHeight(req.GetHeight())
	if err != nil {
		panic(err)
	}
	// Do other logic
	return res
}

func (app *dbBandApp) Commit() (res abci.ResponseCommit) {
	res = app.bandApp.Commit()

	app.dbBand.Commit()
	return res
}
