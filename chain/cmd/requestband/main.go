package main

import (
	"encoding/hex"
	"net/http"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/bandlib"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/gin-gonic/gin"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Request struct {
	Address string `json:"address" binding:"required"`
	Amount  int64  `json:"amount" binding:"required"`
}
type Response struct {
	TxHash string `json:"txHash"`
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
	priv    = getEnv("PRIVATE_KEY", "27313aa3fd8286b54d5dbe16a4fbbc55c7908e844e37a737997fc2ba74403812")
	chainID = getEnv("CHAIN_ID", "bandchain")
)

var pk secp256k1.PrivKeySecp256k1
var bandClient bandlib.BandStatefulClient

func handleRequest(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	to, _ := sdk.AccAddressFromBech32(req.Address)
	result, err := bandClient.SendTransaction(bank.MsgSend{
		FromAddress: bandClient.Sender(),
		ToAddress:   to,
		Amount:      sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(req.Amount))),
	}, 1000000, "", "")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, Response{
		TxHash: result.TxHash,
	})
}

func main() {
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(config)
	config.Seal()

	privBytes, _ := hex.DecodeString(priv)
	copy(pk[:], privBytes)

	var err error
	bandClient, err = bandlib.NewBandStatefulClient(nodeURI, pk, 100, 10, "Band faucet", chainID)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	// Currently gin-contrib/cors not work so add header manually
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	})

	r.POST("/request", handleRequest)

	r.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
