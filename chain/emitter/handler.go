package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

// handleMsg handles the given message by publishing relevant events and populates accounts
// that need balance update in 'app.accs'. Also fills in extra info for this message.
func (app *App) handleMsg(txHash []byte, msg sdk.Msg, log sdk.ABCIMessageLog, extra JsDict) {
	evMap := make(EvMap)
	for _, event := range log.Events {
		for _, kv := range event.Attributes {
			key := event.Type + "." + kv.Key
			evMap[key] = append(evMap[key], kv.Value)
		}
	}
	switch msg := msg.(type) {
	case oracle.MsgRequestData:
		app.handleMsgRequestData(txHash, msg, evMap, extra)
	case oracle.MsgReportData:
		app.handleMsgReportData(txHash, msg, evMap, extra)
	case oracle.MsgCreateDataSource:
		app.handleMsgCreateDataSource(txHash, msg, evMap, extra)
	case oracle.MsgCreateOracleScript:
		app.handleMsgCreateOracleScript(txHash, msg, evMap, extra)
  case staking.MsgCreateValidator:
		app.handleMsgCreateValidator(txHash, msg, evMap, extra)
	}
}
