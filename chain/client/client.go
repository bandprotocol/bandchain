package rpc

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
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

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/d3n/blocks/latest", LatestBlocksRequestHandlerFn(cliCtx)).Methods("GET")
}
