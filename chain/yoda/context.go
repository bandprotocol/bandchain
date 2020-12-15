package yoda

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/bandprotocol/bandchain/chain/yoda/executor"
)

type FeeEstimationData struct {
	askCount    int64
	minCount    int64
	callData    []byte
	rawRequests []rawRequest
	clientID    string
}

type ReportMsgWithKey struct {
	msg               types.MsgReportData
	execVersion       []string
	keyIndex          int64
	feeEstimationData FeeEstimationData
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

	pendingMsgs        chan ReportMsgWithKey
	freeKeys           chan int64
	keyRoundRobinIndex int64 // Must use in conjunction with sync/atomic

	dataSourceCache *sync.Map
	pendingRequests map[types.RequestID]bool

	metricsEnabled bool
	handlingGauge  int64
	pendingGauge   int64
	submittedCount int64
}

func (c *Context) nextKeyIndex() int64 {
	keyIndex := atomic.AddInt64(&c.keyRoundRobinIndex, 1) % int64(len(c.keys))
	return keyIndex
}

func (c *Context) updateHandlingGauge(amount int64) {
	if c.metricsEnabled {
		atomic.AddInt64(&c.handlingGauge, amount)
	}
}

func (c *Context) updatePendingGauge(amount int64) {
	if c.metricsEnabled {
		atomic.AddInt64(&c.pendingGauge, amount)
	}
}

func (c *Context) updateSubmittedCount(amount int64) {
	if c.metricsEnabled {
		atomic.AddInt64(&c.submittedCount, amount)
	}
}
