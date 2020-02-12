package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Request is a type to store detail of request.
type Request struct {
	OracleScriptID           int64            `json:"oracleScriptID"`
	Calldata                 []byte           `json:"calldata"`
	ValidatorList            []sdk.ValAddress `json:"validatorList"`
	SufficientValidatorCount int64            `json:"sufficientValidatorCount"`
	SubmittedValidatorList   []sdk.ValAddress `json:"submittedValidatorList"`
	RequestedBlock           int64            `json:"requestedBlock"`
	RequestedTime            int64            `json:"requestedTime"`
	ExpiredBlock             int64            `json:"expiredBlock"`
	IsResolved               bool             `json:"isResolved"`
}

// NewRequest creates a new Request instance.
func NewRequest(
	oracleScriptID int64,
	calldata []byte,
	validatorList []sdk.ValAddress,
	sufficientValidatorCount int64,
	requestedBlock int64,
	requestedTime int64,
	expiredBlock int64,
) Request {
	return Request{
		OracleScriptID:           oracleScriptID,
		Calldata:                 calldata,
		ValidatorList:            validatorList,
		SufficientValidatorCount: sufficientValidatorCount,
		RequestedBlock:           requestedBlock,
		RequestedTime:            requestedTime,
		ExpiredBlock:             expiredBlock,
	}
}
