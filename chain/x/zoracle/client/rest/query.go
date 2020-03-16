package rest

import (
	"fmt"
	"net/http"

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
	ctx context.CLIContext, queryRequest types.RequestQuerierInfo, withRequestTx, withReportTx bool,
) (RequestRESTInfo, error) {
	var request RequestRESTInfo

	request.ID = queryRequest.ID
	request.OracleScriptID = queryRequest.Request.OracleScriptID
	request.Calldata = queryRequest.Request.Calldata
	request.RequestedValidators = queryRequest.Request.RequestedValidators
	request.SufficientValidatorCount = queryRequest.Request.SufficientValidatorCount
	request.ExpirationHeight = queryRequest.Request.ExpirationHeight
	request.ResolveStatus = queryRequest.Request.ResolveStatus
	request.RawDataRequests = queryRequest.RawDataRequests

	request.Result = queryRequest.Result

	if withRequestTx {
		// Get request detail
		searchRequest, err := utils.QueryTxsByEvents(
			ctx,
			[]string{fmt.Sprintf("%s.%s='%d'",
				types.EventTypeRequest, types.AttributeKeyID, queryRequest.ID)},
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
	}

	txReportMap := make(map[string]TxDetail)

	if withReportTx {
		// Save report tx
		searchReports, err := utils.QueryTxsByEvents(
			ctx,
			[]string{fmt.Sprintf("%s.%s='%d'",
				types.EventTypeReport, types.AttributeKeyRequestID, queryRequest.ID)},
			1,
			10000, // Estimated validator reports
		)
		if err != nil {
			return RequestRESTInfo{}, err
		}
		for _, report := range searchReports.Txs {
			for _, msg := range report.Tx.GetMsgs() {
				msgReport, ok := msg.(types.MsgReportData)
				if ok {
					txReportMap[string(msgReport.Sender)] = buildTxDetail(&report)
					break
				}
			}
		}
	}

	request.Reports = make([]ReportDetail, 0)

	for _, queryReport := range queryRequest.Reports {
		report, ok := txReportMap[string(queryReport.Validator)]
		if !ok {
			report = TxDetail{}
		}
		request.Reports = append(
			request.Reports,
			ReportDetail{
				Reporter: queryReport.Validator,
				Value:    queryReport.RawDataReports,
				Tx:       report,
			},
		)
	}

	return request, nil
}

func getDataSourceByIDHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		dataSourceID := vars[dataSourceIDTag]
		var queryDataSource types.DataSourceQuerierInfo
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/data_source/%s", storeName, dataSourceID), nil)
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
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 100)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		var queryDataSources []types.DataSourceQuerierInfo
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/data_sources/%d/%d", storeName, 1+(page-1)*limit, limit), nil)
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

func getOracleScriptByIDHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		oracleScriptID := vars[oracleScriptIDTag]
		var queryOracleScript types.OracleScriptQuerierInfo
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/oracle_script/%s", storeName, oracleScriptID))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(res, &queryOracleScript)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, queryOracleScript)
	}
}

func getOracleScriptsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 100)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		var queryOracleScripts []types.OracleScriptQuerierInfo
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/oracle_scripts/%d/%d", storeName, 1+(page-1)*limit, limit), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(res, &queryOracleScripts)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, queryOracleScripts)
	}
}

func getRequestByIDHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestID := vars[requestIDTag]
		var queryRequest types.RequestQuerierInfo
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/request/%s", storeName, requestID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(res, &queryRequest)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		request, err := buildRequestRESTInfo(cliCtx, queryRequest, true, true)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, request)
	}
}

func getRequestsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 100)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/request_number", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var requestNumber int64
		err = cliCtx.Codec.UnmarshalJSON(res, &requestNumber)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		firstID := requestNumber - int64(page*limit) + 1
		lastID := requestNumber - int64((page-1)*limit)
		if lastID < 1 {
			rest.PostProcessResponse(w, cliCtx, []types.RequestQuerierInfo{})
			return
		}

		if firstID < 1 {
			firstID = 1
		}

		var queryRequests []types.RequestQuerierInfo
		res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/requests/%d/%d", storeName, firstID, lastID-firstID+1), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(res, &queryRequests)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var requests []RequestRESTInfo

		for idx := len(queryRequests) - 1; idx >= 0; idx-- {
			request, err := buildRequestRESTInfo(cliCtx, queryRequests[idx], true, true)
			if err == nil {
				requests = append(requests, request)
			}
		}
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, requests)
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
