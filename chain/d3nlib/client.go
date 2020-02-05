package d3nlib

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto"
)

type BandStatefulClient struct {
	sequenceNumber uint64

	mtx      sync.RWMutex
	provider BandProvider
}

func NewBandStatefulClient(nodeURI string, privKey crypto.PrivKey) (BandStatefulClient, error) {
	provider, err := NewBandProvider(nodeURI, privKey)
	if err != nil {
		return BandStatefulClient{}, err
	}
	_, seq, err := authtypes.NewAccountRetriever(provider.cliCtx).GetAccountNumberSequence(provider.Sender())
	if err != nil {
		return BandStatefulClient{}, err
	}

	return BandStatefulClient{
		sequenceNumber: seq,
		provider:       provider,
	}, nil
}

func (client *BandStatefulClient) SendTransaction(
	msg sdk.Msg, gas uint64,
	memo, fees, gasPrices, broadcastMode string,
) (sdk.TxResponse, error) {
	// Ask current sequence number
	_, seq, err := authtypes.NewAccountRetriever(client.provider.cliCtx).
		GetAccountNumberSequence(client.provider.Sender())
	if err != nil {
		return sdk.TxResponse{}, err
	}
	client.mtx.Lock()
	if seq > client.sequenceNumber {
		client.sequenceNumber = seq
	}
	nonce := client.sequenceNumber
	client.sequenceNumber++
	client.mtx.Unlock()

	tx, err := client.provider.SendTransaction(
		[]sdk.Msg{msg}, nonce, gas, memo, fees, gasPrices, broadcastMode,
	)

	if err != nil {
		// Reset sequence number to 0 make next request use new sequence number
		client.mtx.Lock()
		client.sequenceNumber = 0
		client.mtx.Unlock()
	}
	return tx, err
}

func (client *BandStatefulClient) Sender() sdk.AccAddress {
	return client.provider.Sender()
}
