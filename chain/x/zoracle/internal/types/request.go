package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Request is a data structure that stores the detail of a request to an oracle script.
type Request struct {
	OracleScriptID           int64            `json:"oracleScriptID"`
	Calldata                 []byte           `json:"calldata"`
	Validators               []sdk.ValAddress `json:"validators"`
	SufficientValidatorCount int64            `json:"sufficientValidatorCount"`
	SubmittedValidatorList   []sdk.ValAddress `json:"submittedValidatorList"`
	RequestHeight            int64            `json:"requestHeight"`
	RequestTime              int64            `json:"requestTime"`
	ExpirationHeight         int64            `json:"expirationHeight"`
	DataSourceCount          int64            `json:"dataSourceCount"`
	IsResolved               bool             `json:"isResolved"`
}

// NewRequest creates a new Request instance.
func NewRequest(
	oracleScriptID int64,
	calldata []byte,
	validators []sdk.ValAddress,
	sufficientValidatorCount int64,
	requestHeight int64,
	requestTime int64,
	expirationHeight int64,
) Request {
	return Request{
		OracleScriptID:           oracleScriptID,
		Calldata:                 calldata,
		Validators:               validators,
		SufficientValidatorCount: sufficientValidatorCount,
		RequestHeight:            requestHeight,
		RequestTime:              requestTime,
		ExpirationHeight:         expirationHeight,
	}
}

// RequestDetail is a data structure that store what datasource and calldata will be used in request.
type RequestDetail struct {
	ScriptID int64  `json:"scriptID"`
	Calldata []byte `json:"calldata"`
}

// NewRequestDetail creates a new RequestDetail instance
func NewRequestDetail(scriptID int64, calldata []byte) RequestDetail {
	return RequestDetail{
		ScriptID: scriptID,
		Calldata: calldata,
	}
}
