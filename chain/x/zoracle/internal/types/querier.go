package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints
const (
	QueryRequestByID      = "request"
	QueryRequests         = "requests"
	QueryPending          = "pending_request"
	QueryScript           = "script"
	QueryAllScripts       = "scripts"
	SerializeParams       = "serialize_params"
	QueryRequestNumber    = "request_number"
	QueryDataSourceByID   = "data_source"
	QueryDataSources      = "data_sources"
	QueryOracleScriptByID = "oracle_script"
	QueryOracleScripts    = "oracle_scripts"
)

type U64Array []uint64

func (u64a U64Array) String() string {
	return fmt.Sprintf("%v", []uint64(u64a))
}

type RequestQuerierInfo struct {
	ID              int64                          `json:"id"`
	Request         Request                        `json:"request"`
	RawDataRequests []RawDataRequestWithExternalID `json:"rawDataRequests"`
	Reports         []ReportWithValidator          `json:"reports"`
	Result          []byte                         `json:"result"`
}

func NewRequestQuerierInfo(
	id int64,
	request Request,
	rawDataRequests []RawDataRequestWithExternalID,
	reports []ReportWithValidator,
	result []byte,
) RequestQuerierInfo {
	return RequestQuerierInfo{
		ID:              id,
		Request:         request,
		RawDataRequests: rawDataRequests,
		Reports:         reports,
		Result:          result,
	}
}

type RawBytes []byte

func (rb RawBytes) String() string {
	// TODO: Not actual used (just for compiled)
	return "NOT_USED"
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func ParseFields(raw []byte) ([]Field, error) {
	var data [][]string
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}
	fields := make([]Field, len(data))
	for idx, row := range data {
		if len(row) != 2 {
			return nil, fmt.Errorf("Invalid field format")
		}
		fields[idx] = Field{
			Name: row[0],
			Type: row[1],
		}
	}
	return fields, nil
}

func MustParseFields(raw []byte) []Field {
	fields, err := ParseFields(raw)
	if err != nil {
		panic(err)
	}
	return fields
}

type ScriptInfo struct {
	Name        string         `json:"name"`
	CodeHash    []byte         `json:"codeHash"`
	Params      []Field        `json:"params"`
	DataSources []Field        `json:"dataSources"`
	Result      []Field        `json:"result"`
	Creator     sdk.AccAddress `json:"creator"`
}

func NewScriptInfo(
	name string,
	codeHash []byte,
	rawParams, rawDataSources, result []Field,
	creator sdk.AccAddress,
) ScriptInfo {
	return ScriptInfo{
		Name:        name,
		CodeHash:    codeHash,
		Params:      rawParams,
		DataSources: rawDataSources,
		Result:      result,
		Creator:     creator,
	}
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
