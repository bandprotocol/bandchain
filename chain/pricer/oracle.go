package pricer

import (
	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	"github.com/bandprotocol/bandchain/chain/pkg/pricecache"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func (app *App) handleEventRequestExecute(evMap EvMap) {
	reqID := types.RequestID(atoi(evMap[types.EventTypeResolve+"."+types.AttributeKeyID][0]))
	result := app.OracleKeeper.MustGetResult(app.DeliverContext, reqID)

	if result.ResponsePacketData.ResolveStatus == types.ResolveStatus_Success {
		// Store latest successful request to request cache file
		app.reqCache.SaveLatestRequest(
			result.RequestPacketData.OracleScriptID, result.RequestPacketData.Calldata, result.RequestPacketData.AskCount, result.RequestPacketData.MinCount, reqID,
		)
		// Check that we need to store data to price cache file
		if app.StdOs[result.RequestPacketData.OracleScriptID] {
			var input Input
			var output Output
			obi.MustDecode(result.RequestPacketData.Calldata, &input)
			obi.MustDecode(result.ResponsePacketData.Result, &output)
			for idx, symbol := range input.Symbols {
				price := pricecache.NewPrice(input.Multiplier, output.Pxs[idx], result.ResponsePacketData.ResolveTime)
				err := app.priceCache.SetPrice(pricecache.GetFilename(symbol, result.RequestPacketData.MinCount, result.RequestPacketData.AskCount), price)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
