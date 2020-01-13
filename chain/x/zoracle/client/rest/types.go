package rest

import (
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
	Reporter       sdk.ValAddress `json:"reporter"`
	TxHash         string         `json:"txhash"`
	ReportAtHeight int64          `json:"reportAtHeight"`
	ReportAtTime   string         `json:"reportAtTime"`
	Value          types.RawJson  `json:"value"`
}

type RequestQueryInfo struct {
	ScriptInfo      types.ScriptInfo `json:"scriptInfo"`
	CodeHash        cmn.HexBytes     `json:"codeHash"`
	Params          types.RawJson    `json:"params"`
	TargetBlock     int64            `json:"targetBlock"`
	Requester       sdk.AccAddress   `json:"requester"`
	TxHash          string           `json:"txhash"`
	RequestAtHeight int64            `json:"requestAtHeight"`
	RequestAtTime   string           `json:"requestAtTime"`
	Reports         []ReportDetail   `json:"reports"`
	Result          []byte           `json:"result"`
}
