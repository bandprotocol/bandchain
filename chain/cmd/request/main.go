package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
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
		log.Fatal("Missing private key")
	}
	nodeURI, ok := os.LookupEnv("NODE_URI")
	if !ok {
		log.Fatal("Missing node uri")
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
	fmt.Println(tx.SendTransaction(zoracle.NewMsgStoreCode(bytes, tx.Sender()), flags.BroadcastBlock))

	codeHash, _ := hex.DecodeString("33cefd052eb9b0cda3d38d4d87313295e45fda5c5b0d6ab9e54870866f62fc80")

	// BTC parameter
	params, _ := hex.DecodeString("0000000000000007626974636f696e0000000000000003425443")
	// Send request by code hash
	fmt.Println(tx.SendTransaction(zoracle.NewMsgRequest(codeHash, params, 4, tx.Sender()), flags.BroadcastBlock))

	// ETH parameter
	params, _ = hex.DecodeString("0000000000000008657468657265756d0000000000000003455448")
	// Send request by same code hash with new parameter
	fmt.Println(tx.SendTransaction(zoracle.NewMsgRequest(codeHash, params, 4, tx.Sender()), flags.BroadcastBlock))
}
