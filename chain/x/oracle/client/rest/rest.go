package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	"github.com/bandprotocol/bandchain/chain/x/oracle/client/common/proof"
)

const (
	idTag               = "idTag"
	dataHashTag         = "dataHashTag"
	validatorAddressTag = "validatorAddressTag"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/params", storeName), getParamsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/counts", storeName), getCountsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/data/{%s}", storeName, dataHashTag), getDataByHashHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/data_sources/{%s}", storeName, idTag), getDataSourceByIDHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/oracle_scripts/{%s}", storeName, idTag), getOracleScriptByIDHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/requests/{%s}", storeName, idTag), getRequestByIDHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/request_search", storeName), getRequestSearchHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/request_prices", storeName), getRequestsPricesHandler(cliCtx, storeName)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/price_symbols", storeName), getRequestsPriceSymbolsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/multi_request_search", storeName), getMultiRequestSearchHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/validators/{%s}", storeName, validatorAddressTag), getValidatorStatusHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/reporters/{%s}", storeName, validatorAddressTag), getReportersHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/proof/{%s}", storeName, proof.RequestIDTag), proof.GetProofHandlerFn(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/multi_proof", storeName), proof.GetMutiProofHandlerFn(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/active_validators", storeName), getActiveValidatorsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/verify_request", storeName), verifyRequest(cliCtx, storeName)).Methods("POST")
}
