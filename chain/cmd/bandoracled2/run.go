package main

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	libclient "github.com/tendermint/tendermint/rpc/lib/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

func newRPCClient(addr string) (*rpcclient.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}

	rpcClient, err := rpcclient.NewHTTPWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}

	return rpcClient, nil
}

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the oracle process",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Read URL from chain configuration
			client, err := newRPCClient("tcp://localhost:26657")
			if err != nil {
				return err
			}

			err = client.Start()
			if err != nil {
				return err
			}

			ctx, cxl := context.WithTimeout(context.Background(), 5*time.Second)
			defer cxl()

			// TODO: We can subscribe only for txs that contain request messages
			eventChan, err := client.Subscribe(ctx, "", "tm.event = 'Tx'")
			if err != nil {
				return err
			}

			logger.Info("👂 Start listening to transaction events...")
			for {
				select {
				case ev := <-eventChan:
					tx := ev.Data.(tmtypes.EventDataTx).TxResult
					// TODO: Get real transaction hash here
					logger.Debug("👀 Inspecting transaction %s", "HASH")

					logs, err := sdk.ParseABCILogs(tx.Result.Log)
					if err != nil {
						logger.Error("❌ Failed to parse transaction logs: %s", err.Error())
						return err
					}

					for _, log := range logs {
						// TODO: Remove magic string
						messageType, err := GetEventValue(log, "message", "action")
						if err != nil {
							logger.Error("❌ Failed to get message type: %s", err.Error())
						}

						if messageType != "request" {
							logger.Debug("⏭️ Skipping non-request message type: %s", messageType)
							continue
						}

						go func(log sdk.ABCIMessageLog) {
							err := handleRequestLog(log)
							if err != nil {
								logger.Error("❌ Failed to handle request: %s", err.Error())
							}
						}(log)
					}
				}
			}
		},
	}
	return cmd
}
