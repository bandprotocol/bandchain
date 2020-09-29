package common

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client/context"
)

func QuerySearchLatestResult(
	route string, cliCtx context.CLIContext, oid, calldata, askCount, minCount string, limit int,
) ([]byte, int64, error) {
	requests, h, err := queryRequestsByLatestTxs(route, cliCtx, oid, calldata, askCount, minCount, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(requests) == 0 {
		bz, err := types.QueryNotFound("request with specified specification not found")
		return bz, 0, err
	}
	bz, err := types.QueryOK(requests)
	return bz, h, err
}
