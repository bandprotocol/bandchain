package rest

import (
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
		paramType := vars[requestID]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/request/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
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
