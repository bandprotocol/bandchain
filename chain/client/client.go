package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	// unnamed import of statik for swagger UI support
	_ "github.com/bandprotocol/d3n/chain/client/lcd/statik"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

const (
	requestIDTag = "requestID"
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

		// limit maximum must be 100.
		if limit > 100 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "limit must be less than or equal 100.")
			return
		}

		// TODO: (1) Sort result in desc order after tendermint/tendermint:#4253 is released
		// TODO: (2) Perform binary search on 'tx.height>?' to optimize the performance

		// Temporary implementation to get latest tx sort by descending timestamp
		// Pull request at bandprotocol/d3n:#224
		node, _ := cliCtx.GetNode()
		block, _ := node.Block(nil)
		totalTxs := int(block.Block.Header.TotalTxs)
		endIndex := totalTxs - (page-1)*limit
		startIndex := totalTxs - page*limit + 1

		var result sdk.SearchTxsResult
		result.PageNumber = page
		result.TotalCount = totalTxs
		result.Limit = limit
		result.PageTotal = (totalTxs-1)/limit + 1
		result.Txs = make([]sdk.TxResponse, 0)

		if startIndex < 1 {
			startIndex = 1
		}

		if endIndex < 1 {
			result.Count = 0
			rest.PostProcessResponseBare(w, cliCtx, result)
			return
		}

		todoIndex := startIndex
		for todoIndex <= endIndex {
			pageOfTodoIndex := (todoIndex-1)/limit + 1
			startIndexOfPage := (pageOfTodoIndex-1)*limit + 1
			searchResult, err := utils.QueryTxsByEvents(cliCtx, []string{"tx.height>0"}, pageOfTodoIndex, limit)

			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			for i, tx := range searchResult.Txs {
				indexOfTx := startIndexOfPage + i
				if indexOfTx >= startIndex && indexOfTx <= endIndex {
					result.Txs = append(result.Txs, tx)
					todoIndex = indexOfTx + 1
				}
			}
		}

		result.Count = len(result.Txs)
		for i, j := 0, len(result.Txs)-1; i < j; i, j = i+1, j-1 {
			result.Txs[i], result.Txs[j] = result.Txs[j], result.Txs[i]
		}

		rest.PostProcessResponseBare(w, cliCtx, result)
	}
}

func GetHealthStatus(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		node, err := cliCtx.GetNode()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		block, err := node.Block(nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		result := "UP"
		if time.Now().Sub(block.Block.Header.Time) > 3*time.Minute {
			result = "DOWN"
		}
		w.Write([]byte(result))
	}
}

func GetProviderStatus(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqNumberResp, _, err := cliCtx.Query("custom/zoracle/request_number")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var requestID string
		err = json.Unmarshal(reqNumberResp, &requestID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, _, err := cliCtx.Query(fmt.Sprintf("custom/zoracle/request/%s", requestID))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var request zoracle.RequestQuerierInfo
		err = cliCtx.Codec.UnmarshalJSON(res, &request)

		block, err := cliCtx.Client.Block(nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		numReporters := len(request.Reports)
		status := "GOOD"
		// TODO: Remove hard-coded provider count threshold
		if block.Block.Height > request.Request.ExpirationHeight && numReporters < 3 {
			fmt.Printf(`BAD ------- requestId: %s, reports: %d`, requestID, numReporters)
			status = "BAD"
		}

		rest.PostProcessResponseBare(w, cliCtx, struct {
			Height       int64  `json:"height"`
			RequestID    string `json:"id"`
			NumReporters int    `json:"num_reporters"`
			Status       string `json:"status"`
		}{
			Height:       block.Block.Height,
			RequestID:    requestID,
			NumReporters: numReporters,
			Status:       status,
		})
	}
}

func ServeSwaggerUI() http.Handler {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)

	return http.StripPrefix("/swagger-ui/", staticServer)
}

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/d3n/blocks/latest", LatestBlocksRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/d3n/txs/latest", LatestTxsRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/d3n/proof/{%s}", requestIDTag), GetProofHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/d3n/health_check", GetHealthStatus(cliCtx)).Methods("GET")
	r.HandleFunc("/d3n/provider_status", GetProviderStatus(cliCtx)).Methods("GET")
	r.PathPrefix("/swagger-ui/").Handler(ServeSwaggerUI())
}
