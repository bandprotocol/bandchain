package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
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
	case oracle.EventTypeRequestExecute:
		app.handleEventRequestExecute(evMap)
	default:
		break
	}
}
