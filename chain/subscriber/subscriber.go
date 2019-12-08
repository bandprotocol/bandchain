package subscriber

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"

	tmclient "github.com/tendermint/tendermint/rpc/client"
)

type Subscriber struct {
	url        string
	wsEndPoint string
	handlers   map[string]func(*abci.Event)
}

func NewSubscriber(url string, wsEndPoint string) *Subscriber {
	return &Subscriber{
		url:        url,
		wsEndPoint: wsEndPoint,
		handlers:   make(map[string]func(*abci.Event)),
	}
}

func (s *Subscriber) AddHandler(event string, handler func(*abci.Event)) {
	s.handlers[event] = handler
}

func (s *Subscriber) Run() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	client := tmclient.NewHTTP(s.url, s.wsEndPoint)
	client.SetLogger(logger)
	err := client.Start()
	if err != nil {
		logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}
	defer client.Stop()

	newBlockResult, err := client.Subscribe(context.Background(), "test", "tm.event = 'NewBlock'", 1000)
	if err != nil {
		logger.Error("Failed to subscribe to query", "err", err, "query", "tm.event = 'NewBlock'")
		os.Exit(1)
	}

	txResult, err := client.Subscribe(context.Background(), "test", "tm.event = 'Tx'", 1000)
	if err != nil {
		logger.Error("Failed to subscribe to query", "err", err, "query", "tm.event = 'Tx'")
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case result := <-newBlockResult:
			data := result.Data.(types.EventDataNewBlock)
			events := data.ResultEndBlock.Events
			for _, event := range events {
				if handler, ok := s.handlers[event.GetType()]; ok {
					handler(&event)
				}
			}
		case result := <-txResult:
			data := result.Data.(types.EventDataTx)
			events := data.Result.Events
			for _, event := range events {
				if handler, ok := s.handlers[event.GetType()]; ok {
					handler(&event)
				}
			}
		case <-quit:
			os.Exit(0)
		}
	}
}
