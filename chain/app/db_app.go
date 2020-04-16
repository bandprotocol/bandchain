package app

import (
	"io"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/bandprotocol/bandchain/chain/db"
	// "github.com/bandprotocol/bandchain/chain/x/oracle"
)

type dbBandApp struct {
	*bandApp
	dbBand *db.BandDB
	txNum  int64
}

func NewDBBandApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, dbBand *db.BandDB, baseAppOptions ...func(*bam.BaseApp),
) *dbBandApp {
	app := NewBandApp(logger, db, traceStore, loadLatest, invCheckPeriod, baseAppOptions...)
	dbBand.BankKeeper = app.BankKeeper
	dbBand.DistrKeeper = app.DistrKeeper
	dbBand.StakingKeeper = app.StakingKeeper
	// dbBand.OracleKeeper = app.OracleKeeper
	// dbBand.IBCKeeper = app.IBCKeeper
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

	// Genaccount genesis
	var genaccountsState genaccounts.GenesisState
	genaccounts.ModuleCdc.MustUnmarshalJSON(genesisState[genaccounts.ModuleName], &genaccountsState)

	for _, account := range genaccountsState {
		err := app.dbBand.SetAccountBalance(account.Address, account.Coins, 0)
		if err != nil {
			panic(err)
		}
	}

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
	app.cdc.MustUnmarshalJSON(genesisState[genutil.ModuleName], &genutilState)

	for _, genTx := range genutilState.GenTxs {
		var tx auth.StdTx
		app.cdc.MustUnmarshalJSON(genTx, &tx)
		for _, msg := range tx.Msgs {
			if createMsg, ok := msg.(staking.MsgCreateValidator); ok {
				app.dbBand.HandleMessage(nil, createMsg, nil)
				app.dbBand.DecreaseAccountBalance(
					createMsg.DelegatorAddress,
					sdk.NewCoins(createMsg.Value),
					0,
				)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	// // Oracle genesis
	// var oracleState oracle.GenesisState
	// app.cdc.MustUnmarshalJSON(genesisState[oracle.ModuleName], &oracleState)

	// // Save data source
	// for idx, dataSource := range oracleState.DataSources {
	// 	err := app.dbBand.AddDataSource(
	// 		int64(idx+1),
	// 		dataSource.Name,
	// 		dataSource.Description,
	// 		dataSource.Owner,
	// 		dataSource.Fee,
	// 		dataSource.Executable,
	// 		time.Now(),
	// 		0,
	// 		nil,
	// 	)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// // Save oracle script
	// for idx, oracleScript := range oracleState.OracleScripts {
	// 	err := app.dbBand.AddOracleScript(
	// 		int64(idx+1),
	// 		oracleScript.Name,
	// 		oracleScript.Description,
	// 		oracleScript.Owner,
	// 		oracleScript.Code,
	// 		time.Now(),
	// 		0,
	// 		nil,
	// 	)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	app.dbBand.Commit()

	return app.bandApp.InitChain(req)
}

func (app *dbBandApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	app.txNum++
	res = app.bandApp.DeliverTx(req)
	lastProcessHeight, err := app.dbBand.GetLastProcessedHeight()
	if err != nil {
		panic(err)
	}
	if lastProcessHeight+1 != app.DeliverContext.BlockHeight() {
		return res
	}

	tx, err := app.TxDecoder(req.Tx)
	if err != nil {
		panic(err)
	}
	if stdTx, ok := tx.(auth.StdTx); ok {
		// Add involved accounts
		involvedAccounts := stdTx.GetSigners()
		txHash := tmhash.Sum(req.Tx)
		app.dbBand.AddTransaction(
			app.txNum,
			txHash,
			app.DeliverContext.BlockTime(),
			res.GasUsed,
			stdTx.Fee.Gas,
			stdTx.Fee.Amount,
			stdTx.GetSigners()[0],
			res.IsOK(),
			app.DeliverContext.BlockHeight(),
		)
		if !res.IsOK() {
			app.dbBand.HandleTransactionFail(stdTx, txHash)
		} else {
			logs, err := sdk.ParseABCILogs(res.Log)
			if err != nil {
				panic(err)
			}

			app.dbBand.HandleTransaction(stdTx, txHash, logs)
			involvedAccounts = append(
				involvedAccounts, app.dbBand.GetInvolvedAccountsFromTx(stdTx)...,
			)
			involvedAccounts = append(
				involvedAccounts, app.dbBand.GetInvolvedAccountsFromTransferEvents(logs)...,
			)
		}
		updatedAccounts := make(map[string]bool)
		for _, account := range involvedAccounts {
			if found := updatedAccounts[account.String()]; !found {
				updatedAccounts[account.String()] = true
				err := app.dbBand.SetAccountBalance(
					account,
					app.dbBand.BankKeeper.GetCoins(app.DeliverContext, account),
					app.DeliverContext.BlockHeight(),
				)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return res
}

func (app *dbBandApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {
	app.txNum = 0
	res = app.bandApp.BeginBlock(req)
	// Begin transaction
	app.dbBand.BeginTransaction()
	app.dbBand.SetContext(app.DeliverContext)
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
		validator := app.StakingKeeper.ValidatorByConsAddr(app.DeliverContext, val.GetValidator().Address)
		reward := app.DistrKeeper.GetValidatorCurrentRewards(app.DeliverContext, validator.GetOperator())

		rewarValue := "0"
		if !reward.Rewards.IsZero() {
			rewarValue = reward.Rewards[0].Amount.String()
		}

		app.dbBand.UpdateValidator(
			validator.GetOperator(),
			&db.Validator{
				CurrentReward: rewarValue,
			},
		)
	}

	err = app.dbBand.ClearOldVotes(req.Header.GetHeight() - 1)
	if err != nil {
		panic(err)
	}

	app.dbBand.AddBlock(
		req.Header.GetHeight(),
		app.DeliverContext.BlockTime(),
		req.Header.GetProposerAddress(),
		req.GetHash(),
	)

	// Handle Begin block event
	events := res.GetEvents()
	for _, event := range events {
		app.dbBand.HandleBeginblockEvent(event)
	}

	return res
}

func (app *dbBandApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
	res = app.bandApp.EndBlock(req)
	inflation := app.bandApp.MintKeeper.GetMinter(app.bandApp.DeliverContext).Inflation.String()
	err := app.dbBand.SetInflationRate(inflation)
	if err != nil {
		panic(err)
	}
	err = app.dbBand.SetLastProcessedHeight(req.GetHeight())
	if err != nil {
		panic(err)
	}

	events := res.GetEvents()
	for _, event := range events {
		app.dbBand.HandleEndblockEvent(event)
	}

	app.dbBand.SetContext(sdk.Context{})
	return res
}

func (app *dbBandApp) Commit() (res abci.ResponseCommit) {
	res = app.bandApp.Commit()

	app.dbBand.Commit()
	return res
}
