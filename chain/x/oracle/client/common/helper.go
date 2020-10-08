package common

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func PostProcessQueryResponse(w http.ResponseWriter, cliCtx context.CLIContext, bz []byte) {
	var result types.QueryResult
	if err := json.Unmarshal(bz, &result); err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.Status)
	rest.PostProcessResponse(w, cliCtx, result.Result)
}
