package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func parseEvents(events sdk.StringEvents) EvMap {
	evMap := make(EvMap)
	for _, event := range events {
		for _, kv := range event.Attributes {
			key := event.Type + "." + kv.Key
			evMap[key] = append(evMap[key], kv.Value)
		}
	}
	return evMap
}

// handleMsg handles the given message by publishing relevant events and populates accounts
// that need balance update in 'app.accs'. Also fills in extra info for this message.
func (app *App) handleMsg(txHash []byte, msg sdk.Msg, log sdk.ABCIMessageLog, extra JsDict) {
	evMap := parseEvents(log.Events)
	switch msg := msg.(type) {
	case oracle.MsgRequestData:
		app.handleMsgRequestData(txHash, msg, evMap, extra)
	case oracle.MsgReportData:
		app.handleMsgReportData(txHash, msg, evMap, extra)
	case oracle.MsgCreateDataSource:
		app.handleMsgCreateDataSource(txHash, msg, evMap, extra)
	case oracle.MsgCreateOracleScript:
		app.handleMsgCreateOracleScript(txHash, msg, evMap, extra)
	}
}

func (app *App) handleEndBlock(event abci.Event) {
	events := sdk.StringifyEvents([]abci.Event{event})
	evMap := parseEvents(events)
	switch event.Type {
	case types.EventTypeRequestExecute:
		resolveStatus := types.ResolveStatus(atoi(evMap[types.EventTypeRequestExecute+"."+types.AttributeKeyResolveStatus][0]))
		dict := JsDict{
			"id":             atoi(evMap[types.EventTypeRequestExecute+"."+types.AttributeKeyRequestID][0]),
			"request_time":   atoi(evMap[types.EventTypeRequestExecute+"."+types.AttributeKeyRequestTime][0]),
			"resolve_time":   atoi(evMap[types.EventTypeRequestExecute+"."+types.AttributeKeyResolveTime][0]),
			"resolve_status": resolveStatus,
		}
		app.Write("UPDATE_REQUEST", dict)
		if resolveStatus == types.ResolveStatus_Success {
			result := []byte(evMap[types.EventTypeRequestExecute+"."+types.AttributeKeyResult][0])
			dict["result"] = result
		}
	default:
		break
	}
}
