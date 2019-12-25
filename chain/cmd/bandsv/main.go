package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/bandprotocol/d3n/chain/cmtx"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
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
}

type OracleInfoResp struct {
	Request zoracle.RequestWithReport `json:"request"`
	Proof   Proof                     `json:"proof"`
}

type IAVLMerklePath struct {
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
	Version     uint64           `json:"version"`
	RequestId   uint64           `json:"requestId"`
	CodeHash    []byte           `json:"codeHash"`
	Params      []byte           `json:"params"`
	Data        []byte           `json:"data"`
	MerklePaths []IAVLMerklePath `json:"merklePaths"`
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

const priv = "06be35b56b048c5a6810a47e2ef612eaed735ccb0d7ea4fc409f23f1d1a16e0b"
const port = "5001"
const nodeURI = "http://localhost:26657"
const queryURI = "http://localhost:1317"

var rpcClient *rpc.HTTP
var pk secp256k1.PrivKeySecp256k1
var txSender cmtx.TxSender
var cdc = codec.New()

func MakeOtherStoresMerkleHash(mspo rootmulti.MultiStoreProofOp) []byte {
	m := map[string][]byte{}
	for _, si := range mspo.Proof.StoreInfos {
		m[si.Name] = tmhash.Sum(tmhash.Sum(si.Core.CommitID.Hash))
	}

	keys := []string{}
	for k := range m {
		if k != "zoracle" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	bs := [][]byte{}
	for _, k := range keys {
		bs = append(bs, m[k])
	}

	h1 := tmhash.Sum(
		append(
			[]byte{1},
			append(
				tmhash.Sum(append([]byte{0}, append(append([]byte{3}, []byte("acc")...), append([]byte{32}, bs[0]...)...)...)),
				tmhash.Sum(append([]byte{0}, append(append([]byte{12}, []byte("distribution")...), append([]byte{32}, bs[1]...)...)...))...,
			)...,
		),
	)

	h2 := tmhash.Sum(
		append(
			[]byte{1},
			append(
				tmhash.Sum(append([]byte{0}, append(append([]byte{3}, []byte("gov")...), append([]byte{32}, bs[2]...)...)...)),
				tmhash.Sum(append([]byte{0}, append(append([]byte{4}, []byte("main")...), append([]byte{32}, bs[3]...)...)...))...,
			)...,
		),
	)

	h3 := tmhash.Sum(
		append(
			[]byte{1},
			append(
				tmhash.Sum(append([]byte{0}, append(append([]byte{6}, []byte("params")...), append([]byte{32}, bs[4]...)...)...)),
				tmhash.Sum(append([]byte{0}, append(append([]byte{8}, []byte("slashing")...), append([]byte{32}, bs[5]...)...)...))...,
			)...,
		),
	)

	h4 := tmhash.Sum(
		append(
			[]byte{1},
			append(
				tmhash.Sum(append([]byte{0}, append(append([]byte{7}, []byte("staking")...), append([]byte{32}, bs[6]...)...)...)),
				tmhash.Sum(append([]byte{0}, append(append([]byte{6}, []byte("supply")...), append([]byte{32}, bs[7]...)...)...))...,
			)...,
		),
	)

	h5 := tmhash.Sum(append([]byte{1}, append(h1, h2...)...))

	h6 := tmhash.Sum(append([]byte{1}, append(h3, h4...)...))

	h7 := tmhash.Sum(append([]byte{1}, append(h5, h6...)...))

	return h7
}

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

	if len(requestData.Code) > 0 {
		requestData.CodeHash = zoracle.NewStoredCode(requestData.Code, txSender.Sender()).GetCodeHash()
		hasCode, err := HasCode(requestData.CodeHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// If codeHash not found then store the code
		if !hasCode {
			_, err := txSender.SendTransaction(zoracle.NewMsgStoreCode(requestData.Code, txSender.Sender()), flags.BroadcastBlock)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
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
	}

	txr, err := txSender.SendTransaction(zoracle.NewMsgRequest(requestData.CodeHash, requestData.Params, 4, txSender.Sender()), flags.BroadcastBlock)
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
		CodeHash:  requestData.CodeHash,
	})
}

func handleGetRequest(c *gin.Context) {
	requestId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key, err := hex.DecodeString(fmt.Sprintf("01%016x", requestId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := rpcClient.ABCIQueryWithOptions(
		"/store/zoracle/key",
		key,
		rpc.ABCIQueryOptions{Height: 0, Prove: true},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	proof := resp.Response.GetProof()
	if proof == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "proof not found"})
		return
	}

	ops := proof.GetOps()
	if ops == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "proof ops not found"})
		return
	}

	var iavlOpData []byte
	var multistoreOpData []byte
	for _, op := range ops {
		opType := op.GetType()
		if opType == "iavl:v" {
			iavlOpData = op.GetData()
		} else if opType == "multistore" {
			multistoreOpData = op.GetData()
		}
	}
	if iavlOpData == nil || multistoreOpData == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "proof was corrupted"})
		return
	}

	var opiavl iavl.IAVLValueOp
	if iavlOpData == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "proof was corrupted"})
		return
	}
	err = cdc.UnmarshalBinaryLengthPrefixed(iavlOpData, &opiavl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	odp := OracleDataProof{}
	odp.RequestId = requestId
	odp.Data = resp.Response.GetValue()

	for i := len(opiavl.Proof.LeftPath) - 1; i >= 0; i-- {
		path := opiavl.Proof.LeftPath[i]
		imp := IAVLMerklePath{}
		imp.SubtreeHeight = uint8(path.Height)
		imp.SubtreeSize = uint64(path.Size)
		imp.SubtreeVersion = uint64(path.Version)
		if len(path.Right) == 0 {
			imp.SiblingHash = path.Left
			imp.IsDataOnRight = true
		} else {
			imp.SiblingHash = path.Right
			imp.IsDataOnRight = false
		}
		odp.MerklePaths = append(odp.MerklePaths, imp)
	}

	var opms rootmulti.MultiStoreProofOp
	err = cdc.UnmarshalBinaryLengthPrefixed(multistoreOpData, &opms)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	osmh := MakeOtherStoresMerkleHash(opms)
	brp := BlockRelayProof{}
	brp.OtherStoresMerkleHash = osmh

	type T struct {
		types.ResponseQuery
		Proof
	}

	c.JSON(200, T{
		resp.Response,
		Proof{
			BlockHeight:     0,
			OracleDataProof: odp,
			BlockRelayProof: brp,
		},
	})
}

func main() {
	viper.Set("nodeURI", nodeURI)
	privBytes, _ := hex.DecodeString(priv)
	copy(pk[:], privBytes)

	txSender = cmtx.NewTxSender(pk)
	rpcClient = rpc.NewHTTP(nodeURI, "/websocket")

	r := gin.Default()

	r.POST("/request", handleRequestData)
	r.GET("/request/:id", handleGetRequest)

	r.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
