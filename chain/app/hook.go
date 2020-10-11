package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Hook is an interface of hook that can process along with abci application
type Hook interface {
	AfterInitChain(ctx sdk.Context, req abci.RequestInitChain, res abci.ResponseInitChain)
	AfterBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock)
	AfterDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx, res abci.ResponseDeliverTx)
	AfterEndBlock(ctx sdk.Context, req abci.RequestEndBlock, res abci.ResponseEndBlock)
	ApplyQuery(req abci.RequestQuery) (res abci.ResponseQuery, stop bool)
	BeforeCommit()
}
