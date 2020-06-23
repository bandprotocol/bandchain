package common

import (
	"encoding/json"
	"net/http"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func PostProcessQueryResponse(w http.ResponseWriter, cliCtx context.CLIContext, bz []byte) {
	var result types.QueryResult
	if err := json.Unmarshal(bz, &result); err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(result.Status)
	rest.PostProcessResponse(w, cliCtx, result.Result)
}
