package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
)

func checkValidReportMsg(ctx sdk.Context, oracleKeeper oracle.Keeper, msg sdk.Msg) bool {
	rep, ok := msg.(oracle.MsgReportData)
	if !ok {
		return false
	}
	if !oracleKeeper.IsReporter(ctx, rep.Validator, rep.Reporter) {
		return false
	}
	if rep.RequestID <= oracleKeeper.GetRequestLastExpired(ctx) {
		return false
	}

	req, err := oracleKeeper.GetRequest(ctx, rep.RequestID)
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

func BandWrapAnteHandler(ante sdk.AnteHandler, oracleKeeper oracle.Keeper, withFeeReportTx bool) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
		newCtx = ctx
		if !withFeeReportTx {
			isValidReportTx := true
			for _, msg := range tx.GetMsgs() {
				isValidReportTx = checkValidReportMsg(ctx, oracleKeeper, msg)
				if !isValidReportTx {
					break
				}
			}
			if ctx.IsCheckTx() && !simulate && isValidReportTx {
				newCtx = newCtx.WithMinGasPrices(sdk.DecCoins{})
			}
		}
		return ante(newCtx, tx, simulate)
	}
}
