package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints
const (
	QueryRequest    = "request"
	QueryPending    = "pending_request"
	QueryScript     = "script"
	QueryAllScripts = "scripts"
	SerializeParams = "serialize_params"
)

type U64Array []uint64

func (u64a U64Array) String() string {
	return fmt.Sprintf("%v", []uint64(u64a))
}

type RequestInfo struct {
	CodeHash    []byte            `json:"codeHash"`
	Params      json.RawMessage   `json:"params"`
	ParamsRaw   []byte            `json:"paramsRaw"`
	TargetBlock uint64            `json:"targetBlock"`
	Reports     []ValidatorReport `json:"reports"`
	Result      json.RawMessage   `json:"result"`
	ResultRaw   []byte            `json:"resultRaw"`
}

func NewRequestInfo(
	codeHash []byte,
	params json.RawMessage,
	paramsRaw []byte,
	targetBlock uint64,
	reports []ValidatorReport,
	result json.RawMessage,
	resultRaw []byte,
) RequestInfo {
	return RequestInfo{
		CodeHash:    codeHash,
		ParamsRaw:   paramsRaw,
		Params:      params,
		TargetBlock: targetBlock,
		Reports:     reports,
		Result:      result,
		ResultRaw:   resultRaw,
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
