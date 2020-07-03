package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

// handleEventSlash implements emitter handler for Slashing event.
func (app *App) handleEventSlash(event EvMap) {
	if raw, ok := event[slashing.AttributeKeyJailed]; ok && len(raw) == 1 {
		consAddress, err := sdk.ConsAddressFromBech32(raw[0])
		if err != nil {
			panic(err)
		}
		validator, found := app.StakingKeeper.GetValidatorByConsAddr(app.DeliverContext, consAddress)
		if !found {
			panic("handleEventSlash: validator not found")
		}
		app.Write("UPDATE_VALIDATOR", JsDict{
			"operator_address": validator.OperatorAddress,
			"tokens":           validator.Tokens.Uint64(),
			"jailed":           validator.Jailed,
		})
	}
}

// handleMsgUnjail implements emitter handler for MsgUnjail.
func (app *App) handleMsgUnjail(msg slashing.MsgUnjail) {
	validator, found := app.StakingKeeper.GetValidator(app.DeliverContext, msg.ValidatorAddr)
	if !found {
		panic("handleMsgUnjail: validator not found")
	}
	app.Write("UPDATE_VALIDATOR", JsDict{
		"operator_address": msg.ValidatorAddr,
		"jailed":           validator.Jailed,
	})

}
