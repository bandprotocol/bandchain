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
	ID         int64          `json:"id"`
	Owner      sdk.AccAddress `json:"owner"`
	Name       string         `json:"name"`
	Fee        sdk.Coins      `json:"fee"`
	Executable []byte         `json:"executable"`
}

func NewDataSourceQuerierInfo(
	id int64,
	owner sdk.AccAddress,
	name string,
	fee sdk.Coins,
	executable []byte,
) DataSourceQuerierInfo {
	return DataSourceQuerierInfo{
		ID:         id,
		Owner:      owner,
		Name:       name,
		Fee:        fee,
		Executable: executable,
	}
}

type OracleScriptQuerierInfo struct {
	ID    int64          `json:"id"`
	Owner sdk.AccAddress `json:"owner"`
	Name  string         `json:"name"`
	Code  []byte         `json:"code"`
}

func NewOracleScriptQuerierInfo(
	id int64,
	owner sdk.AccAddress,
	name string,
	code []byte,
) OracleScriptQuerierInfo {
	return OracleScriptQuerierInfo{
		ID:    id,
		Owner: owner,
		Name:  name,
		Code:  code,
	}
}

type RequestQuerierInfo struct {
	ID              RequestID                      `json:"id"`
	Request         Request                        `json:"request"`
	RawDataRequests []RawDataRequestWithExternalID `json:"rawDataRequests"`
	Reports         []ReportWithValidator          `json:"reports"`
	Result          Result                         `json:"result"`
}

func NewRequestQuerierInfo(
	id RequestID,
	request Request,
	rawDataRequests []RawDataRequestWithExternalID,
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
