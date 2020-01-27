package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"

	"github.com/bandprotocol/d3n/chain/wasm"
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
		request.RequestedAtHeight = searchRequest.Txs[0].Height
		request.RequestedAtTime = searchRequest.Txs[0].Timestamp
		// TODO: Find the correct message and not assume the first message is the one
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
			10000, // Estimated validator reports
		)

		for _, report := range searchReports.Txs {
			// TODO: Find validator address from tx not assume in first log of tx
			reportTxs[report.Logs[0].Events[1].Attributes[2].Value] = ReportDetail{
				TxHash:         report.TxHash,
				ReportedAtTime: report.Timestamp,
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
				Reporter:         report.Validator,
				TxHash:           reportTx.TxHash,
				ReportedAtHeight: int64(report.ReportedAt),
				ReportedAtTime:   reportTx.ReportedAtTime,
				Value:            report.Value,
			}
		}

		rest.PostProcessResponse(w, cliCtx, request)
	}
}

func getScriptInfoHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
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
		err = getStoreTxInfo(cliCtx, &scriptInfo, hash)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, scriptInfo)
	}
}

func getScriptsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 100)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/scripts/%d/%d", storeName, page, limit), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var rawScripts []types.ScriptInfo
		err = cliCtx.Codec.UnmarshalJSON(res, &rawScripts)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		scripts := make([]ScriptInfoWithTx, len(rawScripts))
		for i, _ := range scripts {
			scripts[i].Info = rawScripts[i]
			err := getStoreTxInfo(cliCtx, &scripts[i], hex.EncodeToString(scripts[i].Info.CodeHash))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		rest.PostProcessResponse(w, cliCtx, scripts)
	}
}

func getStoreTxInfo(cliCtx context.CLIContext, script *ScriptInfoWithTx, hash string) error {
	// TODO: Get latest store tx as tx hash (wait tendermint release get result in desc order)
	searchResult, err := utils.QueryTxsByEvents(
		cliCtx,
		[]string{fmt.Sprintf("%s.%s='%s'", types.EventTypeStoreCode, types.AttributeKeyCodeHash, hash)},
		1,
		1,
	)
	if err != nil {
		return err
	}
	script.TxHash = searchResult.Txs[0].TxHash
	script.CreatedAtHeight = searchResult.Txs[0].Height
	script.CreatedAtTime = searchResult.Txs[0].Timestamp
	return nil
}

func getSerializeParams(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramsVar := vars[params]
		codeHashVar := vars[codeHash]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/script/%s", storeName, codeHashVar), nil)
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

		txHash := scriptInfo.TxHash

		res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/txs/%s", txHash), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
		return

		rawParams, err := wasm.SerializeParams([]byte("aaaa"), []byte(paramsVar))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, []byte(fmt.Sprintf("%v \n %v", vars, rawParams)))
	}
}
