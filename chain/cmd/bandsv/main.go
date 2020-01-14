package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	app "github.com/bandprotocol/d3n/chain"
	"github.com/bandprotocol/d3n/chain/cmtx"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

type OracleRequest struct {
	CodeHash []byte `json:"codeHash" binding:"len=0|len=32"`
	Code     []byte `json:"code"`
	Params   []byte `json:"params" binding:"required"`
}

type OracleRequestResp struct {
	RequestId uint64       `json:"id"`
	CodeHash  cmn.HexBytes `json:"codeHash"`
}

type OracleInfoResp struct {
	Request zoracle.RequestInfo `json:"request"`
	Proof   Proof               `json:"proof"`
}

type IAVLMerklePath struct {
	SubtreeHeight  uint8        `json:"subtreeHeight"`
	SubtreeSize    uint64       `json:"subtreeSize"`
	SubtreeVersion uint64       `json:"subtreeVersion"`
	IsDataOnRight  bool         `json:"isDataOnRight"`
	SiblingHash    cmn.HexBytes `json:"siblingHash"`
}

type BlockHeaderMerkleParts struct {
	VersionAndChainIdHash       cmn.HexBytes `json:"versionAndChainIdHash"`
	TimeHash                    cmn.HexBytes `json:"timeHash"`
	TxCountAndLastBlockInfoHash cmn.HexBytes `json:"txCountAndLastBlockInfoHash"`
	ConsensusDataHash           cmn.HexBytes `json:"consensusDataHash"`
	LastResultsHash             cmn.HexBytes `json:"lastResultsHash"`
	EvidenceAndProposerHash     cmn.HexBytes `json:"evidenceAndProposerHash"`
}

type BlockRelayProof struct {
	OracleIAVLStateHash    cmn.HexBytes           `json:"oracleIAVLStateHash"`
	OtherStoresMerkleHash  cmn.HexBytes           `json:"otherStoresMerkleHash"`
	SupplyStoresMerkleHash cmn.HexBytes           `json:"supplyStoresMerkleHash"`
	BlockHeaderMerkleParts BlockHeaderMerkleParts `json:"blockHeaderMerkleParts"`
	SignedDataPrefix       cmn.HexBytes           `json:"signedDataPrefix"`
	Signatures             []TMSignature          `json:"signatures"`
}

type OracleDataProof struct {
	Version     uint64           `json:"version"`
	RequestId   uint64           `json:"requestId"`
	CodeHash    cmn.HexBytes     `json:"codeHash"`
	Params      cmn.HexBytes     `json:"params"`
	Data        cmn.HexBytes     `json:"data"`
	MerklePaths []IAVLMerklePath `json:"merklePaths"`
}
type TMSignature struct {
	R                cmn.HexBytes `json:"r"`
	S                cmn.HexBytes `json:"s"`
	V                uint8        `json:"v"`
	SignedDataSuffix cmn.HexBytes `json:"signedDataSuffix"`
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
var cliCtx context.CLIContext
var cdc *codec.Codec

func cdcEncode(item interface{}) []byte {
	if item != nil && !cmn.IsTypedNil(item) && !cmn.IsEmpty(item) {
		return cdc.MustMarshalBinaryBare(item)
	}
	return nil
}

func RecoverETHAddress(msg, sig, signer []byte) ([]byte, uint8, error) {
	for i := uint8(0); i < 2; i++ {
		pubuc, err := crypto.SigToPub(tmhash.Sum(msg), append(sig, byte(i)))
		if err != nil {
			return nil, 0, err
		}
		pub := crypto.CompressPubkey(pubuc)
		var tmp [33]byte

		copy(tmp[:], pub)
		if string(signer) == string(secp256k1.PubKeySecp256k1(tmp).Address()) {
			return crypto.PubkeyToAddress(*pubuc).Bytes(), 27 + i, nil
		}
	}
	return nil, 0, fmt.Errorf("No match address found")
}

func GetBlockRelayProof(blockId uint64) (BlockRelayProof, error) {
	bp := BlockRelayProof{}
	_blockId := int64(blockId)
	rc, err := rpcClient.Commit(&_blockId)
	if err != nil {
		return BlockRelayProof{}, err
	}
	sh := rc.SignedHeader

	block := *sh.Header
	commit := *sh.Commit

	bhmp := BlockHeaderMerkleParts{}
	bhmp.VersionAndChainIdHash = merkle.SimpleHashFromByteSlices([][]byte{
		cdcEncode(block.Version),
		cdcEncode(block.ChainID),
	})
	bhmp.TimeHash = merkle.SimpleHashFromByteSlices([][]byte{
		cdcEncode(block.Time),
	})
	bhmp.TxCountAndLastBlockInfoHash = merkle.SimpleHashFromByteSlices([][]byte{
		cdcEncode(block.NumTxs),
		cdcEncode(block.TotalTxs),
		cdcEncode(block.LastBlockID),
		cdcEncode(block.LastCommitHash),
	})
	bhmp.ConsensusDataHash = merkle.SimpleHashFromByteSlices([][]byte{
		cdcEncode(block.DataHash),
		cdcEncode(block.ValidatorsHash),
		cdcEncode(block.NextValidatorsHash),
		cdcEncode(block.ConsensusHash),
	})
	bhmp.LastResultsHash = merkle.SimpleHashFromByteSlices([][]byte{
		cdcEncode(block.LastResultsHash),
	})
	bhmp.EvidenceAndProposerHash = merkle.SimpleHashFromByteSlices([][]byte{
		cdcEncode(block.EvidenceHash),
		cdcEncode(block.ProposerAddress),
	})

	bp.BlockHeaderMerkleParts = bhmp

	leftMsgCount := map[string]int{}
	bp.Signatures = []TMSignature{}
	addrs := []string{}
	mapAddrs := map[string]TMSignature{}
	for i := 0; i < len(commit.Precommits); i++ {
		if commit.Precommits[i] == nil {
			continue
		}

		vote := types.Vote(*commit.Precommits[i])

		msg := vote.SignBytes("bandchain")
		lr := strings.Split(hex.EncodeToString(msg), hex.EncodeToString(block.Hash()))

		val, ok := leftMsgCount[lr[0]]
		if ok {
			leftMsgCount[lr[0]] = val + 1
		} else {
			leftMsgCount[lr[0]] = 0
		}

		lr1, err := hex.DecodeString(lr[1])
		if err != nil {
			continue
		}

		addr, v, err := RecoverETHAddress(msg, vote.Signature, vote.ValidatorAddress)
		if err != nil {
			continue
		}
		addrs = append(addrs, string(addr))
		mapAddrs[string(addr)] = TMSignature{
			vote.Signature[:32],
			vote.Signature[32:],
			v,
			lr1,
		}
		// bp.Signatures = append(bp.Signatures, TMSignature{
		// 	vote.Signature[:32],
		// 	vote.Signature[32:],
		// 	v,
		// 	lr1,
		// })
	}
	if len(addrs) < 4 {
		return BlockRelayProof{}, fmt.Errorf("Too many invalid precommits")
	}

	sort.Strings(addrs)
	for _, addr := range addrs {
		bp.Signatures = append(bp.Signatures, mapAddrs[addr])
	}

	maxCount := 0
	for k, count := range leftMsgCount {
		if count > maxCount {
			kb, err := hex.DecodeString(k)
			if err != nil {
				return BlockRelayProof{}, err
			}
			bp.SignedDataPrefix = kb
			maxCount = count
		}
	}

	if err != nil {
		return BlockRelayProof{}, err
	}

	return bp, nil
}

func MakeOtherStoresMerkleHash(mspo rootmulti.MultiStoreProofOp) (cmn.HexBytes, cmn.HexBytes, cmn.HexBytes) {
	var zoracleHash []byte
	m := map[string][]byte{}
	for _, si := range mspo.Proof.StoreInfos {
		m[si.Name] = tmhash.Sum(tmhash.Sum(si.Core.CommitID.Hash))
		if si.Name == "zoracle" {
			zoracleHash = si.Core.CommitID.Hash
		}
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
				tmhash.Sum(append([]byte{0}, append(append([]byte{3}, []byte("acc")...), append([]byte{32}, m["acc"]...)...)...)),
				tmhash.Sum(append([]byte{0}, append(append([]byte{12}, []byte("distribution")...), append([]byte{32}, m["distribution"]...)...)...))...,
			)...,
		),
	)

	h2 := tmhash.Sum(
		append(
			[]byte{1},
			append(
				tmhash.Sum(append([]byte{0}, append(append([]byte{3}, []byte("gov")...), append([]byte{32}, m["gov"]...)...)...)),
				tmhash.Sum(append([]byte{0}, append(append([]byte{4}, []byte("main")...), append([]byte{32}, m["main"]...)...)...))...,
			)...,
		),
	)

	h3 := tmhash.Sum(
		append(
			[]byte{1},
			append(
				tmhash.Sum(append([]byte{0}, append(append([]byte{4}, []byte("mint")...), append([]byte{32}, m["mint"]...)...)...)),
				tmhash.Sum(append([]byte{0}, append(append([]byte{6}, []byte("params")...), append([]byte{32}, m["params"]...)...)...))...,
			)...,
		),
	)

	h4 := tmhash.Sum(
		append(
			[]byte{1},
			append(
				tmhash.Sum(append([]byte{0}, append(append([]byte{8}, []byte("slashing")...), append([]byte{32}, m["slashing"]...)...)...)),
				tmhash.Sum(append([]byte{0}, append(append([]byte{7}, []byte("staking")...), append([]byte{32}, m["staking"]...)...)...))...,
			)...,
		),
	)

	h5 := tmhash.Sum(append([]byte{0}, append(append([]byte{6}, []byte("supply")...), append([]byte{32}, m["supply"]...)...)...))

	h6 := tmhash.Sum(append([]byte{1}, append(h1, h2...)...))

	h7 := tmhash.Sum(append([]byte{1}, append(h3, h4...)...))

	h8 := tmhash.Sum(append([]byte{1}, append(h6, h7...)...))

	return h5, h8, zoracleHash
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

	var rwr zoracle.RequestInfo
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/zoracle/request/%d", requestId), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = json.Unmarshal(res, &rwr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	key := zoracle.ResultStoreKey(requestId, rwr.CodeHash, rwr.Params)

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

	height := uint64(resp.Response.Height) + 1
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
	odp.CodeHash = rwr.CodeHash
	odp.Params = rwr.ParamsHex
	odp.Version = uint64(opiavl.Proof.Leaves[0].Version)

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

	ssmh, osmh, oiavlsh := MakeOtherStoresMerkleHash(opms)
	var brp BlockRelayProof
	for i := 0; i < 50; i++ {
		brp, err = GetBlockRelayProof(height)
		if err != nil {
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	brp.OtherStoresMerkleHash = osmh
	brp.SupplyStoresMerkleHash = ssmh
	brp.OracleIAVLStateHash = oiavlsh

	c.JSON(200, Proof{
		BlockHeight:     height,
		OracleDataProof: odp,
		BlockRelayProof: brp,
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
	r.GET("/request/:id", handleGetRequest)

	r.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
