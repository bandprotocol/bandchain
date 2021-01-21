package yoda

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

type VerificationMessage struct {
	ChainID    string           `json:"chain_id"`
	Validator  sdk.ValAddress   `json:"validator"`
	RequestID  types.RequestID  `json:"request_id"`
	ExternalID types.ExternalID `json:"external_id"`
}

func NewVerificationMessage(
	chainID string, validator sdk.ValAddress, requestID types.RequestID, externalID types.ExternalID,
) VerificationMessage {
	return VerificationMessage{
		ChainID:    chainID,
		Validator:  validator,
		RequestID:  requestID,
		ExternalID: externalID,
	}
}

func (msg VerificationMessage) GetSignBytes() []byte {
	return sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg))
}
