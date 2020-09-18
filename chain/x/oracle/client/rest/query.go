package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	clientcmn "github.com/bandprotocol/bandchain/chain/x/oracle/client/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func getParamsHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", route, types.QueryParams))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

func getCountsHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", route, types.QueryCounts))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

func getDataByHashHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		vars := mux.Vars(r)
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryData, vars[dataHashTag]))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Disposition", "attachment;")
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.Write(res)
	}
}

func getDataSourceByIDHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		vars := mux.Vars(r)
		bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryDataSources, vars[idTag]))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

func getOracleScriptByIDHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		vars := mux.Vars(r)
		bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryOracleScripts, vars[idTag]))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

func getRequestByIDHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		vars := mux.Vars(r)
		bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryRequests, vars[idTag]))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

func getRequestSearchHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		bz, height, err := clientcmn.QuerySearchLatestRequest(
			route, cliCtx,
			r.FormValue("oid"), r.FormValue("calldata"), r.FormValue("ask_count"), r.FormValue("min_count"),
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

func getValidatorStatusHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		vars := mux.Vars(r)
		bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryValidatorStatus, vars[validatorAddressTag]))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

func getReportersHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		vars := mux.Vars(r)
		bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryReporters, vars[validatorAddressTag]))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

func getActiveValidatorsHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", route, types.QueryActiveValidators))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

type requestDetail struct {
	ChainID    string           `json:"chain_id"`
	Validator  sdk.ValAddress   `json:"validator"`
	RequestID  types.RequestID  `json:"request_id"`
	ExternalID types.ExternalID `json:"external_id"`
	Reporter   string           `json:"reporter"`
	Signature  []byte           `json:"signature"`
}

func verifyRequest(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var detail requestDetail
		err := json.NewDecoder(r.Body).Decode(&detail)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		bz, height, err := clientcmn.VerifyRequest(
			route, cliCtx, detail.ChainID, detail.Reporter, detail.Validator,
			detail.RequestID, detail.ExternalID, detail.Signature,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}
