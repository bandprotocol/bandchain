package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bandprotocol/d3n/chain/cmtx"
	sub "github.com/bandprotocol/d3n/chain/subscriber"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

var txSender cmtx.TxSender

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

	txSender = cmtx.NewTxSender(priv)
	s := sub.NewSubscriber(viper.GetString("nodeURI"), "/websocket")

	// Tx events
	s.AddHandler(zoracle.EventTypeRequest, handleRequest)

	// start subscription
	s.Run()
}

type Command struct {
	Cmd       string   `json:"cmd"`
	Arguments []string `json:"args"`
}

func handleRequest(event *abci.Event) {
	var requestID uint64
	var commands []Command

	for _, kv := range event.GetAttributes() {
		switch string(kv.Key) {
		case zoracle.AttributeKeyRequestID:
			var err error
			requestID, err = strconv.ParseUint(string(kv.Value), 10, 64)
			if err != nil {
				fmt.Printf("handleRequest %s", err)
				return
			}
		case zoracle.AttributeKeyPrepare:
			byteValue, err := hex.DecodeString(string(kv.Value))
			if err != nil {
				fmt.Printf("handleRequest %s", err)
				return
			}
			err = json.Unmarshal(byteValue, &commands)
			if err != nil {
				fmt.Printf("handleRequest %s", err)
				return
			}
		}
	}
	var answer []string
	for _, command := range commands {
		if command.Cmd != "curl" {
			fmt.Printf("handleRequest unknown command %s", command.Cmd)
			return
		}
		cmd := exec.Command(command.Cmd, command.Arguments...)
		query, err := cmd.Output()
		if err != nil {
			fmt.Printf("handleRequest query err with command %s %v", command.Cmd, command.Arguments)
			return
		}
		answer = append(answer, string(query))
	}
	b, _ := json.Marshal(answer)

	tx, err := txSender.SendTransaction(zoracle.NewMsgReport(requestID, b, sdk.ValAddress(txSender.Sender())), flags.BroadcastSync)
	if err != nil {
		fmt.Printf("handleRequest %s", err)
		return
	}
	fmt.Println("Tx:", tx)
}
