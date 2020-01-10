package rest

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

type ScriptInfoWithTx struct {
	Info      types.ScriptInfo `json:"info"`
	TxHash    string           `json:"txhash"`
	CreatedAt int64            `json:"createdAt"`
}
