package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func (app *App) emitSetValidator(addrs sdk.ValAddress) {
	val, _ := app.StakingKeeper.GetValidator(app.DeliverContext, addrs)
	app.Write("SET_VALIDATOR", JsDict{
		"operator_address":      addrs.String(),
		"consensus_address":     sdk.ConsAddress(val.ConsPubKey.Address()).String(),
		"consensus_pubkey":      sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, val.ConsPubKey),
		"moniker":               val.Description.Moniker,
		"identity":              val.Description.Identity,
		"website":               val.Description.Website,
		"details":               val.Description.Details,
		"commission_rate":       val.Commission.Rate.String(),
		"commission_max_rate":   val.Commission.MaxRate.String(),
		"commission_max_change": val.Commission.MaxChangeRate.String(),
		"min_self_delegation":   val.MinSelfDelegation.String(),
		"tokens":                val.Tokens.Uint64(),
		"jailed":                val.Jailed,
		"delegator_shares":      val.DelegatorShares.String(),
	})
}

// handleMsgCreateValidator implements emitter handler for MsgCreateValidator.
func (app *App) handleMsgCreateValidator(msg staking.MsgCreateValidator) {
	app.emitSetValidator(msg.ValidatorAddress)
}
