package proof

import (
	"testing"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func TestCalculateResultHash(t *testing.T) {
	// RawByte is d9c589270a0962616e64207465737410011a1e30333030303030303432353434333634303030303030303030303030303020012801
	reqPacket := oracle.OracleRequestPacketData{
		ClientID:       "band test",
		OracleScriptID: 1,
		Calldata:       "030000004254436400000000000000",
		AskCount:       1,
		MinCount:       1,
	}
	// RawByte is 79b5957c0a0962616e6420746573741006180120d3d482f50528d6d482f50530023a1032316438306130303030303030303030
	resPacket := oracle.OracleResponsePacketData{
		ClientID:      "band test",
		RequestID:     6,
		RequestTime:   1587587667,
		ResolveTime:   1587587670,
		ResolveStatus: oracle.ResolveStatus(1),
		AnsCount:      1,
		Result:        "21d80a0000000000",
	}
	expectedResultHash := hexToBytes("eed54e012643cecb617e9530fa251c92fbf1fe8db6bb5134d57efd421297b74b")

	// TODO: Find a way to encode protobuf to chain
	require.Equal(t, expectedResultHash, tmhash.Sum(
		append(
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(reqPacket)),
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(resPacket))...,
		)))
}
