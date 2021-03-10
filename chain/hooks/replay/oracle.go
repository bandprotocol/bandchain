package replay

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func (h *Hook) emitSetDataSourceTxHash(id types.DataSourceID, txHash []byte) {
	h.Write("SET_DATA_SOURCE", common.JsDict{
		"id":      id,
		"tx_hash": txHash,
	})
}

func (h *Hook) emitSetOracleScriptTxHash(id types.OracleScriptID, txHash []byte) {
	h.Write("SET_ORACLE_SCRIPT", common.JsDict{
		"id":      id,
		"tx_hash": txHash,
	})
}

// handleMsgRequestData implements emitter handler for MsgRequestData.
func (h *Hook) handleMsgRequestData(
	ctx sdk.Context, txHash []byte, msg oracle.MsgRequestData, evMap common.EvMap, extra common.JsDict,
) {
	id := types.RequestID(common.Atoi(evMap[types.EventTypeRequest+"."+types.AttributeKeyID][0]))
	req := h.oracleKeeper.MustGetRequest(ctx, id)
	h.Write("UPDATE_REQUEST", common.JsDict{
		"id":      id,
		"tx_hash": txHash,
		"sender":  msg.Sender.String(),
	})
	common.EmitSetRequestCountPerDay(h, ctx.BlockTime().UnixNano())
	common.EmitUpdateOracleScriptRequest(h, msg.OracleScriptID)
	for _, raw := range req.RawRequests {
		common.EmitUpdateDataSourceRequest(h, raw.DataSourceID)
		common.EmitUpdateRelatedDsOs(h, raw.DataSourceID, msg.OracleScriptID)
	}
	os := h.oracleKeeper.MustGetOracleScript(ctx, msg.OracleScriptID)
	extra["id"] = id
	extra["name"] = os.Name
	extra["schema"] = os.Schema
}

func (h *Hook) handleMsgReportData(
	txHash []byte, msg oracle.MsgReportData, evMap common.EvMap, extra common.JsDict,
) {
	h.Write("UPDATE_REPORT", common.JsDict{
		"tx_hash":    txHash,
		"request_id": msg.RequestID,
		"validator":  msg.Validator.String(),
		"reporter":   msg.Reporter.String(),
	})
}

func (h *Hook) emitHistoricalValidatorStatus(ctx sdk.Context, operatorAddress sdk.ValAddress) {
	status := h.oracleKeeper.GetValidatorStatus(ctx, operatorAddress).IsActive
	common.EmitHistoricalValidatorStatus(h, operatorAddress, status, ctx.BlockTime().UnixNano())
}

// handleMsgCreateDataSource implements emitter handler for MsgCreateDataSource.
func (h *Hook) handleMsgCreateDataSource(
	txHash []byte, msg oracle.MsgCreateDataSource, evMap common.EvMap, extra common.JsDict,
) {
	id := types.DataSourceID(common.Atoi(evMap[types.EventTypeCreateDataSource+"."+types.AttributeKeyID][0]))
	h.emitSetDataSourceTxHash(id, txHash)
	extra["id"] = id
}

// handleMsgCreateOracleScript implements emitter handler for MsgCreateOracleScript.
func (h *Hook) handleMsgCreateOracleScript(
	txHash []byte, msg oracle.MsgCreateOracleScript, evMap common.EvMap, extra common.JsDict,
) {
	id := types.OracleScriptID(common.Atoi(evMap[types.EventTypeCreateOracleScript+"."+types.AttributeKeyID][0]))
	h.emitSetOracleScriptTxHash(id, txHash)
	extra["id"] = id
}

// handleMsgEditDataSource implements emitter handler for MsgEditDataSource.
func (h *Hook) handleMsgEditDataSource(
	txHash []byte, msg oracle.MsgEditDataSource, evMap common.EvMap, extra common.JsDict,
) {
	id := msg.DataSourceID
	h.emitSetDataSourceTxHash(id, txHash)
	extra["id"] = id
}

// handleMsgEditOracleScript implements emitter handler for MsgEditOracleScript.
func (h *Hook) handleMsgEditOracleScript(
	txHash []byte, msg oracle.MsgEditOracleScript, evMap common.EvMap, extra common.JsDict,
) {
	id := msg.OracleScriptID
	h.emitSetOracleScriptTxHash(id, txHash)
	extra["id"] = id
}

// handleMsgAddReporter implements emitter handler for MsgAddReporter.
func (h *Hook) handleMsgAddReporter(
	ctx sdk.Context, txHash []byte, msg oracle.MsgAddReporter, evMap common.EvMap, extra common.JsDict,
) {
	val, _ := h.stakingKeeper.GetValidator(ctx, msg.Validator)
	extra["validator_moniker"] = val.GetMoniker()
	h.AddAccountsInTx(msg.Reporter)
}

// handleMsgRemoveReporter implements emitter handler for MsgRemoveReporter.
func (h *Hook) handleMsgRemoveReporter(
	ctx sdk.Context, txHash []byte, msg oracle.MsgRemoveReporter, evMap common.EvMap, extra common.JsDict,
) {
	val, _ := h.stakingKeeper.GetValidator(ctx, msg.Validator)
	extra["validator_moniker"] = val.GetMoniker()
	h.AddAccountsInTx(msg.Reporter)
}

// handleMsgActivate implements emitter handler for handleMsgActivate.
func (h *Hook) handleMsgActivate(
	ctx sdk.Context, txHash []byte, msg oracle.MsgActivate, evMap common.EvMap, extra common.JsDict,
) {
	h.emitHistoricalValidatorStatus(ctx, msg.Validator)
}

// handleEventDeactivate implements emitter handler for EventDeactivate.
func (h *Hook) handleEventDeactivate(ctx sdk.Context, evMap common.EvMap) {
	addr, _ := sdk.ValAddressFromBech32(evMap[types.EventTypeDeactivate+"."+types.AttributeKeyValidator][0])
	h.emitHistoricalValidatorStatus(ctx, addr)
}
