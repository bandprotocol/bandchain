package types

import (
	"fmt"
)

type OracleRequestPacketData struct {
	OracleScriptID           OracleScriptID `json:"oracleScriptID"`
	Calldata                 []byte         `json:"calldata"`
	RequestedValidatorCount  int64          `json:"requestedValidatorCount"`
	SufficientValidatorCount int64          `json:"sufficientValidatorCount"`
	Expiration               int64          `json:"expiration"`
	PrepareGas               uint64         `json:"prepareGas"`
	ExecuteGas               uint64         `json:"executeGas"`
}

func NewOracleRequestPacketData(
	oracleScriptID OracleScriptID, calldata []byte, requestedValidatorCount int64,
	sufficientValidatorCount int64, expiration int64, prepareGas uint64, executeGas uint64,
) OracleRequestPacketData {
	return OracleRequestPacketData{
		OracleScriptID:           oracleScriptID,
		Calldata:                 calldata,
		RequestedValidatorCount:  requestedValidatorCount,
		SufficientValidatorCount: sufficientValidatorCount,
		Expiration:               expiration,
		PrepareGas:               prepareGas,
		ExecuteGas:               executeGas,
	}
}

func (p OracleRequestPacketData) String() string {
	return fmt.Sprintf(`OracleRequestPacketData:
	OracleScriptID:           %d
	Calldata:                 %x
	RequestedValidatorCount:  %d
	SufficientValidatorCount: %d
	Expiration:               %d
	PrepareGas:               %d
	ExecuteGas:               %d`,
		p.OracleScriptID,
		p.Calldata,
		p.RequestedValidatorCount,
		p.SufficientValidatorCount,
		p.Expiration,
		p.PrepareGas,
		p.ExecuteGas,
	)
}

func (p OracleRequestPacketData) ValidateBasic() error {
	// TODO: Validate oracle request packet
	return nil
}

func (p OracleRequestPacketData) GetBytes() []byte {
	return ModuleCdc.MustMarshalBinaryBare(p)
}
