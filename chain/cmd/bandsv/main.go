package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

type OracleRequest struct {
	CodeHash []byte `json:"codeHash" binding:"len=0|len=32"`
	Code     []byte `json:"code"`
	Params   []byte `json:"params" binding:"required"`
}

type OracleRequestResp struct {
	RequestId uint64 `json:"id"`
	CodeHash  []byte `json:"codeHash"`
}

type OracleInfoResp struct {
	Request zoracle.RequestWithReport `json:"request"`
	Proof   Proof                     `json:"proof"`
}

type IAVLPath struct {
	SubtreeHeight  uint8  `json:"subtreeHeight"`
	SubtreeSize    uint64 `json:"subtreeSize"`
	SubtreeVersion uint64 `json:"subtreeVersion"`
	IsDataOnRight  bool   `json:"isDataOnRight"`
	SiblingHash    []byte `json:"siblingHash"`
}

type BlockHeaderMerkleParts struct {
	VersionAndChainIdHash       []byte `json:"versionAndChainIdHash"`
	TimeHash                    []byte `json:"timeHash"`
	TxCountAndLastBlockInfoHash []byte `json:"txCountAndLastBlockInfoHash"`
	ConsensusDataHash           []byte `json:"consensusDataHash"`
	LastResultsHash             []byte `json:"lastResultsHash"`
	EvidenceAndProposerHash     []byte `json:"evidenceAndProposerHash"`
}

type BlockRelayProof struct {
	OracleIAVLStateHash    []byte                 `json:"oracleIAVLStateHash"`
	OtherStoresMerkleHash  []byte                 `json:"otherStoresMerkleHash"`
	BlockHeaderMerkleParts BlockHeaderMerkleParts `json:"blockHeaderMerkleParts"`
	SignedDataPrefix       []byte                 `json:"signedDataPrefix"`
	Signatures             []TMSignature          `json:"signatures"`
}

type OracleDataProof struct {
	Version     uint64     `json:"version"`
	RequestId   uint64     `json:"requestId"`
	CodeHash    []byte     `json:"codeHash"`
	Params      []byte     `json:"params"`
	Value       []byte     `json:"value"`
	MerklePaths []IAVLPath `json:"merklePaths"`
}
type TMSignature struct {
	R                []byte `json:"r"`
	S                []byte `json:"s"`
	V                uint8  `json:"v"`
	SignedDataSuffix []byte `json:"signedDataSuffix"`
}

type Proof struct {
	BlockHeight     uint64          `json:"blockHeight"`
	OracleDataProof OracleDataProof `json:"oracleDataProof"`
	BlockRelayProof BlockRelayProof `json:"blockRelayProof"`
}

func main() {
	r := gin.Default()

	r.POST("/request", func(c *gin.Context) {
		var json OracleRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(json.CodeHash) == 0 && len(json.Code) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code/codeHash"})
			return
		}
		if len(json.CodeHash) > 0 && len(json.Code) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only one of code/codeHash can be sent"})
			return
		}
		if len(json.Code) > 0 {
			// TODO:
			//   1. Compute code hash.
			//   2. Check if code hash already exist on chain. If not deploy it.
			//   3. Save code hash to `json`
		}
		// TODO:
		//   1. Send the request to the blockchain
		//   2. Wait for the tx to confirm and get back the request ID
		c.JSON(200, OracleRequestResp{
			RequestId: 10,
			CodeHash:  nil,
		})
	})

	r.GET("/request/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": c.Param("id"),
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
