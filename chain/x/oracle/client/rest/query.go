package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/bandprotocol/bandchain/chain/hooks/price"
	clientcmn "github.com/bandprotocol/bandchain/chain/x/oracle/client/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type RequestPrices struct {
	Symbols  []string `json:"symbols"`
	MinCount uint64   `json:"min_count"`
	AskCount uint64   `json:"ask_count"`
}

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

func getRequestsPricesHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var requestPrices RequestPrices
		err := decoder.Decode(&requestPrices)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		prices := make([]price.Price, len(requestPrices.Symbols))
		height := int64(0)
		for idx, symbol := range requestPrices.Symbols {
			bz, h, err := cliCtx.Query(fmt.Sprintf("band/prices/%s/%d/%d", symbol, requestPrices.AskCount, requestPrices.MinCount))
			if h > height {
				height = h
			}
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			var price price.Price
			err = cliCtx.Codec.UnmarshalBinaryBare(bz, &price)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			prices[idx] = price
		}
		bz, err := types.QueryOK(prices)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

func getRequestsPriceSymbolsHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		bz, height, err := cliCtx.Query(fmt.Sprintf("band/price_symbols/%s/%s", r.FormValue("ask_count"), r.FormValue("min_count")))

		var symbols []string
		if err := cliCtx.Codec.UnmarshalBinaryBare(bz, &symbols); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		bz, err = types.QueryOK(symbols)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

func getMultiRequestSearchHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		limit := 1
		if rawLimit := r.FormValue("limit"); rawLimit != "" {
			var err error
			limit, err = strconv.Atoi(rawLimit)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			}
		}
		bz, height, err := clientcmn.QueryMultiSearchLatestRequest(
			route, cliCtx,
			r.FormValue("oid"), r.FormValue("calldata"), r.FormValue("ask_count"), r.FormValue("min_count"), limit,
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
	RequestID  types.RequestID  `json:"request_id,string"`
	ExternalID types.ExternalID `json:"external_id,string"`
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
		reporterPubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, detail.Reporter)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		bz, height, err := clientcmn.VerifyRequest(
			route, cliCtx, detail.ChainID, detail.RequestID, detail.ExternalID,
			detail.Validator, reporterPubkey, detail.Signature,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}
