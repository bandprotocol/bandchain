package proof

import (
	"math/big"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

type RequestPacketEthereum struct {
	ClientID       string
	OracleScriptID *big.Int
	Calldata       string
	AskCount       *big.Int
	MinCount       *big.Int
}

func transformRequestPacket(p oracle.OracleRequestPacketData) RequestPacketEthereum {
	return RequestPacketEthereum{
		ClientID:       p.ClientID,
		OracleScriptID: big.NewInt(int64(p.OracleScriptID)),
		Calldata:       p.Calldata,
		AskCount:       big.NewInt(int64(p.AskCount)),
		MinCount:       big.NewInt(int64(p.MinCount)),
	}
}

type ResponsePacketEthereum struct {
	ClientID      string
	RequestID     *big.Int
	AnsCount      *big.Int
	RequestTime   *big.Int
	ResolveTime   *big.Int
	ResolveStatus uint8
	Result        string
}

func transformResponsePacket(p oracle.OracleResponsePacketData) ResponsePacketEthereum {
	return ResponsePacketEthereum{
		ClientID:      p.ClientID,
		RequestID:     big.NewInt(int64(p.RequestID)),
		AnsCount:      big.NewInt(int64(p.AnsCount)),
		RequestTime:   big.NewInt(int64(p.RequestTime)),
		ResolveTime:   big.NewInt(int64(p.ResolveTime)),
		ResolveStatus: uint8(p.ResolveStatus),
		Result:        p.Result,
	}
}
