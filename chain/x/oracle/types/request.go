package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewRequest creates a new Request instance.
func NewRequest(
	oracleScriptID OracleScriptID,
	calldata []byte,
	requestedValidators []sdk.ValAddress,
	minCount int64,
	requestHeight int64,
	requestTime int64,
	clientID string,
	ibc *RequestIBC,
) Request {
	return Request{
		OracleScriptID:      oracleScriptID,
		Calldata:            calldata,
		RequestedValidators: requestedValidators,
		MinCount:            minCount,
		RequestHeight:       requestHeight,
		RequestTime:         requestTime,
		ClientID:            clientID,
		IBC:                 ibc,
	}
}

// NewRawRequest creates a new RawRequest instance.
func NewRawRequest(externalID ExternalID, did DataSourceID, calldata []byte) RawRequest {
	return RawRequest{
		ExternalID:   externalID,
		DataSourceID: did,
		Calldata:     calldata,
	}
}
