package replay

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	dist "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
)

// handleMsgWithdrawDelegatorReward implements emitter handler for MsgWithdrawDelegatorReward.
func (h *Hook) handleMsgWithdrawDelegatorReward(
	ctx sdk.Context, txHash []byte, msg dist.MsgWithdrawDelegatorReward, evMap common.EvMap, extra common.JsDict,
) {
	withdrawAddr := h.distrKeeper.GetDelegatorWithdrawAddr(ctx, msg.DelegatorAddress)
	h.AddAccountsInTx(withdrawAddr)
	val, _ := h.stakingKeeper.GetValidator(ctx, msg.ValidatorAddress)
	extra["moniker"] = val.Description.Moniker
	extra["identity"] = val.Description.Identity
	extra["reward_amount"] = evMap[dist.EventTypeWithdrawRewards+"."+sdk.AttributeKeyAmount][0]
}

// handleMsgSetWithdrawAddress implements emitter handler for MsgSetWithdrawAddress.
func (h *Hook) handleMsgSetWithdrawAddress(
	txHash []byte, msg dist.MsgSetWithdrawAddress, evMap common.EvMap, extra common.JsDict,
) {
	h.AddAccountsInTx(msg.WithdrawAddress)
	extra["delegator_address"] = msg.DelegatorAddress
	extra["withdraw_address"] = msg.WithdrawAddress
}

// handleMsgWithdrawValidatorCommission implements emitter handler for MsgWithdrawValidatorCommission.
func (h *Hook) handleMsgWithdrawValidatorCommission(
	ctx sdk.Context, txHash []byte, msg dist.MsgWithdrawValidatorCommission, evMap common.EvMap, extra common.JsDict,
) {
	withdrawAddr := h.distrKeeper.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(msg.ValidatorAddress))
	h.AddAccountsInTx(withdrawAddr)
	val, _ := h.stakingKeeper.GetValidator(ctx, msg.ValidatorAddress)
	extra["moniker"] = val.Description.Moniker
	extra["identity"] = val.Description.Identity
	extra["commission_amount"] = evMap[types.EventTypeWithdrawCommission+"."+sdk.AttributeKeyAmount][0]
}
