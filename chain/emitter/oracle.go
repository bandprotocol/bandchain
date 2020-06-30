package emitter

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// handleMsgRequestData implements emitter handler for MsgRequestData.
func (app *App) handleMsgRequestData(
	txHash []byte, msg oracle.MsgRequestData, evMap EvMap, extra JsDict,
) {
	id := types.RequestID(atoi(evMap[types.EventTypeRequest+"."+types.AttributeKeyID][0]))
	req := app.OracleKeeper.MustGetRequest(app.DeliverContext, id)
	app.Write("NEW_REQUEST", JsDict{
		"id":               id,
		"tx_hash":          txHash,
		"oracle_script_id": msg.OracleScriptID,
		"calldata":         msg.Calldata,
		"ask_count":        msg.AskCount,
		"min_count":        msg.MinCount,
		"sender":           msg.Sender.String(),
		"client_id":        msg.ClientID,
	})
	es := evMap[types.EventTypeRawRequest+"."+types.AttributeKeyExternalID]
	ds := evMap[types.EventTypeRawRequest+"."+types.AttributeKeyDataSourceID]
	cs := evMap[types.EventTypeRawRequest+"."+types.AttributeKeyCalldata]
	for idx, _ := range es {
		app.Write("NEW_RAW_REQUEST", JsDict{
			"request_id":     id,
			"external_id":    es[idx],
			"data_source_id": ds[idx],
			"calldata":       []byte(cs[idx]),
		})
	}
	for _, val := range req.RequestedValidators {
		app.Write("NEW_VAL_REQUEST", JsDict{
			"request_id": id,
			"validator":  val.String(),
		})
	}
	os := app.OracleKeeper.MustGetOracleScript(app.DeliverContext, msg.OracleScriptID)
	extra["id"] = id
	extra["name"] = os.Name
	extra["schema"] = os.Schema
}

// handleMsgReportData implements emitter handler for MsgReportData.
func (app *App) handleMsgReportData(
	txHash []byte, msg oracle.MsgReportData, evMap EvMap, extra JsDict,
) {
	app.Write("NEW_REPORT", JsDict{
		"tx_hash":    txHash,
		"request_id": msg.RequestID,
		"validator":  msg.Validator.String(),
		"reporter":   msg.Reporter.String(),
	})
	for _, data := range msg.RawReports {
		app.Write("NEW_RAW_REPORT", JsDict{
			"request_id":  msg.RequestID,
			"validator":   msg.Validator.String(),
			"external_id": data.ExternalID,
			"data":        data.Data,
			"exit_code":   data.ExitCode,
		})
	}
}
