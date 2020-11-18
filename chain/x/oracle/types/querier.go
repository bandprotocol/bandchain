package types

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Query endpoints supported by the oracle Querier.
const (
	QueryParams           = "params"
	QueryCounts           = "counts"
	QueryData             = "data"
	QueryDataSources      = "data_sources"
	QueryOracleScripts    = "oracle_scripts"
	QueryRequests         = "requests"
	QueryValidatorStatus  = "validator_status"
	QueryReporters        = "reporters"
	QueryActiveValidators = "active_validators"
	QueryPendingRequests  = "pending_requests"
)

// QueryResult wraps querier result with HTTP status to return to application.
type QueryResult struct {
	Status int             `json:"status"`
	Result json.RawMessage `json:"result"`
}

// QueryOK creates and marshals a QueryResult instance with HTTP status OK.
func QueryOK(result interface{}) ([]byte, error) {
	return json.MarshalIndent(QueryResult{
		Status: http.StatusOK,
		Result: codec.MustMarshalJSONIndent(ModuleCdc, result),
	}, "", "  ")
}

// QueryBadRequest creates and marshals a QueryResult instance with HTTP status BadRequest.
func QueryBadRequest(result interface{}) ([]byte, error) {
	return json.MarshalIndent(QueryResult{
		Status: http.StatusBadRequest,
		Result: codec.MustMarshalJSONIndent(ModuleCdc, result),
	}, "", "  ")
}

// QueryNotFound creates and marshals a QueryResult instance with HTTP status NotFound.
func QueryNotFound(result interface{}) ([]byte, error) {
	return json.MarshalIndent(QueryResult{
		Status: http.StatusNotFound,
		Result: codec.MustMarshalJSONIndent(ModuleCdc, result),
	}, "", "  ")
}

// QueryCountsResult is the struct for the result of query counts.
type QueryCountsResult struct {
	DataSourceCount   int64 `json:"data_source_count"`
	OracleScriptCount int64 `json:"oracle_script_count"`
	RequestCount      int64 `json:"request_count"`
}

// QueryRequestResult is the struct for the result of request query.
type QueryRequestResult struct {
	Request Request  `json:"request"`
	Reports []Report `json:"reports"`
	Result  *Result  `json:"result"`
}

// QueryActiveValidatorResult is the struct for the result of request active validators.
type QueryActiveValidatorResult struct {
	Address sdk.ValAddress `json:"address"`
	Power   uint64         `json:"power"`
}
