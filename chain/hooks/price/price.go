package price

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/syndtr/goleveldb/leveldb"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// Hook uses levelDB to store the latest price of standard price reference.
type Hook struct {
	cdc          *codec.Codec
	stdOs        map[types.OracleScriptID]bool
	oracleKeeper keeper.Keeper
	db           *leveldb.DB
}

// NewHook creates a price hook instance that will be added in Band App.
func NewHook(cdc *codec.Codec, oracleKeeper keeper.Keeper, oids []types.OracleScriptID, priceDBDir string) *Hook {
	stdOs := make(map[types.OracleScriptID]bool)
	for _, oid := range oids {
		stdOs[oid] = true
	}
	db, err := leveldb.OpenFile(priceDBDir, nil)
	if err != nil {
		panic(err)
	}
	return &Hook{
		cdc:          cdc,
		stdOs:        stdOs,
		oracleKeeper: oracleKeeper,
		db:           db,
	}
}

// AfterInitChain specify actions need to do after chain initialization (app.Hook interface).
func (h *Hook) AfterInitChain(ctx sdk.Context, req abci.RequestInitChain, res abci.ResponseInitChain) {
}

// AfterBeginBlock specify actions need to do after begin block period (app.Hook interface).
func (h *Hook) AfterBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock) {
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
				// Check that we need to store data to db
				if h.stdOs[result.RequestPacketData.OracleScriptID] {
					var input Input
					var output Output
					obi.MustDecode(result.RequestPacketData.Calldata, &input)
					obi.MustDecode(result.ResponsePacketData.Result, &output)
					for idx, symbol := range input.Symbols {
						price := NewPrice(symbol, input.Multiplier, output.Pxs[idx], result.ResponsePacketData.RequestID, result.ResponsePacketData.ResolveTime)
						err := h.db.Put([]byte(fmt.Sprintf("%d,%d,%s", result.RequestPacketData.AskCount, result.RequestPacketData.MinCount, symbol)),
							h.cdc.MustMarshalBinaryBare(price), nil)
						if err != nil {
							panic(err)
						}
					}
				}
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
		case "prices":
			if len(paths) < 5 {
				return common.QueryResultError(errors.New("no route for prices query specified")), true
			}
			symbol := paths[2]
			askCount := common.Atoui(paths[3])
			minCount := common.Atoui(paths[4])
			bz, err := h.db.Get([]byte(fmt.Sprintf("%d,%d,%s", askCount, minCount, symbol)), nil)
			if err != nil {
				return common.QueryResultError(fmt.Errorf(
					"Cannot get price of %s with %d/%d counts with error: %s",
					symbol, minCount, askCount, err.Error(),
				)), true
			}
			return common.QueryResultSuccess(bz, req.Height), true
		case "price_symbols":
			if len(paths) < 4 {
				return common.QueryResultError(errors.New("no route for symbol prices query specified")), true
			}
			askCount := common.Atoui(paths[2])
			minCount := common.Atoui(paths[3])
			prefix := []byte(fmt.Sprintf("%d,%d,", askCount, minCount))

			it := h.db.NewIterator(nil, nil)
			it.Seek(prefix)

			symbols := []string{}
			for ; it.Valid() && bytes.HasPrefix(it.Key(), prefix); it.Next() {
				var p Price
				h.cdc.MustUnmarshalBinaryBare(it.Value(), &p)
				symbols = append(symbols, p.Symbol)
			}

			bz := h.cdc.MustMarshalBinaryBare(symbols)
			return common.QueryResultSuccess(bz, req.Height), true
		default:
			return abci.ResponseQuery{}, false
		}
	} else {
		return abci.ResponseQuery{}, false
	}
}

// BeforeCommit specify actions need to do before commit block (app.Hook interface).
func (h *Hook) BeforeCommit() {}
