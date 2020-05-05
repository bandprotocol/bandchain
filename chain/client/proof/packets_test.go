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
	// RawByte is 79b5957c0a0962616e6420746573741001180420f8cb8bf50528fccb8bf50530013a1064383732306230303030303030303030
	resPacket := oracle.OracleResponsePacketData{
		ClientID:      "band test",
		RequestID:     1,
		AnsCount:      4,
		RequestTime:   1587734008,
		ResolveTime:   1587734012,
		ResolveStatus: oracle.ResolveStatus(1),
		Result:        "d8720b0000000000",
	}
	expectedResultHash := hexToBytes("63d30f34c4b3439a95386912ec9ee2e9c01666685b6a25b11c96d46d47f37a42")

	require.Equal(t, expectedResultHash, tmhash.Sum(
		append(
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(reqPacket)),
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(resPacket))...,
		)))
}

func TestEmptyClientID(t *testing.T) {
	// RawByte is d9c5892710011a20303430303030303034323431346534343430343230663030303030303030303020042804
	reqPacket := oracle.OracleRequestPacketData{
		ClientID:       "",
		OracleScriptID: 1,
		Calldata:       "0400000042414e4440420f0000000000",
		AskCount:       4,
		MinCount:       4,
	}
	// RawByte is 79b5957c1008180420f0d0b0f50528f4d0b0f50530013a1063653730313330303030303030303030
	resPacket := oracle.OracleResponsePacketData{
		ClientID:      "",
		RequestID:     8,
		AnsCount:      4,
		RequestTime:   1588340848,
		ResolveTime:   1588340852,
		ResolveStatus: oracle.ResolveStatus(1),
		Result:        "ce70130000000000",
	}
	expectedResultHash := hexToBytes("d1f7ae1f21a04d5500b30ccfffe58a10f2376790edd617e563c8607cea0cd1c5")

	require.Equal(t, expectedResultHash, tmhash.Sum(
		append(
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(reqPacket)),
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(resPacket))...,
		)))
}
