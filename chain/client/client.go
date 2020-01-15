package rpc

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	requestID = "requestID"
)

func getLatestBlocks(cliCtx context.CLIContext, r *http.Request) ([]byte, error) {
	_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 10)
	if err != nil {
		return nil, err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	status, err := node.Status()
	if err != nil {
		return nil, err
	}

	blockHeight := status.SyncInfo.LatestBlockHeight
	res := make([]*ctypes.ResultBlock, 0)

	for i := (page - 1) * limit; i < page*limit && i < int(blockHeight); i++ {
		height := blockHeight - int64(i)
		block, err := node.Block(&height)
		if err != nil {
			return nil, err
		}
		res = append(res, block)
	}

	if cliCtx.Indent {
		return codec.Cdc.MarshalJSONIndent(res, "", "  ")
	}

	return codec.Cdc.MarshalJSON(res)
}

func LatestBlocksRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := getLatestBlocks(cliCtx, r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}

func LatestTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		_, page, limit, err := rest.ParseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// TODO: (1) Sort result in desc order after tendermint/tendermint:#4253 is released
		// TODO: (2) Perform binary search on 'tx.height>?' to optimize the performance

		// Temporary implementation to get latest tx sort by descending timestamp
		// Pull request at bandprotocol/d3n:#224
		searchResult, err := utils.QueryTxsByEvents(cliCtx, []string{"tx.height>0"}, 1, 1000000)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		noTx := searchResult.TotalCount
		startIdx := noTx - (page-1)*limit - 1 // Latest tx in page
		endIdx := noTx - page*limit           // Oldest tx in page

		var result sdk.SearchTxsResult
		result.PageNumber = page
		result.TotalCount = noTx
		result.Limit = limit
		result.Txs = make([]sdk.TxResponse, 0)
		for ; startIdx >= endIdx && startIdx >= 0; startIdx-- {
			result.Txs = append(result.Txs, searchResult.Txs[startIdx])
		}
		result.Count = len(result.Txs)
		result.PageTotal = (noTx + limit - 1) / limit

		rest.PostProcessResponseBare(w, cliCtx, result)
	}
}

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/d3n/blocks/latest", LatestBlocksRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/d3n/txs/latest", LatestTxsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/d3n/proof/{%s}", requestID), GetProofHandlerFn(cliCtx)).Methods("GET")
}
