package request

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type Request struct {
	RequestID      types.RequestID      `db:"request_id, primarykey" json:"request_id"`
	OracleScriptID types.OracleScriptID `db:"oracle_script_id" json:"oracle_script_id"`
	Calldata       []byte               `db:"calldata" json:"calldata"`
	MinCount       uint64               `db:"min_count" json:"min_count"`
	AskCount       uint64               `db:"ask_count" json:"ask_count"`
	ResolveTime    int64                `db:"resolve_time" json:"resolve_time"`
}

func (h *RequestHook) InsertRequest(requestID types.RequestID, oracleScriptID types.OracleScriptID, calldata []byte, minCount uint64, askCount uint64, resolveTime int64) {
	err := h.dbMap.Insert(&Request{
		RequestID:      requestID,
		OracleScriptID: oracleScriptID,
		Calldata:       calldata,
		MinCount:       minCount,
		AskCount:       askCount,
		ResolveTime:    resolveTime,
	})
	if err != nil {
		panic(err)
	}
}

func (h *RequestHook) GetLatestRequestID(oid types.OracleScriptID, calldata []byte, minCount uint64, askCount uint64) types.RequestID {
	var latestRequest Request
	h.dbMap.SelectOne(&latestRequest,
		`select * from request
where oracle_script_id = ? and calldata = ? and min_count = ? and ask_count = ?
order by resolve_time desc limit 1`,
		oid, calldata, minCount, askCount)
	return types.RequestID(latestRequest.RequestID)
}

func (h *RequestHook) GetMultiRequestID(oid types.OracleScriptID, calldata []byte, minCount uint64, askCount uint64, limit int64) []types.RequestID {
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

func (h *RequestHook) DeleteExpiredData(resolveTime int64) {
	h.dbMap.Exec(`delete from request
where resolve_time < ?
		`, resolveTime)
}
