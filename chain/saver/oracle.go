package saver

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/pkg/obi"
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
		fmt.Println("->", result.ResponsePacketData.ResolveTime)
		for idx, symbol := range input.Symbols {
			err := app.cahce.AddFile(symbol, result.RequestPacketData.MinCount, result.RequestPacketData.MinCount, input.Multiplier, output.Pxs[idx], result.ResponsePacketData.ResolveTime)
			if err != nil {
				panic(err)
			}
		}
	}
}
