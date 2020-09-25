package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/emitter/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func (app *App) emitSetDataSourceTxHash(id types.DataSourceID, txHash []byte) {
	app.Write("SET_DATA_SOURCE", common.JsDict{
		"id":      id,
		"tx_hash": txHash,
	})
}

func (app *App) emitSetOracleScriptTxHash(id types.OracleScriptID, txHash []byte) {
	app.Write("SET_ORACLE_SCRIPT", common.JsDict{
		"id":      id,
		"tx_hash": txHash,
	})
}

// handleMsgRequestData implements emitter handler for MsgRequestData.
func (app *App) handleMsgRequestData(
	txHash []byte, msg oracle.MsgRequestData, evMap common.EvMap, extra common.JsDict,
) {
	id := types.RequestID(common.Atoi(evMap[types.EventTypeRequest+"."+types.AttributeKeyID][0]))
	req := app.OracleKeeper.MustGetRequest(app.DeliverContext, id)
	app.Write("UPDATE_REQUEST", common.JsDict{
		"id":      id,
		"tx_hash": txHash,
		"sender":  msg.Sender.String(),
	})
	common.EmitSetRequestCountPerDay(app, app.DeliverContext.BlockTime().UnixNano())
	common.EmitUpdateOracleScriptRequest(app, msg.OracleScriptID)
	for _, raw := range req.RawRequests {
		common.EmitUpdateDataSourceRequest(app, raw.DataSourceID)
		common.EmitUpdateRelatedDsOs(app, raw.DataSourceID, msg.OracleScriptID)
	}
	os := app.OracleKeeper.MustGetOracleScript(app.DeliverContext, msg.OracleScriptID)
	extra["id"] = id
	extra["name"] = os.Name
	extra["schema"] = os.Schema
}

func (app *App) handleMsgReportData(
	txHash []byte, msg oracle.MsgReportData, evMap common.EvMap, extra common.JsDict,
) {
	app.Write("UPDATE_REPORT", common.JsDict{
		"tx_hash":    txHash,
		"request_id": msg.RequestID,
		"validator":  msg.Validator.String(),
		"reporter":   msg.Reporter.String(),
	})
}

func (app *App) emitHistoricalValidatorStatus(operatorAddress sdk.ValAddress) {
	status := app.OracleKeeper.GetValidatorStatus(app.DeliverContext, operatorAddress).IsActive
	common.EmitHistoricalValidatorStatus(app, operatorAddress, status, app.DeliverContext.BlockTime().UnixNano())
}

// handleMsgCreateDataSource implements emitter handler for MsgCreateDataSource.
func (app *App) handleMsgCreateDataSource(
	txHash []byte, msg oracle.MsgCreateDataSource, evMap common.EvMap, extra common.JsDict,
) {
	id := types.DataSourceID(common.Atoi(evMap[types.EventTypeCreateDataSource+"."+types.AttributeKeyID][0]))
	app.emitSetDataSourceTxHash(id, txHash)
	extra["id"] = id
}

// handleMsgCreateOracleScript implements emitter handler for MsgCreateOracleScript.
func (app *App) handleMsgCreateOracleScript(
	txHash []byte, msg oracle.MsgCreateOracleScript, evMap common.EvMap, extra common.JsDict,
) {
	id := types.OracleScriptID(common.Atoi(evMap[types.EventTypeCreateOracleScript+"."+types.AttributeKeyID][0]))
	app.emitSetOracleScriptTxHash(id, txHash)
	extra["id"] = id
}

// handleMsgEditDataSource implements emitter handler for MsgEditDataSource.
func (app *App) handleMsgEditDataSource(
	txHash []byte, msg oracle.MsgEditDataSource, evMap common.EvMap, extra common.JsDict,
) {
	id := msg.DataSourceID
	app.emitSetDataSourceTxHash(id, txHash)
	extra["id"] = id
}

// handleMsgEditOracleScript implements emitter handler for MsgEditOracleScript.
func (app *App) handleMsgEditOracleScript(
	txHash []byte, msg oracle.MsgEditOracleScript, evMap common.EvMap, extra common.JsDict,
) {
	id := msg.OracleScriptID
	app.emitSetOracleScriptTxHash(id, txHash)
	extra["id"] = id
}

// handleMsgAddReporter implements emitter handler for MsgAddReporter.
func (app *App) handleMsgAddReporter(
	txHash []byte, msg oracle.MsgAddReporter, evMap common.EvMap, extra common.JsDict,
) {
	val, _ := app.StakingKeeper.GetValidator(app.DeliverContext, msg.Validator)
	extra["validator_moniker"] = val.GetMoniker()
	app.AddAccountsInTx(msg.Reporter)
}

// handleMsgRemoveReporter implements emitter handler for MsgRemoveReporter.
func (app *App) handleMsgRemoveReporter(
	txHash []byte, msg oracle.MsgRemoveReporter, evMap common.EvMap, extra common.JsDict,
) {
	val, _ := app.StakingKeeper.GetValidator(app.DeliverContext, msg.Validator)
	extra["validator_moniker"] = val.GetMoniker()
	app.AddAccountsInTx(msg.Reporter)
}

// handleMsgActivate implements emitter handler for handleMsgActivate.
func (app *App) handleMsgActivate(
	txHash []byte, msg oracle.MsgActivate, evMap common.EvMap, extra common.JsDict,
) {
	app.emitHistoricalValidatorStatus(msg.Validator)
}

// handleEventDeactivate implements emitter handler for EventDeactivate.
func (app *App) handleEventDeactivate(evMap common.EvMap) {
	addr, _ := sdk.ValAddressFromBech32(evMap[types.EventTypeDeactivate+"."+types.AttributeKeyValidator][0])
	app.emitHistoricalValidatorStatus(addr)
}
