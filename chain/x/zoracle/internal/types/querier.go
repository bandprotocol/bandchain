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

type RequestInfo struct {
	CodeHash    []byte            `json:"codeHash"`
	ParamsHex   []byte            `json:"paramsHex"`
	Params      RawJson           `json:"params"`
	TargetBlock uint64            `json:"targetBlock"`
	Reports     []ValidatorReport `json:"reports"`
	Result      []byte            `json:"result"`
}

func NewRequestInfo(
	codeHash []byte,
	paramsHex []byte,
	params RawJson,
	targetBlock uint64,
	reports []ValidatorReport,
	result []byte,
) RequestInfo {
	return RequestInfo{
		CodeHash:    codeHash,
		ParamsHex:   paramsHex,
		Params:      params,
		TargetBlock: targetBlock,
		Reports:     reports,
		Result:      result,
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
