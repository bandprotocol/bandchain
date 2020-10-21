package request

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type RequestHook struct {
	cdc          *codec.Codec
	oracleKeeper keeper.Keeper
	dbMap        *gorp.DbMap
	limit        uint64
}

func initDb(sqliteDir string) *gorp.DbMap {
	db, err := sql.Open("sqlite3", sqliteDir)
	if err != nil {
		panic(err)
	}
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTableWithName(Request{}, "request")
	err = dbMap.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}
	return dbMap
}

func NewRequestHook(cdc *codec.Codec, oracleKeeper keeper.Keeper, sqliteDir string) *RequestHook {
	return &RequestHook{
		cdc:          cdc,
		oracleKeeper: oracleKeeper,
		dbMap:        initDb(sqliteDir),
		limit:        100,
	}
}

func (h RequestHook) AfterInitChain(ctx sdk.Context, req abci.RequestInitChain, res abci.ResponseInitChain) {
}

func (h RequestHook) AfterBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock) {
}

func (h RequestHook) AfterDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx, res abci.ResponseDeliverTx) {
}

func (h RequestHook) AfterEndBlock(ctx sdk.Context, req abci.RequestEndBlock, res abci.ResponseEndBlock) {
	for _, event := range res.Events {
		events := sdk.StringifyEvents([]abci.Event{event})
		evMap := common.ParseEvents(events)
		_ = evMap
		_ = events
		switch event.Type {
		case types.EventTypeResolve:
			reqID := types.RequestID(common.Atoi(evMap[types.EventTypeResolve+"."+types.AttributeKeyID][0]))
			result := h.oracleKeeper.MustGetResult(ctx, reqID)

			if result.ResponsePacketData.ResolveStatus == types.ResolveStatus_Success {
				h.UpsertRequest(Request{
					OracleScriptID: result.RequestPacketData.OracleScriptID,
					Calldata:       result.RequestPacketData.Calldata,
					MinCount:       result.RequestPacketData.MinCount,
					AskCount:       result.RequestPacketData.AskCount,
					RequestIDs:     fmt.Sprintf("%d", reqID),
				})
			}
		default:
			break
		}
	}
}

func (h RequestHook) ApplyQuery(req abci.RequestQuery) (res abci.ResponseQuery, stop bool) {
	paths := strings.Split(req.Path, "/")
	if paths[0] == "band" {
		switch paths[1] {
		case "latest_request":
			if len(paths) != 6 {
				return common.QueryResultError(errors.New(fmt.Sprintf("expect 6 arguments given %d", len(paths)))), true
			}
			oid := types.OracleScriptID(common.Atoi(paths[2]))
			calldata, err := hex.DecodeString(paths[3])
			if err != nil {
				return abci.ResponseQuery{}, false
			}
			askCount := common.Atoi(paths[4])
			minCount := common.Atoi(paths[5])
			// TODO: fix this
			requestID := h.GetLatestRequestID(oid, calldata, uint64(askCount), uint64(minCount))
			bz, err := h.cdc.MarshalBinaryBare(requestID)
			return common.QueryResultSuccess(bz, req.Height), true

		default:
			return abci.ResponseQuery{}, false
		}
	} else {
		return abci.ResponseQuery{}, false
	}
}

func (h RequestHook) BeforeCommit() {}
