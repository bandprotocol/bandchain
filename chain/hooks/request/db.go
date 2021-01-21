package request

import (
	"encoding/hex"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

type Request struct {
	RequestID      types.RequestID      `db:"request_id, primarykey" json:"request_id"`
	OracleScriptID types.OracleScriptID `db:"oracle_script_id" json:"oracle_script_id"`
	Calldata       string               `db:"calldata" json:"calldata"`
	MinCount       uint64               `db:"min_count" json:"min_count"`
	AskCount       uint64               `db:"ask_count" json:"ask_count"`
	ResolveTime    int64                `db:"resolve_time" json:"resolve_time"`
}

func (h *Hook) insertRequest(requestID types.RequestID, oracleScriptID types.OracleScriptID, calldata []byte, askCount uint64, minCount uint64, resolveTime int64) {
	err := h.trans.Insert(&Request{
		RequestID:      requestID,
		OracleScriptID: oracleScriptID,
		Calldata:       hex.EncodeToString(calldata),
		MinCount:       minCount,
		AskCount:       askCount,
		ResolveTime:    resolveTime,
	})
	if err != nil {
		panic(err)
	}
}

func (h *Hook) getMultiRequestID(oid types.OracleScriptID, calldata string, askCount uint64, minCount uint64, limit int64) []types.RequestID {
	var requests []Request
	h.dbMap.Select(&requests,
		`select * from request
where oracle_script_id = ? and calldata = ? and min_count = ? and ask_count = ?
order by resolve_time desc limit ?`,
		oid, calldata, minCount, askCount, limit)
	requestIDs := make([]types.RequestID, len(requests))
	for idx, request := range requests {
		requestIDs[idx] = request.RequestID
	}
	return requestIDs
}
