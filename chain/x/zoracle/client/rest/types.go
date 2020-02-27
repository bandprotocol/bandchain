package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
)

type ReportDetail struct {
	Reporter sdk.ValAddress        `json:"reporter"`
	Value    []types.RawDataReport `json:"value"`
	Tx       TxDetail              `json:"tx,omitempty"`
}

type RequestRESTInfo struct {
	ID                       int64                                `json:"id"`
	OracleScriptID           int64                                `json:"oracleScriptID"`
	Calldata                 []byte                               `json:"calldata"`
	RequestedValidators      []sdk.ValAddress                     `json:"requestedValidators"`
	SufficientValidatorCount int64                                `json:"sufficientValidatorCount"`
	ExpirationHeight         int64                                `json:"expirationHeight"`
	IsResolved               bool                                 `json:"isResolved"`
	Requester                sdk.AccAddress                       `json:"requester"`
	RequestTx                TxDetail                             `json:"requestTx,omitempty"`
	RawDataRequests          []types.RawDataRequestWithExternalID `json:"rawDataRequests"`
	Reports                  []ReportDetail                       `json:"reports"`
	Result                   types.Result                         `json:"result"`
}

type TxDetail struct {
	Hash      string `json:"hash"`
	Height    int64  `json:"height"`
	Timestamp string `json:"timestamp"`
}
