package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
)

// MempoolFeeDecorator will check if the transaction's fee is at least as large
// as the local validator's minimum gasFee (defined in validator config).
// If fee is too low, decorator returns error and tx is rejected from mempool.
// Note this only applies when ctx.CheckTx = true
// If fee is high enough or not CheckTx, then call next AnteHandler
// CONTRACT: Tx must implement FeeTx to use MempoolFeeDecorator
type MempoolFeeDecorator struct {
	oracleKeeper oracle.Keeper
	mempool      ante.MempoolFeeDecorator
}

func NewMempoolFeeDecorator(ok oracle.Keeper) MempoolFeeDecorator {
	return MempoolFeeDecorator{oracleKeeper: ok, mempool: ante.NewMempoolFeeDecorator()}
}

func (mfd MempoolFeeDecorator) checkValidReportMsg(ctx sdk.Context, msg sdk.Msg) bool {
	rep, ok := msg.(oracle.MsgReportData)
	if !ok {
		return false
	}
	if !mfd.oracleKeeper.IsReporter(ctx, rep.Validator, rep.Reporter) {
		return false
	}
	if rep.RequestID <= mfd.oracleKeeper.GetRequestLastExpired(ctx) {
		return false
	}

	req, err := mfd.oracleKeeper.GetRequest(ctx, rep.RequestID)
	if err != nil {
		return false
	}
	if !keeper.ContainsVal(req.RequestedValidators, rep.Validator) {
		return false
	}
	if len(rep.RawReports) != len(req.RawRequests) {
		return false
	}
	for _, report := range rep.RawReports {
		if !keeper.ContainsEID(req.RawRequests, report.ExternalID) {
			return false
		}
	}
	return true
}

func (mfd MempoolFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isValidReportTx := true
	for _, msg := range tx.GetMsgs() {
		if isValidReportTx = mfd.checkValidReportMsg(ctx, msg); !isValidReportTx {
			break
		}
	}
	newCtx = ctx
	if ctx.IsCheckTx() && !simulate && isValidReportTx {
		newCtx = newCtx.WithMinGasPrices(sdk.DecCoins{})
	}
	return mfd.mempool.AnteHandle(newCtx, tx, simulate, next)
}
