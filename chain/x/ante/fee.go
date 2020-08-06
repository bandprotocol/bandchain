package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (mfd MempoolFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	_, ok := tx.(ante.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	isValidReportTx := true
	for _, msg := range tx.GetMsgs() {
		report, ok := msg.(oracle.MsgReportData)
		if !ok {
			isValidReportTx = false
			break
		}

		if !mfd.oracleKeeper.GetValidatorStatus(ctx, report.Validator).IsActive {
			isValidReportTx = false
			break
		}

		request := mfd.oracleKeeper.MustGetRequest(ctx, report.RequestID)
		if !keeper.ContainsVal(request.RequestedValidators, report.Validator) {
			isValidReportTx = false
			break
		}

		reports := mfd.oracleKeeper.GetReports(ctx, report.RequestID)
		vals := make([]sdk.ValAddress, len(reports))
		for idx, rp := range reports {
			vals[idx] = rp.Validator
		}
		if keeper.ContainsVal(vals, report.Validator) {
			isValidReportTx = false
			break
		}
	}
	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() && !simulate && !isValidReportTx {
		return mfd.mempool.AnteHandle(ctx, tx, simulate, next)
	}

	return next(ctx, tx, simulate)
}
