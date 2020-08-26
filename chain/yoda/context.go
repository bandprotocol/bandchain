package yoda

import (
	"sync/atomic"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/bandprotocol/bandchain/chain/yoda/executor"
)

type ReportMsgWithKey struct {
	msg         types.MsgReportData
	execVersion []string
	keyIndex    int64
}

type Context struct {
	client           rpcclient.Client
	validator        sdk.ValAddress
	gasPrices        sdk.DecCoins
	keys             []keys.Info
	executor         executor.Executor
	fileCache        filecache.Cache
	broadcastTimeout time.Duration
	maxTry           uint64
	rpcPollInterval  time.Duration
	maxReport        uint64

	pendingMsgs chan ReportMsgWithKey
	freeKeys    chan int64
	keyIndex    int64
}

func (c *Context) getKeyIndex() int64 {
	keyIndex := atomic.AddInt64(&c.keyIndex, 1) % int64(len(c.keys))
	return keyIndex
}
