package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	requestIDTag      = "requestIDTag"
	dataSourceIDTag   = "dataSourceIDTag"
	oracleScriptIDTag = "oracleScriptIDTag"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/data_source/{%s}", storeName, dataSourceIDTag), getDataSourceByIDHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/data_sources", storeName), getDataSourcesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/oracle_script/{%s}", storeName, oracleScriptIDTag), getOracleScriptByIDHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/oracle_scripts", storeName), getOracleScriptsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/request/{%s}", storeName, requestIDTag), getRequestByIDHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/requests", storeName), getRequestsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/request_number", storeName), getRequestNumberHandler(cliCtx, storeName)).Methods("GET")
}
