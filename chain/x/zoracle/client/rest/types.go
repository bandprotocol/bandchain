package rest

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

type ScriptInfoWithTx struct {
	Info            types.ScriptInfo `json:"info"`
	TxHash          string           `json:"txhash"`
	CreatedAtHeight int64            `json:"createdAtHeight"`
	CreatedAtTime   string           `json:"createdAtTime"`
}

type ReportDetail struct {
	Reporter         sdk.ValAddress  `json:"reporter"`
	TxHash           string          `json:"txhash"`
	ReportedAtHeight int64           `json:"reportedAtHeight"`
	ReportedAtTime   string          `json:"reportedAtTime"`
	Value            json.RawMessage `json:"value"`
}

type RequestQueryInfo struct {
	ScriptInfo        types.ScriptInfo `json:"scriptInfo"`
	CodeHash          cmn.HexBytes     `json:"codeHash"`
	Params            json.RawMessage  `json:"params"`
	TargetBlock       int64            `json:"targetBlock"`
	Requester         sdk.AccAddress   `json:"requester"`
	TxHash            string           `json:"txhash"`
	RequestedAtHeight int64            `json:"requestedAtHeight"`
	RequestedAtTime   string           `json:"requestedAtTime"`
	Reports           []ReportDetail   `json:"reports"`
	Result            json.RawMessage  `json:"result"`
}

type SerializeParams struct {
	Result string `json:"result"`
}
