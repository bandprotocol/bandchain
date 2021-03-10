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
	val, _ := h.stakingKeeper.GetValidator(ctx, msg.ValidatorAddress)
	h.emitSetHistoricalBondedTokenOnValidator(ctx, msg.ValidatorAddress)
	extra["moniker"] = val.Description.Moniker
	extra["identity"] = val.Description.Identity
}

// handleMsgEditValidator implements emitter handler for MsgEditValidator.
func (h *Hook) handleMsgEditValidator(
	ctx sdk.Context, txHash []byte, msg staking.MsgEditValidator, evMap common.EvMap, extra common.JsDict,
) {
	val, _ := h.stakingKeeper.GetValidator(ctx, msg.ValidatorAddress)
	extra["moniker"] = val.Description.Moniker
	extra["identity"] = val.Description.Identity
}

// handleMsgDelegate implements emitter handler for MsgDelegate
func (h *Hook) handleMsgDelegate(
	ctx sdk.Context, txHash []byte, msg staking.MsgDelegate, evMap common.EvMap, extra common.JsDict,
) {
	val, _ := h.stakingKeeper.GetValidator(ctx, msg.ValidatorAddress)
	h.emitSetHistoricalBondedTokenOnValidator(ctx, msg.ValidatorAddress)
	extra["moniker"] = val.Description.Moniker
	extra["identity"] = val.Description.Identity
}

// handleMsgUndelegate implements emitter handler for MsgUndelegate
func (h *Hook) handleMsgUndelegate(
	ctx sdk.Context, txHash []byte, msg staking.MsgUndelegate, evMap common.EvMap, extra common.JsDict,
) {
	val, _ := h.stakingKeeper.GetValidator(ctx, msg.ValidatorAddress)
	h.emitSetHistoricalBondedTokenOnValidator(ctx, msg.ValidatorAddress)
	extra["moniker"] = val.Description.Moniker
	extra["identity"] = val.Description.Identity
}

// handleMsgBeginRedelegate implements emitter handler for MsgBeginRedelegate
func (h *Hook) handleMsgBeginRedelegate(
	ctx sdk.Context, txHash []byte, msg staking.MsgBeginRedelegate, evMap common.EvMap, extra common.JsDict,
) {
	valSrc, _ := h.stakingKeeper.GetValidator(ctx, msg.ValidatorSrcAddress)
	valDst, _ := h.stakingKeeper.GetValidator(ctx, msg.ValidatorDstAddress)
	h.emitSetHistoricalBondedTokenOnValidator(ctx, msg.ValidatorSrcAddress)
	h.emitSetHistoricalBondedTokenOnValidator(ctx, msg.ValidatorDstAddress)
	extra["val_src_moniker"] = valSrc.Description.Moniker
	extra["val_src_identity"] = valSrc.Description.Identity
	extra["val_dst_moniker"] = valDst.Description.Moniker
	extra["val_dst_identity"] = valDst.Description.Identity
}
