package common

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/client/context"

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
	bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%d", route, types.QueryRequests, rid))
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

func queryMultitRequest(cliCtx context.CLIContext, oid, calldata, askCount, minCount, limit string) ([]types.RequestID, error) {
	bz, _, err := cliCtx.Query(fmt.Sprintf("band/multi_request/%s/%s/%s/%s/%s", oid, calldata, askCount, minCount, limit))
	if err != nil {
		return nil, err
	}
	var reqIDs []types.RequestID
	err = cliCtx.Codec.UnmarshalBinaryBare(bz, &reqIDs)
	if err != nil {
		return nil, err
	}
	return reqIDs, nil
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
	requestIDs, err := queryMultitRequest(cliCtx, oid, calldata, askCount, minCount, fmt.Sprintf("%s", limit))
	if err != nil {
		return nil, 0, err
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
