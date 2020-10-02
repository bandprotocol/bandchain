package pricer

import (
	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	"github.com/bandprotocol/bandchain/chain/pkg/pricecache"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func (app *App) handleEventRequestExecute(evMap EvMap) {
	reqID := types.RequestID(atoi(evMap[types.EventTypeResolve+"."+types.AttributeKeyID][0]))
	result := app.OracleKeeper.MustGetResult(app.DeliverContext, reqID)
	if _, ok := app.StdOs[result.RequestPacketData.OracleScriptID]; !ok {
		return
	}
	if result.ResponsePacketData.ResolveStatus == types.ResolveStatus_Success {
		var input Input
		var output Output
		obi.MustDecode(result.RequestPacketData.Calldata, &input)
		obi.MustDecode(result.ResponsePacketData.Result, &output)
		for idx, symbol := range input.Symbols {
			price := pricecache.NewPrice(input.Multiplier, output.Pxs[idx], result.ResponsePacketData.ResolveTime)
			err := app.cache.SetPrice(pricecache.GetFilename(symbol, result.RequestPacketData.MinCount, result.RequestPacketData.MinCount), price)
			if err != nil {
				panic(err)
			}
		}
	}
}
