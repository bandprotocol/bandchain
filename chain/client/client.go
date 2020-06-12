package rpc

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	"github.com/bandprotocol/bandchain/chain/client/proof"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/bandchain/genesis", GetGenesisHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bandchain/evm-validators", GetEVMValidators(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/bandchain/proof/{%s}", proof.RequestIDTag), proof.GetProofHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/bandchain/getfile/{%s}", Filename), GetFile()).Methods("GET")
}
