package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
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
		bz, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", route, types.QueryParams), nil)
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
		bz, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", route, types.QueryCounts), nil)
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
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryData, vars[dataHashTag]), nil)
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
		bz, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryDataSources, vars[idTag]), nil)
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
		bz, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryOracleScripts, vars[idTag]), nil)
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
		bz, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryRequests, vars[idTag]), nil)
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

func getReportersHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		vars := mux.Vars(r)
		bz, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryReporters, vars[validatorAddressTag]), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}
