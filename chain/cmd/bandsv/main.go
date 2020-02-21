package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bandprotocol/d3n/chain/app"
	"github.com/bandprotocol/d3n/chain/d3nlib"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/client"

	"github.com/bandprotocol/d3n/chain/byteexec"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

const (
	Asynchronous = "ASYNCHRONOUS"
	Synchronous  = "SYNCHRONOUS"
	Full         = "FULL"

	DefaultRequestedValidatorCount  = 1
	DefaultSufficientValidatorCount = 1
	DefaultExpiration               = 100
)

type OracleRequest struct {
	Type                     string `json:"type" binding:"required"`
	OracleScriptID           int64  `json:"oracleScriptID,string" binding:"required"`
	Calldata                 []byte `json:"calldata" binding:"required"`
	RequestedValidatorCount  int64  `json:"requestedValidatorCount,string"`
	SufficientValidatorCount int64  `json:"sufficientValidatorCount,string"`
	Expiration               int64  `json:"expiration,string"`
}

type OracleRequestResp struct {
	TxHash    string          `json:"txHash"`
	RequestID int64           `json:"id,omitempty"`
	Proof     json.RawMessage `json:"proof,omitempty"`
}

type ExecuteRequest struct {
	Executable cmn.HexBytes `json:"executable" binding:"required"`
	Calldata   string       `json:"calldata" binding:"required"`
}

type ExecuteResponse struct {
	Result string `json:"result"`
}

func getEnv(key, def string) string {
	tmp := os.Getenv(key)
	if tmp == "" {
		return def
	}
	return tmp
}

var (
	port     = getEnv("PORT", "5001")
	nodeURI  = getEnv("NODE_URI", "http://localhost:26657")
	queryURI = getEnv("QUERY_URI", "http://localhost:1317")
	priv     = getEnv("PRIVATE_KEY", "eedda7a96ad35758f2ffc404d6ccd7be913f149a530c70e95e2e3ee7a952a877")
)

var rpcClient *rpc.HTTP
var pk secp256k1.PrivKeySecp256k1
var bandClient d3nlib.BandStatefulClient
var cdc *codec.Codec

type serializeResponse struct {
	Result cmn.HexBytes `json:"result"`
}

func handleRequestData(c *gin.Context) {
	var req OracleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqType := req.Type
	if reqType != Asynchronous && reqType != Synchronous && reqType != Full {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type not match"})
		return
	}

	if req.RequestedValidatorCount == 0 {
		req.RequestedValidatorCount = DefaultRequestedValidatorCount
	}

	if req.SufficientValidatorCount == 0 {
		req.SufficientValidatorCount = DefaultSufficientValidatorCount
	}

	if req.Expiration == 0 {
		req.Expiration = DefaultExpiration
	}

	// unconfirmed respond
	if reqType == Asynchronous {
		txr, err := bandClient.SendTransaction(
			zoracle.NewMsgRequestData(
				req.OracleScriptID,
				req.Calldata,
				req.RequestedValidatorCount,
				req.SufficientValidatorCount,
				req.Expiration,
				bandClient.Sender(),
			),
			20000000, "", "", "",
			flags.BroadcastAsync,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, OracleRequestResp{
			TxHash: txr.TxHash,
		})
		return
	}

	txr, err := bandClient.SendTransaction(
		zoracle.NewMsgRequestData(
			req.OracleScriptID,
			req.Calldata,
			req.RequestedValidatorCount,
			req.SufficientValidatorCount,
			req.Expiration,
			bandClient.Sender(),
		),
		20000000, "", "", "",
		flags.BroadcastBlock,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	requestID := int64(0)
	for _, event := range txr.Events {
		if event.Type == "request" {
			for _, attr := range event.Attributes {
				if string(attr.Key) == "id" {
					requestID, err = strconv.ParseInt(attr.Value, 10, 64)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					break
				}
			}
		}
	}
	if requestID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("cannot find requestID: %v", txr)})
		return
	}

	// confirmed respond
	if reqType == Synchronous {
		c.JSON(200, OracleRequestResp{
			TxHash:    txr.TxHash,
			RequestID: requestID,
		})
		return
	}

	for i := 0; i < 10; i++ {
		resp, err := grequests.Get(fmt.Sprintf(`%s/d3n/proof/%d`, queryURI, requestID), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// confirmed + proof respond
		if resp.StatusCode == 200 {
			var proof struct {
				Result json.RawMessage `json:"result"`
			}

			err = resp.JSON(&proof)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, OracleRequestResp{
				TxHash:    txr.TxHash,
				RequestID: requestID,
				Proof:     proof.Result,
			})
			return
		}

		time.Sleep(3 * time.Second)
	}

	// finding proof timeout
	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf(`Cannot find proof in this TxHash %s`, txr.TxHash)})
}

func handleExecute(c *gin.Context) {
	var req ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := byteexec.RunOnDocker(req.Executable, 1*time.Minute, req.Calldata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, ExecuteResponse{
		Result: strings.TrimSpace(string(result)),
	})
}

// relayResponse is the function that query and response
// please call this function last one.
func relayResponse(c *gin.Context, url string, options *grequests.RequestOptions) {
	resp, err := grequests.Get(url, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), resp.Bytes())
}

func handleQueryRequest(c *gin.Context) {
	requestID := c.Param("requestID")
	relayResponse(c, fmt.Sprintf(`%s/zoracle/request/%s`, queryURI, requestID), nil)
}

func handleQueryTx(c *gin.Context) {
	txHash := c.Param("txHash")
	relayResponse(c, fmt.Sprintf(`%s/txs/%s`, queryURI, txHash), nil)
}

func handleQueryProof(c *gin.Context) {
	requestID := c.Param("requestID")
	relayResponse(c, fmt.Sprintf(`%s/d3n/proof/%s`, queryURI, requestID), nil)
}

func main() {
	privBytes, _ := hex.DecodeString(priv)
	copy(pk[:], privBytes)

	var err error
	bandClient, err = d3nlib.NewBandStatefulClient(nodeURI, pk)
	if err != nil {
		panic(err)
	}
	cdc = app.MakeCodec()
	rpcClient = rpc.NewHTTP(nodeURI, "/websocket")

	r := gin.Default()
	// Currently gin-contrib/cors not work so add header manually
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	})

	r.POST("/request", handleRequestData)
	r.POST("/execute", handleExecute)

	r.GET("/request/:requestID", handleQueryRequest)
	r.GET("/txs/:txHash", handleQueryTx)
	r.GET("/proof/:requestID", handleQueryProof)

	r.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
