package pricer

import (
	"errors"
	"io"
	"strings"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/pkg/pricecache"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// App extends the standard Band Cosmos-SDK application with Price cache
// functionality to act as an event producer for all events in the blockchains.
type App struct {
	*bandapp.BandApp
	StdOs map[types.OracleScriptID]bool
	cache pricecache.Cache
}

func NewBandAppWithPricer(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string,
	disableFeelessReports bool, oids []types.OracleScriptID, fileDir string, baseAppOptions ...func(*bam.BaseApp),
) *App {
	app := bandapp.NewBandApp(
		logger, db, traceStore, loadLatest, invCheckPeriod, skipUpgradeHeights,
		home, disableFeelessReports, baseAppOptions...,
	)
	stdOs := make(map[types.OracleScriptID]bool)
	for _, oid := range oids {
		stdOs[oid] = true
	}
	return &App{
		BandApp: app,
		StdOs:   stdOs,
		cache:   pricecache.New(fileDir),
	}
}

// EndBlock calls into the underlying EndBlock and save price to cache
func (app *App) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	res := app.BandApp.EndBlock(req)
	for _, event := range res.Events {
		app.handleBeginBlockEndBlockEvent(event)
	}
	return res
}

func QueryResultError(err error) abci.ResponseQuery {
	space, code, log := sdkerrors.ABCIInfo(err, true)
	return abci.ResponseQuery{
		Code:      code,
		Codespace: space,
		Log:       log,
	}
}

func QueryResultSuccess(value []byte, height int64) abci.ResponseQuery {
	space, code, log := sdkerrors.ABCIInfo(nil, true)
	return abci.ResponseQuery{
		Code:      code,
		Codespace: space,
		Log:       log,
		Height:    height,
		Value:     value,
	}
}

func (app *App) Query(req abci.RequestQuery) abci.ResponseQuery {
	paths := strings.Split(req.Path, "/")
	if paths[0] == "prices" {
		if len(paths) < 2 {
			return QueryResultError(errors.New("no route for prices query specified"))
		}
		price, err := app.cache.GetPrice(paths[1])
		if err != nil {
			return QueryResultError(err)
		}
		bz, err := app.Codec().MarshalBinaryBare(price)
		if err != nil {
			return QueryResultError(err)
		}
		return QueryResultSuccess(bz, req.Height)
	}
	return app.BandApp.Query(req)
}
