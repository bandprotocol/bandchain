package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewOracleRequestPacketData creates a new OracleRequestPacketData instance.
func NewOracleRequestPacketData(
	clientID string, oracleScriptID OracleScriptID, calldata string,
	askCount int64, minCount int64,
) OracleRequestPacketData {
	return OracleRequestPacketData{
		ClientID:       clientID,
		OracleScriptID: oracleScriptID,
		Calldata:       calldata,
		AskCount:       askCount,
		MinCount:       minCount,
	}
}

// func (p OracleRequestPacketData) String() string {
// 	return fmt.Sprintf(`OracleRequestPacketData:
//     ClientID:       %s
//     OracleScriptID: %d
//     Calldata:       %s
//     AskCount:       %d
//     MinCount:       %d`,
// 		p.ClientID,
// 		p.OracleScriptID,
// 		p.Calldata,
// 		p.AskCount,
// 		p.MinCount,
// 	)
// }

func (p OracleRequestPacketData) ValidateBasic() error {
	// TODO: Validate oracle request packet
	return nil
}

func (p OracleRequestPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}

func NewOracleResponsePacketData(
	clientID string, requestID RequestID, ansCount int64, requestTime int64, resolveTime int64,
	resolveStatus ResolveStatus, result string,
) OracleResponsePacketData {
	return OracleResponsePacketData{
		ClientID:      clientID,
		RequestID:     requestID,
		AnsCount:      ansCount,
		RequestTime:   requestTime,
		ResolveTime:   resolveTime,
		ResolveStatus: resolveStatus,
		Result:        result,
	}
}

// func (p OracleResponsePacketData) String() string {
// 	return fmt.Sprintf(`OracleResponsePacketData:
// 	ClientID: %s
// 	RequestID: %d
// 	AnsCount: %d
// 	RequestTime: %d
// 	ResolveTime: %d
// 	ResolveStatus: %d
// 	Result: %s`,
// 		p.ClientID,
// 		p.RequestID,
// 		p.AnsCount,
// 		p.RequestTime,
// 		p.ResolveTime,
// 		p.ResolveStatus,
// 		p.Result,
// 	)
// }

func (p OracleResponsePacketData) ValidateBasic() error {
	// TODO: Validate oracle request packet
	return nil
}

func (p OracleResponsePacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(p))
}
