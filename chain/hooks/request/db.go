package request

import (
	"fmt"
	"strings"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type Request struct {
	OracleScriptID types.OracleScriptID `db:"oracle_script_id, primarykey" json:"oracle_script_id"`
	Calldata       []byte               `db:"calldata, primarykey" json:"calldata"`
	MinCount       uint64               `db:"min_count, primarykey" json:"min_count"`
	AskCount       uint64               `db:"ask_count, primarykey" json:"ask_count"`
	RequestIDs     string               `db:"request_ids" json:"request_ids"`
}

func (h *RequestHook) UpsertRequest(request Request) {
	err := h.dbMap.Insert(&request)
	if err != nil {
		h.UpdateRequest(request)
	}
}

func (h *RequestHook) UpdateRequest(request Request) {
	obj, err := h.dbMap.Get(Request{}, request.OracleScriptID, request.Calldata, request.MinCount, request.AskCount)
	if err != nil {
		panic(err)
	}
	data := obj.(*Request)
	data.RequestIDs = fmt.Sprintf("%s,%s", data.RequestIDs, request.RequestIDs)
	_, err = h.dbMap.Update(data)
	if err != nil {
		panic(err)
	}
}

func (h *RequestHook) GetLatestRequestID(oid types.OracleScriptID, calldata []byte, minCount uint64, askCount uint64) types.RequestID {
	obj, err := h.dbMap.Get(Request{}, oid, calldata, minCount, askCount)
	if err != nil {
		panic(err)
	}
	data := obj.(*Request)
	raws := strings.Split(data.RequestIDs, ",")
	raw := raws[len(raws)-1]
	fmt.Println("!@", data)
	fmt.Println("--->", raw, raws)
	return types.RequestID(common.Atoi(raw))
}
