package emitter

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func (app *App) emitSetDataSource(id types.DataSourceID, ds types.DataSource, txHash []byte) {
	app.Write("SET_DATA_SOURCE", JsDict{
		"id":          id,
		"name":        ds.Name,
		"description": ds.Description,
		"owner":       ds.Owner.String(),
		"executable":  app.OracleKeeper.GetFile(ds.Filename),
		"tx_hash":     txHash,
	})
}

func (app *App) emitSetOracleScript(id types.OracleScriptID, os types.OracleScript, txHash []byte) {
	app.Write("SET_ORACLE_SCRIPT", JsDict{
		"id":              id,
		"name":            os.Name,
		"description":     os.Description,
		"owner":           os.Owner.String(),
		"schema":          os.Schema,
		"codehash":        os.Filename,
		"source_code_url": os.SourceCodeURL,
		"tx_hash":         txHash,
	})
}

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
		"resolve_status":   types.ResolveStatus_Open,
	})
	es := evMap[types.EventTypeRawRequest+"."+types.AttributeKeyExternalID]
	ds := evMap[types.EventTypeRawRequest+"."+types.AttributeKeyDataSourceID]
	cs := evMap[types.EventTypeRawRequest+"."+types.AttributeKeyCalldata]
	for idx := range es {
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

// handleMsgCreateDataSource implements emitter handler for MsgCreateDataSource.
func (app *App) handleMsgCreateDataSource(
	txHash []byte, msg oracle.MsgCreateDataSource, evMap EvMap, extra JsDict,
) {
	id := types.DataSourceID(atoi(evMap[types.EventTypeCreateDataSource+"."+types.AttributeKeyID][0]))
	ds := app.BandApp.OracleKeeper.MustGetDataSource(app.DeliverContext, id)
	app.emitSetDataSource(id, ds, txHash)
	extra["id"] = id
}

// handleMsgCreateOracleScript implements emitter handler for MsgCreateOracleScript.
func (app *App) handleMsgCreateOracleScript(
	txHash []byte, msg oracle.MsgCreateOracleScript, evMap EvMap, extra JsDict,
) {
	id := types.OracleScriptID(atoi(evMap[types.EventTypeCreateOracleScript+"."+types.AttributeKeyID][0]))
	os := app.BandApp.OracleKeeper.MustGetOracleScript(app.DeliverContext, id)
	app.emitSetOracleScript(id, os, txHash)
	extra["id"] = id
}

// handleMsgEditDataSource implements emitter handler for MsgEditDataSource.
func (app *App) handleMsgEditDataSource(
	txHash []byte, msg oracle.MsgEditDataSource, evMap EvMap, extra JsDict,
) {
	id := msg.DataSourceID
	ds := app.BandApp.OracleKeeper.MustGetDataSource(app.DeliverContext, id)
	app.emitSetDataSource(id, ds, txHash)
}

// handleMsgEditOracleScript implements emitter handler for MsgEditOracleScript.
func (app *App) handleMsgEditOracleScript(
	txHash []byte, msg oracle.MsgEditOracleScript, evMap EvMap, extra JsDict,
) {
	id := msg.OracleScriptID
	os := app.BandApp.OracleKeeper.MustGetOracleScript(app.DeliverContext, id)
	app.emitSetOracleScript(id, os, txHash)
}

// handleEventRequestExecute implements emitter handler for EventRequestExecute.
func (app *App) handleEventRequestExecute(evMap EvMap) {
	id := types.RequestID(atoi(evMap[types.EventTypeResolve+"."+types.AttributeKeyID][0]))
	result := app.OracleKeeper.MustGetResult(app.DeliverContext, id)
	app.Write("UPDATE_REQUEST", JsDict{
		"id":             id,
		"request_time":   result.ResponsePacketData.RequestTime,
		"resolve_time":   result.ResponsePacketData.ResolveTime,
		"resolve_status": result.ResponsePacketData.ResolveStatus,
		"result":         result.ResponsePacketData.Result,
	})
}
