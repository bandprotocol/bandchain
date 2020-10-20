package price

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type PriceHook struct {
	cdc          *codec.Codec
	stdOs        map[types.OracleScriptID]bool
	oracleKeeper keeper.Keeper
	dbMap        *gorp.DbMap
}

func initDb(sqliteDir string) *gorp.DbMap {
	db, err := sql.Open("sqlite3", sqliteDir)
	if err != nil {
		panic(err)
	}
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTableWithName(Price{}, "price")
	err = dbMap.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}
	return dbMap
}

func NewPriceHook(cdc *codec.Codec, oracleKeeper keeper.Keeper, oids []types.OracleScriptID, sqliteDir string) *PriceHook {
	stdOs := make(map[types.OracleScriptID]bool)
	for _, oid := range oids {
		stdOs[oid] = true
	}
	return &PriceHook{
		cdc:          cdc,
		stdOs:        stdOs,
		oracleKeeper: oracleKeeper,
		dbMap:        initDb(sqliteDir),
	}
}

func (h PriceHook) AfterInitChain(ctx sdk.Context, req abci.RequestInitChain, res abci.ResponseInitChain) {
}

func (h PriceHook) AfterBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock) {
}

func (h PriceHook) AfterDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx, res abci.ResponseDeliverTx) {

}

func (h PriceHook) AfterEndBlock(ctx sdk.Context, req abci.RequestEndBlock, res abci.ResponseEndBlock) {
	for _, event := range res.Events {
		events := sdk.StringifyEvents([]abci.Event{event})
		evMap := common.ParseEvents(events)
		switch event.Type {
		case types.EventTypeResolve:
			reqID := types.RequestID(common.Atoi(evMap[types.EventTypeResolve+"."+types.AttributeKeyID][0]))
			result := h.oracleKeeper.MustGetResult(ctx, reqID)

			if result.ResponsePacketData.ResolveStatus == types.ResolveStatus_Success {
				// Check that we need to store data to sqlite db
				if h.stdOs[result.RequestPacketData.OracleScriptID] {
					var input Input
					var output Output
					obi.MustDecode(result.RequestPacketData.Calldata, &input)
					obi.MustDecode(result.ResponsePacketData.Result, &output)
					for idx, symbol := range input.Symbols {
						h.UpsertPrice(Price{
							Symbol:      symbol,
							MinCount:    result.RequestPacketData.MinCount,
							AskCount:    result.RequestPacketData.AskCount,
							Multiplier:  input.Multiplier,
							Px:          output.Pxs[idx],
							RequestID:   result.ResponsePacketData.RequestID,
							ResolveTime: result.ResponsePacketData.ResolveTime,
						})
					}
				}
			}
		default:
			break
		}
	}
}

func (h PriceHook) ApplyQuery(req abci.RequestQuery) (res abci.ResponseQuery, stop bool) {
	paths := strings.Split(req.Path, "/")
	if paths[0] == "band" {
		switch paths[1] {
		case "prices":
			if len(paths) < 5 {
				return common.QueryResultError(errors.New("no route for prices query specified")), true
			}
			symbol := paths[2]
			minCount := common.Atoi(paths[3])
			askCount := common.Atoi(paths[4])
			price, err := h.dbMap.Get(Price{}, symbol, minCount, askCount)
			if err != nil {
				return common.QueryResultError(err), true
			}
			bz, err := h.cdc.MarshalBinaryBare(price)
			if err != nil {
				return common.QueryResultError(err), true
			}
			return common.QueryResultSuccess(bz, req.Height), true
		default:
			return abci.ResponseQuery{}, false
		}
	} else {
		return abci.ResponseQuery{}, false
	}
}

func (h PriceHook) BeforeCommit() {}
