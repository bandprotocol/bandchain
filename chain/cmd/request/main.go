package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/bandprotocol/bandx/oracle/cmtx"
	"github.com/bandprotocol/bandx/oracle/x/oracle"
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

	file, err := os.Open("../../wasm/res/awesome_oracle_bg.wasm")
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

	fmt.Println(tx.SendTransaction(oracle.NewMsgRequest(bytes, 4, tx.Sender())))
}
