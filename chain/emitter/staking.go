package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func (app *App) SetValidator(addrs sdk.ValAddress, blockHeight int64) {
	val, found := app.StakingKeeper.GetValidator(app.DeliverContext, addrs)
	if !found {
		panic("expected validator, not found")
	}
	app.Write("SET_VALIDATOR", JsDict{
		"operator_address":      addrs.String(),
		"consensus_address":     val.ConsPubKey.Address().Bytes(),
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
		"jailed":                false,
		"delegator_shares":      val.DelegatorShares.String(),
		"bonded_height":         blockHeight,
		"current_reward":        "0",
		"current_ratio":         "0",
	})
}

// handleMsgCreateValidator implements emitter handler for MsgCreateValidator.
func (app *App) handleMsgCreateValidator(msg staking.MsgCreateValidator) {
	app.SetValidator(msg.ValidatorAddress, app.DeliverContext.BlockHeight())
}
