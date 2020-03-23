package app

import (
	"io"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
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

	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	// Staking genesis (Not used in our chain)
	// var stakingState staking.GenesisState
	// staking.ModuleCdc.MustUnmarshalJSON(genesisState[staking.ModuleName], &stakingState)

	// for _, val := range stakingState.Validators {
	// 	err := app.dbBand.AddValidator(val.GetOperator().String(), val.GetConsAddr().String())
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// Genutil genesis
	var genutilState genutil.GenesisState
	genutil.ModuleCdc.MustUnmarshalJSON(genesisState[genutil.ModuleName], &genutilState)

	for _, genTx := range genutilState.GenTxs {
		var tx authtypes.StdTx
		genutil.ModuleCdc.MustUnmarshalJSON(genTx, &tx)
		for _, msg := range tx.Msgs {
			if createMsg, ok := msg.(staking.MsgCreateValidator); ok {
				err := app.dbBand.AddValidator(
					createMsg.ValidatorAddress,
					createMsg.PubKey,
				)
				if err != nil {
					panic(err)
				}
			}
		}
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
	if res.IsOK() && lastProcessHeight+1 == app.DeliverContext.BlockHeight() {
		for _, event := range res.Events {
			kvMap := make(map[string]string)
			for _, kv := range event.Attributes {
				kvMap[string(kv.GetKey())] = string(kv.GetValue())
			}
			app.dbBand.HandleEvent(event.Type, kvMap)
		}
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

	for _, val := range req.GetLastCommitInfo().Votes {
		app.dbBand.AddValidatorUpTime(
			val.GetValidator().Address,
			req.Header.GetHeight()-1,
			val.GetSignedLastBlock(),
		)
	}

	err = app.dbBand.ClearOldVotes(req.Header.GetHeight() - 1)
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
