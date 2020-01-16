package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/bandprotocol/d3n/chain/cmtx"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// File to send new request to bandchain
func main() {
	// Get environment variable
	privS, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		// Default private key is eedda7a96ad35758f2ffc404d6ccd7be913f149a530c70e95e2e3ee7a952a877
		privS = "eedda7a96ad35758f2ffc404d6ccd7be913f149a530c70e95e2e3ee7a952a877"
	}
	nodeURI, ok := os.LookupEnv("NODE_URI")
	if !ok {
		// Default node uri is tcp://localhost:26657
		nodeURI = "tcp://localhost:26657"
	}
	viper.Set("nodeURI", nodeURI)

	privB, _ := hex.DecodeString(privS)
	var priv secp256k1.PrivKeySecp256k1
	copy(priv[:], privB)

	tx := cmtx.NewTxSender(priv)

	file, err := os.Open("../../wasm/res/test_u64.wasm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		panic(statsErr)
	}

	size := stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	if err != nil {
		panic(err)
	}

	// Send token
	// to, _ := sdk.AccAddressFromBech32("band13zmknvkq2sj920spz90g4r9zjan8g584x8qalj")
	// fmt.Println(tx.SendTransaction(bank.MsgSend{
	// 	FromAddress: tx.Sender(),
	// 	ToAddress:   to,
	// 	Amount:      sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(10))),
	// }, flags.BroadcastBlock))

	// Send transaction to store code first (commend it if already stored code)
	fmt.Println(tx.SendTransaction(zoracle.NewMsgStoreCode(bytes, "Crypto price", tx.Sender()), flags.BroadcastBlock))

	codeHash, _ := hex.DecodeString("089a092741d2bbe10b1cfaa8e48d1512a51ef183e579ae29f89af59db3e72c85")

	// BTC parameter
	params, _ := hex.DecodeString("0000000000000007626974636f696e0000000000000003425443")
	// Send request by code hash
	fmt.Println(tx.SendTransaction(zoracle.NewMsgRequest(codeHash, params, 4, tx.Sender()), flags.BroadcastBlock))

	// ETH parameter
	params, _ = hex.DecodeString("0000000000000008657468657265756d0000000000000003455448")
	// Send request by same code hash with new parameter
	fmt.Println(tx.SendTransaction(zoracle.NewMsgRequest(codeHash, params, 4, tx.Sender()), flags.BroadcastBlock))
}
