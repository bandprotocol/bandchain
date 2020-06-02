package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	dataSourceIDTag   = "dataSourceIDTag"
	oracleScriptIDTag = "oracleScriptIDTag"
	dataHashTag       = "dataHashTag"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/params", storeName), getParamsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/counts", storeName), getCountsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/data/{%s}", storeName, dataHashTag), getDataByHashHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/data_sources/{%s}", storeName, dataSourceIDTag), getDataSourceByIDHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/oracle_scripts/{%s}", storeName, oracleScriptIDTag), getOracleScriptByIDHandler(cliCtx, storeName)).Methods("GET")
}
