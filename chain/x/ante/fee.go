package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	types "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
)

var (
	_ FeeTx = (*types.StdTx)(nil) // assert StdTx implements FeeTx
)

// FeeTx defines the interface to be implemented by Tx to use the FeeDecorators
type FeeTx interface {
	sdk.Tx
	GetGas() uint64
	GetFee() sdk.Coins
	FeePayer() sdk.AccAddress
}

// MempoolFeeDecorator will check if the transaction's fee is at least as large
// as the local validator's minimum gasFee (defined in validator config).
// If fee is too low, decorator returns error and tx is rejected from mempool.
// Note this only applies when ctx.CheckTx = true
// If fee is high enough or not CheckTx, then call next AnteHandler
// CONTRACT: Tx must implement FeeTx to use MempoolFeeDecorator
type MempoolFeeDecorator struct {
	oracleKeeper oracle.Keeper
}

func NewMempoolFeeDecorator(ok oracle.Keeper) MempoolFeeDecorator {
	return MempoolFeeDecorator{oracleKeeper: ok}
}

func (mfd MempoolFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}
	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

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
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdk.NewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !feeCoins.IsAnyGTE(requiredFees) {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}

	return next(ctx, tx, simulate)
}
