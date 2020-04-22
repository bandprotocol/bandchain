package rpc

import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/tendermint/iavl"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	types "github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type GetResultResponse struct {
	RequestPacket  types.OracleRequestPacketData  `json:"request_packet"`
	ResponsePacket types.OracleResponsePacketData `json:"response_packet"`
}

func mustParseInt64(b []byte) int64 {
	i64, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		panic(err)
	}
	return i64
}

func GetResult(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rawRID, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		commit, err := cliCtx.Client.Commit(nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rid := oracle.RequestID(rawRID)
		resp, err := cliCtx.Client.ABCIQueryWithOptions(
			"/store/oracle/key",
			oracle.ResultStoreKey(rid),
			rpcclient.ABCIQueryOptions{Height: commit.Height - 1, Prove: true},
		)

		proof := resp.Response.GetProof()
		var iavlOp iavl.ValueOp
		for _, op := range proof.GetOps() {
			switch op.GetType() {
			case "iavl:v":
				cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(op.GetData(), &iavlOp)
			}
		}

		eventHeight := iavlOp.Proof.Leaves[0].Version
		blockResults, err := cliCtx.Client.BlockResults(&eventHeight)

		var reqPacket types.OracleRequestPacketData
		var resPacket types.OracleResponsePacketData
		for _, ev := range blockResults.EndBlockEvents {
			if ev.GetType() != types.EventTypeRequestExecute {
				continue
			}
			for _, kv := range ev.GetAttributes() {
				switch string(kv.Key) {
				case types.AttributeKeyRequestID:
					resPacket.RequestID = oracle.RequestID(mustParseInt64(kv.Value))
				case types.AttributeKeyClientID:
					reqPacket.ClientID = string(kv.Value)
					resPacket.ClientID = string(kv.Value)
				case types.AttributeKeyOracleScriptID:
					reqPacket.OracleScriptID = oracle.OracleScriptID(mustParseInt64(kv.Value))
				case types.AttributeKeyCalldata:
					reqPacket.Calldata = string(kv.Value)
				case types.AttributeKeyAskCount:
					reqPacket.AskCount = mustParseInt64(kv.Value)
				case types.AttributeKeyMinCount:
					reqPacket.MinCount = mustParseInt64(kv.Value)
				case types.AttributeKeyAnsCount:
					resPacket.AnsCount = mustParseInt64(kv.Value)
				case types.AttributeKeyRequestTime:
					resPacket.RequestTime = mustParseInt64(kv.Value)
				case types.AttributeKeyResolveTime:
					resPacket.ResolveTime = mustParseInt64(kv.Value)
				case types.AttributeKeyResolveStatus:
					resPacket.ResolveStatus = oracle.ResolveStatus(mustParseInt64(kv.Value))
				case types.AttributeKeyResult:
					resPacket.Result = string(kv.Value)
				}
			}
		}

		rest.PostProcessResponseBare(w, cliCtx, GetResultResponse{
			RequestPacket:  reqPacket,
			ResponsePacket: resPacket,
		})
	}
}
