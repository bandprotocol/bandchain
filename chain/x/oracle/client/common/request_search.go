package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func queryRequest(route string, cliCtx context.CLIContext, rid string) ([]byte, int64, error) {
	res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryRequests, rid), nil)
	if err != nil {
		return nil, 0, err
	}
	return res, height, nil
}

func QuerySearchLatestRequest(
	route string, cliCtx context.CLIContext, oid, calldata, askCount, minCount string,
) ([]byte, int64, error) {

	events := []string{
		fmt.Sprintf("%s.%s='%s'", types.EventTypeRequest, types.AttributeKeyOracleScriptID, oid),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeRequest, types.AttributeKeyCalldata, calldata),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeRequest, types.AttributeKeyAskCount, askCount),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeRequest, types.AttributeKeyMinCount, minCount),
	}
	searchResult, err := authclient.QueryTxsByEvents(cliCtx, events, 1, 30, "desc")
	if err != nil {
		return nil, 0, err
	}
	for _, tx := range searchResult.Txs {
		for _, log := range tx.Logs {
			for _, ev := range log.Events {
				if ev.Type != types.EventTypeRequest {
					continue
				}
				rid := ""
				ok := true
				for _, attr := range ev.Attributes {
					if attr.Key == types.AttributeKeyID {
						rid = attr.Value
					}
					if attr.Key == types.AttributeKeyOracleScriptID && attr.Value != oid ||
						attr.Key == types.AttributeKeyCalldata && attr.Value != calldata ||
						attr.Key == types.AttributeKeyAskCount && attr.Value != askCount ||
						attr.Key == types.AttributeKeyMinCount && attr.Value != minCount {
						ok = false
						break
					}
				}
				if ok && rid != "" {
					res, h, err := queryRequest(route, cliCtx, rid)
					if err != nil {
						return nil, 0, err
					}
					var out types.QueryRequestResult
					cliCtx.Codec.MustUnmarshalJSON(res, &out)
					if out.Result != nil {
						return res, h, nil
					}
				}
			}
		}
	}
	return nil, 0, fmt.Errorf("request not found")
}
