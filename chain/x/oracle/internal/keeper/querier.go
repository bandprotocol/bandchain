package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandx/oracle/x/oracle/internal/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryRequest:
			return queryRequest(ctx, path[1:], req, keeper)
		case types.QueryPending:
			return queryPending(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

// queryRequest is a query function to get request by id
func queryRequest(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) == 0 {
		return nil, sdk.ErrInternal("must specify the requestid")
	}
	reqID, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("wrong format for requestid %s", err.Error()))
	}
	requestFromState, sdkErr := keeper.GetRequest(ctx, reqID)
	if err != nil {
		return nil, sdkErr
	}
	reports := keeper.GetValidatorReports(ctx, reqID)

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.NewRequestWithReport(requestFromState, reports))
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return res, nil
}

// queryPending is a query function to get request on pending period
func queryPending(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	res, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetPending(ctx))
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return res, nil
}
