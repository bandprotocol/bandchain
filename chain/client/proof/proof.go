package proof

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gorilla/mux"
	"github.com/tendermint/iavl"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

var (
	relayArguments  abi.Arguments
	verifyArguments abi.Arguments
)

const (
	RequestIDTag = "requestID"
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

type OracleDataProof struct {
	RequestPacket  oracle.OracleRequestPacketData  `json:"requestPacket"`
	ResponsePacket oracle.OracleResponsePacketData `json:"responsePacket"`
	Version        uint64                          `json:"version"`
	MerklePaths    []IAVLMerklePath                `json:"merklePaths"`
}

func (o *OracleDataProof) encodeToEthData(blockHeight uint64) ([]byte, error) {
	parsePaths := make([]IAVLMerklePathEthereum, len(o.MerklePaths))
	for i, path := range o.MerklePaths {
		parsePaths[i] = path.encodeToEthFormat()
	}
	return verifyArguments.Pack(
		big.NewInt(int64(blockHeight)),
		transformRequestPacket(o.RequestPacket),
		transformResponsePacket(o.ResponsePacket),
		big.NewInt(int64(o.Version)),
		parsePaths,
	)
}

type JsonProof struct {
	BlockHeight     uint64          `json:"blockHeight"`
	OracleDataProof OracleDataProof `json:"oracleDataProof"`
	BlockRelayProof BlockRelayProof `json:"blockRelayProof"`
}

type Proof struct {
	JsonProof     JsonProof        `json:"jsonProof"`
	EVMProofBytes tmbytes.HexBytes `json:"evmProofBytes"`
}

func GetProofHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		intRequestID, err := strconv.ParseUint(vars[RequestIDTag], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		requestID := otypes.RequestID(intRequestID)

		commit, err := cliCtx.Client.Commit(nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		resp, err := cliCtx.Client.ABCIQueryWithOptions(
			"/store/oracle/key",
			otypes.ResultStoreKey(requestID),
			rpcclient.ABCIQueryOptions{Height: commit.Height - 1, Prove: true},
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

		var iavlProof iavl.ValueOp
		var multiStoreProof rootmulti.MultiStoreProofOp
		for _, op := range ops {
			opType := op.GetType()
			if opType == "iavl:v" {
				err := cliCtx.Codec.UnmarshalBinaryLengthPrefixed(op.GetData(), &iavlProof)
				if err != nil {
					rest.WriteErrorResponse(w, http.StatusInternalServerError,
						fmt.Sprintf("iavl: %s", err.Error()),
					)
					return
				}
			} else if opType == "multistore" {
				mp, err := rootmulti.MultiStoreProofOpDecoder(op)
				multiStoreProof = mp.(rootmulti.MultiStoreProofOp)
				if err != nil {
					rest.WriteErrorResponse(w, http.StatusInternalServerError,
						fmt.Sprintf("multiStore: %s", err.Error()),
					)
					return
				}
			}
		}

		eventHeight := iavlProof.Proof.Leaves[0].Version
		signatures, prefix, err := GetSignaturesAndPrefix(&commit.SignedHeader)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		blockRelay := BlockRelayProof{
			MultiStoreProof:        GetMultiStoreProof(multiStoreProof),
			BlockHeaderMerkleParts: GetBlockHeaderMerkleParts(cliCtx.Codec, commit.Header),
			Signatures:             signatures,
			SignedDataPrefix:       prefix,
		}
		resValue := resp.Response.GetValue()

		type result struct {
			oracle.OracleRequestPacketData
			oracle.OracleResponsePacketData
		}
		var rs result
		obi.MustDecode(resValue, &rs)
		reqPacket := otypes.NewOracleRequestPacketData(
			rs.OracleRequestPacketData.ClientID, rs.OracleRequestPacketData.OracleScriptID, rs.OracleRequestPacketData.Calldata,
			rs.OracleRequestPacketData.AskCount, rs.OracleRequestPacketData.MinCount)
		resPacket := otypes.NewOracleResponsePacketData(
			rs.OracleResponsePacketData.ClientID, rs.OracleResponsePacketData.RequestID, rs.OracleResponsePacketData.AnsCount,
			rs.OracleResponsePacketData.RequestTime, rs.OracleResponsePacketData.ResolveTime, rs.OracleResponsePacketData.ResolveStatus,
			rs.OracleResponsePacketData.Result)

		oracleData := OracleDataProof{
			RequestPacket:  reqPacket,
			ResponsePacket: resPacket,
			Version:        uint64(eventHeight),
			MerklePaths:    GetIAVLMerklePaths(&iavlProof),
		}

		// Calculate byte for proofbytes
		var relayAndVerifyArguments abi.Arguments
		format := `[{"type":"bytes"},{"type":"bytes"}]`
		err = json.Unmarshal([]byte(format), &relayAndVerifyArguments)
		if err != nil {
			panic(err)
		}

		blockRelayBytes, err := blockRelay.encodeToEthData(uint64(commit.Height))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		oracleDataBytes, err := oracleData.encodeToEthData(uint64(commit.Height))
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
				BlockHeight:     uint64(commit.Height),
				OracleDataProof: oracleData,
				BlockRelayProof: blockRelay,
			},
			EVMProofBytes: evmProofBytes,
		})
	}
}
