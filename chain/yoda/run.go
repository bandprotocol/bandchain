package yoda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	httpclient "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/bandprotocol/bandchain/chain/yoda/executor"
)

const (
	// TODO: We can subscribe only for txs that contain request messages
	TxQuery = "tm.event = 'Tx'"
	// EventChannelCapacity is a buffer size of channel between node and this program
	EventChannelCapacity = 2000
)

func runImpl(c *Context, l *Logger) error {
	l.Info(":rocket: Starting WebSocket subscriber")
	err := c.client.Start()
	if err != nil {
		return err
	}

	ctx, cxl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cxl()

	l.Info(":ear: Subscribing to events with query: %s...", TxQuery)
	eventChan, err := c.client.Subscribe(ctx, "", TxQuery, EventChannelCapacity)
	if err != nil {
		return err
	}

	availiableKeys := make([]bool, len(c.keys))
	waitingMsgs := make([][]ReportMsgWithKey, len(c.keys))
	for i := range availiableKeys {
		availiableKeys[i] = true
		waitingMsgs[i] = []ReportMsgWithKey{}
	}

	// Get pending requests and handle them
	rawPendingRequests, err := c.client.ABCIQueryWithOptions(fmt.Sprintf("custom/%s/%s/%s", types.StoreKey, types.QueryPendingRequests, c.validator.String()), nil, rpcclient.ABCIQueryOptions{})
	if err != nil {
		return err
	}

	var result types.QueryResult
	if err := json.Unmarshal(rawPendingRequests.Response.GetValue(), &result); err != nil {
		return err
	}

	var pendingRequests []types.RequestID
	cdc.MustUnmarshalJSON(result.Result, &pendingRequests)

	for _, id := range pendingRequests {
		go handlePendingRequest(c, l.With("rid", id), id)
	}

	for {
		select {
		case ev := <-eventChan:
			go handleTransaction(c, l, ev.Data.(tmtypes.EventDataTx).TxResult)
		case keyIndex := <-c.freeKeys:
			if len(waitingMsgs[keyIndex]) != 0 {
				if uint64(len(waitingMsgs[keyIndex])) > c.maxReport {
					go SubmitReport(c, l, keyIndex, waitingMsgs[keyIndex][:c.maxReport])
					waitingMsgs[keyIndex] = waitingMsgs[keyIndex][c.maxReport:]
				} else {
					go SubmitReport(c, l, keyIndex, waitingMsgs[keyIndex])
					waitingMsgs[keyIndex] = []ReportMsgWithKey{}
				}
			} else {
				availiableKeys[keyIndex] = true
			}
		case pm := <-c.pendingMsgs:
			if availiableKeys[pm.keyIndex] {
				availiableKeys[pm.keyIndex] = false
				go SubmitReport(c, l, pm.keyIndex, []ReportMsgWithKey{pm})
			} else {
				waitingMsgs[pm.keyIndex] = append(waitingMsgs[pm.keyIndex], pm)
			}
		}
	}
}

func runCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"r"},
		Short:   "Run the oracle process",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.ChainID == "" {
				return errors.New("Chain ID must not be empty")
			}
			keys, err := keybase.List()
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				return errors.New("No key available")
			}
			c.keys = keys
			c.validator, err = sdk.ValAddressFromBech32(cfg.Validator)
			if err != nil {
				return err
			}
			err = sdk.VerifyAddressFormat(c.validator)
			if err != nil {
				return err
			}
			c.gasPrices, err = sdk.ParseDecCoins(cfg.GasPrices)
			if err != nil {
				return err
			}
			allowLevel, err := log.AllowLevel(cfg.LogLevel)
			if err != nil {
				return err
			}
			l := NewLogger(allowLevel)
			c.executor, err = executor.NewExecutor(cfg.Executor)
			if err != nil {
				return err
			}
			l.Info(":star: Creating HTTP client with node URI: %s", cfg.NodeURI)
			c.client, err = httpclient.New(cfg.NodeURI, "/websocket")
			if err != nil {
				return err
			}
			c.fileCache = filecache.New(filepath.Join(viper.GetString(flags.FlagHome), "files"))
			c.broadcastTimeout, err = time.ParseDuration(cfg.BroadcastTimeout)
			if err != nil {
				return err
			}
			c.maxTry = cfg.MaxTry
			c.maxReport = cfg.MaxReport
			c.rpcPollInterval, err = time.ParseDuration(cfg.RPCPollInterval)
			if err != nil {
				return err
			}
			c.pendingMsgs = make(chan ReportMsgWithKey)
			c.freeKeys = make(chan int64, len(keys))
			c.keyRoundRobinIndex = -1
			c.dataSourceCache = new(sync.Map)
			return runImpl(c, l)
		},
	}
	cmd.Flags().String(flags.FlagChainID, "", "chain ID of BandChain network")
	cmd.Flags().String(flags.FlagNode, "tcp://localhost:26657", "RPC url to BandChain node")
	cmd.Flags().String(flagValidator, "", "validator address")
	cmd.Flags().String(flagExecutor, "", "executor name and url for executing the data source script")
	cmd.Flags().String(flags.FlagGasPrices, "", "gas prices for report transaction")
	cmd.Flags().String(flagLogLevel, "info", "set the logger level")
	cmd.Flags().String(flagBroadcastTimeout, "5m", "The time that Yoda will wait for tx commit")
	cmd.Flags().String(flagRPCPollInterval, "1s", "The duration of rpc poll interval")
	cmd.Flags().Uint64(flagMaxTry, 5, "The maximum number of tries to submit a report transaction")
	cmd.Flags().Uint64(flagMaxReport, 10, "The maximum number of reports in one transaction")
	viper.BindPFlag(flags.FlagChainID, cmd.Flags().Lookup(flags.FlagChainID))
	viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
	viper.BindPFlag(flagValidator, cmd.Flags().Lookup(flagValidator))
	viper.BindPFlag(flags.FlagGasPrices, cmd.Flags().Lookup(flags.FlagGasPrices))
	viper.BindPFlag(flagLogLevel, cmd.Flags().Lookup(flagLogLevel))
	viper.BindPFlag(flagExecutor, cmd.Flags().Lookup(flagExecutor))
	viper.BindPFlag(flagBroadcastTimeout, cmd.Flags().Lookup(flagBroadcastTimeout))
	viper.BindPFlag(flagRPCPollInterval, cmd.Flags().Lookup(flagRPCPollInterval))
	viper.BindPFlag(flagMaxTry, cmd.Flags().Lookup(flagMaxTry))
	viper.BindPFlag(flagMaxReport, cmd.Flags().Lookup(flagMaxReport))
	return cmd
}
