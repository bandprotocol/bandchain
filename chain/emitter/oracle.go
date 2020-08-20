package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func parseBytes(b []byte) []byte {
	if len(b) == 0 {
		return []byte{}
	}
	return b
}

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

func (app *App) emitHistoricalValidatorStatus(operatorAddress sdk.ValAddress) {
	status := app.OracleKeeper.GetValidatorStatus(app.DeliverContext, operatorAddress).IsActive
	app.Write("SET_HISTORICAL_VALIDATOR_STATUS", JsDict{
		"operator_address": operatorAddress,
		"status":           status,
		"timestamp":        app.DeliverContext.BlockTime().UnixNano(),
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
		"calldata":         parseBytes(msg.Calldata),
		"ask_count":        msg.AskCount,
		"min_count":        msg.MinCount,
		"sender":           msg.Sender.String(),
		"client_id":        msg.ClientID,
		"resolve_status":   types.ResolveStatus_Open,
	})
	for _, raw := range req.RawRequests {
		app.Write("NEW_RAW_REQUEST", JsDict{
			"request_id":     id,
			"external_id":    raw.ExternalID,
			"data_source_id": raw.DataSourceID,
			"calldata":       parseBytes(raw.Calldata),
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
			"data":        parseBytes(data.Data),
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
		"result":         parseBytes(result.ResponsePacketData.Result),
	})
}

// handleMsgAddReporter implements emitter handler for MsgAddReporter.
func (app *App) handleMsgAddReporter(
	txHash []byte, msg oracle.MsgAddReporter, evMap EvMap, extra JsDict,
) {
	val, _ := app.StakingKeeper.GetValidator(app.DeliverContext, msg.Validator)
	extra["validator_moniker"] = val.GetMoniker()
	app.AddAccountsInTx(msg.Reporter)
	app.Write("SET_REPORTER", JsDict{
		"reporter":  msg.Reporter,
		"validator": msg.Validator,
	})
}

// handleMsgRemoveReporter implements emitter handler for MsgRemoveReporter.
func (app *App) handleMsgRemoveReporter(
	txHash []byte, msg oracle.MsgRemoveReporter, evMap EvMap, extra JsDict,
) {
	val, _ := app.StakingKeeper.GetValidator(app.DeliverContext, msg.Validator)
	extra["validator_moniker"] = val.GetMoniker()
	app.AddAccountsInTx(msg.Reporter)
	app.Write("REMOVE_REPORTER", JsDict{
		"reporter":  msg.Reporter,
		"validator": msg.Validator,
	})
}

// handleMsgActivate implements emitter handler for handleMsgActivate.
func (app *App) handleMsgActivate(
	txHash []byte, msg oracle.MsgActivate, evMap EvMap, extra JsDict,
) {
	app.emitUpdateValidatorStatus(msg.Validator)
	app.emitHistoricalValidatorStatus(msg.Validator)
}

// handleEventDeactivate implements emitter handler for EventDeactivate.
func (app *App) handleEventDeactivate(evMap EvMap) {
	addr, _ := sdk.ValAddressFromBech32(evMap[types.EventTypeDeactivate+"."+types.AttributeKeyValidator][0])
	app.emitUpdateValidatorStatus(addr)
	app.emitHistoricalValidatorStatus(addr)
}
