package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

func getRequestHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		reqID := vars[requestID]
		var request RequestQueryInfo

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/request/%s", storeName, reqID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		var queryRequest types.RequestInfo
		err = cliCtx.Codec.UnmarshalJSON(res, &queryRequest)

		request.CodeHash = queryRequest.CodeHash
		request.Params = queryRequest.Params
		request.TargetBlock = int64(queryRequest.TargetBlock)
		request.Result = queryRequest.Result

		// Get request detail
		searchRequest, err := utils.QueryTxsByEvents(
			cliCtx,
			[]string{fmt.Sprintf("%s.%s='%s'", types.EventTypeRequest, types.AttributeKeyRequestID, reqID)},
			1,
			1,
		)

		request.TxHash = searchRequest.Txs[0].TxHash
		request.RequestAtHeight = searchRequest.Txs[0].Height
		request.RequestAtTime = searchRequest.Txs[0].Timestamp
		request.Requester = searchRequest.Txs[0].Tx.GetMsgs()[0].(types.MsgRequest).Sender
		// Script detail
		scriptInfoRaw, _, err := cliCtx.QueryWithData(
			fmt.Sprintf(
				"custom/%s/script/%s",
				storeName,
				hex.EncodeToString(queryRequest.CodeHash),
			),
			nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(scriptInfoRaw, &request.ScriptInfo)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// Save report tx
		reportTxs := make(map[string]ReportDetail)
		searchReports, err := utils.QueryTxsByEvents(
			cliCtx,
			[]string{fmt.Sprintf("%s.%s='%s'", types.EventTypeReport, types.AttributeKeyRequestID, reqID)},
			1,
			30, // Estimated validator reports
		)

		for _, report := range searchReports.Txs {
			reportTxs[report.Logs[0].Events[1].Attributes[1].Value] = ReportDetail{
				TxHash:       report.TxHash,
				ReportAtTime: report.Timestamp,
			}
		}

		request.Reports = make([]ReportDetail, len(queryRequest.Reports))
		for i, report := range queryRequest.Reports {
			reportTx, ok := reportTxs[report.Validator.String()]
			if !ok {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, "Cannot find report tx")
				return
			}
			request.Reports[i] = ReportDetail{
				Reporter:       report.Validator,
				ReportAtHeight: int64(report.ReportAt),
				TxHash:         reportTx.TxHash,
				ReportAtTime:   reportTx.ReportAtTime,
				Value:          report.Value,
			}
		}

		rest.PostProcessResponse(w, cliCtx, request)
	}
}

func getScriptHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars[codeHash]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/script/%s", storeName, hash), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		var scriptInfo ScriptInfoWithTx
		err = cliCtx.Codec.UnmarshalJSON(res, &scriptInfo.Info)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		// Find TxHash and height that of transaction
		// TODO: Get latest store tx as tx hash (wait tendermint release get result in desc order)
		searchResult, err := utils.QueryTxsByEvents(
			cliCtx,
			[]string{fmt.Sprintf("%s.%s='%s'", types.EventTypeStoreCode, types.AttributeKeyCodeHash, hash)},
			1,
			1,
		)
		scriptInfo.TxHash = searchResult.Txs[0].TxHash
		scriptInfo.CreatedAtHeight = searchResult.Txs[0].Height
		scriptInfo.CreatedAtTime = searchResult.Txs[0].Timestamp

		rest.PostProcessResponse(w, cliCtx, scriptInfo)
	}
}
