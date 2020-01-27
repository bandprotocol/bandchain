package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	app "github.com/bandprotocol/d3n/chain"
	"github.com/bandprotocol/d3n/chain/cmtx"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
	cmc "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/client"

	"github.com/bandprotocol/d3n/chain/wasm"
)

type HexString []byte

func (hexstr *HexString) UnmarshalJSON(b []byte) error {
	b, err := hex.DecodeString(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*hexstr = b
	return nil
}

type RawJson []byte

var EmptyMap = []byte("{}")

func (j RawJson) MarshalJSON() ([]byte, error) {
	return []byte(j), nil
}

type OracleRequest struct {
	CodeHash HexString `json:"codeHash" binding:"len=0|len=32"`
	Code     HexString `json:"code"`
	Name     string    `json:"name" binding:"required"`
	Params   string    `json:"params" binding:"required"`
}

type OracleRequestResp struct {
	RequestId uint64       `json:"id"`
	CodeHash  cmn.HexBytes `json:"codeHash"`
}

type ExecuteRequest struct {
	Code   HexString       `json:"code" binding:"required"`
	Params json.RawMessage `json:"params" binding:"required"`
}

type ExecuteResponse struct {
	Result RawJson `json:"result"`
}

type ParamsInfoRequest struct {
	Code HexString `json:"code" binding:"required"`
}

type ParamsInfoResponse struct {
	Params RawJson `json:"params"`
}

type Command struct {
	Cmd       string   `json:"cmd"`
	Arguments []string `json:"args"`
}

var allowedCommands = map[string]bool{"curl": true}

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

const priv = "06be35b56b048c5a6810a47e2ef612eaed735ccb0d7ea4fc409f23f1d1a16e0b"

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
	if len(requestData.CodeHash) == 0 && len(requestData.Code) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code/codeHash"})
		return
	}
	if len(requestData.CodeHash) > 0 && len(requestData.Code) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only one of code/codeHash can be sent"})
		return
	}

	// TODO
	// Need some work around to make params can be empty bytes
	if len(requestData.Params) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Params should not be empty bytes"})
		return
	}

	var params []byte

	if len(requestData.Code) > 0 {
		if len(requestData.Name) <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name should not be empty string"})
			return
		}
		requestData.CodeHash = zoracle.NewStoredCode(requestData.Code, requestData.Name, txSender.Sender()).GetCodeHash()
		hasCode, err := HasCode(requestData.CodeHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// If codeHash not found then store the code
		if !hasCode {
			_, err := txSender.SendTransaction(zoracle.NewMsgStoreCode(requestData.Code, requestData.Name, txSender.Sender()), flags.BroadcastBlock)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// Parse params
		params, err = wasm.SerializeParams(requestData.Code, []byte(requestData.Params))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if len(requestData.CodeHash) > 0 {
		hasCode, err := HasCode(requestData.CodeHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !hasCode {
			c.JSON(http.StatusBadRequest, gin.H{"error": "codeHash not found"})
			return
		}

		params, err = hex.DecodeString(requestData.Params)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	txr, err := txSender.SendTransaction(zoracle.NewMsgRequest(requestData.CodeHash, params, 4, txSender.Sender()), flags.BroadcastBlock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	requestId := uint64(0)
	events := txr.Events
	for _, event := range events {
		if event.Type == "request" {
			for _, attr := range event.Attributes {
				if string(attr.Key) == "id" {
					requestId, err = strconv.ParseUint(attr.Value, 10, 64)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					break
				}
			}
		}
	}
	if requestId == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("cannot find requestId: %v", txr)})
		return
	}

	c.JSON(200, OracleRequestResp{
		RequestId: requestId,
		CodeHash:  cmn.HexBytes(requestData.CodeHash),
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

	rawParams, err := wasm.SerializeParams(req.Code, []byte(req.Params))
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

	r.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
