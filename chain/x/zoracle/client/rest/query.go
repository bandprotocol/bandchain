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

func newRequestQueryInfo(
	ctx context.CLIContext, id string, queryRequest types.RequestQuerierInfo,
) (RequestQueryInfo, error) {
	var request RequestQueryInfo

	request.OracleScriptID = queryRequest.Request.OracleScriptID
	request.Calldata = queryRequest.Request.Calldata
	request.RequestedValidators = queryRequest.Request.RequestedValidators
	request.SufficientValidatorCount = queryRequest.Request.SufficientValidatorCount
	request.ExpirationHeight = queryRequest.Request.ExpirationHeight
	request.IsResolved = queryRequest.Request.IsResolved
	request.RawDataRequests = queryRequest.RawDataRequests

	request.Result = queryRequest.Result

	// Get request detail
	searchRequest, err := utils.QueryTxsByEvents(
		ctx,
		[]string{fmt.Sprintf("%s.%s='%s'", types.EventTypeRequest, types.AttributeKeyRequestID, id)},
		1,
		1,
	)
	if err != nil {
		return RequestQueryInfo{}, err
	}

	request.TxHash = searchRequest.Txs[0].TxHash
	request.RequestedAtHeight = searchRequest.Txs[0].Height
	request.RequestedAtTime = searchRequest.Txs[0].Timestamp
	// TODO: Find the correct message and not assume the first message is the one
	request.Requester = searchRequest.Txs[0].Tx.GetMsgs()[0].(types.MsgRequestData).Sender

	// Save report tx
	searchReports, err := utils.QueryTxsByEvents(
		ctx,
		[]string{fmt.Sprintf("%s.%s='%s'", types.EventTypeReport, types.AttributeKeyRequestID, id)},
		1,
		10000, // Estimated validator reports
	)

	request.Reports = make([]ReportDetail, 0)

	for _, report := range searchReports.Txs {
		// TODO: Find validator address from tx not assume in first log of tx
		validatorAddress := report.Logs[0].Events[1].Attributes[1].Value
		for _, queryReport := range queryRequest.Reports {
			if queryReport.Validator.String() == validatorAddress {
				request.Reports = append(
					request.Reports,
					ReportDetail{
						Reporter:         queryReport.Validator,
						TxHash:           report.TxHash,
						ReportedAtHeight: report.Height,
						ReportedAtTime:   report.Timestamp,
						Value:            queryReport.RawDataReports,
					},
				)
				continue
			}
		}
	}

	return request, nil

}

func getRequestHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		reqID := vars[requestID]
		var queryRequest types.RequestQuerierInfo
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/request/%s", storeName, reqID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(res, &queryRequest)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		request, err := newRequestQueryInfo(cliCtx, reqID, queryRequest)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
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
	if len(searchResult.Txs) == 1 {
		script.TxHash = searchResult.Txs[0].TxHash
		script.CreatedAtHeight = searchResult.Txs[0].Height
		script.CreatedAtTime = searchResult.Txs[0].Timestamp
	} else {
		script.TxHash = "0000000000000000000000000000000000000000000000000000000000000000"
		script.CreatedAtHeight = 0
		script.CreatedAtTime = "0"
	}
	return nil
}

func getSerializeParams(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		params := r.URL.Query()["params"][0]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/serialize_params/%s/%s", storeName, vars[codeHash], params), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if string(res) == "null" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid params")
			return
		}

		var serializeParamsBytes []byte
		err = cliCtx.Codec.UnmarshalJSON(res, &serializeParamsBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, hex.EncodeToString(serializeParamsBytes))
	}
}

func getRequestNumberHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/request_number", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var requestNumber uint64
		err = cliCtx.Codec.UnmarshalJSON(res, &requestNumber)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, requestNumber)
	}
}
