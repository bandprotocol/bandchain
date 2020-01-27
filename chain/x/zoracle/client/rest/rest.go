package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	requestID = "requestID"
	codeHash  = "codeHash"
	params    = "params"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/request/{%s}", storeName, requestID), getRequestHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/script/{%s}", storeName, codeHash), getScriptInfoHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/scripts", storeName), getScriptsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{%s}/serialize-params/{%s}", storeName, codeHash, params), getSerializeParams(cliCtx, storeName)).Methods("GET")
}
