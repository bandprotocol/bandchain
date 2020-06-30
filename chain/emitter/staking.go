package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func createValidator(msg staking.MsgCreateValidator, BlockHeight int64) JsDict {
	return JsDict{
		"operator_address":      msg.ValidatorAddress.String(),
		"consensus_address":     msg.PubKey.Address().Bytes(),
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

// handleMsgCreateValidator implements emitter handler for MsgCreateValidator.
func (app *App) handleMsgCreateValidator(
	txHash []byte, msg staking.MsgCreateValidator, evMap EvMap, extra JsDict,
) {
	app.Write("NEW_VALIDATOR", createValidator(msg, app.DeliverContext.BlockHeight()))
}
