package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	app "github.com/bandprotocol/d3n/chain"
	"github.com/bandprotocol/d3n/chain/cmtx"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
	cmc "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/client"

	"github.com/bandprotocol/d3n/chain/wasm"
)

type OracleRequest struct {
	CodeHash cmn.HexBytes    `json:"codeHash" binding:"len=32"`
	Params   json.RawMessage `json:"params" binding:"required"`
}

type OracleRequestResp struct {
	RequestId uint64       `json:"id"`
	TxHash    cmn.HexBytes `json:"txHash"`
}

type ExecuteRequest struct {
	Code   cmn.HexBytes    `json:"code" binding:"required"`
	Params json.RawMessage `json:"params" binding:"required"`
}

type ExecuteResponse struct {
	Result json.RawMessage `json:"result"`
}

type ParamsInfoRequest struct {
	Code cmn.HexBytes `json:"code" binding:"required"`
}

type ParamsInfoResponse struct {
	Params json.RawMessage `json:"params"`
}

type StoreRequest struct {
	Code cmn.HexBytes `json:"code" binding:"required"`
	Name string       `json:"name" binding:"required"`
}

type StoreResponse struct {
	TxHash   cmn.HexBytes `json:"txHash"`
	CodeHash cmn.HexBytes `json:"codeHash"`
}

type Command struct {
	Cmd       string   `json:"cmd"`
	Arguments []string `json:"args"`
}

var allowedCommands = map[string]bool{"curl": true, "date": true}

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

func getEnv(key, def string) string {
	tmp := os.Getenv(key)
	if tmp == "" {
		return def
	}
	return tmp
}

var (
	port    = getEnv("PORT", "5001")
	nodeURI = getEnv("NODE_URI", "http://localhost:26657")
	priv    = getEnv("PRIVATE_KEY", "eedda7a96ad35758f2ffc404d6ccd7be913f149a530c70e95e2e3ee7a952a877")
)

var rpcClient *rpc.HTTP
var pk secp256k1.PrivKeySecp256k1
var txSender cmtx.TxSender
var cliCtx cmc.CLIContext
var cdc *codec.Codec

// TODO
// - Add query from rest client and ask via that endpoint
func HasCode(codeHash []byte) (bool, error) {
	key := zoracle.CodeHashStoreKey(codeHash)
	resp, err := rpcClient.ABCIQuery("/store/zoracle/key", key)
	if err != nil {
		return false, err
	}

	return len(resp.Response.Value) > 0, nil
}

func handleRequestData(c *gin.Context) {
	var requestData OracleRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Mock this endpoint for front-end for now

	// if len(requestData.CodeHash) == 0 && len(requestData.Code) == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code/codeHash"})
	// 	return
	// }
	// if len(requestData.CodeHash) > 0 && len(requestData.Code) > 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Only one of code/codeHash can be sent"})
	// 	return
	// }

	// // TODO
	// // Need some work around to make params can be empty bytes
	// if len(requestData.Params) <= 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Params should not be empty bytes"})
	// 	return
	// }

	// var params []byte

	// if len(requestData.Code) > 0 {
	// 	if len(requestData.Name) <= 0 {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Name should not be empty string"})
	// 		return
	// 	}
	// 	requestData.CodeHash = zoracle.NewStoredCode(requestData.Code, requestData.Name, txSender.Sender()).GetCodeHash()
	// 	hasCode, err := HasCode(requestData.CodeHash)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	// If codeHash not found then store the code
	// 	if !hasCode {
	// 		_, err := txSender.SendTransaction(zoracle.NewMsgStoreCode(requestData.Code, requestData.Name, txSender.Sender()), flags.BroadcastBlock)
	// 		if err != nil {
	// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 			return
	// 		}
	// 	}

	// 	// Parse params
	// 	params, err = wasm.SerializeParams(requestData.Code, []byte(requestData.Params))
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// } else if len(requestData.CodeHash) > 0 {
	// 	hasCode, err := HasCode(requestData.CodeHash)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	if !hasCode {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "codeHash not found"})
	// 		return
	// 	}

	// 	params, err = hex.DecodeString(requestData.Params)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// }

	// txr, err := txSender.SendTransaction(zoracle.NewMsgRequest(requestData.CodeHash, params, 4, txSender.Sender()), flags.BroadcastBlock)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// requestId := uint64(0)
	// events := txr.Events
	// for _, event := range events {
	// 	if event.Type == "request" {
	// 		for _, attr := range event.Attributes {
	// 			if string(attr.Key) == "id" {
	// 				requestId, err = strconv.ParseUint(attr.Value, 10, 64)
	// 				if err != nil {
	// 					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 					return
	// 				}
	// 				break
	// 			}
	// 		}
	// 	}
	// }
	// if requestId == 0 {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("cannot find requestId: %v", txr)})
	// 	return
	// }

	mockTxHash, _ := hex.DecodeString("A5A8482E19F434FD7083B79B51270527243DB1B4EAAD2CEBB3AA75915719589A")

	c.JSON(200, OracleRequestResp{
		RequestId: 1,
		TxHash:    mockTxHash,
	})
}

func handleParamsInfo(c *gin.Context) {
	var req ParamsInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := wasm.ParamsInfo(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, ParamsInfoResponse{
		Params: res,
	})
}

func handleExecute(c *gin.Context) {
	var req ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawParams, err := wasm.SerializeParams(req.Code, req.Params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Call prepare function
	prepare, err := wasm.Prepare(req.Code, rawParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var commands []Command
	err = json.Unmarshal(prepare, &commands)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var answer []string
	for _, command := range commands {
		if !allowedCommands[command.Cmd] {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("handleRequest unknown command %s", command.Cmd)})
			return
		}
		dockerCommand := Command{
			Cmd: "docker",
			Arguments: append([]string{
				"run", "--rm", "band-provider",
				command.Cmd,
			}, command.Arguments...),
		}
		query, err := execWithTimeout(dockerCommand, 10*time.Second)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		answer = append(answer, string(query))
	}

	b, _ := json.Marshal(answer)

	rawResult, err := wasm.Execute(req.Code, rawParams, [][]byte{b})
	result, err := wasm.ParseResult(req.Code, rawResult)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, ExecuteResponse{
		Result: result,
	})
}

func handleStore(c *gin.Context) {
	var req StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	codeHash := zoracle.NewStoredCode(req.Code, req.Name, txSender.Sender()).GetCodeHash()
	tx, err := txSender.SendTransaction(zoracle.NewMsgStoreCode(req.Code, req.Name, txSender.Sender()), flags.BroadcastBlock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	txHash, err := hex.DecodeString(tx.TxHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, StoreResponse{
		TxHash:   txHash,
		CodeHash: cmn.HexBytes(codeHash),
	})
}

func main() {
	viper.Set("nodeURI", nodeURI)
	privBytes, _ := hex.DecodeString(priv)
	copy(pk[:], privBytes)

	txSender = cmtx.NewTxSender(pk)
	cdc = app.MakeCodec()
	rpcClient = rpc.NewHTTP(nodeURI, "/websocket")
	cliCtx = cmtx.NewCLIContext(txSender.Sender()).WithCodec(cdc)

	r := gin.Default()

	r.POST("/request", handleRequestData)
	r.POST("/params-info", handleParamsInfo)
	r.POST("/execute", handleExecute)
	r.POST("/store", handleStore)

	// Allows all origins
	r.Use(cors.Default())
	r.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
