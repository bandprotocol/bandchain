package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

// buildTxDetail takes a TxResponse instance and builds new TxDetail contains only necessary fields.
func buildTxDetail(tx *sdk.TxResponse) TxDetail {
	return TxDetail{
		Hash:      tx.TxHash,
		Height:    tx.Height,
		Timestamp: tx.Timestamp,
	}
}

// buildRequestRESTInfo takes a RequestQuerierInfo instance and builds a more comprehensive version of it.
func buildRequestRESTInfo(
	ctx context.CLIContext, id string, queryRequest types.RequestQuerierInfo,
) (RequestRESTInfo, error) {
	var request RequestRESTInfo

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
	if err != nil || len(searchRequest.Txs) != 1 {
		return RequestRESTInfo{}, err
	}
	request.RequestTx = buildTxDetail(&searchRequest.Txs[0])
	for _, msg := range searchRequest.Txs[0].Tx.GetMsgs() {
		msgRequest, ok := msg.(types.MsgRequestData)
		if ok {
			request.Requester = msgRequest.Sender
			break
		}
	}

	// Save report tx
	searchReports, err := utils.QueryTxsByEvents(
		ctx,
		[]string{fmt.Sprintf("%s.%s='%s'", types.EventTypeReport, types.AttributeKeyRequestID, id)},
		1,
		10000, // Estimated validator reports
	)

	request.Reports = make([]ReportDetail, 0)

	for _, report := range searchReports.Txs {
		var validatorAddress sdk.ValAddress
		for _, msg := range report.Tx.GetMsgs() {
			msgReport, ok := msg.(types.MsgReportData)
			if ok {
				validatorAddress = msgReport.Sender
				break
			}
		}
		for _, queryReport := range queryRequest.Reports {
			if queryReport.Validator.Equals(validatorAddress) {
				request.Reports = append(
					request.Reports,
					ReportDetail{
						Reporter: queryReport.Validator,
						Value:    queryReport.RawDataReports,
						Tx:       buildTxDetail(&report),
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

		request, err := buildRequestRESTInfo(cliCtx, reqID, queryRequest)
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

func getDataSourceHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		dsID := vars[dataSourceID]
		var queryDataSource types.DataSourceQuerierInfo
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/data_source/%s", storeName, dsID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(res, &queryDataSource)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, queryDataSource)
	}
}

func getDataSourcesHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		pageIDStr := vars[dataSourcePageIndex]
		pageSizeStr := vars[dataSourcePageSize]

		pageID, err := strconv.ParseInt(pageIDStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if pageID < 1 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "data source page id should be >= 1")
			return
		}

		pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if pageSize < 1 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "data source page size should be >= 1 and <= 100")
			return
		}

		var queryDataSources []types.DataSourceQuerierInfo
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/data_sources/%d/%d", storeName, 1+(pageID-1)*pageSize, pageSize), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(res, &queryDataSources)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, queryDataSources)
	}
}
