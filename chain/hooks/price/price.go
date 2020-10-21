package price

import (
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

type PriceHook struct {
	cdc          *codec.Codec
	stdOs        map[types.OracleScriptID]bool
	oracleKeeper keeper.Keeper
	db           *leveldb.DB
}

func NewPriceHook(cdc *codec.Codec, oracleKeeper keeper.Keeper, oids []types.OracleScriptID, priceDBDir string) *PriceHook {
	stdOs := make(map[types.OracleScriptID]bool)
	for _, oid := range oids {
		stdOs[oid] = true
	}
	db, err := leveldb.OpenFile(priceDBDir, nil)
	if err != nil {
		panic(err)
	}
	return &PriceHook{
		cdc:          cdc,
		stdOs:        stdOs,
		oracleKeeper: oracleKeeper,
		db:           db,
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
				// Check that we need to store data to db
				if h.stdOs[result.RequestPacketData.OracleScriptID] {
					var input Input
					var output Output
					obi.MustDecode(result.RequestPacketData.Calldata, &input)
					obi.MustDecode(result.ResponsePacketData.Result, &output)
					for idx, symbol := range input.Symbols {
						price := NewPrice(symbol, input.Multiplier, output.Pxs[idx], result.ResponsePacketData.RequestID, result.ResponsePacketData.ResolveTime)
						err := h.db.Put([]byte(fmt.Sprintf("%s,%d,%d", symbol, result.RequestPacketData.MinCount, result.RequestPacketData.AskCount)), h.cdc.MustMarshalBinaryBare(price), nil)
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
func (h PriceHook) ApplyQuery(req abci.RequestQuery) (res abci.ResponseQuery, stop bool) {
	paths := strings.Split(req.Path, "/")
	if paths[0] == "band" {
		switch paths[1] {
		case "prices":
			if len(paths) < 3 {
				return common.QueryResultError(errors.New("no route for prices query specified")), true
			}
			bz, err := h.db.Get([]byte(paths[2]), nil)
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
