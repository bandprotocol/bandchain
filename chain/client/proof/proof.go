package proof

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

var (
	relayArguments  abi.Arguments
	verifyArguments abi.Arguments
)

// TODO: Remove this variable
const (
	requestIDTag = "requestID"
)

func init() {
	err := json.Unmarshal(relayFormat, &relayArguments)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(verifyFormat, &verifyArguments)
	if err != nil {
		panic(err)
	}
}

type IAVLMerklePath struct {
	IsDataOnRight  bool             `json:"isDataOnRight"`
	SubtreeHeight  uint8            `json:"subtreeHeight"`
	SubtreeSize    uint64           `json:"subtreeSize"`
	SubtreeVersion uint64           `json:"subtreeVersion"`
	SiblingHash    tmbytes.HexBytes `json:"siblingHash"`
}

type IAVLMerklePathEthereum struct {
	IsDataOnRight  bool
	SubtreeHeight  uint8
	SubtreeSize    *big.Int
	SubtreeVersion *big.Int
	SiblingHash    common.Hash
}

func (merklePath *IAVLMerklePath) encodeToEthFormat() IAVLMerklePathEthereum {
	return IAVLMerklePathEthereum{
		merklePath.IsDataOnRight,
		merklePath.SubtreeHeight,
		big.NewInt(int64(merklePath.SubtreeSize)),
		big.NewInt(int64(merklePath.SubtreeVersion)),
		common.BytesToHash(merklePath.SiblingHash),
	}
}

type TMSignature struct {
	R                tmbytes.HexBytes `json:"r"`
	S                tmbytes.HexBytes `json:"s"`
	V                uint8            `json:"v"`
	SignedDataSuffix tmbytes.HexBytes `json:"signedDataSuffix"`
}

type TMSignatureEthereum struct {
	R                common.Hash
	S                common.Hash
	V                uint8
	SignedDataSuffix []byte
}

func (signature *TMSignature) encodeToEthFormat() TMSignatureEthereum {
	return TMSignatureEthereum{
		common.BytesToHash(signature.R),
		common.BytesToHash(signature.S),
		signature.V,
		signature.SignedDataSuffix,
	}
}

type BlockRelayProof struct {
	MultiStoreProof        MultiStoreProof        `json:"multiStoreProof"`
	BlockHeaderMerkleParts BlockHeaderMerkleParts `json:"blockHeaderMerkleParts"`
	SignedDataPrefix       tmbytes.HexBytes       `json:"signedDataPrefix"`
	Signatures             []TMSignature          `json:"signatures"`
}

func (blockRelay *BlockRelayProof) encodeToEthData(blockHeight uint64) ([]byte, error) {
	parseSignatures := make([]TMSignatureEthereum, len(blockRelay.Signatures))
	for i, sig := range blockRelay.Signatures {
		parseSignatures[i] = sig.encodeToEthFormat()
	}
	return relayArguments.Pack(
		big.NewInt(int64(blockHeight)),
		blockRelay.MultiStoreProof.encodeToEthFormat(),
		blockRelay.BlockHeaderMerkleParts.encodeToEthFormat(),
		blockRelay.SignedDataPrefix,
		parseSignatures,
	)
}

// TODO: It will be fixed on #1325
type OracleDataProof struct {
	Version        uint64                `json:"version"`
	RequestID      oracle.RequestID      `json:"requestID"`
	OracleScriptID oracle.OracleScriptID `json:"oracleScriptID"`
	Calldata       tmbytes.HexBytes      `json:"calldata"`
	Data           tmbytes.HexBytes      `json:"data"`
	MerklePaths    []IAVLMerklePath      `json:"merklePaths"`
}

func (oracleData *OracleDataProof) encodeToEthData(blockHeight uint64) ([]byte, error) {
	parsePaths := make([]IAVLMerklePathEthereum, len(oracleData.MerklePaths))
	for i, path := range oracleData.MerklePaths {
		parsePaths[i] = path.encodeToEthFormat()
	}
	return verifyArguments.Pack(
		big.NewInt(int64(blockHeight)),
		oracleData.Data,
		uint64(oracleData.RequestID),
		uint64(oracleData.OracleScriptID),
		oracleData.Calldata,
		big.NewInt(int64(oracleData.Version)),
		parsePaths,
	)
}

type Proof struct {
	JsonProof     JsonProof        `json:"jsonProof"`
	EVMProofBytes tmbytes.HexBytes `json:"evmProofBytes"`
}

type JsonProof struct {
	BlockHeight     uint64          `json:"blockHeight"`
	OracleDataProof OracleDataProof `json:"oracleDataProof"`
	BlockRelayProof BlockRelayProof `json:"blockRelayProof"`
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

// Combine MultiStoreProof + BlockHeaderMerkleParts
func GetBlockRelayProof(cliCtx context.CLIContext, blockId uint64) (BlockRelayProof, error) {
	bp := BlockRelayProof{}
	_blockId := int64(blockId)
	rc, err := cliCtx.Client.Commit(&_blockId)
	if err != nil {
		return BlockRelayProof{}, err
	}
	sh := rc.SignedHeader

	block := *sh.Header
	commit := *sh.Commit

	bp.BlockHeaderMerkleParts = GetBlockHeaderMerkleParts(cliCtx.Codec, sh.Header)

	leftMsg := ""
	bp.Signatures = []TMSignature{}
	addrs := []string{}
	mapAddrs := map[string]TMSignature{}
	for i, vote := range commit.Signatures {
		msg := commit.VoteSignBytes(block.ChainID, i)
		lr := strings.Split(hex.EncodeToString(msg), hex.EncodeToString(block.Hash()))

		if leftMsg != "" && leftMsg != lr[0] {
			return BlockRelayProof{}, fmt.Errorf("Inconsistent prefix signature bytes")
		}
		leftMsg = lr[0]

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
	}
	if len(addrs) == 0 {
		return BlockRelayProof{}, fmt.Errorf("Too many invalid precommits")
	}

	sort.Strings(addrs)
	for _, addr := range addrs {
		bp.Signatures = append(bp.Signatures, mapAddrs[addr])
	}

	kb, err := hex.DecodeString(leftMsg)
	if err != nil {
		return BlockRelayProof{}, err
	}
	bp.SignedDataPrefix = kb

	if err != nil {
		return BlockRelayProof{}, err
	}

	return bp, nil
}

// TODO: It will be fixed #1326
func GetProofHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		intRequestID, err := strconv.ParseUint(vars[requestIDTag], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		requestID := oracle.RequestID(intRequestID)

		rc, err := cliCtx.Client.Commit(nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		brp, err := GetBlockRelayProof(cliCtx, uint64(rc.Height))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var queryRequest oracle.RequestQuerierInfo
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/oracle/request/%d", requestID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(res, &queryRequest)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		key := oracle.ResultStoreKey(requestID)

		resp, err := cliCtx.Client.ABCIQueryWithOptions(
			"/store/oracle/key",
			key,
			rpcclient.ABCIQueryOptions{Height: rc.Height - 1, Prove: true},
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		proof := resp.Response.GetProof()
		if proof == nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "proof not found")
			return
		}

		ops := proof.GetOps()
		if ops == nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "proof ops not found")
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
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "proof was corrupted")
			return
		}

		var opiavl iavl.ValueOp
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(iavlOpData, &opiavl)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		odp := OracleDataProof{}
		odp.RequestID = requestID
		odp.Data = resp.Response.GetValue()
		odp.OracleScriptID = queryRequest.Request.OracleScriptID
		odp.Calldata = queryRequest.Request.Calldata
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

		var multiStoreProof rootmulti.MultiStoreProofOp
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(multistoreOpData, &multiStoreProof)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		brp.MultiStoreProof = GetMultiStoreProof(multiStoreProof)

		// Calculate byte for proofbytes
		var relayAndVerifyArguments abi.Arguments
		format := `[{"type":"bytes"},{"type":"bytes"}]`
		err = json.Unmarshal([]byte(format), &relayAndVerifyArguments)
		if err != nil {
			panic(err)
		}

		blockRelayBytes, err := brp.encodeToEthData(uint64(rc.Height))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		oracleDataBytes, err := odp.encodeToEthData(uint64(rc.Height))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		evmProofBytes, err := relayAndVerifyArguments.Pack(blockRelayBytes, oracleDataBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, Proof{
			JsonProof: JsonProof{
				BlockHeight:     uint64(rc.Height),
				OracleDataProof: odp,
				BlockRelayProof: brp},
			EVMProofBytes: evmProofBytes,
		})
	}
}
