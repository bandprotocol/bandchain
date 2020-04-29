package proof

import (
	"testing"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func TestCalculateResultHash(t *testing.T) {
	// RawByte is d9c589270a0962616e64207465737410011a1e30333030303030303432353434333634303030303030303030303030303020042804
	reqPacket := oracle.OracleRequestPacketData{
		ClientID:       "band test",
		OracleScriptID: 1,
		Calldata:       "030000004254436400000000000000",
		AskCount:       4,
		MinCount:       4,
	}
	// RawByte is 79b5957c0a0962616e6420746573741001180420f8cb8bf50528fccb8bf50530023a1064383732306230303030303030303030
	resPacket := oracle.OracleResponsePacketData{
		ClientID:      "band test",
		RequestID:     1,
		AnsCount:      4,
		RequestTime:   1587734008,
		ResolveTime:   1587734012,
		ResolveStatus: oracle.ResolveStatus(1),
		Result:        "d8720b0000000000",
	}
	expectedResultHash := hexToBytes("090dcd42c7a6729dfbe719a9cc6c78ed94a88f720420d5614e92918ff3567077")

	require.Equal(t, expectedResultHash, tmhash.Sum(
		append(
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(reqPacket)),
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(resPacket))...,
		)))
}
