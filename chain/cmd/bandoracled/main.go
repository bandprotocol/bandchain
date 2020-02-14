package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bandprotocol/d3n/chain/d3nlib"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

type Command struct {
	Cmd       string   `json:"cmd"`
	Arguments []string `json:"args"`
}

const limitTimeOut = 10 * time.Second

func getEnv(key, defaultValue string) string {
	tmp := os.Getenv(key)
	if tmp == "" {
		return defaultValue
	}
	return tmp
}

var (
	bandClient      d3nlib.BandStatefulClient
	allowedCommands = map[string]bool{"curl": true, "date": true}
	nodeURI         = getEnv("NODE_URI", "http://localhost:26657")
	privS           = getEnv("PRIVATE_KEY", "06be35b56b048c5a6810a47e2ef612eaed735ccb0d7ea4fc409f23f1d1a16e0b")
)

func getLatestRequestID() (int64, error) {
	cliCtx := bandClient.GetContext()
	res, _, err := cliCtx.Query("custom/zoracle/request_number")
	if err != nil {
		return 0, err
	}
	var requestID int64
	err = cliCtx.Codec.UnmarshalJSON(res, &requestID)
	if err != nil {
		return 0, err
	}
	return requestID, nil
}

func main() {
	privB, _ := hex.DecodeString(privS)
	var priv secp256k1.PrivKeySecp256k1
	copy(priv[:], privB)

	var err error
	bandClient, err = d3nlib.NewBandStatefulClient(nodeURI, priv)
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
			go handleRequestAndLog(currentRequestID)
		}
		time.Sleep(1 * time.Second)
	}
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

func getPrepareBytes(searchResult *sdk.SearchTxsResult) ([]byte, error) {
	for _, tx := range searchResult.Txs {
		// Stringevents (type of tx.Events) are deprecated in next cosmos-release.
		for _, event := range tx.Events {
			if event.Type == "request" {
				for _, kv := range event.Attributes {
					if string(kv.Key) == zoracle.AttributeKeyPrepare {
						return hex.DecodeString(string(kv.Value))
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("Cannot find prepare bytes")
}

func handleRequest(requestID int64) (sdk.TxResponse, error) {
	searchResult, err := utils.QueryTxsByEvents(
		bandClient.GetContext(),
		[]string{fmt.Sprintf("request.id='%d'", requestID)},
		1,
		100,
	)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	byteValue, err := getPrepareBytes(searchResult)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var commands []Command
	err = json.Unmarshal(byteValue, &commands)
	if err != nil {
		return sdk.TxResponse{}, err
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

	// b, _ := json.Marshal(answers)

	// TODO: Send new report
	return sdk.TxResponse{}, nil

	// return bandClient.SendTransaction(
	// 	zoracle.NewMsgReport(requestID, b, sdk.ValAddress(bandClient.Sender())),
	// 	10000000, "", "", "",
	// 	flags.BroadcastSync,
	// )
}

func handleRequestAndLog(requestID int64) {
	fmt.Println(handleRequest(requestID))
}
