package main

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"
)

func runImpl(client *rpchttp.HTTP, key keyring.Info) error {
	logger.Info("🚀 Starting WebSocket subscriber")
	err := client.Start()
	if err != nil {
		return err
	}

	ctx, cxl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cxl()

	// TODO: We can subscribe only for txs that contain request messages
	query := "tm.event = 'Tx'"
	logger.Info("👂 Subscribing to events with query: %s...", query)
	eventChan, err := client.Subscribe(ctx, "", query)
	if err != nil {
		return err
	}

	for {
		select {
		case ev := <-eventChan:
			go handleTransaction(client, key, ev.Data.(tmtypes.EventDataTx).TxResult)
		}
	}
}

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run [key-name]",
		Aliases: []string{"r"},
		Short:   "Run the oracle process",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := keybase.Key(args[0])
			if err != nil {
				return err
			}

			logger.Info("⭐️ Creating HTTP client with node URI %s", nodeURI)
			client, err := rpchttp.New(nodeURI, "/websocket")
			if err != nil {
				return err
			}
			return runImpl(client, key)
		},
	}
	return cmd
}
