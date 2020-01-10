package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints
const (
	QueryRequest = "request"
	QueryPending = "pending_request"
	QueryScript  = "script"
)

type U64Array []uint64

func (u64a U64Array) String() string {
	return fmt.Sprintf("%v", []uint64(u64a))
}

type ReportDetail struct {
	Reporter sdk.AccAddress `json:"reporter"`
	TxHash   string         `json:"txhash"`
	ReportAt uint64         `json:"reportAt"`
	Value    interface{}    `json:"value"`
}

type RequestWithReport struct {
	// Script
	CodeHash    []byte         `json:"codeHash"`
	Params      interface{}    `json:"params"`
	TargetBlock uint64         `json:"targetBlock"`
	Requester   sdk.AccAddress `json:"requester"`
	RequestAt   uint64         `json:"requestAt"`
	Reports     []ReportDetail `json:"reports"`
	Result      []byte         `json:"result"`
}

func NewRequestWithReport(request Request, result []byte, reports []ValidatorReport) RequestWithReport {
	return RequestWithReport{
		Request: request,
		Result:  result,
		Reports: reports,
	}
}

func (re RequestWithReport) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(re.Request)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(d, &data); err != nil {
		return nil, err
	}

	data["result"] = re.Result
	data["reports"] = re.Reports

	return json.Marshal(data)
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
	Params      []Field        `json:"params"`
	DataSources []Field        `json:"dataSources"`
	Creator     sdk.AccAddress `json:"creator"`
}

func NewScriptInfo(name string, rawParams, rawDataSources []Field, creator sdk.AccAddress) ScriptInfo {
	return ScriptInfo{
		Name:        name,
		Params:      rawParams,
		DataSources: rawDataSources,
		Creator:     creator,
	}
}
