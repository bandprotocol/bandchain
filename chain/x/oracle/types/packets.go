package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
