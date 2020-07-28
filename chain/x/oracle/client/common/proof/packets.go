package proof

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type RequestPacketEthereum struct {
	ClientId       string
	OracleScriptId uint64
	Params         []byte
	AskCount       uint64
	MinCount       uint64
}

func transformRequestPacket(p types.OracleRequestPacketData) RequestPacketEthereum {
	return RequestPacketEthereum{
		ClientId:       p.ClientID,
		OracleScriptId: uint64(p.OracleScriptID),
		Params:         p.Calldata,
		AskCount:       uint64(p.AskCount),
		MinCount:       uint64(p.MinCount),
	}
}

type ResponsePacketEthereum struct {
	ClientId      string
	RequestId     uint64
	AnsCount      uint64
	RequestTime   uint64
	ResolveTime   uint64
	ResolveStatus uint8
	Result        []byte
}

func transformResponsePacket(p types.OracleResponsePacketData) ResponsePacketEthereum {
	return ResponsePacketEthereum{
		ClientId:      p.ClientID,
		RequestId:     uint64(p.RequestID),
		AnsCount:      uint64(p.AnsCount),
		RequestTime:   uint64(p.RequestTime),
		ResolveTime:   uint64(p.ResolveTime),
		ResolveStatus: uint8(p.ResolveStatus),
		Result:        p.Result,
	}
}
