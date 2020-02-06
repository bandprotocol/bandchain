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
	"github.com/levigross/grequests"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bandprotocol/d3n/chain/d3nlib"
	sub "github.com/bandprotocol/d3n/chain/subscriber"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

type Command struct {
	Cmd       string   `json:"cmd"`
	Arguments []string `json:"args"`
}

const limitTimeOut = 10 * time.Second

// TODO: Replace by `BandStatefulClient` after implementation finished.
var bandProvider d3nlib.BandProvider
var allowedCommands = map[string]bool{"curl": true, "date": true}

func getEnv(key, defaultValue string) string {
	tmp := os.Getenv(key)
	if tmp == "" {
		return defaultValue
	}
	return tmp
}

var (
	nodeURI  = getEnv("NODE_URI", "http://localhost:26657")
	queryURI = getEnv("QUERY_URI", "http://localhost:1317")
	privS    = getEnv("PRIVATE_KEY", "eedda7a96ad35758f2ffc404d6ccd7be913f149a530c70e95e2e3ee7a952a877")
)

func getLatestRequestID() (uint64, error) {
	resp, err := grequests.Get(fmt.Sprintf("%s/zoracle/request_number", queryURI), nil)
	if err != nil {
		return 0, err
	}

	var responseStruct struct {
		Result string `json:"result"`
	}
	if err := resp.JSON(&responseStruct); err != nil {
		return 0, err
	}

	return strconv.ParseUint(responseStruct.Result, 10, 64)
}

func main() {
	privB, _ := hex.DecodeString(privS)
	var priv secp256k1.PrivKeySecp256k1
	copy(priv[:], privB)

	var err error
	bandProvider, err = d3nlib.NewBandProvider(nodeURI, priv)
	if err != nil {
		panic(err)
	}

	currentRequestID, err := getLatestRequestID()
	if err != nil {
		panic(err)
	}

	// Setup poll loop
	for {
		newRequestID, err := getLatestRequestID()
		if err != nil {
			log.Println("Cannot get request number error: ", err.Error())
		}

		for currentRequestID < newRequestID {
			currentRequestID++
			go newHandleRequest(currentRequestID)
		}
		time.Sleep(1 * time.Second)
	}

	s := sub.NewSubscriber(nodeURI, "/websocket")

	// Tx events
	s.AddHandler(zoracle.EventTypeRequest, handleRequestAndLog)

	// start subscription
	s.Run()
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

func newHandleRequest(requestID uint64) {
	fmt.Println("Have new request", requestID)
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

	type queryParallelInfo struct {
		index  int
		answer string
		err    error
	}
	chanQueryParallelInfo := make(chan queryParallelInfo, len(commands))
	for i, command := range commands {
		go func(index int, command Command) {
			info := queryParallelInfo{index: index, answer: "", err: nil}
			if !allowedCommands[command.Cmd] {
				info.err = fmt.Errorf("handleRequest unknown command %s", command.Cmd)
				chanQueryParallelInfo <- info
				return
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
				info.err = fmt.Errorf("handleRequest query err with command: %s %v, error: %v", command.Cmd, command.Arguments, err)
				chanQueryParallelInfo <- info
				return
			}

			info.answer = string(query)
			chanQueryParallelInfo <- info
		}(i, command)
	}

	answers := make([]string, len(commands))
	for i := 0; i < len(commands); i++ {
		info := <-chanQueryParallelInfo
		if info.err != nil {
			return sdk.TxResponse{}, info.err
		}
		answers[info.index] = info.answer
	}

	b, _ := json.Marshal(answers)

	tx, err := bandProvider.SendTransaction(
		[]sdk.Msg{zoracle.NewMsgReport(requestID, b, sdk.ValAddress(bandProvider.Sender()))},
		0, 10000000, "", "", "",
		flags.BroadcastSync,
	)
	if err != nil {
		return sdk.TxResponse{}, fmt.Errorf("handleRequest send tx fail : %s", err)
	}
	return tx, nil
}

func handleRequestAndLog(event *abci.Event) {
	fmt.Println(handleRequest(event))
}
