package pricer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

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

func (app *App) handleBeginBlockEndBlockEvent(event abci.Event) {
	events := sdk.StringifyEvents([]abci.Event{event})
	evMap := parseEvents(events)
	switch event.Type {
	case types.EventTypeResolve:
		app.handleEventRequestExecute(evMap)
	default:
		break
	}
}
