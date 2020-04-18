package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type ReportDetail struct {
	Reporter sdk.ValAddress              `json:"reporter"`
	Value    []types.RawDataReportWithID `json:"value"`
	Tx       TxDetail                    `json:"tx,omitempty"`
}

type RequestRESTInfo struct {
	ID                       types.RequestID                      `json:"id"`
	OracleScriptID           types.OracleScriptID                 `json:"oracleScriptID"`
	Calldata                 []byte                               `json:"calldata"`
	RequestedValidators      []sdk.ValAddress                     `json:"requestedValidators"`
	SufficientValidatorCount int64                                `json:"sufficientValidatorCount"`
	ExpirationHeight         int64                                `json:"expirationHeight"`
	ResolveStatus            types.ResolveStatus                  `json:"resolveStatus"`
	Requester                sdk.AccAddress                       `json:"requester"`
	RequestTx                TxDetail                             `json:"requestTx,omitempty"`
	RawDataRequests          []types.RawRequest `json:"rawDataRequests"`
	Reports                  []ReportDetail                       `json:"reports"`
	Result                   types.Result                         `json:"result"`
}

type TxDetail struct {
	Hash      string `json:"hash"`
	Height    int64  `json:"height"`
	Timestamp string `json:"timestamp"`
}
