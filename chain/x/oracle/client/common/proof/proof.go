package proof

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gorilla/mux"
	"github.com/tendermint/iavl"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	clientcmn "github.com/bandprotocol/bandchain/chain/x/oracle/client/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
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
		parseSignatures,
	)
}

type OracleDataProof struct {
	RequestPacket  types.OracleRequestPacketData  `json:"requestPacket"`
	ResponsePacket types.OracleResponsePacketData `json:"responsePacket"`
	Version        uint64                         `json:"version"`
	MerklePaths    []IAVLMerklePath               `json:"merklePaths"`
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

type JsonMultiProof struct {
	BlockHeight          uint64            `json:"blockHeight"`
	OracleDataMultiProof []OracleDataProof `json:"oracleDataMultiProof"`
	BlockRelayProof      BlockRelayProof   `json:"blockRelayProof"`
}

type Proof struct {
	JsonProof     JsonProof        `json:"jsonProof"`
	EVMProofBytes tmbytes.HexBytes `json:"evmProofBytes"`
}

type MultiProof struct {
	JsonProof     JsonMultiProof   `json:"jsonProof"`
	EVMProofBytes tmbytes.HexBytes `json:"evmProofBytes"`
}

func GetProofHandlerFn(ctx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, ctx, r)
		if !ok {
			return
		}
		height := &cliCtx.Height
		if cliCtx.Height == 0 {
			height = nil
		}

		vars := mux.Vars(r)
		intRequestID, err := strconv.ParseUint(vars[RequestIDTag], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		requestID := types.RequestID(intRequestID)
		bz, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%d", route, types.QueryRequests, requestID))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var qResult types.QueryResult
		if err := json.Unmarshal(bz, &qResult); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if qResult.Status != http.StatusOK {
			clientcmn.PostProcessQueryResponse(w, cliCtx, bz)
			return
		}
		var request types.QueryRequestResult
		if err := cliCtx.Codec.UnmarshalJSON(qResult.Result, &request); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if request.Result == nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "Result has not been resolved")
			return
		}

		commit, err := cliCtx.Client.Commit(height)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		resp, err := cliCtx.Client.ABCIQueryWithOptions(
			"/store/oracle/key",
			types.ResultStoreKey(requestID),
			rpcclient.ABCIQueryOptions{Height: commit.Height - 1, Prove: true},
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		proof := resp.Response.GetProof()
		if proof == nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "Proof not found")
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
		if iavlProof.Proof == nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "Proof has not been ready.")
			return
		}
		eventHeight := iavlProof.Proof.Leaves[0].Version
		signatures, err := GetSignaturesAndPrefix(&commit.SignedHeader)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		blockRelay := BlockRelayProof{
			MultiStoreProof:        GetMultiStoreProof(multiStoreProof),
			BlockHeaderMerkleParts: GetBlockHeaderMerkleParts(cliCtx.Codec, commit.Header),
			Signatures:             signatures,
		}
		resValue := resp.Response.GetValue()

		type result struct {
			Req types.OracleRequestPacketData
			Res types.OracleResponsePacketData
		}
		var rs result
		obi.MustDecode(resValue, &rs)

		oracleData := OracleDataProof{
			RequestPacket:  rs.Req,
			ResponsePacket: rs.Res,
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

func GetMutiProofHandlerFn(ctx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, ctx, r)
		if !ok {
			return
		}
		height := &cliCtx.Height
		if cliCtx.Height == 0 {
			height = nil
		}

		vars := mux.Vars(r)
		commitHeight := uint64(0)
		requestIDs := strings.Split(vars[RequestIDTag], ",")
		blockRelay := BlockRelayProof{}
		blockRelayBytes := []byte{}
		arrayOfOracleDataBytes := [][]byte{}
		arrayOfOracleData := []OracleDataProof{}

		for _, requestID := range requestIDs {
			intRequestID, err := strconv.ParseUint(requestID, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			requestID := types.RequestID(intRequestID)
			bz, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%d", route, types.QueryRequests, requestID))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			var qResult types.QueryResult
			if err := json.Unmarshal(bz, &qResult); err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			if qResult.Status != http.StatusOK {
				clientcmn.PostProcessQueryResponse(w, cliCtx, bz)
				return
			}
			var request types.QueryRequestResult
			if err := cliCtx.Codec.UnmarshalJSON(qResult.Result, &request); err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			if request.Result == nil {
				rest.WriteErrorResponse(w, http.StatusNotFound, "Result has not been resolved")
				return
			}

			commit, err := cliCtx.Client.Commit(height)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			commitHeight = uint64(commit.Height)

			resp, err := cliCtx.Client.ABCIQueryWithOptions(
				"/store/oracle/key",
				types.ResultStoreKey(requestID),
				rpcclient.ABCIQueryOptions{Height: commit.Height - 1, Prove: true},
			)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			proof := resp.Response.GetProof()
			if proof == nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, "Proof not found")
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
			if iavlProof.Proof == nil {
				rest.WriteErrorResponse(w, http.StatusNotFound, "Proof has not been ready.")
				return
			}
			eventHeight := iavlProof.Proof.Leaves[0].Version
			signatures, err := GetSignaturesAndPrefix(&commit.SignedHeader)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			blockRelay := BlockRelayProof{
				MultiStoreProof:        GetMultiStoreProof(multiStoreProof),
				BlockHeaderMerkleParts: GetBlockHeaderMerkleParts(cliCtx.Codec, commit.Header),
				Signatures:             signatures,
			}
			resValue := resp.Response.GetValue()

			type result struct {
				Req types.OracleRequestPacketData
				Res types.OracleResponsePacketData
			}
			var rs result
			obi.MustDecode(resValue, &rs)

			oracleData := OracleDataProof{
				RequestPacket:  rs.Req,
				ResponsePacket: rs.Res,
				Version:        uint64(eventHeight),
				MerklePaths:    GetIAVLMerklePaths(&iavlProof),
			}

			blockRelayBytes, err = blockRelay.encodeToEthData(uint64(commit.Height))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			oracleDataBytes, err := oracleData.encodeToEthData(uint64(commit.Height))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			arrayOfOracleDataBytes = append(arrayOfOracleDataBytes, oracleDataBytes)
			arrayOfOracleData = append(arrayOfOracleData, oracleData)
		}

		// Calculate byte for MultiProofbytes
		var relayAndVerifyArguments abi.Arguments
		format := `[{"type":"bytes"},{"type":"bytes[]"}]`
		err := json.Unmarshal([]byte(format), &relayAndVerifyArguments)
		if err != nil {
			panic(err)
		}

		evmProofBytes, err := relayAndVerifyArguments.Pack(blockRelayBytes, arrayOfOracleDataBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, MultiProof{
			JsonProof: JsonMultiProof{
				BlockHeight:          commitHeight,
				OracleDataMultiProof: arrayOfOracleData,
				BlockRelayProof:      blockRelay,
			},
			EVMProofBytes: evmProofBytes,
		})
	}
}
