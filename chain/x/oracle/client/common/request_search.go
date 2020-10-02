package common

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func queryLatestRequest(cliCtx context.CLIContext, oid, calldata, askCount, minCount string) (types.RequestID, error) {
	bz, _, err := cliCtx.Query(fmt.Sprintf("band/latest_request/%s/%s/%s/%s", oid, calldata, askCount, minCount))
	if err != nil {
		return 0, err
	}
	var reqID types.RequestID
	err = cliCtx.Codec.UnmarshalBinaryBare(bz, &reqID)
	if err != nil {
		return 0, err
	}
	return reqID, nil
}

func queryRequest(route string, cliCtx context.CLIContext, rid types.RequestID) (types.QueryRequestResult, int64, error) {
	bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryRequests, rid))
	if err != nil {
		return types.QueryRequestResult{}, 0, err
	}
	var result types.QueryResult
	if err := json.Unmarshal(bz, &result); err != nil {
		return types.QueryRequestResult{}, 0, err
	}
	var reqResult types.QueryRequestResult
	cliCtx.Codec.MustUnmarshalJSON(result.Result, &reqResult)
	return reqResult, height, nil
}

func QuerySearchLatestRequest(
	route string, cliCtx context.CLIContext, oid, calldata, askCount, minCount string,
) ([]byte, int64, error) {
	id, err := queryLatestRequest(cliCtx, oid, calldata, askCount, minCount)
	if err != nil {
		bz, err := types.QueryNotFound("request with specified specification not found")
		return bz, 0, err
	}
	out, h, err := queryRequest(route, cliCtx, id)
	bz, err := types.QueryOK(out)
	return bz, h, err
}

func queryRequests(
	route string, cliCtx context.CLIContext, requestIDs []types.RequestID,
) ([]types.QueryRequestResult, int64, error) {
	type queryResult struct {
		result types.QueryRequestResult
		err    error
	}
	queryResultsChan := make(chan queryResult, len(requestIDs))
	for _, rid := range requestIDs {
		go func(rid types.RequestID) {
			out, _, err := queryRequest(route, cliCtx, rid)
			if err != nil {
				queryResultsChan <- queryResult{err: err}
				return
			}
			queryResultsChan <- queryResult{result: out}
		}(rid)
	}
	requests := make([]types.QueryRequestResult, 0)
	for idx := 0; idx < len(requestIDs); idx++ {
		select {
		case req := <-queryResultsChan:
			if req.err != nil {
				return nil, 0, req.err
			}
			if req.result.Result != nil {
				requests = append(requests, req.result)
			}
		}
	}

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].Result.ResponsePacketData.ResolveTime > requests[j].Result.ResponsePacketData.ResolveTime
	})

	return requests, cliCtx.Height, nil
}

func QueryMultiSearchLatestRequest(
	route string, cliCtx context.CLIContext, oid, calldata, askCount, minCount string, limit int,
) ([]byte, int64, error) {
	query := fmt.Sprintf("%s.%s='%s' AND %s.%s='%s' AND %s.%s='%s' AND %s.%s='%s'",
		types.EventTypeRequest, types.AttributeKeyOracleScriptID, oid,
		types.EventTypeRequest, types.AttributeKeyCalldata, calldata,
		types.EventTypeRequest, types.AttributeKeyAskCount, askCount,
		types.EventTypeRequest, types.AttributeKeyMinCount, minCount,
	)
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, 0, err
	}
	resTxs, err := node.TxSearch(query, !cliCtx.TrustNode, 1, 30+limit, "desc")
	if err != nil {
		return nil, 0, err
	}
	requestIDs := make([]types.RequestID, 0)
	for _, tx := range resTxs.Txs {
		if !cliCtx.TrustNode {
			err := utils.ValidateTxResult(cliCtx, tx)
			if err != nil {
				return nil, 0, err
			}
		}
		logs, _ := sdk.ParseABCILogs(tx.TxResult.Log)
		for _, log := range logs {
			for _, ev := range log.Events {
				if ev.Type != types.EventTypeRequest {
					continue
				}
				rid := types.RequestID(0)
				ok := true
				for _, attr := range ev.Attributes {
					if attr.Key == types.AttributeKeyID {
						id, err := strconv.ParseUint(attr.Value, 10, 64)
						if err != nil {
							return nil, 0, err
						}
						rid = types.RequestID(id)
					}
					if attr.Key == types.AttributeKeyOracleScriptID && attr.Value != oid ||
						attr.Key == types.AttributeKeyCalldata && attr.Value != calldata ||
						attr.Key == types.AttributeKeyAskCount && attr.Value != askCount ||
						attr.Key == types.AttributeKeyMinCount && attr.Value != minCount {
						ok = false
						break
					}
				}
				if ok && rid != 0 {
					requestIDs = append(requestIDs, rid)
				}
			}
		}
	}
	queryRequestResults, h, err := queryRequests(route, cliCtx, requestIDs)
	if err != nil {
		return nil, 0, err
	}
	if len(queryRequestResults) == 0 {
		bz, err := types.QueryNotFound("request with specified specification not found")
		return bz, 0, err
	}
	if len(queryRequestResults) > limit {
		queryRequestResults = queryRequestResults[:limit]
	}
	bz, err := types.QueryOK(queryRequestResults)
	return bz, h, err
}
