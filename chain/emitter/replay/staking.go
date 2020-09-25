package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/bandprotocol/bandchain/chain/emitter/common"
)

func (app *App) emitSetHistoricalBondedTokenOnValidator(addr sdk.ValAddress) {
	val, _ := app.StakingKeeper.GetValidator(app.DeliverContext, addr)
	common.EmitSetHistoricalBondedTokenOnValidator(app, addr, val.Tokens.Uint64(), app.DeliverContext.BlockTime().UnixNano())
}

// handleMsgCreateValidator implements emitter handler for MsgCreateValidator.
func (app *App) handleMsgCreateValidator(
	txHash []byte, msg staking.MsgCreateValidator, evMap common.EvMap, extra common.JsDict,
) {
	app.emitSetHistoricalBondedTokenOnValidator(msg.ValidatorAddress)
}

// handleMsgDelegate implements emitter handler for MsgDelegate
func (app *App) handleMsgDelegate(
	txHash []byte, msg staking.MsgDelegate, evMap common.EvMap, extra common.JsDict,
) {
	app.emitSetHistoricalBondedTokenOnValidator(msg.ValidatorAddress)
}

// handleMsgUndelegate implements emitter handler for MsgUndelegate
func (app *App) handleMsgUndelegate(
	txHash []byte, msg staking.MsgUndelegate, evMap common.EvMap, extra common.JsDict,
) {
	app.emitSetHistoricalBondedTokenOnValidator(msg.ValidatorAddress)
}

// handleMsgBeginRedelegate implements emitter handler for MsgBeginRedelegate
func (app *App) handleMsgBeginRedelegate(
	txHash []byte, msg staking.MsgBeginRedelegate, evMap common.EvMap, extra common.JsDict,
) {
	app.emitSetHistoricalBondedTokenOnValidator(msg.ValidatorSrcAddress)
	app.emitSetHistoricalBondedTokenOnValidator(msg.ValidatorDstAddress)
}
