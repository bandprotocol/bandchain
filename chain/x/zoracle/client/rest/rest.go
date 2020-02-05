package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	requestID = "requestID"
	codeHash  = "codeHash"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/request/{%s}", storeName, requestID), getRequestHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/script/{%s}", storeName, codeHash), getScriptInfoHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/scripts", storeName), getScriptsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/serialize_params/{%s}", storeName, codeHash), getSerializeParams(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/request_number", storeName), getRequestNumberHandler(cliCtx, storeName)).Methods("GET")
}
