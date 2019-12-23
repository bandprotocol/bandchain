package main

import (
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/bandprotocol/d3n/chain/cmtx"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	rpc "github.com/tendermint/tendermint/rpc/client"
)

type OracleRequest struct {
	CodeHash []byte `json:"codeHash" binding:"len=0|len=32"`
	Code     []byte `json:"code"`
	Params   []byte `json:"params" binding:"required"`
}

type OracleRequestResp struct {
	RequestId uint64 `json:"id"`
	CodeHash  []byte `json:"codeHash"`
	Params    []byte `json:"params"`
}

const priv = "a96e62ed3955e65be3aaa3f12d87b6b5cf26039ecfa948dc5107a495418e5430"
const port = "5001"
const nodeURI = "http://localhost:26657"
const queryURI = "http://localhost:1317"

var rpcClient *rpc.HTTP
var pk secp256k1.PrivKeySecp256k1

// TODO
// - Add query from rest client and ask via that endpoint
func HasCode(codeHash []byte) bool {
	key := zoracle.CodeHashStoreKey(codeHash)
	resp, err := rpcClient.ABCIQuery("/store/zoracle/key", key)
	if err != nil {
		return false
	}

	return len(resp.Response.Value) > 0
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

	if len(requestData.Code) > 0 {
		requestData.CodeHash = crypto.Sha256(requestData.Code)
		hasCode := HasCode(requestData.CodeHash)
		// If codeHash not found then store the code
		if !hasCode {
			txSender := cmtx.NewTxSender(pk)
			_, err := txSender.SendTransaction(zoracle.NewMsgStoreCode(requestData.Code, txSender.Sender()), flags.BroadcastBlock)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	} else if len(requestData.CodeHash) > 0 {
		hasCode := HasCode(requestData.CodeHash)
		if !hasCode {
			c.JSON(http.StatusBadRequest, gin.H{"error": "codeHash not found"})
			return
		}
	}

	txSender := cmtx.NewTxSender(pk)
	txr, err := txSender.SendTransaction(zoracle.NewMsgRequest(requestData.CodeHash, requestData.Params, 5, txSender.Sender()), flags.BroadcastBlock)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var requestId uint64
	events := txr.Events
	for _, event := range events {
		if event.Type == "request" {
			for _, attr := range event.Attributes {
				if string(attr.Key) == "id" {
					requestId, err = strconv.ParseUint(string(attr.Value), 10, 64)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				}
			}
		}
	}

	c.JSON(200, OracleRequestResp{
		RequestId: requestId,
		CodeHash:  requestData.CodeHash,
		Params:    requestData.Params,
	})
}

func handleGetRequest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": c.Param("id"),
	})
}

func main() {
	viper.Set("nodeURI", nodeURI)
	privBytes, _ := hex.DecodeString(priv)
	copy(pk[:], privBytes)

	rpcClient = rpc.NewHTTP(nodeURI, "/websocket")

	r := gin.Default()

	r.POST("/request", handleRequestData)
	r.GET("/request/:id", handleGetRequest)

	r.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
