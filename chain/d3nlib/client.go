package d3nlib

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type BandStatefulClient struct {
	sequenceNumber uint64

	mtx      sync.Mutex
	provider BandProvider
}

func NewBandStatefulClient(nodeURI string, privKey crypto.PrivKey) (BandStatefulClient, error) {
	provider, err := NewBandProvider(nodeURI, privKey)
	if err != nil {
		return BandStatefulClient{}, err
	}
	seq, err := provider.GetSequenceFromChain()
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
	seq, err := client.provider.GetSequenceFromChain()
	if err != nil {
		return sdk.TxResponse{}, err
	}
	var nonce uint64
	{
		client.mtx.Lock()
		defer client.mtx.Unlock()
		if seq > client.sequenceNumber {
			client.sequenceNumber = seq
		}
		nonce = client.sequenceNumber
		client.sequenceNumber++
	}

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
