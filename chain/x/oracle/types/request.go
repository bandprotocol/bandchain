package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ResolveStatus int8

const (
	Open ResolveStatus = iota
	Success
	Failure
)

// Request is a data structure that stores the detail of a request to an oracle script.
type Request struct {
	OracleScriptID           OracleScriptID   `json:"oracleScriptID"`
	Calldata                 []byte           `json:"calldata"`
	RequestedValidators      []sdk.ValAddress `json:"requestedValidators"`
	SufficientValidatorCount int64            `json:"sufficientValidatorCount"`
	ReceivedValidators       []sdk.ValAddress `json:"receivedValidators"`
	RequestHeight            int64            `json:"requestHeight"`
	RequestTime              int64            `json:"requestTime"`
	ExpirationHeight         int64            `json:"expirationHeight"`
	ResolveStatus            ResolveStatus    `json:"resolveStatus"`
	ClientID                 string           `json:"clientID"`

	SourcePort    string `json:"soucePort"`
	SourceChannel string `json:"sourceChannel"`
}

// NewRequest creates a new Request instance.
func NewRequest(
	oracleScriptID OracleScriptID,
	calldata []byte,
	requestedValidators []sdk.ValAddress,
	sufficientValidatorCount int64,
	requestHeight int64,
	requestTime int64,
	expirationHeight int64,
	clientID string,
) Request {
	return Request{
		OracleScriptID:           oracleScriptID,
		Calldata:                 calldata,
		RequestedValidators:      requestedValidators,
		SufficientValidatorCount: sufficientValidatorCount,
		RequestHeight:            requestHeight,
		RequestTime:              requestTime,
		ExpirationHeight:         expirationHeight,
		ResolveStatus:            Open,
		SourcePort:               "",
		SourceChannel:            "",
		ClientID:                 clientID,
	}
}

// RawDataRequest is a data structure that store what datasource and calldata will be used in request.
type RawDataRequest struct {
	DataSourceID DataSourceID `json:"dataSourceID"`
	Calldata     []byte       `json:"calldata"`
}

// NewRawDataRequest creates a new RawDataRequest instance
func NewRawDataRequest(dataSourceID DataSourceID, calldata []byte) RawDataRequest {
	return RawDataRequest{
		DataSourceID: dataSourceID,
		Calldata:     calldata,
	}
}

// RawDataRequestWithExternalID is a raw data request that contain external id.
type RawDataRequestWithExternalID struct {
	ExternalID     ExternalID     `json:"externalID"`
	RawDataRequest RawDataRequest `json:"detail"`
}

// NewRawDataRequestWithExternalID creates a new RawDataRequestWithExternalID instance.
func NewRawDataRequestWithExternalID(externalID ExternalID, rawDataRequest RawDataRequest) RawDataRequestWithExternalID {
	return RawDataRequestWithExternalID{
		ExternalID:     externalID,
		RawDataRequest: rawDataRequest,
	}
}
