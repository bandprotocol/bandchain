package bandlib

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type msgDetail struct {
	msg     sdk.Msg
	gas     uint64
	fees    sdk.Coins
	txChan  chan sdk.TxResponse
	errChan chan error
}

// BandStatefulClient contains state client
type BandStatefulClient struct {
	memo            string
	maximumMsgPerTx int
	provider        BandProvider
	msgChan         chan msgDetail
	msgs            []sdk.Msg
	totalGas        uint64
	totalFees       sdk.Coins
	txChans         []chan sdk.TxResponse
	errChans        []chan error
	readMode        <-chan struct{}
	readModeBackUp  <-chan struct{}
}

// NewBandStatefulClient creates new instance of BandStatefulClient.
func NewBandStatefulClient(
	nodeURI string, privKey crypto.PrivKey, msgsCap, maximumMsgPerTx int, memo string,
) (BandStatefulClient, error) {
	provider, err := NewBandProvider(nodeURI, privKey)
	if err != nil {
		return BandStatefulClient{}, err
	}
	ch := make(chan struct{})
	close(ch)
	client := BandStatefulClient{
		provider:        provider,
		msgChan:         make(chan msgDetail, msgsCap),
		maximumMsgPerTx: maximumMsgPerTx,
		memo:            memo,
		readMode:        ch,
		readModeBackUp:  ch,
	}
	go client.loop()

	return client, nil
}

func (client *BandStatefulClient) SendTransaction(
	msg sdk.Msg, gas uint64, fees string,
) (sdk.TxResponse, error) {
	// Add msg to channel
	parsedFees, err := sdk.ParseCoins(fees)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	txChan := make(chan sdk.TxResponse, 0)
	errChan := make(chan error, 0)

	client.msgChan <- msgDetail{
		msg: msg, gas: gas, fees: parsedFees,
		txChan: txChan, errChan: errChan,
	}

	select {
	case txResponse := <-txChan:
		return txResponse, nil
	case err := <-errChan:
		return sdk.TxResponse{}, err
	}
}

func (client *BandStatefulClient) Sender() sdk.AccAddress {
	return client.provider.Sender()
}

func (client *BandStatefulClient) GetContext() context.CLIContext {
	return client.provider.cliCtx
}

func (client *BandStatefulClient) loop() {
	for {
		select {
		case <-client.readMode:
			{
				select {
				case msg := <-client.msgChan:
					{
						client.msgs = append(client.msgs, msg.msg)
						client.totalGas += msg.gas
						client.totalFees = client.totalFees.Add(msg.fees)
						client.txChans = append(client.txChans, msg.txChan)
						client.errChans = append(client.errChans, msg.errChan)
						if len(client.msgs) == client.maximumMsgPerTx {
							client.readMode = nil
						}
					}
				case <-time.After(100 * time.Millisecond):
					{
						if len(client.msgs) != 0 {
							client.readMode = nil
						}
					}
				}
			}
		default:
			{
				seq, err := client.provider.GetSequenceFromChain()
				if err != nil {
					// TODO: error handler
					fmt.Println(err)
					continue
				}
				tx, err := client.provider.SendTransaction(
					client.msgs,
					seq,
					client.totalGas,
					client.memo,
					client.totalFees.String(),
					flags.BroadcastBlock,
				)
				if err != nil {
					for _, errChan := range client.errChans {
						errChan <- err
					}
				} else {
					for _, txChan := range client.txChans {
						txChan <- tx
					}
				}

				// Clear state
				client.readMode = client.readModeBackUp
				client.msgs = []sdk.Msg{}
				client.totalGas = uint64(0)
				client.totalFees = sdk.Coins{}
				client.txChans = []chan sdk.TxResponse{}
				client.errChans = []chan error{}
			}
		}
	}
}
