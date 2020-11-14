package replay

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
)

func (h *Hook) emitSetHistoricalBondedTokenOnValidator(ctx sdk.Context, addr sdk.ValAddress) {
	val, _ := h.stakingKeeper.GetValidator(ctx, addr)
	common.EmitSetHistoricalBondedTokenOnValidator(h, addr, val.Tokens.Uint64(), ctx.BlockTime().UnixNano())
}

// handleMsgCreateValidator implements emitter handler for MsgCreateValidator.
func (h *Hook) handleMsgCreateValidator(
	ctx sdk.Context, txHash []byte, msg staking.MsgCreateValidator, evMap common.EvMap, extra common.JsDict,
) {
	h.emitSetHistoricalBondedTokenOnValidator(ctx, msg.ValidatorAddress)
}

// handleMsgDelegate implements emitter handler for MsgDelegate
func (h *Hook) handleMsgDelegate(
	ctx sdk.Context, txHash []byte, msg staking.MsgDelegate, evMap common.EvMap, extra common.JsDict,
) {
	h.emitSetHistoricalBondedTokenOnValidator(ctx, msg.ValidatorAddress)
}

// handleMsgUndelegate implements emitter handler for MsgUndelegate
func (h *Hook) handleMsgUndelegate(
	ctx sdk.Context, txHash []byte, msg staking.MsgUndelegate, evMap common.EvMap, extra common.JsDict,
) {
	h.emitSetHistoricalBondedTokenOnValidator(ctx, msg.ValidatorAddress)
}

// handleMsgBeginRedelegate implements emitter handler for MsgBeginRedelegate
func (h *Hook) handleMsgBeginRedelegate(
	ctx sdk.Context, txHash []byte, msg staking.MsgBeginRedelegate, evMap common.EvMap, extra common.JsDict,
) {
	h.emitSetHistoricalBondedTokenOnValidator(ctx, msg.ValidatorSrcAddress)
	h.emitSetHistoricalBondedTokenOnValidator(ctx, msg.ValidatorDstAddress)
}
