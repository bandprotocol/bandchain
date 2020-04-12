package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OracleRequestPacketData struct {
	OracleScriptID           OracleScriptID `json:"oracleScriptID"`
	Calldata                 string         `json:"calldata"`
	RequestedValidatorCount  int64          `json:"requestedValidatorCount"`
	SufficientValidatorCount int64          `json:"sufficientValidatorCount"`
	ClientID                 string         `json:"clientID"`
}

func NewOracleRequestPacketData(
	oracleScriptID OracleScriptID, calldata string, requestedValidatorCount int64,
	sufficientValidatorCount int64, clientID string,
) OracleRequestPacketData {
	return OracleRequestPacketData{
		OracleScriptID:           oracleScriptID,
		Calldata:                 calldata,
		RequestedValidatorCount:  requestedValidatorCount,
		SufficientValidatorCount: sufficientValidatorCount,
		ClientID:                 clientID,
	}
}

func (p OracleRequestPacketData) String() string {
	return fmt.Sprintf(`OracleRequestPacketData:
	OracleScriptID:           %d
	Calldata:                 %s
	RequestedValidatorCount:  %d
	SufficientValidatorCount: %d
	ClientID:                 %s`,
		p.OracleScriptID,
		p.Calldata,
		p.RequestedValidatorCount,
		p.SufficientValidatorCount,
		p.ClientID,
	)
}

func (p OracleRequestPacketData) ValidateBasic() error {
	// TODO: Validate oracle request packet
	return nil
}

func (p OracleRequestPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}

type OracleResponsePacketData struct {
	RequestID     RequestID     `json:"requestID"`
	ClientID      string        `json:"clientID"`
	ResolveStatus ResolveStatus `json:"resolveStatus"`
	Result        string        `json:"result"`
}

func NewOracleResponsePacketData(
	requestID RequestID,
	clientID string,
	resolveStatus ResolveStatus,
	result string,
) OracleResponsePacketData {
	return OracleResponsePacketData{
		RequestID:     requestID,
		ClientID:      clientID,
		ResolveStatus: resolveStatus,
		Result:        result,
	}
}

func (p OracleResponsePacketData) String() string {
	return fmt.Sprintf(`OracleResponsePacketData:
	RequestID: %d
	ClientID: %s
	ResolveStatus: %d
	Result: %s`,
		p.RequestID,
		p.ClientID,
		p.ResolveStatus,
		p.Result,
	)
}

func (p OracleResponsePacketData) ValidateBasic() error {
	// TODO: Validate oracle request packet
	return nil
}

func (p OracleResponsePacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}
