package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

// handleEventSlash implements emitter handler for Slashing event.
func (app *App) handleEventSlash(event EvMap) {
	if raw, ok := event[slashing.EventTypeSlash+"."+slashing.AttributeKeyJailed]; ok && len(raw) == 1 {
		consAddress, _ := sdk.ConsAddressFromBech32(raw[0])
		validator, _ := app.StakingKeeper.GetValidatorByConsAddr(app.DeliverContext, consAddress)
		app.Write("UPDATE_VALIDATOR", JsDict{
			"operator_address": validator.OperatorAddress.String(),
			"tokens":           validator.Tokens.Uint64(),
			"jailed":           validator.Jailed,
		})
	}
}

// handleMsgUnjail implements emitter handler for MsgUnjail.
func (app *App) handleMsgUnjail(
	txHash []byte, msg slashing.MsgUnjail, evMap EvMap, extra JsDict,
) {
	validator, _ := app.StakingKeeper.GetValidator(app.DeliverContext, msg.ValidatorAddr)
	app.Write("UPDATE_VALIDATOR", JsDict{
		"operator_address": msg.ValidatorAddr.String(),
		"jailed":           validator.Jailed,
	})
	extra["moniker"] = validator.Description.Moniker
	extra["identity"] = validator.Description.Identity
}
