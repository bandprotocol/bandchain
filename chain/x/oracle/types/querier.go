package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints
const (
	QueryDataSourceByID   = "data_source"
	QueryDataSources      = "data_sources"
	QueryOracleScriptByID = "oracle_script"
	QueryOracleScripts    = "oracle_scripts"
	QueryRequestByID      = "request"
	QueryRequests         = "requests"
	QueryPending          = "pending_request"
	QueryRequestNumber    = "request_number"
)

type RawBytes []byte

func (rb RawBytes) String() string {
	// TODO: Not actual used (just for compiled)
	return "NOT_USED"
}

type U64Array []uint64

func (u64a U64Array) String() string {
	return fmt.Sprintf("%v", []uint64(u64a))
}

type DataSourceQuerierInfo struct {
	ID          DataSourceID   `json:"id"`
	Owner       sdk.AccAddress `json:"owner"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Fee         sdk.Coins      `json:"fee"`
	Executable  []byte         `json:"executable"`
}

func NewDataSourceQuerierInfo(
	id DataSourceID,
	owner sdk.AccAddress,
	name string,
	description string,
	fee sdk.Coins,
	executable []byte,
) DataSourceQuerierInfo {
	return DataSourceQuerierInfo{
		ID:          id,
		Owner:       owner,
		Name:        name,
		Description: description,
		Fee:         fee,
		Executable:  executable,
	}
}

type OracleScriptQuerierInfo struct {
	ID            OracleScriptID `json:"id"`
	Owner         sdk.AccAddress `json:"owner"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Code          []byte         `json:"code"`
	Schema        string         `json:"schema"`
	SourceCodeURL string         `json:"source_code_url"`
}

func NewOracleScriptQuerierInfo(
	id OracleScriptID,
	owner sdk.AccAddress,
	name string,
	description string,
	code []byte,
	schema string,
	sourceCodeURL string,
) OracleScriptQuerierInfo {
	return OracleScriptQuerierInfo{
		ID:            id,
		Owner:         owner,
		Description:   description,
		Name:          name,
		Code:          code,
		Schema:        schema,
		SourceCodeURL: sourceCodeURL,
	}
}

type RequestQuerierInfo struct {
	ID              RequestID                      `json:"id"`
	Request         Request                        `json:"request"`
	RawDataRequests []RawRequest `json:"rawDataRequests"`
	Reports         []ReportWithValidator          `json:"reports"`
	Result          Result                         `json:"result"`
}

func NewRequestQuerierInfo(
	id RequestID,
	request Request,
	rawDataRequests []RawRequest,
	reports []ReportWithValidator,
	result Result,
) RequestQuerierInfo {
	return RequestQuerierInfo{
		ID:              id,
		Request:         request,
		RawDataRequests: rawDataRequests,
		Reports:         reports,
		Result:          result,
	}
}
