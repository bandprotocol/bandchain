package yoda

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Constant used to estimate gas price of reports transaction.
const (
	// cosmos
	baseFixedGas        = uint64(37764)
	baseTransactionSize = uint64(200)
	txCostPerByte       = uint64(5) // Using DefaultTxSizeCostPerByte of BandChain

	readingBaseCost = uint64(1000)
	writingBaseCost = uint64(2000)

	readingCostPerByte = uint64(3)
	writingCostPerByte = uint64(30)

	payingFeeCost = uint64(15200)

	// band
	baseReportCost        = uint64(4006)
	countingPerReportCost = uint64(1479)
	addingPendingCost     = uint64(3400)

	baseRequestSize = uint64(32)
	addressSize     = uint64(20)

	baseRawRequestSize = uint64(16)
)

func estimateTxSize(msgs []sdk.Msg) uint64 {
	// base tx + reports
	size := baseTransactionSize

	for _, msg := range msgs {
		msg, ok := msg.(types.MsgReportData)
		if !ok {
			panic("Don't support non-report data message")
		}

		ser := cdc.MustMarshalBinaryBare(msg)
		size += uint64(len(ser))
	}

	return size
}

func estimateStoringReportCost(msg sdk.Msg) uint64 {
	cost := writingBaseCost
	cost += uint64(len(cdc.MustMarshalBinaryBare(msg))) * writingCostPerByte

	return cost
}

func estimateReadingRequestCost(f FeeEstimationData) uint64 {
	cost := readingBaseCost

	size := baseRequestSize
	size += uint64(len(f.callData))
	size += uint64(f.validators) * addressSize
	size += uint64(len(f.clientId))

	for _, r := range f.rawRequests {
		size += baseRawRequestSize + uint64(len(r.calldata))
	}

	cost += size * readingCostPerByte

	return cost
}

func estimateReportHandleCost(msg sdk.Msg, f FeeEstimationData) uint64 {
	cost := baseReportCost

	// read request twice
	cost += 2 * estimateReadingRequestCost(f)

	// write report once
	cost += estimateStoringReportCost(msg)

	// count report
	cost += countingPerReportCost * uint64(f.minCount+1)

	// add pending
	cost += addingPendingCost

	return cost
}

func estimateGas(c *Context, msgs []sdk.Msg, feeEstimations []FeeEstimationData) uint64 {
	gas := baseFixedGas

	txSize := estimateTxSize(msgs)
	gas += txCostPerByte * txSize

	// process paying fee
	if len(c.gasPrices) > 0 {
		gas += payingFeeCost
	}

	for i := range msgs {
		gas += estimateReportHandleCost(msgs[i], feeEstimations[i])
	}

	return gas
}
