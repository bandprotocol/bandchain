package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ResolveStatus int8

const (
	Open ResolveStatus = iota
	Success
	Failure
	Expired
)

// Request is a data structure that stores the detail of a request to an oracle script.
type Request struct {
	OracleScriptID           OracleScriptID   `json:"oracle_script_id"`
	Calldata                 []byte           `json:"calldata"`
	RequestedValidators      []sdk.ValAddress `json:"requested_validators"`
	SufficientValidatorCount int64            `json:"sufficient_validator_count"`
	ReceivedValidators       []sdk.ValAddress `json:"received_validators"`
	RequestHeight            int64            `json:"request_height"`
	RequestTime              int64            `json:"request_time"`
	ExpirationHeight         int64            `json:"expiration_height"`
	ClientID                 string           `json:"client_id"`

	SourcePort    string `json:"source_port"`
	SourceChannel string `json:"source_channel"`
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

type RawRequest = RawDataRequest

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
