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

	file, err := os.Open("../../wasm/res/owasm_example_bg.wasm")
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

	// Send transaction to store code first (commend it if already stored code)
	fmt.Println(tx.SendTransaction(zoracle.NewMsgStoreCode(bytes, tx.Sender()), flags.BroadcastBlock))

	codeHash, _ := hex.DecodeString("549ffdeccba87459c61fb2929b678f4c6cd671d5f41aa8576d3f4528c553f272")

	// BTC parameter
	params, _ := hex.DecodeString("0000000000000007626974636f696e0000000000000003425443")
	// Send request by code hash
	fmt.Println(tx.SendTransaction(zoracle.NewMsgRequest(codeHash, params, 4, tx.Sender()), flags.BroadcastBlock))

	// ETH parameter
	params, _ = hex.DecodeString("0000000000000008657468657265756d0000000000000003455448")
	// Send request by same code hash with new parameter
	fmt.Println(tx.SendTransaction(zoracle.NewMsgRequest(codeHash, params, 4, tx.Sender()), flags.BroadcastBlock))
}
