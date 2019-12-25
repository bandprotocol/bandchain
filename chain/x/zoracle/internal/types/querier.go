package types

import (
	"encoding/json"
	"fmt"
)

// query endpoints
const (
	QueryRequest = "request"
	QueryPending = "pending_request"
)

type U64Array []uint64

func (u64a U64Array) String() string {
	return fmt.Sprintf("%v", []uint64(u64a))
}

type RequestWithReport struct {
	Request
	Result  []byte            `json:"result"`
	Reports []ValidatorReport `json:"reports"`
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
