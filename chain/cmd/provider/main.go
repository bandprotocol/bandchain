package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bandprotocol/d3n/chain/cmtx"
	sub "github.com/bandprotocol/d3n/chain/subscriber"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

const limitTimeOut = 10 * time.Second

var txSender cmtx.TxSender
var allowedCommands = map[string]bool{"curl": true}

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
	s.AddHandler(zoracle.EventTypeRequest, handleRequestAndLog)

	// start subscription
	s.Run()
}

type Command struct {
	Cmd       string   `json:"cmd"`
	Arguments []string `json:"args"`
}

func execWithTimeout(command Command, limit time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	cmd := exec.CommandContext(ctx, command.Cmd, command.Arguments...)
	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return []byte{}, fmt.Errorf("Command timed out")
	}
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func handleRequest(event *abci.Event) (sdk.TxResponse, error) {
	var requestID uint64
	var commands []Command

	for _, kv := range event.GetAttributes() {
		switch string(kv.Key) {
		case zoracle.AttributeKeyRequestID:
			var err error
			requestID, err = strconv.ParseUint(string(kv.Value), 10, 64)
			if err != nil {
				return sdk.TxResponse{}, fmt.Errorf("handleRequest %s", err)
			}
		case zoracle.AttributeKeyPrepare:
			byteValue, err := hex.DecodeString(string(kv.Value))
			if err != nil {
				return sdk.TxResponse{}, fmt.Errorf("handleRequest %s", err)
			}
			err = json.Unmarshal(byteValue, &commands)
			if err != nil {
				return sdk.TxResponse{}, fmt.Errorf("handleRequest %s", err)
			}
		}
	}
	var answer []string
	for _, command := range commands {
		if !allowedCommands[command.Cmd] {
			return sdk.TxResponse{}, fmt.Errorf("handleRequest unknown command %s", command.Cmd)
		}
		dockerCommand := Command{
			Cmd: "docker",
			Arguments: append([]string{
				"run", "--rm", "band-provider",
				command.Cmd,
			}, command.Arguments...),
		}
		query, err := execWithTimeout(dockerCommand, limitTimeOut)
		if err != nil {
			return sdk.TxResponse{}, fmt.Errorf("handleRequest query err with command: %s %v, error: %v", command.Cmd, command.Arguments, err)
		}
		answer = append(answer, string(query))
	}
	b, _ := json.Marshal(answer)

	tx, err := txSender.SendTransaction(zoracle.NewMsgReport(requestID, b, sdk.ValAddress(txSender.Sender())), flags.BroadcastSync)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("handleRequest send tx fail : %s", err)
	}
	return tx, nil
}

func handleRequestAndLog(event *abci.Event) {
	fmt.Println(handleRequest(event))
}
