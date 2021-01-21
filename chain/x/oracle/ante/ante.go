package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	lru "github.com/hashicorp/golang-lru"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/keeper"
)

var (
	repTxCount       *lru.Cache
	nextRepOnlyBlock int64
)

func init() {
	var err error
	repTxCount, err = lru.New(20000)
	if err != nil {
		panic(err)
	}
}

func checkValidReportMsg(ctx sdk.Context, oracleKeeper oracle.Keeper, rep oracle.MsgReportData) bool {
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

// NewFeelessReportsAnteHandler returns a new ante handler that waives minimum gas price
// requirement if the incoming tx is a valid report transaction.
func NewFeelessReportsAnteHandler(ante sdk.AnteHandler, oracleKeeper oracle.Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
		if ctx.IsCheckTx() && !simulate {
			// TODO: Move this out of "FeelessReports" ante handler.
			isRepOnlyBlock := ctx.BlockHeight() == nextRepOnlyBlock
			isValidReportTx := true
			for _, msg := range tx.GetMsgs() {
				rep, ok := msg.(oracle.MsgReportData)
				if !ok || !checkValidReportMsg(ctx, oracleKeeper, rep) {
					isValidReportTx = false
					break
				}
				if !isRepOnlyBlock {
					key := fmt.Sprintf("%s:%d", rep.Validator.String(), rep.RequestID)
					val, ok := repTxCount.Get(key)
					nextVal := 1
					if ok {
						nextVal = val.(int) + 1
					}
					repTxCount.Add(key, nextVal)
					if nextVal > 20 {
						nextRepOnlyBlock = ctx.BlockHeight() + 1
					}
				}
			}
			if isRepOnlyBlock && !isValidReportTx {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Block reserved for report txs")
			}
			if isValidReportTx {
				minGas := ctx.MinGasPrices()
				newCtx, err := ante(ctx.WithMinGasPrices(sdk.DecCoins{}), tx, simulate)
				// Set minimum gas price context and return context to caller.
				return newCtx.WithMinGasPrices(minGas), err
			}
		}
		return ante(ctx, tx, simulate)
	}
}
