package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	requestID = "requestID"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/request/{%s}", storeName, requestID), getRequestHandler(cliCtx, storeName)).Methods("GET")
}
