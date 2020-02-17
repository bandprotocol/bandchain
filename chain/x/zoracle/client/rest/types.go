package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

type ScriptInfoWithTx struct {
	Info            types.ScriptInfo `json:"info"`
	TxHash          string           `json:"txhash"`
	CreatedAtHeight int64            `json:"createdAtHeight"`
	CreatedAtTime   string           `json:"createdAtTime"`
}

type ReportDetail struct {
	Reporter         sdk.ValAddress        `json:"reporter"`
	TxHash           string                `json:"txhash"`
	ReportedAtHeight int64                 `json:"reportedAtHeight"`
	ReportedAtTime   string                `json:"reportedAtTime"`
	Value            []types.RawDataReport `json:"value"`
}

type RequestQueryInfo struct {
	OracleScriptID           int64                  `json:"oracleScriptID"`
	Calldata                 []byte                 `json:"calldata"`
	RequestedValidators      []sdk.ValAddress       `json:"requestedValidators"`
	SufficientValidatorCount int64                  `json:"sufficientValidatorCount"`
	ExpirationHeight         int64                  `json:"expirationHeight"`
	IsResolved               bool                   `json:"isResolved"`
	Requester                sdk.AccAddress         `json:"requester"`
	TxHash                   string                 `json:"txhash"`
	RequestedAtHeight        int64                  `json:"requestedAtHeight"`
	RequestedAtTime          string                 `json:"requestedAtTime"`
	RawDataRequests          []types.RawDataRequest `json:"rawDataRequests"`
	Reports                  []ReportDetail         `json:"reports"`
	Result                   []byte                 `json:"result"`
}
