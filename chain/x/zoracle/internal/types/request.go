package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	ExecuteGas               uint64           `json:"executeGas"`
	IsResolved               bool             `json:"isResolved"`
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
	executeGas uint64,
) Request {
	return Request{
		OracleScriptID:           oracleScriptID,
		Calldata:                 calldata,
		RequestedValidators:      requestedValidators,
		SufficientValidatorCount: sufficientValidatorCount,
		RequestHeight:            requestHeight,
		RequestTime:              requestTime,
		ExpirationHeight:         expirationHeight,
		ExecuteGas:               executeGas,
	}
}

// RawDataRequest is a data structure that store what datasource and calldata will be used in request.
type RawDataRequest struct {
	DataSourceID int64  `json:"dataSourceID"`
	Calldata     []byte `json:"calldata"`
}

// NewRawDataRequest creates a new RawDataRequest instance
func NewRawDataRequest(dataSourceID int64, calldata []byte) RawDataRequest {
	return RawDataRequest{
		DataSourceID: dataSourceID,
		Calldata:     calldata,
	}
}

// RawDataRequestWithExternalID is a raw data request that contain external id.
type RawDataRequestWithExternalID struct {
	ExternalID     int64          `json:"externalID"`
	RawDataRequest RawDataRequest `json:"detail"`
}

// NewRawDataRequestWithExternalID creates a new RawDataRequestWithExternalID instance.
func NewRawDataRequestWithExternalID(externalID int64, rawDataRequest RawDataRequest) RawDataRequestWithExternalID {
	return RawDataRequestWithExternalID{
		ExternalID:     externalID,
		RawDataRequest: rawDataRequest,
	}
}
