package emitter

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
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
	}
}

func (app *App) handleEndBlock(event abci.Event) {
	kvMap := make(map[string]string)
	for _, kv := range event.Attributes {
		kvMap[string(kv.Key)] = string(kv.Value)
	}
	switch event.Type {
	case types.EventTypeRequestExecute:
		id, err := strconv.ParseInt(kvMap[types.AttributeKeyRequestID], 10, 64)
		if err != nil {
			panic(err)
		}

		numResolveStatus, err := strconv.ParseInt(kvMap[types.AttributeKeyResolveStatus], 10, 8)
		if err != nil {
			panic(err)
		}

		resolveStatus := types.ResolveStatus(numResolveStatus)

		requestTime, err := strconv.ParseInt(kvMap[types.AttributeKeyRequestTime], 10, 64)
		if err != nil {
			panic(err)
		}

		resolveTime, err := strconv.ParseInt(kvMap[types.AttributeKeyResolveTime], 10, 64)
		if err != nil {
			panic(err)
		}

		dict := JsDict{
			"id":             id,
			"request_time":   requestTime,
			"resolve_time":   resolveTime,
			"resolve_status": resolveStatus,
		}

		app.Write("UPDATE_REQUEST", dict)
		if resolveStatus == types.ResolveStatus_Success {
			result := []byte(kvMap[types.AttributeKeyResult])
			dict["result"] = result
		}
	default:
		break
	}
}
