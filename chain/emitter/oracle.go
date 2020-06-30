package emitter

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func createValidator(msg staking.MsgCreateValidator, BlockHeight int64) JsDict {
	return JsDict{
		"operator_address":      msg.ValidatorAddress.String(),
		"consensus_address":     msg.PubKey.Address().String(),
		"consensus_pubkey":      sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, msg.PubKey),
		"moniker":               msg.Description.Moniker,
		"identity":              msg.Description.Identity,
		"website":               msg.Description.Website,
		"details":               msg.Description.Details,
		"commission_rate":       msg.Commission.Rate.String(),
		"commission_max_rate":   msg.Commission.MaxRate.String(),
		"commission_max_change": msg.Commission.MaxChangeRate.String(),
		"min_self_delegation":   msg.MinSelfDelegation.String(),
		"tokens":                msg.Value.Amount.Uint64(),
		"jailed":                false,
		"delegator_shares":      msg.Value.Amount.String(),
		"bonded_height":         BlockHeight,
		"current_reward":        "0",
		"current_ratio":         "0",
	}
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

// handleMsgCreateValidator implements emitter handler for MsgCreateValidator.
func (app *App) handleMsgCreateValidator(
	txHash []byte, msg staking.MsgCreateValidator, evMap EvMap, extra JsDict,
) {
	app.Write("NEW_VALIDATOR", createValidator(msg, app.DeliverContext.BlockHeight()))
}
