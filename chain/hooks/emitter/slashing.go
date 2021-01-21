package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"

	"github.com/GeoDB-Limited/odincore/chain/hooks/common"
)

// handleEventSlash implements emitter handler for Slashing event.
func (h *Hook) handleEventSlash(ctx sdk.Context, event common.EvMap) {
	if raw, ok := event[slashing.EventTypeSlash+"."+slashing.AttributeKeyJailed]; ok && len(raw) == 1 {
		consAddress, _ := sdk.ConsAddressFromBech32(raw[0])
		validator, _ := h.stakingKeeper.GetValidatorByConsAddr(ctx, consAddress)
		h.Write("UPDATE_VALIDATOR", common.JsDict{
			"operator_address": validator.OperatorAddress.String(),
			"tokens":           validator.Tokens.Uint64(),
			"jailed":           validator.Jailed,
		})
	}
}

// handleMsgUnjail implements emitter handler for MsgUnjail.
func (h *Hook) handleMsgUnjail(
	ctx sdk.Context, msg slashing.MsgUnjail,
) {
	validator, _ := h.stakingKeeper.GetValidator(ctx, msg.ValidatorAddr)
	h.Write("UPDATE_VALIDATOR", common.JsDict{
		"operator_address": msg.ValidatorAddr.String(),
		"jailed":           validator.Jailed,
	})
}
