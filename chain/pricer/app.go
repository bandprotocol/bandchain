package pricer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bandapp "github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/pkg/pricecache"
	"github.com/bandprotocol/bandchain/chain/pkg/requestcache"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// App extends the standard Band Cosmos-SDK application with Price cache.
type App struct {
	*bandapp.BandApp
	StdOs      map[types.OracleScriptID]bool
	priceCache pricecache.Cache
	reqCache   requestcache.Cache
}

func NewBandAppWithPricer(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string,
	disableFeelessReports bool, oids []types.OracleScriptID, priceCacheDir string, reqCacheDir string, baseAppOptions ...func(*bam.BaseApp),
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
		BandApp:    app,
		StdOs:      stdOs,
		priceCache: pricecache.New(priceCacheDir),
		reqCache:   requestcache.New(reqCacheDir),
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

func queryResultError(err error) abci.ResponseQuery {
	space, code, log := sdkerrors.ABCIInfo(err, true)
	return abci.ResponseQuery{
		Code:      code,
		Codespace: space,
		Log:       log,
	}
}

func queryResultSuccess(value []byte, height int64) abci.ResponseQuery {
	space, code, log := sdkerrors.ABCIInfo(nil, true)
	return abci.ResponseQuery{
		Code:      code,
		Codespace: space,
		Log:       log,
		Height:    height,
		Value:     value,
	}
}

// Query returns response query if the route is prices else calls into the underlying Query
func (app *App) Query(req abci.RequestQuery) abci.ResponseQuery {
	paths := strings.Split(req.Path, "/")
	if paths[0] == "band" {
		switch paths[1] {
		case "price":
			if len(paths) < 3 {
				return queryResultError(errors.New("no route for prices query specified"))
			}
			price, err := app.priceCache.GetPrice(paths[2])
			if err != nil {
				return queryResultError(err)
			}
			bz, err := app.Codec().MarshalBinaryBare(price)
			if err != nil {
				return queryResultError(err)
			}
			return queryResultSuccess(bz, req.Height)
		case "latest_request":
			if len(paths) != 6 {
				return queryResultError(fmt.Errorf("expect 6 arguments given %d", len(paths)))
			}
			oid := types.OracleScriptID(atoi(paths[2]))
			calldata, err := hex.DecodeString(paths[3])
			if err != nil {
				return queryResultError(err)
			}
			askCount := atoi(paths[4])
			minCount := atoi(paths[4])
			reqID, err := app.reqCache.GetLatestRequest(oid, calldata, askCount, minCount)
			if err != nil {
				return queryResultError(err)
			}
			bz, err := app.Codec().MarshalBinaryBare(reqID)
			return queryResultSuccess(bz, req.Height)
		default:
			return queryResultError(fmt.Errorf("no route for %s", paths[1]))
		}
	} else {
		return app.BandApp.Query(req)
	}
}
