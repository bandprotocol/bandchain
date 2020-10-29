package request

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-gorp/gorp"

	// DB driver
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// Hook inherits from Band app hook to save latest request into SQL database.
type Hook struct {
	cdc          *codec.Codec
	oracleKeeper keeper.Keeper
	dbMap        *gorp.DbMap
	trans        *gorp.Transaction
}

func getDB(driverName string, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	return db
}

func initDb(connStr string) *gorp.DbMap {
	connStrs := strings.Split(connStr, "://")
	if len(connStrs) != 2 {
		panic("failed to parse connection string")
	}
	var dbMap *gorp.DbMap
	switch connStrs[0] {
	case "sqlite3":
		dbMap = &gorp.DbMap{Db: getDB(connStrs[0], connStrs[1]), Dialect: gorp.SqliteDialect{}}
	case "postgres":
		dbMap = &gorp.DbMap{Db: getDB(connStrs[0], connStrs[1]), Dialect: gorp.PostgresDialect{}}
	case "mysql":
		dbMap = &gorp.DbMap{Db: getDB(connStrs[0], connStrs[1]), Dialect: gorp.MySQLDialect{}}
	default:
		panic(fmt.Sprintf("unknown driver %s", connStrs[0]))
	}
	indexName := "ix_calldata_min_count_ask_count_oracle_script_id_resolve_time"
	dbMap.AddTableWithName(Request{}, "request").AddIndex(indexName, "Btree", []string{"calldata", "min_count", "ask_count", "oracle_script_id", "resolve_time"})
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}
	err = dbMap.CreateIndex()
	// Check error if it's not creating existed index, panic the process.
	if err != nil && err.Error() != fmt.Sprintf("index %s already exists", indexName) {
		panic(err)
	}
	return dbMap
}

// NewHook creates a request hook instance that will be added in Band App.
func NewHook(cdc *codec.Codec, oracleKeeper keeper.Keeper, connStr string) *Hook {
	return &Hook{
		cdc:          cdc,
		oracleKeeper: oracleKeeper,
		dbMap:        initDb(connStr),
	}
}

// AfterInitChain specify actions need to do after chain initialization (app.Hook interface).
func (h *Hook) AfterInitChain(ctx sdk.Context, req abci.RequestInitChain, res abci.ResponseInitChain) {
}

// AfterBeginBlock specify actions need to do after begin block period (app.Hook interface).
func (h *Hook) AfterBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock) {
	trans, err := h.dbMap.Begin()
	if err != nil {
		panic(err)
	}
	h.trans = trans
}

// AfterDeliverTx specify actions need to do after transaction has been processed (app.Hook interface).
func (h *Hook) AfterDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx, res abci.ResponseDeliverTx) {
}

// AfterEndBlock specify actions need to do after end block period (app.Hook interface).
func (h *Hook) AfterEndBlock(ctx sdk.Context, req abci.RequestEndBlock, res abci.ResponseEndBlock) {
	for _, event := range res.Events {
		events := sdk.StringifyEvents([]abci.Event{event})
		evMap := common.ParseEvents(events)
		switch event.Type {
		case types.EventTypeResolve:
			reqID := types.RequestID(common.Atoi(evMap[types.EventTypeResolve+"."+types.AttributeKeyID][0]))
			result := h.oracleKeeper.MustGetResult(ctx, reqID)
			if result.ResponsePacketData.ResolveStatus == types.ResolveStatus_Success {
				h.insertRequest(reqID, result.RequestPacketData.OracleScriptID, result.RequestPacketData.Calldata,
					result.RequestPacketData.AskCount, result.RequestPacketData.MinCount,
					result.ResponsePacketData.ResolveTime)
			}
		default:
			break
		}
	}
}

// ApplyQuery catch the custom query that matches specific paths (app.Hook interface).
func (h *Hook) ApplyQuery(req abci.RequestQuery) (res abci.ResponseQuery, stop bool) {
	paths := strings.Split(req.Path, "/")
	if paths[0] == "band" {
		switch paths[1] {
		case "latest_request":
			if len(paths) != 7 {
				return common.QueryResultError(fmt.Errorf("expect 7 arguments given %d", len(paths))), true
			}
			oid := types.OracleScriptID(common.Atoi(paths[2]))
			calldata := paths[3]
			askCount := common.Atoui(paths[4])
			minCount := common.Atoui(paths[5])
			limit := common.Atoi(paths[6])
			requestIDs := h.getMultiRequestID(oid, calldata, askCount, minCount, limit)
			bz, err := h.cdc.MarshalBinaryBare(requestIDs)
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

// BeforeCommit specify actions need to do before commit block (app.Hook interface).
func (h *Hook) BeforeCommit() {
	err := h.trans.Commit()
	if err != nil {
		h.trans.Rollback()
	}
}
